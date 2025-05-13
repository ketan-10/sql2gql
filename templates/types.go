package templates

type TemplateType uint16

const (
	ENUM TemplateType = iota
	TABLE
	REPO
	GO_WIRE
	GRAPH_SCHEMA
	RLTS
	GQLGEN
	ENUM_SCALAR
	GO_RESOLVER
)

func (tt *TemplateType) String() string {
	switch *tt {
	case ENUM:
		return "enum"
	case TABLE:
		return "table"
	case REPO:
		return "repo"
	case GO_WIRE:
		return "go_wire"
	case RLTS:
		return "rlts"
	case GRAPH_SCHEMA:
		return "schema"
	case GQLGEN:
		return "gqlgen"
	case ENUM_SCALAR:
		return "scalar"
	case GO_RESOLVER:
		return "go_resolver"
	}

	return ""
}

func (tt *TemplateType) Extension() string {
	switch *tt {
	case ENUM, TABLE, REPO, GO_WIRE, RLTS, GO_RESOLVER:
		return "go"
	case GRAPH_SCHEMA, ENUM_SCALAR:
		return "graphql"
	case GQLGEN:
		return "yml"
	}

	return ""
}
func (tt *TemplateType) PlaceAtRoot() bool {
	switch *tt {
	case GO_WIRE, GQLGEN, ENUM_SCALAR, GO_RESOLVER:
		return true
	}

	return false
}
