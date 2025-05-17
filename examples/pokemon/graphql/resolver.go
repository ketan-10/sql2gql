package graphql

import (
	"github.com/ketan-10/sql2gql/examples/pokemon/graphql/gen"
	"github.com/ketan-10/sql2gql/examples/pokemon/sql2gql"
)


type Resolver struct {
	sql2gql.GoResolver
}

func (r *Resolver) Query() gen.QueryResolver {
	return r
}

func (r *Resolver) Mutation() gen.MutationResolver {
	return r
}
