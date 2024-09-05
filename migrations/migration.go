package migrations

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

const driver = "postgres"

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) < 3 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	dir, dbstring, command := args[0], args[1], args[2]

	if err := goose.SetDialect(driver); err != nil {
		log.Fatal(err)
	}

	if dbstring == "" {
		log.Fatalf("-dbstring=%q not supported\n", dbstring)
	}

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", dbstring, err)
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, dir, arguments...); err != nil {
		log.Fatalf("failed run migration. %s", err)
	}
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

var (
	usagePrefix = `Usage: goose MIGRATION_DIR CONN_STRING COMMAND
Drivers:
    postgres
Examples:
    goose ./sql "user=postgres dbname=postgres sslmode=disable" status
    goose ./sql "user=postgres dbname=postgres sslmode=disable" up
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)
