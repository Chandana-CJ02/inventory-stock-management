// Package migrations embeds the numbered SQL migration files
// (NNNN_name.up.sql / NNNN_name.down.sql) so a single binary can migrate any
// database without external files.
package migrations

import "embed"

// FS holds every *.sql file in this directory.
//
//go:embed *.sql
var FS embed.FS
