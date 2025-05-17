package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql" // empty import, to load drivers
	"github.com/ketan-10/sql2gql/examples/pokemon/graphql/gen"
	"github.com/ketan-10/sql2gql/examples/pokemon/internal/context_manager"
	"github.com/ketan-10/sql2gql/examples/pokemon/wire_app"
)

func main() {
	port := "8080"

	var connectionString string

	flag.StringVar(&connectionString, "connection", "", "database connection string")
	flag.Parse()
	if connectionString == "" {
		fmt.Fprintln(os.Stderr, "Error: --connection is required")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context_manager.WithConnection(context.Background(), connectionString)

	app, _, err := wire_app.GetApp(ctx)

	if err != nil {
		panic(err)
	}

	c := gen.Config{Resolvers: app.Resolver}

	srv := handler.NewDefaultServer(gen.NewExecutableSchema(c))

	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
