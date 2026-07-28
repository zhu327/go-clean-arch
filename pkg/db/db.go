package db

import (
	"context"
	"database/sql"
	"fmt"

	"go-clean-arch/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBTimezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}
	if err := NewRunner(sqlMigrationFS(), migrationPath()).Up(context.Background(), sqlMigrationDB{sqlDB}); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

type sqlMigrationDB struct {
	db *sql.DB
}

func (db sqlMigrationDB) BeginTx(ctx context.Context) (MigrationTx, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return sqlMigrationTx{tx: tx}, nil
}

type sqlMigrationTx struct {
	tx *sql.Tx
}

func (tx sqlMigrationTx) ExecContext(ctx context.Context, query string, args ...any) error {
	_, err := tx.tx.ExecContext(ctx, query, args...)
	return err
}

func (tx sqlMigrationTx) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx sqlMigrationTx) Commit() error   { return tx.tx.Commit() }
func (tx sqlMigrationTx) Rollback() error { return tx.tx.Rollback() }
