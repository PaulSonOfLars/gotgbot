package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func contains(s string, ss []string) bool {
	for _, str := range ss {
		if s == str {
			return true
		}
	}
	return false
}

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

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func titleToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func titleToCamelCase(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

var tgToGoTypeMap = map[string]string{
	tgTypeInteger: "int64",
	tgTypeFloat:   "float64",
	tgTypeBoolean: "bool",
	tgTypeString:  "string",
}

func toGoType(s string) string {
	pref := ""
	for isTgArray(s) {
		pref += "[]"
		s = strings.TrimPrefix(s, "Array of ")
	}

	if t, ok := tgToGoTypeMap[s]; ok {
		return pref + t
	}
	return pref + s
}

func stripPointersAndArrays(retType string) string {
	for isPointer(retType) {
		retType = strings.TrimPrefix(retType, "*")
	}

	for isArray(retType) {
		retType = strings.TrimPrefix(retType, "[]")
	}
	return retType
}

func isTgArray(s string) bool {
	return strings.HasPrefix(s, "Array of ")
}

func isPointer(s string) bool {
	return strings.HasPrefix(s, "*")
}

func isArray(s string) bool {
	return strings.HasPrefix(s, "[]")
}

func getDefaultTypeVal(d APIDescription, s string) string {
	if isPointer(s) || isArray(s) {
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
	default:
		if _, ok := d.Types[s]; ok {
			return "nil"
		}

		// this isnt great
		return s
	}
}

func getDefaultReturnVals(d APIDescription, types []string) []string {
	var retVals []string
	for _, retType := range types {
		retVals = append(retVals, getDefaultTypeVal(d, retType))
	}
	return retVals
}

// goTypeStringer provides us with the fmt strings that allow us to convert basic types into strings.
// For example, we define how to handle ints, bools, and strings.
// More complicated types should be handled separately.
func goTypeStringer(t string) string {
	switch t {
	case "int64":
		return "strconv.FormatInt(%s, 10)"
	case "*int64":
		return "strconv.FormatInt(*%s, 10)"
	case "float64":
		return "strconv.FormatFloat(%s, 'f', -1, 64)"
	case "bool":
		return "strconv.FormatBool(%s)"
	case "*bool":
		return "strconv.FormatBool(*%s)"
	case "string":
		return "%s"
	case "*string":
		return "*%s"
	default:
		return ""
	}
}

// getAllFields merges all the fields from list of types.
func getAllFields(d APIDescription, types []TypeDescription, parentType string) ([]Field, error) {
	if len(types) == 0 {
		return nil, nil
	}

	var fieldNames []string            // Ordered, unique list of names for all fields across all types
	presentIn := map[string][]string{} // Map of fields -> list of types which use it
	fields := map[string]Field{}       // Map of fieldnames -> fields to use (using correct pointers)

	for _, t := range types {
		for _, f := range t.Fields {
			typeName := t.getTypeNameFromParent(parentType)
			presentIn[f.Name] = append(presentIn[f.Name], typeName)

			// Some fields need to be cleaned up to be agnostic across all fields.
			// eg: 	The BotCommandScopeDefault has a "Type" field saying 'Scope type, must be default'
			// This is clearly only valid for the "default" scope; not the others - hence, clean it.
			for _, replace := range []string{
				fmt.Sprintf(", must be %s", typeName),
				fmt.Sprintf(", always \"%s\"", typeName),
			} {
				f.Description = strings.Replace(f.Description, replace, "", 1)
			}

			// If this is the first time a field is found, add it to the list of field names
			if len(presentIn[f.Name]) == 1 {
				fieldNames = append(fieldNames, f.Name)
				fields[f.Name] = f
			}

			// Check the current type info for this field, to handle cases where one is a pointer
			// (eg "User" and "*User" in ChatBoostSourceGiveaway)
			currField := fields[f.Name]
			newType, err := f.getPreferredType(d)
			if err != nil {
				return nil, err
			}

			currType, err := currField.getPreferredType(d)
			if err != nil {
				return nil, err
			}

			if currType != newType {
				if isPointer(newType) {
					fields[f.Name] = f
				}
			}
		}
	}

	var retField []Field
	for _, n := range fieldNames {
		f, ok := fields[n]
		if !ok {
			return nil, errors.New("missing types for " + n)
		}

		typesUsingField := presentIn[f.Name]
		if len(typesUsingField) != len(types) {
			// If not all subtypes use it, then its optional; update description.
			if f.Required {
				f.Description = "Optional. " + f.Description
			}

			f.Required = false
			f.Description = fmt.Sprintf("%s (Only for %s)", f.Description, strings.Join(typesUsingField, ", "))
		} else if parentType == "ChatBoostSource" && f.Name == "user" {
			// Special edge case for odd typing situation
			f.Description = "Optional. User that provided the boost (may be empty for ChatBoostSourceGiveaway)"
		}

		retField = append(retField, f)
	}

	return retField, nil
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

// getFieldNames turns a list of fields into a list of the field's names, as described.
func getFieldNames(fs []Field) (out []string) {
	for _, t := range fs {
		out = append(out, t.Name)
	}
	return out
}

// getReplyMarkupTypes gets all the different types which are used in "reply_markup" fields.
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

// extractQuotedValues is a very basic quote extraction method. It only works on normal double quotes ("), it does not
// handle any sort of escaping, and it expects an even number of quotes to function.
// But that's all we need for this package, and so we're happy.
func extractQuotedValues(in string) ([]string, error) {
	if strings.Count(in, "\"")%2 != 0 {
		return nil, errors.New("uneven number of quotes in string")
	}

	var out []string
	startIdx := -1
	for idx, r := range in {
		if r == '"' {
			if startIdx == -1 {
				// This is an opening quote
				startIdx = idx
				continue
			}

			// This is a closing quote, so append to outputs and reset startIdx.
			out = append(out, in[startIdx+1:idx])
			startIdx = -1
		}
	}
	return out, nil
}
