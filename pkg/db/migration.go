package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strconv"
)

var migrationFileName = regexp.MustCompile(`^(\d+)_.+\.up\.sql$`)

// MigrationDB begins transactions for schema migrations.
type MigrationDB interface {
	BeginTx(context.Context) (MigrationTx, error)
}

// MigrationTx executes and records schema migrations atomically.
type MigrationTx interface {
	ExecContext(context.Context, string, ...any) error
	QueryContext(context.Context, string, ...any) (Rows, error)
	Commit() error
	Rollback() error
}

// Rows is the subset of SQL rows used by the migration runner.
type Rows interface {
	Next() bool
	Scan(...any) error
	Err() error
	Close() error
}

// Runner applies versioned SQL migrations from an fs.FS.
type Runner struct {
	fs   fs.FS
	path string
}

// NewRunner creates a migration runner for migrations under path.
func NewRunner(migrations fs.FS, path string) Runner {
	return Runner{fs: migrations, path: path}
}

type migration struct {
	version  int64
	name     string
	checksum string
	sql      string
}

type appliedMigration struct {
	version  int64
	name     string
	checksum string
}

// Up applies every migration that is not recorded in schema_migrations.
func (r Runner) Up(ctx context.Context, db MigrationDB) (err error) {
	migrations, err := r.load()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(824631927)`); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	if err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT PRIMARY KEY, name TEXT NOT NULL, checksum TEXT NOT NULL, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	if err = tx.ExecContext(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS name TEXT`); err != nil {
		return fmt.Errorf("add migration name: %w", err)
	}
	if err = tx.ExecContext(ctx, `ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum TEXT`); err != nil {
		return fmt.Errorf("add migration checksum: %w", err)
	}
	applied, err := appliedMigrations(ctx, tx)
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if previous, exists := applied[migration.version]; exists {
			if migration.version == 1 && (previous.name == "" || previous.checksum == "") &&
				(previous.name == "" || previous.name == migration.name) &&
				(previous.checksum == "" || previous.checksum == migration.checksum) {
				if err = verifyLegacyUsersSchema(ctx, tx); err != nil {
					return err
				}
				if err = tx.ExecContext(ctx, `UPDATE schema_migrations SET name = $2, checksum = $3 WHERE version = $1 AND (name IS NULL OR checksum IS NULL)`, migration.version, migration.name, migration.checksum); err != nil {
					return fmt.Errorf("baseline migration %s metadata: %w", migration.name, err)
				}
				continue
			}
			if previous.name != migration.name || previous.checksum != migration.checksum {
				return fmt.Errorf("migration %s checksum mismatch", migration.name)
			}
			continue
		}
		if migration.version == 1 {
			legacyUsers, checkErr := usersTableExists(ctx, tx)
			if checkErr != nil {
				return checkErr
			}
			if legacyUsers {
				if err = verifyLegacyUsersSchema(ctx, tx); err != nil {
					return err
				}
				if err = tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, name, checksum) VALUES ($1, $2, $3)`, migration.version, migration.name, migration.checksum); err != nil {
					return fmt.Errorf("record legacy baseline %s: %w", migration.name, err)
				}
				continue
			}
		}
		if err = tx.ExecContext(ctx, migration.sql); err != nil {
			return fmt.Errorf("apply migration %s: %w", migration.name, err)
		}
		if err = tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, name, checksum) VALUES ($1, $2, $3)`, migration.version, migration.name, migration.checksum); err != nil {
			return fmt.Errorf("record migration %s: %w", migration.name, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit migrations: %w", err)
	}
	return nil
}

func (r Runner) load() ([]migration, error) {
	entries, err := fs.ReadDir(r.fs, r.path)
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}
	migrations := make([]migration, 0, len(entries))
	versions := make(map[int64]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := migrationFileName.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration version %q: %w", entry.Name(), err)
		}
		if previous, exists := versions[version]; exists {
			return nil, fmt.Errorf("duplicate migration version %d: %s and %s", version, previous, entry.Name())
		}
		contents, err := fs.ReadFile(r.fs, path.Join(r.path, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		versions[version] = entry.Name()
		digest := sha256.Sum256(contents)
		migrations = append(
			migrations,
			migration{
				version:  version,
				name:     entry.Name(),
				checksum: hex.EncodeToString(digest[:]),
				sql:      string(contents),
			},
		)
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	return migrations, nil
}

func usersTableExists(ctx context.Context, tx MigrationTx) (bool, error) {
	return queryBoolean(
		ctx,
		tx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'users')`,
	)
}

func verifyLegacyUsersSchema(ctx context.Context, tx MigrationTx) error {
	emailUnique, err := queryBoolean(
		ctx,
		tx,
		`SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'users'::regclass AND contype = 'u' AND conname = 'users_email_key')`,
	)
	if err != nil {
		return fmt.Errorf("verify legacy users email constraint: %w", err)
	}
	usernameUnique, err := queryBoolean(
		ctx,
		tx,
		`SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'users'::regclass AND contype = 'u' AND conname = 'users_username_key')`,
	)
	if err != nil {
		return fmt.Errorf("verify legacy users username constraint: %w", err)
	}
	if !emailUnique || !usernameUnique {
		return fmt.Errorf("legacy users schema does not have required email and username unique constraints")
	}
	return nil
}

func queryBoolean(ctx context.Context, tx MigrationTx, query string) (bool, error) {
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		return false, fmt.Errorf("query returned no rows")
	}
	var value bool
	if err := rows.Scan(&value); err != nil {
		return false, err
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	return value, nil
}

func appliedMigrations(ctx context.Context, tx MigrationTx) (map[int64]appliedMigration, error) {
	rows, err := tx.QueryContext(
		ctx,
		`SELECT version, COALESCE(name, ''), COALESCE(checksum, '') FROM schema_migrations`,
	)
	if err != nil {
		return nil, fmt.Errorf("query applied migrations: %w", err)
	}
	defer func() { _ = rows.Close() }()
	migrations := make(map[int64]appliedMigration)
	for rows.Next() {
		var migration appliedMigration
		if err := rows.Scan(&migration.version, &migration.name, &migration.checksum); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		migrations[migration.version] = migration
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}
	return migrations, nil
}

func sqlMigrationFS() fs.FS { return migrationsFS }

func migrationPath() string { return "migrations" }
