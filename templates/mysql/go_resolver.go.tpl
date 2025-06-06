// Code generated by sql2gql. DO NOT EDIT.

package sql2gql

import (
	"github.com/google/wire"
)


type GoResolver struct {
{{ range . }}
    repo.I{{ camelCase .Table.TableName }}Repository
    rlts.I{{ camelCase .Table.TableName }}RltsRepository
{{- end }}
}

//type IGoResolver interface {
    {{- range . }}
//    {{ camelCase .Table.TableName }}() gen.{{ camelCase .Table.TableName }}Resolver
    {{- end}}
//}

var NewGoResolver = wire.NewSet(
    wire.Struct(new(GoResolver), "*"),
//    wire.Bind(new(IGoResolver), new(GoResolver)),
)

{{ range . }}
func (r *GoResolver) {{ camelCase .Table.TableName }}() gen.{{ camelCase .Table.TableName }}Resolver {
	return r.I{{ camelCase .Table.TableName }}RltsRepository
}
{{- end }}
