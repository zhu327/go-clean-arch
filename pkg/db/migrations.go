package db

import "embed"

//go:embed migrations/*.up.sql
var migrationsFS embed.FS
