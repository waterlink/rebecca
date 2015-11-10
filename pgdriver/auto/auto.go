// Package auto allows one to setup rebecca/pgdriver in a convenient way from
// environment variables.
//
// To start using auto use this nameless import:
//
// 				import _ "github.com/waterlink/rebecca/pgdriver/auto"
//
// Empty import will fetch connection options from respective environment
// variables:
//
// - `REBECCA_PG_USER`
// - `REBECCA_PG_PASS`
// - `REBECCA_PG_HOST`
// - `REBECCA_PG_PORT`
// - `REBECCA_PG_DATABASE`
// - `REBECCA_PG_SSLMODE`
//
// or, if present: `REBECCA_PG_URL`
package auto

import (
	"fmt"
	"os"

	"github.com/waterlink/rebecca"
	"github.com/waterlink/rebecca/pgdriver"
)

func init() {
	d := pgdriver.NewDriver(pgURL())
	rebecca.SetupDriver(d)
}

func pgURL() string {
	if pgURL := os.Getenv("REBECCA_PG_URL"); pgURL != "" {
		return pgURL
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getenv("REBECCA_PG_USER", "postgres"),
		getenv("REBECCA_PG_PASS", ""),
		getenv("REBECCA_PG_HOST", "127.0.0.1"),
		getenv("REBECCA_PG_PORT", "5432"),
		getenv("REBECCA_PG_DATABASE", "postgres"),
		getenv("REBECCA_PG_SSLMODE", "disable"),
	)
}

func getenv(name, def string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return def
}
