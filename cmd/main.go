package cmd

import (
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ketan-10/mysql2graphql/internal"

	// empty import, so that init method will be called, and drivers will be loaded, this has driver import
	_ "github.com/ketan-10/mysql2graphql/loaders"
)

func Execute() {

	var err error
	fmt.Println("Started")

	args := internal.GetDefaultArgs()

	// validate Args
	flag.StringVar(&args.DBC, "connection", "", "database connection string")
	flag.Parse()
	if args.DBC == "" {
		fmt.Fprintln(os.Stderr, "Error: --connection is required")
		flag.Usage()
		os.Exit(1)
	}

	// open DB and save it to args
	err = openDB(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer args.DB.Close()

	// loaders
	err = getLoaderOfDriver(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = args.Loader.LoadSchema(args)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// create files

	generateFiles(args)
}

func generateFiles(args *internal.Args) {

	for _, gen := range args.Generated {

		dirName := "./" + args.GeneratedDir
		if !gen.TemplateType.PlaceAtRoot() {
			dirName += "/" + gen.TemplateType.String()
		}

		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			os.MkdirAll(dirName, 0755)
			// os.MkdirAll(dirName, os.ModeDir)
		}
		file, err := os.Create(dirName + "/" + gen.OutFileName + "." + gen.TemplateType.Extension())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.Write(gen.Buffer.Bytes())
		if err != nil {
			panic(err)
		}

	}
}

// uriToGoDSN converts connection string to a Go DSN.
// ex: mysql://bob:password@127.0.0.1:3306/training?charset=utf8mb4&parseTime=true
// -> bob:password@tcp(127.0.0.1:3306)/training?charset=utf8mb4&parseTime=true
func uriToGoDSN(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	user := u.User.Username()
	pass, _ := u.User.Password()
	host := u.Host
	dbName := strings.TrimPrefix(u.Path, "/")
	params := u.RawQuery
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		user, pass, host, dbName, params,
	)
	return dsn, nil
}

func openDB(args *internal.Args) error {
	url, err := uriToGoDSN(args.DBC)
	if err != nil {
		return err
	}

	args.DatabaseType = internal.MYSQL

	// open mysql connection
	args.DB, err = sql.Open("mysql", url)
	if err != nil {
		return err
	}
	return nil
}

func getLoaderOfDriver(args *internal.Args) error {
	var ok bool
	args.Loader, ok = internal.AllLoaders[args.DatabaseType]
	if !ok {
		return fmt.Errorf("for driver %s, no registered loader found", args.DatabaseType.String())
	}
	return nil
}
