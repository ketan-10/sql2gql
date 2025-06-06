package internal

import (
	"bytes"
	"database/sql"
	"text/template"

	"github.com/ketan-10/sql2gql/templates"
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
		GeneratedDir: "sql2gql",
	}
}

type GeneratedTemplate struct {
	TemplateType templates.TemplateType
	OutFileName  string
	Buffer       *bytes.Buffer
}

func (arg *Args) ExecuteTemplate(tt templates.TemplateType, outFilename string, obj interface{}) error {

	genTmp := &GeneratedTemplate{
		TemplateType: tt,
		OutFileName:  outFilename,
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
