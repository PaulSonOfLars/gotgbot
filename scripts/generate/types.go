package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

var (
	inputParamsTmpl           = template.Must(template.New("inputParamsMethod").Parse(inputParamsMethod))
	customMarshalTmpl         = template.Must(template.New("customMarshal").Parse(customMarshal))
	customUnmarshalTmpl       = template.Must(template.New("customUnmarshal").Parse(customUnmarshal))
	customStructUnmarshalTmpl = template.Must(template.New("customStructUnmarshal").Parse(customStructUnmarshal))
)

func generateTypes(d APIDescription) error {
	file := strings.Builder{}
	file.WriteString(`
// THIS FILE IS AUTOGENERATED. DO NOT EDIT.
// Regen by running 'go generate' in the repo root.

package gotgbot

import (
	"encoding/json"
	"fmt"
	"io"
)
`)

	// the reply_markup field is weird; this allows it to support multiple types.
	replyMarkupInterface, err := generateGenericInterfaceType(d, tgTypeReplyMarkup, getReplyMarkupTypes(d))
	if err != nil {
		return fmt.Errorf("failed to generate reply_markup interface: %w", err)
	}

	file.WriteString(replyMarkupInterface)

	for _, tgTypeName := range orderedTgTypes(d) {
		tgType := d.Types[tgTypeName]

		typeDef, err := generateTypeDef(d, tgType)
		if err != nil {
			return fmt.Errorf("failed to generate type definition of %s: %w", tgTypeName, err)
		}

		file.WriteString(typeDef)
	}

	return writeGenToFile(file, "gen_types.go")
}

func generateTypeDef(d APIDescription, tgType TypeDescription) (string, error) {
	// If interface type, generate interface sections
	if len(tgType.Subtypes) != 0 || tgType.Name == tgTypeInputFile {
		return generateParentType(d, tgType)
	}

	typeFields, err := generateTypeFields(d, tgType)
	if err != nil {
		return "", fmt.Errorf("failed to generate type fields for %s: %w", tgType.Name, err)
	}

	typeDef := strings.Builder{}
	typeDef.WriteString(tgType.docs())

	if typeFields == "" {
		typeDef.WriteString("\ntype " + tgType.Name + " struct{}")
	} else {
		typeDef.WriteString("\ntype " + tgType.Name + " struct {")
		typeDef.WriteString(typeFields)
		typeDef.WriteString("\n}")
	}

	customUnmarshalDef, err := setupCustomUnmarshal(d, tgType)
	if err != nil {
		return "", fmt.Errorf("failed to setup custom unmarshal for %s: %w", tgType.Name, err)
	}
	typeDef.WriteString(customUnmarshalDef)

	interfaces, err := fulfilParentTypeInterfaces(d, tgType)
	if err != nil {
		return "", fmt.Errorf("failed to generate parent type interfaces %s: %w", tgType.Name, err)
	}
	typeDef.WriteString(interfaces)

	ok, fieldName, err := containsInputFile(d, tgType, map[string]bool{})
	if err != nil {
		return "", fmt.Errorf("failed to check if type requires special handling: %w", err)
	}
	if ok {
		err = inputParamsTmpl.Execute(&typeDef, inputParamsMethodData{
			Type:  tgType.Name,
			Field: snakeToTitle(fieldName),
		})
		if err != nil {
			return "", fmt.Errorf("failed to generate %s inputparam methods: %w", tgType.Name, err)
		}
	}

	return typeDef.String(), nil
}

// fieldContainsInputFile checks whether the field's type contains any inputfiles, and thus might be used to send data.
func fieldContainsInputFile(d APIDescription, field Field) (bool, error) {
	goType, err := field.getPreferredType(d)
	if err != nil {
		return false, err
	}
	cleanName := strings.TrimPrefix(goType, "[]")
	tgType, ok := d.Types[cleanName]
	if !ok {
		return false, nil
	}

	ok, _, err = containsInputFile(d, tgType, map[string]bool{})
	return ok, err
}

// containsInputFile returns a boolean to indicate whether or not tgType contains an InputFile.
// If true, it also returns the field name of that inputfile.
func containsInputFile(d APIDescription, tgType TypeDescription, checked map[string]bool) (bool, string, error) {
	// If already checked, we don't need to check again. This avoids infinite recursive loops.
	if checked[tgType.Name] {
		return false, "", nil
	}
	checked[tgType.Name] = true

	if tgType.Name == tgTypeInputMedia {
		return true, "media", nil
	}

	for _, f := range tgType.Fields {
		goType, err := f.getPreferredType(d)
		if err != nil {
			return false, "", err
		}

		if goType == tgTypeInputFile {
			return true, f.Name, nil
		}

		if isTgType(d, goType) {
			ok, _, err := containsInputFile(d, d.Types[goType], checked)
			if err != nil {
				return false, "", fmt.Errorf("failed to check if %s contains inputfiles: %w", goType, err)
			}
			if ok {
				// We return an error, because we can't actually handle this case yet.
				return false, "", errors.New("no support for recursive checks of inputfiles yet")
			}
		}
	}
	return false, "", nil
}

func generateParentType(d APIDescription, tgType TypeDescription) (string, error) {
	subTypes, err := getTypesByName(d, tgType.Subtypes)
	if err != nil {
		return "", fmt.Errorf("failed to get subtypes by name for %s: %w", tgType.Name, err)
	}

	interfaceDefinition, err := generateGenericInterfaceType(d, tgType.Name, subTypes)
	if err != nil {
		return "", fmt.Errorf("failed to generate generic interface type for %s: %w", tgType.Name, err)
	}

	typeDef := strings.Builder{}
	typeDef.WriteString(tgType.docs())
	typeDef.WriteString(interfaceDefinition)

	// If an interface type is sent by the API (eg in an update, or as a return) then we need to define custom
	// UnmarshalJSON methods to handle those edge cases into the right structs.
	if len(tgType.Subtypes) > 0 && tgType.sentByAPI(d) {
		unmarshalFunc, err := interfaceUnmarshalFunc(d, tgType)
		if err != nil {
			return "", fmt.Errorf("unable to generate interface unmarshal function: %w", err)
		}
		if unmarshalFunc != "" {
			typeDef.WriteString("\n\n" + unmarshalFunc)
		}
	}

	return typeDef.String(), nil
}

// Incoming types which marshal into interfaces need special handling to make sure the interfaces are
// populated correctly.
func setupCustomUnmarshal(d APIDescription, tgType TypeDescription) (string, error) {
	var fields []customUnmarshalFieldData
	generateCustomMarshal := false
	for idx, f := range tgType.Fields {
		prefType, err := f.getPreferredType(d)
		if err != nil {
			return "", err
		}

		if isTgType(d, prefType) {
			fieldType, err := getTypeByName(d, prefType)
			if err != nil {
				return "", fmt.Errorf("failed to get type of parameter %s in %s: %w", prefType, tgType.Name, err)
			}

			if len(fieldType.Subtypes) > 0 {
				subtypes, err := getTypesByName(d, fieldType.Subtypes)
				if err != nil {
					return "", fmt.Errorf("failed to get subtypes from %s: %w", fieldType.Name, err)
				}
				if len(getCommonFields(subtypes)) > 0 {
					generateCustomMarshal = true
				}
			}
		}

		if idx == 0 && len(tgType.SubtypeOf) > 0 && f.isConstantField(d, tgType) {
			continue
		}

		if isTgType(d, prefType) && !f.Required {
			prefType = "*" + prefType
		}

		fields = append(fields, customUnmarshalFieldData{
			Name:    snakeToTitle(f.Name),
			Custom:  len(d.Types[prefType].Subtypes) > 0,
			Type:    prefType,
			JSONTag: fmt.Sprintf("`json:\"%s\"`", f.Name),
		})
	}

	if !generateCustomMarshal {
		return "", nil
	}

	bd := strings.Builder{}
	err := customUnmarshalTmpl.Execute(&bd, customUnmarshalData{
		Type:   tgType.Name,
		Fields: fields,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate custom unmarshal: %w", err)
	}
	return bd.String(), nil
}

func fulfilParentTypeInterfaces(d APIDescription, tgType TypeDescription) (string, error) {
	typeInterfaces := strings.Builder{}
	for _, parentTypeName := range tgType.SubtypeOf {
		parentType, err := getTypeByName(d, parentTypeName)
		if err != nil {
			return "", fmt.Errorf("failed to get parent type %s of %s: %w", parentTypeName, tgType.Name, err)
		}

		commonFields, err := commonFieldGenerator(d, tgType, parentType)
		if err != nil {
			return "", err
		}

		typeInterfaces.WriteString(commonFields)

		typeInterfaces.WriteString(generateGenericInterfaceMethod(tgType.Name, parentTypeName))
	}

	for _, t := range getReplyMarkupTypes(d) {
		if tgType.Name == t.Name {
			typeInterfaces.WriteString(generateGenericInterfaceMethod(tgType.Name, tgTypeReplyMarkup))
			break
		}
	}

	return typeInterfaces.String(), nil
}

func generateGenericInterfaceMethod(typeName string, parentType string) string {
	methodName := titleToCamelCase(parentType)
	return fmt.Sprintf(`
// %s.%s is a dummy method to avoid interface implementation.
func (v %s) %s() {}
`, typeName, methodName, typeName, methodName)
}

func interfaceUnmarshalFunc(d APIDescription, tgType TypeDescription) (string, error) {
	constantField, err := tgType.getConstantFieldFromParent(d)
	if err != nil {
		return "", fmt.Errorf("failed to generate custom unmarshaller for %s: %w", tgType.Name, err)
	}
	if constantField == "" {
		return "", nil
	}

	var cases []customStructUnmarshalCaseData
	for _, subTypeName := range tgType.Subtypes {
		shortName := d.Types[subTypeName].getTypeNameFromParent(tgType.Name)
		cases = append(cases, customStructUnmarshalCaseData{
			ConstantFieldValue: shortName,
			TypeName:           subTypeName,
		})
	}

	bd := strings.Builder{}
	err = customStructUnmarshalTmpl.Execute(&bd, customStructUnmarshalData{
		UnmarshalFuncName: "unmarshal" + tgType.Name,
		ParentType:        tgType.Name,
		ConstantFieldName: snakeToTitle(constantField),
		CaseStatements:    cases,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate interface unmarshaller: %w", err)
	}

	return bd.String(), nil
}

func commonFieldGenerator(d APIDescription, tgType TypeDescription, parentType TypeDescription) (string, error) {
	// Some items need a custom marshaller to handle the "type" field
	shortName := tgType.getTypeNameFromParent(parentType.Name)

	subTypes, err := getTypesByName(d, parentType.Subtypes)
	if err != nil {
		return "", fmt.Errorf("failed to get subtypes of parent type %s of %s: %w", parentType.Name, tgType.Name, err)
	}

	commonFields := getCommonFields(subTypes)
	if len(commonFields) == 0 {
		return "", nil
	}

	constantField, err := parentType.getConstantFieldFromParent(d)
	if err != nil {
		return "", fmt.Errorf("failed to get constant field from %s: %w", parentType.Name, err)
	}

	bd := strings.Builder{}
	if len(commonFields) > 0 {
		commonGetMethods, err := generateAllCommonGetMethods(d, tgType.Name, commonFields, constantField, shortName)
		if err != nil {
			return "", err
		}
		bd.WriteString(commonGetMethods)

		// We only generate the merge func if there is a comm
		if constantField != "" {
			mergeFunc, err := generateMergeFunc(d, tgType.Name, shortName, tgType.Fields, parentType.Name, constantField)
			if err != nil {
				return "", err
			}
			bd.WriteString(mergeFunc)
		}
	}

	if constantField != "" {
		err = customMarshalTmpl.Execute(&bd, customMarshalData{
			Type:                  tgType.Name,
			ConstantFieldName:     strings.Title(constantField),
			ConstantJSONFieldName: constantField,
			ConstantValueName:     shortName,
		})
		if err != nil {
			return "", fmt.Errorf("failed to generate custom marshal function for %s: %w", tgType.Name, err)
		}
	}

	return bd.String(), nil
}

func generateAllCommonGetMethods(d APIDescription, typeName string, commonFields []Field, constantField string, shortName string) (string, error) {
	bd := strings.Builder{}
	for _, commonField := range commonFields {
		commonValueName := "v." + snakeToTitle(commonField.Name)
		if commonField.Name == constantField {
			commonValueName = strconv.Quote(shortName)
		}

		prefType, err := commonField.getPreferredType(d)
		if err != nil {
			return "", fmt.Errorf("failed to get preferred type for field %s of %s: %w", commonField.Name, typeName, err)
		}

		bd.WriteString(generateCommonGetMethod(typeName, snakeToTitle(commonField.Name), prefType, commonValueName))
	}
	return bd.String(), nil
}

// generateTypeFields generates the field contents of a telegram "struct" by parsing the fields.
func generateTypeFields(d APIDescription, tgType TypeDescription) (string, error) {
	var constantFields []string
	for _, f := range tgType.Fields {
		if f.isConstantField(d, tgType) {
			constantFields = append(constantFields, f.Name)
		}
	}

	return generateStructFields(d, tgType.Fields, constantFields)
}

// generateStructFields generates the go representation of a list of fields.
func generateStructFields(d APIDescription, fields []Field, constantFields []string) (string, error) {
	typeFields := strings.Builder{}
	for _, f := range fields {
		fieldType, err := f.getPreferredType(d)
		if err != nil {
			return "", fmt.Errorf("failed to get preferred type: %w", err)
		}

		skip := false
		for _, constantField := range constantFields {
			if f.Name == constantField {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if isTgType(d, fieldType) && !f.Required && len(d.Types[fieldType].Subtypes) == 0 {
			fieldType = "*" + fieldType
		}

		typeFields.WriteString("\n// " + f.Description)
		if f.Required && !isArray(fieldType) {
			typeFields.WriteString("\n" + snakeToTitle(f.Name) + " " + fieldType + " `json:\"" + f.Name + "\"`")
		} else {
			typeFields.WriteString("\n" + snakeToTitle(f.Name) + " " + fieldType + " `json:\"" + f.Name + ",omitempty\"`")
		}
	}

	return typeFields.String(), nil
}

func generateGenericInterfaceType(d APIDescription, name string, subtypes []TypeDescription) (string, error) {
	if len(subtypes) == 0 {
		return "\ntype " + name + " interface{}", nil
	}

	commonFields := getCommonFields(subtypes)

	constantField, err := getConstantFieldFromCommons(d, commonFields)
	if err != nil {
		return "", fmt.Errorf("failed to get constant fields: %w", err)
	}

	hasInputFile, fieldName, err := containsInputFile(d, subtypes[0], map[string]bool{})
	if err != nil {
		return "", fmt.Errorf("failed to check if %s types all contain inputfiles: %w", name, err)
	}

	// If the inputfile is a common field, then the interface contains fields.
	hasInputFile = hasInputFile && contains(fieldName, getFieldNames(commonFields))

	bd := strings.Builder{}
	bd.WriteString(fmt.Sprintf("\ntype %s interface{", name))
	for _, f := range commonFields {
		prefType, err := f.getPreferredType(d)
		if err != nil {
			return "", fmt.Errorf("failed to get preferred type for %s: %w", f.Name, err)
		}
		bd.WriteString(fmt.Sprintf("\nGet%s() %s", snakeToTitle(f.Name), prefType))
	}

	// create a dummy func to avoid external types implementing this interface
	bd.WriteString(fmt.Sprintf("\n// %s exists to avoid external types implementing this interface.", titleToCamelCase(name)))
	bd.WriteString(fmt.Sprintf("\n%s()", titleToCamelCase(name)))

	if hasInputFile {
		bd.WriteString("\n// InputParams allows for uploading attachments with files.")
		bd.WriteString("\nInputParams(string, map[string]NamedReader) ([]byte, error)")
	}

	if len(commonFields) > 0 && constantField != "" {
		bd.WriteString(fmt.Sprintf("\n// Merge%s returns a Merged%s struct to simplify working with complex telegram types in a non-generic world.", name, name))
		bd.WriteString(fmt.Sprintf("\nMerge%s() Merged%s", name, name))
	}
	bd.WriteString("\n}")

	if len(commonFields) > 0 && constantField != "" {
		mergedStruct, err := generateMergedStruct(d, name, subtypes)
		if err != nil {
			return "", fmt.Errorf("failed to generate merged struct: %w", err)
		}

		bd.WriteString("\n" + mergedStruct)

		commonGetMethods, err := generateAllCommonGetMethods(d, "Merged"+name, commonFields, "", "")
		if err != nil {
			return "", fmt.Errorf("failed to generate common get methods: %w", err)
		}

		bd.WriteString(commonGetMethods)
		bd.WriteString(generateGenericInterfaceMethod("Merged"+name, name))
		bd.WriteString(fmt.Sprintf(`
// Merge%s returns a Merged%s struct to simplify working with types in a non-generic world.
func (v Merged%s) Merge%s() Merged%s {
	return v
}`, name, name, name, name, name))
	}

	return bd.String(), nil
}

func generateMergedStruct(d APIDescription, name string, subtypes []TypeDescription) (string, error) {
	allFields, err := getAllFields(d, subtypes, name)
	if err != nil {
		return "", fmt.Errorf("failed to get all fields: %w", err)
	}

	fields, err := generateStructFields(d, allFields, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate struct fields: %w", err)
	}

	return fmt.Sprintf(`
// Merged%s is a helper type to simplify interactions with the various %s subtypes.
type Merged%s struct {
	%s
}
`, name, name, name, strings.TrimSpace(fields)), nil
}

func generateCommonGetMethod(t string, commonName string, commonType string, commonValue string) string {
	return fmt.Sprintf(`
// Get%s is a helper method to easily access the common fields of an interface.
func (v %s) Get%s() %s {
	return %s
}
`, commonName, t, commonName, commonType, commonValue)
}

func generateMergeFunc(d APIDescription, typeName string, shortname string, fields []Field, parentType string, constantField string) (string, error) {
	subTypes, err := getTypesByName(d, d.Types[parentType].Subtypes)
	if err != nil {
		return "", fmt.Errorf("failed to get subtypes by name for %s: %w", typeName, err)
	}

	allParentFields, err := getAllFields(d, subTypes, parentType)
	if err != nil {
		return "", fmt.Errorf("failed to get all fields for %s with parent type %s: %w", typeName, parentType, err)
	}

	bd := strings.Builder{}

	bd.WriteString(fmt.Sprintf("\n// Merge%s returns a Merged%s struct to simplify working with types in a non-generic world.", parentType, parentType))
	bd.WriteString(fmt.Sprintf("\nfunc (v %s) Merge%s() Merged%s {", typeName, parentType, parentType))
	bd.WriteString(fmt.Sprintf("\n\treturn Merged%s{", parentType))
	for _, f := range fields {
		if f.Name == constantField {
			bd.WriteString(fmt.Sprintf("\n\t%s: \"%s\",", snakeToTitle(f.Name), shortname))
			continue
		}

		deref := false
		for _, parentField := range allParentFields {
			if parentField.Name == f.Name {
				fieldType, err := f.getPreferredType(d)
				if err != nil {
					return "", fmt.Errorf("failed to get preferred type: %w", err)
				}

				if isTgType(d, fieldType) && f.Required != parentField.Required {
					deref = true
				}
			}
		}

		if deref {
			bd.WriteString(fmt.Sprintf("\n\t%s: &v.%s,", snakeToTitle(f.Name), snakeToTitle(f.Name)))
		} else {
			bd.WriteString(fmt.Sprintf("\n\t%s: v.%s,", snakeToTitle(f.Name), snakeToTitle(f.Name)))
		}
	}

	bd.WriteString("\n\t}")
	bd.WriteString("\n}")
	bd.WriteString("\n")
	return bd.String(), nil
}

type customUnmarshalFieldData struct {
	Name    string
	Custom  bool
	Type    string
	JSONTag string
}

type customUnmarshalData struct {
	Type   string
	Fields []customUnmarshalFieldData
}

const customUnmarshal = `
// UnmarshalJSON is a custom JSON unmarshaller to use the helpers which allow for unmarshalling structs into interfaces.
func (v *{{.Type}}) UnmarshalJSON(b []byte) error {
	// All fields in {{.Type}}, with interface fields as json.RawMessage
	type tmp struct {
        {{ range $f := .Fields }}
		{{ $f.Name }} {{ if $f.Custom }} json.RawMessage {{ else }} {{ $f.Type }} {{ end }} {{ $f.JSONTag -}}
		{{- end }}
	}
	t := tmp{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	{{ range $f := .Fields }}
		{{- if $f.Custom}}
			v.{{ $f.Name }}, err = unmarshal{{ $f.Type }}(t.{{$f.Name}})
			if err != nil {
				return err
			}
		{{- else }}
			v.{{ $f.Name }} = t.{{ $f.Name }}
		{{- end }}
	{{- end }}

	return nil
}
`

type customStructUnmarshalData struct {
	UnmarshalFuncName string
	ParentType        string
	ConstantFieldName string
	CaseStatements    []customStructUnmarshalCaseData
}

type customStructUnmarshalCaseData struct {
	ConstantFieldValue string
	TypeName           string
}

const customStructUnmarshal = `
// {{.UnmarshalFuncName}}Array is a JSON unmarshalling helper which allows unmarshalling an array of interfaces 
// using {{.UnmarshalFuncName}}.
func {{.UnmarshalFuncName}}Array(d json.RawMessage) ([]{{.ParentType}}, error) {
	var ds []json.RawMessage
	err := json.Unmarshal(d, &ds)
	if err != nil {
		return nil, err
	}

	var vs []{{.ParentType}}
	for _, d := range ds {
		v, err := {{.UnmarshalFuncName}}(d)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}

	return vs, nil
}

// {{.UnmarshalFuncName}} is a JSON unmarshal helper to marshal the right structs into a {{.ParentType}} interface
// based on the {{.ConstantFieldName}} field.
func {{.UnmarshalFuncName}}(d json.RawMessage) ({{.ParentType}}, error) {
		if len(d) == 0 {
			return nil, nil
		}

		t := struct {
			{{.ConstantFieldName}} string
		}{}
		err := json.Unmarshal(d, &t)
		if err != nil {
			return nil, err
		}

		switch t.{{.ConstantFieldName}} {
		{{-  range $val := .CaseStatements }}
		case "{{ $val.ConstantFieldValue }}":
			s := {{ $val.TypeName }}{}
			err := json.Unmarshal(d, &s)
			if err != nil {
				return nil, err
			}
			return s, nil
		{{ end }}
		}
		return nil, fmt.Errorf("unknown interface with {{.ConstantFieldName}} %v", t.{{.ConstantFieldName}})
}`

type customMarshalData struct {
	Type                  string
	ConstantFieldName     string
	ConstantJSONFieldName string
	ConstantValueName     string
}

// The alias type is required to avoid infinite MarshalJSON loops.
const customMarshal = `
// MarshalJSON is a custom JSON marshaller to allow for enforcing the {{.ConstantFieldName}} value.
func (v {{.Type}}) MarshalJSON() ([]byte, error) {
	type alias {{.Type}}
	a := struct{
		{{.ConstantFieldName}} string ` + "`json:\"{{.ConstantJSONFieldName}}\"`" + `
		alias
	}{
		{{.ConstantFieldName}}: "{{.ConstantValueName}}",
		alias: (alias)(v),
	}
	return json.Marshal(a)
}
`

type inputParamsMethodData struct {
	Type  string
	Field string
}

const inputParamsMethod = `
func (v {{.Type}}) InputParams(mediaName string, data map[string]NamedReader) ([]byte, error) {
	if v.{{.Field}} != nil {
		switch m := v.{{.Field}}.(type) {
		case string:
			// ok, noop

		case NamedReader:
			v.{{.Field}} = "attach://" + mediaName
			data[mediaName] = m

		case io.Reader:
			v.{{.Field}} = "attach://" + mediaName
			data[mediaName] = NamedFile{File: m}

		default:
			return nil, fmt.Errorf("unknown type: %T", v.{{.Field}})
		}
	}
	
	return json.Marshal(v)
}
`
