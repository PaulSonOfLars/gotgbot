package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func snakeToTitle(s string) string {
	bd := strings.Builder{}

	for _, split := range strings.Split(s, "_") {
		bd.WriteString(strings.Title(split))
	}

	return bd.String()
}

func snakeToCamel(s string) string {
	title := snakeToTitle(s)

	return strings.ToLower(title[:1]) + title[1:]
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func titleToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func toGoType(s string) string {
	pref := ""
	for isTgArray(s) {
		pref += "[]"
		s = strings.TrimPrefix(s, "Array of ")
	}

	switch s {
	case tgTypeInteger:
		return pref + "int64"
	case tgTypeFloat:
		return pref + "float64"
	case tgTypeBoolean:
		return pref + "bool"
	case tgTypeString:
		return pref + "string"
	}

	return pref + s
}

func isTgArray(s string) bool {
	return strings.HasPrefix(s, "Array of ")
}

func isArray(s string) bool {
	return strings.HasPrefix(s, "[]")
}

func getDefaultReturnVal(s string) string {
	if strings.HasPrefix(s, "*") || strings.HasPrefix(s, "[]") {
		return "nil"
	}

	switch s {
	case "int64":
		return "0"
	case "float64":
		return "0.0"
	case "bool":
		return "false"
	case "string":
		return "\"\""
	}

	// this isnt great
	return s
}

func goTypeStringer(t string) string {
	switch t {
	case "int64":
		return "strconv.FormatInt(%s, 10)"
	case "float64":
		return "strconv.FormatFloat(%s, 'f', -1, 64)"
	case "bool":
		return "strconv.FormatBool(%s)"
	case "string":
		return "%s"
	default:
		return ""
	}
}

func getAllFields(types []TypeDescription, parentType string) []Field {
	if len(types) == 0 {
		return nil
	}

	var fields []Field
	isOK := map[string][]string{}

	for _, t := range types {
		for _, f := range t.Fields {
			isOK[f.Name] = append(isOK[f.Name], t.getTypeNameFromParent(parentType))

			if len(isOK[f.Name]) == 1 {
				fields = append(fields, f)
			}
		}
	}

	for idx, f := range fields {
		typesUsingField := isOK[f.Name]
		if len(typesUsingField) == len(types) {
			continue
		}

		// If not all subtypes use it, then its optional; update description.
		if f.Required {
			f.Description = "Optional. " + f.Description
		}

		fields[idx].Required = false
		fields[idx].Description = fmt.Sprintf("%s (Only for %s)", f.Description, strings.Join(typesUsingField, ", "))
	}

	return fields
}

func getCommonFields(types []TypeDescription) []Field {
	if len(types) == 0 {
		return nil
	}

	count := map[string]int{}

	for _, t := range types {
		for _, f := range t.Fields {
			if !f.Required {
				continue
			}

			count[f.Name]++
		}
	}

	var fields []Field

	// only need to iterate on first, since guaranteed overlap
	for _, f := range types[0].Fields {
		if count[f.Name] == len(types) {
			fields = append(fields, f)
		}
	}

	return fields
}

func getReplyMarkupTypes(d APIDescription) []TypeDescription {
	typesMap := map[string]struct{}{}
	for _, m := range d.Methods {
		for _, f := range m.Fields {
			if f.Name == "reply_markup" {
				for _, t := range f.Types {
					typesMap[t] = struct{}{}
				}
			}
		}
	}

	var typeNames []string
	for t := range typesMap {
		typeNames = append(typeNames, t)
	}
	sort.Strings(typeNames)

	var types []TypeDescription
	for _, t := range typeNames {
		types = append(types, d.Types[t])
	}

	return types
}

func getTypeByName(d APIDescription, typeName string) (TypeDescription, error) {
	t, ok := d.Types[typeName]
	if !ok {
		return t, fmt.Errorf("unknown typename %s", typeName)
	}
	return t, nil
}

func getTypesByName(d APIDescription, typeNames []string) ([]TypeDescription, error) {
	var types []TypeDescription

	for _, typeName := range typeNames {
		t, err := getTypeByName(d, typeName)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	return types, nil
}
