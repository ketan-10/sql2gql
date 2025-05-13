package internal

import (
	"bytes"
	"database/sql"
	"text/template"

	"github.com/ketan-10/mysql2graphql/templates"
)

type Args struct {
	// DBC is database connection string
	DBC          string
	DB           *sql.DB
	Loader       ILoader
	DatabaseType DatabaseType
	DatabaseName string
	GeneratedDir string
	Generated    []*GeneratedTemplate
}
func GetDefaultArgs() *Args {
	return &Args{
		GeneratedDir: "mysql2gql",
	}
}

type GeneratedTemplate struct {
	TemplateType templates.TemplateType
	FileName     string
	Buffer       *bytes.Buffer
}

func (arg *Args) ExecuteTemplate(tt templates.TemplateType, fileName string, obj interface{}) error {

	genTmp := &GeneratedTemplate{
		TemplateType: tt,
		FileName:     fileName,
		Buffer:       new(bytes.Buffer),
	}

	// read template file
	templateFileLocation := arg.DatabaseType.String() + "/" + tt.String() + "." + tt.Extension() + ".tpl"
	fileContent, err := templates.TemplatesFS.ReadFile(templateFileLocation)
	if err != nil {
		return err
	}

	t, err := template.
		New(templateFileLocation).
		Funcs(template.FuncMap(templates.HelperFunc)).
		Parse(string(fileContent))
	if err != nil {
		return err
	}
	err = t.Execute(genTmp.Buffer, obj)
	if err != nil {
		return err
	}
	// save the generated buffer
	arg.Generated = append(arg.Generated, genTmp)
	return nil
}
