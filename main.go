package main

import (
	"flag"
	"fmt"
	"github.com/eugeniypetrov/browscap-go/browscap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"reflect"
	"time"
)

const (
	CommandCompile = "compile"
	CommandFind    = "find"
)

func getStorage(storageName string, dsn string) (browscap.BrowserStorage, error) {
	switch storageName {
	case "mysql":
		db, err := sqlx.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("error opening db: %w", err)
		}
		return browscap.NewMysqlBrowserStorage(db), nil
	case "sqlite":
		db, err := sqlx.Open("sqlite3", dsn)
		if err != nil {
			return nil, fmt.Errorf("error opening db: %w", err)
		}
		return browscap.NewSqliteBrowserStorage(db), nil
	case "postgres":
		db, err := sqlx.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("error opening db: %w", err)
		}
		return browscap.NewPostgresBrowserStorage(db), nil
	default:
		return nil, fmt.Errorf("unknown storage %s", storageName)
	}
}

func compile(filename string, storageName string, dsn string) error {
	log.Println("compiling", filename)

	storage, err := getStorage(storageName, dsn)
	if err != nil {
		return fmt.Errorf("error getting storage: %w", err)
	}

	l := browscap.NewLoader(storage)

	err = l.Compile(filename)
	if err != nil {
		return fmt.Errorf("error compiling: %w", err)
	}

	return nil
}

func find(userAgent string, storageName string, dsn string) error {
	log.Println("loading", userAgent)

	storage, err := getStorage(storageName, dsn)
	if err != nil {
		return fmt.Errorf("error getting storage: %w", err)
	}

	l := browscap.NewLoader(storage)

	start := time.Now()
	bc, err := l.Load()
	if err != nil {
		return fmt.Errorf("error loading: %w", err)
	}

	log.Printf("loaded (elapsed %s)", time.Since(start))

	start = time.Now()
	browser, err := bc.GetBrowser(userAgent)
	if err != nil {
		return fmt.Errorf("error finding: %w", err)
	}
	elapsed := time.Since(start)

	b := reflect.ValueOf(browser).Elem()
	for i := 0; i < b.NumField(); i++ {
		field := b.Field(i)
		fmt.Printf(
			"%-30s %v\n",
			b.Type().Field(i).Name,
			field.Interface(),
		)
	}

	log.Printf("elapsed %s", elapsed)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case CommandCompile:
		fs := flag.NewFlagSet(CommandCompile, flag.ExitOnError)
		filename := fs.String("filename", "full_php_browscap.ini", "browscap ini file")
		storage := fs.String("storage", "sqlite", "storage (mysql, sqlite, postgres)")
		dsn := fs.String("dsn", "browscap.sqlite", "data source name")

		err := fs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("error parsing compile command. %s", err)
		}

		err = compile(*filename, *storage, *dsn)
		if err != nil {
			log.Fatalf("error compiling. %s", err)
		}
	case CommandFind:
		fs := flag.NewFlagSet(CommandFind, flag.ExitOnError)
		userAgent := fs.String("user-agent", "", "user agent")
		storage := fs.String("storage", "sqlite", "storage (mysql, sqlite, memory)")
		dsn := fs.String("dsn", "browscap.sqlite", "data source name")

		err := fs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("error parsing find command. %s", err)
		}

		err = find(*userAgent, *storage, *dsn)
		if err != nil {
			log.Fatalf("error finding. %s", err)
		}
	default:
		log.Fatalf("unexpected subcommand %s", cmd)
	}
}
