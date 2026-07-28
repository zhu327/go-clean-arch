package db

import (
	"context"
	"errors"
	"strings"
	"testing"
	"testing/fstest"
)

type fakeMigrationDB struct {
	beginErr error
	tx       *fakeMigrationTx
}

func (db *fakeMigrationDB) BeginTx(context.Context) (MigrationTx, error) {
	if db.beginErr != nil {
		return nil, db.beginErr
	}
	return db.tx, nil
}

type fakeMigrationTx struct {
	executed       []string
	args           [][]any
	applied        []appliedMigration
	usersExist     bool
	emailUnique    bool
	usernameUnique bool
	commit         bool
	rollback       bool
}

func (tx *fakeMigrationTx) ExecContext(_ context.Context, query string, args ...any) error {
	tx.executed = append(tx.executed, query)
	tx.args = append(tx.args, args)
	return nil
}

func (tx *fakeMigrationTx) QueryContext(_ context.Context, query string, _ ...any) (Rows, error) {
	if strings.Contains(query, "information_schema.tables") {
		return &fakeBoolRows{value: tx.usersExist}, nil
	}
	if strings.Contains(query, "users_email_key") {
		return &fakeBoolRows{value: tx.emailUnique}, nil
	}
	if strings.Contains(query, "users_username_key") {
		return &fakeBoolRows{value: tx.usernameUnique}, nil
	}
	return &fakeRows{migrations: tx.applied}, nil
}

func (tx *fakeMigrationTx) Commit() error {
	tx.commit = true
	return nil
}

func (tx *fakeMigrationTx) Rollback() error {
	tx.rollback = true
	return nil
}

type fakeRows struct {
	migrations []appliedMigration
	index      int
}

func (r *fakeRows) Next() bool { return r.index < len(r.migrations) }
func (r *fakeRows) Scan(dest ...any) error {
	if len(dest) != 3 {
		return errors.New("expected version, name, checksum destinations")
	}
	migration := r.migrations[r.index]
	r.index++
	*(dest[0].(*int64)) = migration.version
	*(dest[1].(*string)) = migration.name
	*(dest[2].(*string)) = migration.checksum
	return nil
}
func (r *fakeRows) Err() error   { return nil }
func (r *fakeRows) Close() error { return nil }

type fakeBoolRows struct {
	value bool
	read  bool
}

func (r *fakeBoolRows) Next() bool { return !r.read }
func (r *fakeBoolRows) Scan(dest ...any) error {
	r.read = true
	*(dest[0].(*bool)) = r.value
	return nil
}
func (r *fakeBoolRows) Err() error   { return nil }
func (r *fakeBoolRows) Close() error { return nil }

func TestRunnerUpAppliesPendingMigrationsInVersionOrderUnderAdvisoryLock(t *testing.T) {
	files := fstest.MapFS{
		"migrations/002_add_usernames.up.sql": &fstest.MapFile{
			Data: []byte("ALTER TABLE users ADD COLUMN username TEXT;"),
		},
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
		"migrations/ignore.sql": &fstest.MapFile{Data: []byte("ignored")},
	}
	tx := &fakeMigrationTx{}
	runner := NewRunner(files, "migrations")

	if err := runner.Up(context.Background(), &fakeMigrationDB{tx: tx}); err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	if !tx.commit {
		t.Fatal("expected migration transaction to commit")
	}
	if !strings.Contains(tx.executed[0], "pg_advisory_xact_lock") {
		t.Fatalf("first statement = %q, want transaction advisory lock", tx.executed[0])
	}
	if tx.executed[4] != "CREATE TABLE users (id BIGSERIAL PRIMARY KEY);" {
		t.Fatalf("first migration SQL = %q", tx.executed[4])
	}
	if got := tx.args[5]; len(got) != 3 || got[0] != int64(1) || got[1] != "001_create_users.up.sql" || got[2] == "" {
		t.Fatalf("first migration metadata = %#v", got)
	}
}

func TestRunnerUpRejectsAppliedMigrationChecksumMismatch(t *testing.T) {
	files := fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
	}
	tx := &fakeMigrationTx{
		applied: []appliedMigration{{version: 1, name: "001_create_users.up.sql", checksum: "changed"}},
	}

	err := NewRunner(files, "migrations").Up(context.Background(), &fakeMigrationDB{tx: tx})
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("Up error = %v, want checksum mismatch", err)
	}
	if !tx.rollback {
		t.Fatal("expected rollback after checksum mismatch")
	}
}

func TestRunnerUpAdoptsVerifiedLegacyUsersSchema(t *testing.T) {
	files := fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
	}
	tx := &fakeMigrationTx{usersExist: true, emailUnique: true, usernameUnique: true}

	if err := NewRunner(files, "migrations").Up(context.Background(), &fakeMigrationDB{tx: tx}); err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	for _, statement := range tx.executed {
		if strings.Contains(statement, "CREATE TABLE users") {
			t.Fatalf("legacy users schema must be baselined rather than recreated: %q", statement)
		}
	}
	if len(tx.args) == 0 || len(tx.args[len(tx.args)-1]) != 3 || tx.args[len(tx.args)-1][0] != int64(1) {
		t.Fatalf("expected atomic baseline metadata insert, got %#v", tx.args)
	}
}

func TestRunnerUpRejectsUnverifiedLegacyUsersSchema(t *testing.T) {
	files := fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
	}
	tx := &fakeMigrationTx{usersExist: true, emailUnique: true}
	err := NewRunner(files, "migrations").Up(context.Background(), &fakeMigrationDB{tx: tx})
	if err == nil || !strings.Contains(err.Error(), "legacy users schema") {
		t.Fatalf("Up error = %v, want legacy schema validation failure", err)
	}
	if !tx.rollback {
		t.Fatal("expected rollback")
	}
}

func TestRunnerUpBaselinesLegacyMetadataWithoutNameOrChecksum(t *testing.T) {
	files := fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
	}
	tx := &fakeMigrationTx{
		applied:        []appliedMigration{{version: 1}},
		usersExist:     true,
		emailUnique:    true,
		usernameUnique: true,
	}
	if err := NewRunner(files, "migrations").Up(context.Background(), &fakeMigrationDB{tx: tx}); err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	found := false
	for _, statement := range tx.executed {
		found = found || strings.Contains(statement, "UPDATE schema_migrations")
	}
	if !found {
		t.Fatal("expected controlled metadata baseline update")
	}
}

func TestRunnerUpBaselinesLegacyMetadataWhenOneFieldIsMissing(t *testing.T) {
	files := fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
	}
	tx := &fakeMigrationTx{
		applied:        []appliedMigration{{version: 1, name: "001_create_users.up.sql"}},
		usersExist:     true,
		emailUnique:    true,
		usernameUnique: true,
	}
	if err := NewRunner(files, "migrations").Up(context.Background(), &fakeMigrationDB{tx: tx}); err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
}

func TestRunnerUpRejectsDuplicateMigrationVersions(t *testing.T) {
	runner := NewRunner(fstest.MapFS{
		"migrations/001_create_users.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE users (id BIGSERIAL PRIMARY KEY);"),
		},
		"migrations/001_create_accounts.up.sql": &fstest.MapFile{
			Data: []byte("CREATE TABLE accounts (id BIGSERIAL PRIMARY KEY);"),
		},
	}, "migrations")

	err := runner.Up(context.Background(), &fakeMigrationDB{tx: &fakeMigrationTx{}})
	if err == nil {
		t.Fatal("expected duplicate version error")
	}
}
