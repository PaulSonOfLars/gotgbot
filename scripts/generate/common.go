package main

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const internalTypeRef = "GOTGBOT INTERNAL"

func snakeToTitle(s string) string {
	bd := strings.Builder{}

	for split := range strings.SplitSeq(s, "_") {
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

func isCoreTGType(t string) bool {
	return t == tgTypeString || t == tgTypeBoolean || t == tgTypeInteger || t == tgTypeFloat
}

func isRichTextType(t string) bool {
	return strings.HasPrefix(t, "Rich")
}

func toGoType(s string) string {
	pref := ""
	for isTgArray(s) {
		pref += "[]"
		_, s = typeOfTgArray(s)
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

const tgArrayOfPrefix = "Array of "

func isTgArray(s string) bool {
	return strings.HasPrefix(s, tgArrayOfPrefix)
}

func typeOfTgArray(s string) (bool, string) {
	return isTgArray(s), strings.TrimPrefix(s, tgArrayOfPrefix)
}

func isPointer(s string) bool {
	return strings.HasPrefix(s, "*") ||
		s == tgTypeInputFile || s == typeInputFileOrString || s == typeInputString ||
		s == tgTypeInputMedia || s == tgTypeInputPaidMedia ||
		s == typeReplyMarkup
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

		// this isn't great
		return s + "{}"
	}
}

func getDefaultReturnVals(d APIDescription, types []string) []string {
	var retVals []string
	for _, retType := range types {
		retVals = append(retVals, getDefaultTypeVal(d, retType))
	}
	return retVals
}

// getAllFields merges all the fields from list of types.
func getAllFields(d APIDescription, types []TypeDescription, parentType string) ([]Field, error) {
	if len(types) == 0 {
		return nil, nil
	}

	constantType, err := d.Types[parentType].getConstantFieldFromParent(d)
	if err != nil {
		return nil, fmt.Errorf("error getting field %s: %w", parentType, err)
	}

	var fieldNames []string            // Ordered, unique list of names for all fields across all types
	presentIn := map[string][]string{} // Map of fields -> list of types which use it
	fields := map[string]Field{}       // Map of fieldnames -> fields to use (using correct pointers)

	for _, t := range types {
		shortName, err := t.getTypeNameFromParent(constantType.Name)
		if err != nil {
			return nil, fmt.Errorf("error getting shortname %s: %w", t.Name, err)
		}

		for _, f := range t.Fields {
			presentIn[f.Name] = append(presentIn[f.Name], shortName)

			// Some fields need to be cleaned up to be agnostic across all fields.
			// eg: 	The BotCommandScopeDefault has a "Type" field saying 'Scope type, must be default'
			// This is clearly only valid for the "default" scope; not the others - hence, clean it.
			for _, replace := range []*regexp.Regexp{
				typeAlwaysMatcher, // always "x"
				typeMustBeMatcher, // must be x
			} {
				f.Description = replace.ReplaceAllString(f.Description, "")
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
	fieldTypes := map[string]Field{}

	for _, t := range types {
		for _, f := range t.Fields {
			if !f.Required {
				continue
			}

			count[f.Name]++

			if other, ok := fieldTypes[f.Name]; !ok {
				fieldTypes[f.Name] = f
			} else if !slices.Equal(f.Types, other.Types) {
				// if not equal, make sure to just set the count to a silly value
				count[f.Name] += len(types)
			}
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

func getTypesByName(d APIDescription, parentType string, typeNames []string) ([]TypeDescription, error) {
	var types []TypeDescription

	for _, typeName := range typeNames {
		if isCoreTGType(typeName) {
			newTypeName := parentType + typeName
			d.Types[newTypeName] = newTypeDescription(newTypeName, typeName, parentType)
			typeName = newTypeName

		} else if isTgArray(typeName) {
			newTypeName := parentType + "Array"
			d.Types[newTypeName] = newTypeDescription(newTypeName, typeName, parentType)
			typeName = newTypeName
		}

		t, err := getTypeByName(d, typeName)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	return types, nil
}

func newTypeDescription(newTypeName string, typeName string, parentType string) TypeDescription {
	return TypeDescription{
		Name: newTypeName,
		Fields: []Field{{
			Name:        "type",
			Types:       []string{tgTypeString},
			Required:    true,
			Description: fmt.Sprintf("Type of the %s, always \"%s\"", parentType, strings.ToLower(typeName)),
		}},
		Href:      internalTypeRef,
		SubtypeOf: unique(typeName, parentType),
	}
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

func checkAllChildrenFieldTypes(d APIDescription, parentType string, subtypes []TypeDescription) bool {
	for _, t := range subtypes {
		if !childFieldTypesMatch(d, parentType, t.Fields) {
			return false
		}
	}
	return true
}

func childFieldTypesMatch(d APIDescription, parentType string, fields []Field) bool {
	subTypes, err := getTypesByName(d, parentType, d.Types[parentType].Subtypes)
	if err != nil {
		return false
	}

	allParentFields, err := getAllFields(d, subTypes, parentType)
	if err != nil {
		return false
	}

	for _, f := range fields {
		for _, parentField := range allParentFields {
			if parentField.Name == f.Name && !slices.Equal(f.Types, parentField.Types) {
				return false
			}
		}
	}

	return true
}

func (f *Field) ConstantName(parentName string, shortName string) string {
	return ConstantName(parentName, f.Name, shortName)
}

func ConstantName(parentName string, fieldName string, shortName string) string {
	return snakeToTitle(parentName) + snakeToTitle(fieldName) + snakeToTitle(shortName)
}

func unique(ss ...string) []string {
	seen := make(map[string]struct{})
	unique := make([]string, 0, len(ss))
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			unique = append(unique, s)
		}
	}
	return unique
}
