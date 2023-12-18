package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"log"
)

const dialect = "postgres"

var (
	uri   = ""
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "./migrations", "directory with migration files")
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate db",
	Run: func(cmd *cobra.Command, args []string) {
		flags.Usage = usage

		if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
			flags.Usage()

			return
		}

		flag.StringVar(&uri, "db_dsn", "", "database connection-string.")
		flag.Parse()

		command := args[0]
		switch command {
		case "create":
			if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
				log.Fatalf("migrate run: %v", err)
			}
			return
		case "fix":
			if err := goose.Run("fix", nil, *dir); err != nil {
				log.Fatalf("migrate run: %v", err)
			}
			return
		}

		db, err := sql.Open("postgres", uri)
		if err != nil {
			panic(fmt.Errorf("sql.Open %w", err))
		}

		if err = db.Ping(); err != nil {
			panic(fmt.Errorf("db.Ping %w", err))
		}

		if err = goose.SetDialect(dialect); err != nil {
			log.Fatal(err)
		}

		if err = goose.Run(command, db, *dir, args[1:]...); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
	},
}

func usage() {
	fmt.Print(usagePrefix)
	flags.PrintDefaults()
	fmt.Print(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                   Apply sequential ordering to migrations
`
)
