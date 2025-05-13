package templates

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/ketan-10/mysql2graphql/utils"
)

var removeSpecialChar = regexp.MustCompile(`[^a-zA-Z0-9]`)

var HelperFunc template.FuncMap = template.FuncMap{
	"camelCase": func(input string) string {
		return utils.SnakeToCamel(removeSpecialChar.ReplaceAllLiteralString(input, "_"))
	},
	"camelCaseVar": func(input string) string {
		return utils.LowCaseFirst(utils.SnakeToCamel((removeSpecialChar.ReplaceAllLiteralString(input, "_"))))
	},
	"joinWith":      joinWith,
	"shortName":     shortName,
	"convertToNull": convertToNull,
}

func joinWith(with string, values ...string) string {
	return strings.Join(values, with)
}

// fetch only uppercase values
func shortName(name string) string {
	short := strings.ToLower(strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r
		}
		return -1
	}, name))

	if len(short) == 0 {
		short = strings.ToLower(name[0:1])
	}
	if short == "er" || short == "err" || short == "va" || short == "var" {
		short = short + "_"
	}
	return short
}

func convertToNull(name string, typ string) string {
	switch typ {
	case "int":
		return "sql.NullInt64{Valid: true, Int64: int64(" + name + ")}"
	case "string":
		return "sql.NullString{Valid: true, String: " + name + "}"
	}
	return name
}
