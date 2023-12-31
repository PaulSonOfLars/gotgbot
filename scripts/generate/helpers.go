package main

import (
	"fmt"
	"strings"
	"text/template"
)

func generateHelpers(d APIDescription) error {
	helpers := strings.Builder{}
	helpers.WriteString(`
// THIS FILE IS AUTOGENERATED. DO NOT EDIT.
// Regen by running 'go generate' in the repo root.

package gotgbot

`)

	for _, tgMethodName := range orderedMethods(d) {
		tgMethod := d.Methods[tgMethodName]

		helper, err := generateHelperDef(d, tgMethod)
		if err != nil {
			return fmt.Errorf("failed to generate helpers for %s: %w", tgMethodName, err)
		}

		if helper == "" {
			continue
		}

		helpers.WriteString(helper)
	}

	return writeGenToFile(helpers, "gen_helpers.go")
}

func generateHelperDef(d APIDescription, tgMethod MethodDescription) (string, error) {
	helperDef := strings.Builder{}
	hasFromChat := false

	for _, x := range tgMethod.Fields {
		if x.Name == "from_chat_id" {
			hasFromChat = true
			break
		}
	}

	for _, typeName := range orderedTgTypes(d) {
		if typeName == tgTypeFile {
			continue
		}

		tgType := d.Types[typeName]
		if len(tgType.Subtypes) != 0 {
			// Interfaces can't have methods on them
			continue
		}

		// Get list of fields which match
		fields := getMethodFieldsTypeMatches(tgMethod, tgType)
		if len(fields) == 0 {
			continue
		}

		newMethodName, ok, err := getMethodFieldsSubtypeMatches(d, tgMethod, tgType, hasFromChat, fields)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}

		ret, err := tgMethod.GetReturnTypes(d)
		if err != nil {
			return "", fmt.Errorf("failed to get return type for %s: %w", tgMethod.Name, err)
		}

		receiverName := tgType.receiverName()

		funcCallArgList, funcDefArgList, optsContent, err := generateHelperArguments(d, tgMethod, receiverName, fields)
		if err != nil {
			return "", err
		}

		funcDefArgs := strings.Join(funcDefArgList, ", ")
		funcCallArgs := strings.Join(funcCallArgList, ", ")

		helperDef.WriteString("\n// " + newMethodName + " Helper method for Bot." + strings.Title(tgMethod.Name))

		err = helperFuncTmpl.Execute(&helperDef, helperFuncData{
			Receiver:     receiverName,
			TypeName:     typeName,
			HelperName:   newMethodName,
			ReturnType:   strings.Join(ret, ", "),
			FuncDefArgs:  funcDefArgs,
			Contents:     optsContent,
			OptsName:     tgMethod.optsName(),
			MethodName:   strings.Title(tgMethod.Name),
			FuncCallArgs: funcCallArgs,
		})
		if err != nil {
			return "", fmt.Errorf("failed to execute template to generate %s helper method on %s: %w", newMethodName, typeName, err)
		}
	}

	return helperDef.String(), nil
}

func generateHelperArguments(d APIDescription, tgMethod MethodDescription, receiverName string, fields map[string]string) ([]string, []string, string, error) {
	var funcCallArgList []string
	optsContent := strings.Builder{}
	funcDefArgList := []string{"b *Bot"}
	hasOpts := false

	for _, mf := range tgMethod.Fields {
		hasOpts = hasOpts || !mf.Required

		prefType, err := mf.getPreferredType(d)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to get preferred type for field %s of %s: %w", mf.Name, tgMethod.Name, err)
		}

		if fName, ok := fields[mf.Name]; ok {
			if !mf.Required {
				defaultValue := getDefaultTypeVal(d, prefType)
				optsContent.WriteString("\n	if opts." + snakeToTitle(mf.Name) + " == " + defaultValue + " {")
				if isPointer(prefType) {
					optsContent.WriteString("\n		opts." + snakeToTitle(mf.Name) + " = &" + receiverName + "." + snakeToTitle(fName))
				} else {
					optsContent.WriteString("\n		opts." + snakeToTitle(mf.Name) + " = " + receiverName + "." + snakeToTitle(fName))
				}
				optsContent.WriteString("\n	}")
				continue
			}

			funcCallArgList = append(funcCallArgList, receiverName+"."+snakeToTitle(fName))
			continue
		}

		if !mf.Required {
			continue
		}

		funcDefArgList = append(funcDefArgList, snakeToCamel(mf.Name)+" "+prefType)
		funcCallArgList = append(funcCallArgList, snakeToCamel(mf.Name))
	}

	funcDefArgList = append(funcDefArgList, "opts *"+tgMethod.optsName())
	funcCallArgList = append(funcCallArgList, "opts")

	return funcCallArgList, funcDefArgList, optsContent.String(), nil
}

func getMethodFieldsSubtypeMatches(d APIDescription, tgMethod MethodDescription, tgType TypeDescription, hasFromChat bool, fields map[string]string) (string, bool, error) {
	typeName := tgType.Name
	if typeName == "InaccessibleMessage" {
		typeName = "Message"
	}
	newMethodName := strings.Replace(tgMethod.Name, typeName, "", 1)

	for _, f := range tgType.Fields {
		if f.Name == "reply_to_message" {
			// this subfield just causes confusion; we always want the message_id
			continue
		}

		for _, mf := range tgMethod.Fields {
			prefType, err := f.getPreferredType(d)
			if err != nil {
				return "", false, fmt.Errorf("failed to get preferred type for field %s of %s: %w", mf.Name, tgMethod.Name, err)
			}

			if isTgType(d, prefType) && f.Name+"_id" == mf.Name {
				newMethodName = strings.ReplaceAll(newMethodName, prefType, "")

				if hasFromChat && mf.Name == "chat_id" {
					fields["from_chat_id"] = f.Name + ".Id"
				} else {
					fields[mf.Name] = f.Name + ".Id" // Note: maybe not just assume ID field exists?
				}
			}
		}
	}
	return strings.Title(newMethodName), newMethodName != tgMethod.Name, nil
}

func getMethodFieldsTypeMatches(tgMethod MethodDescription, tgType TypeDescription) map[string]string {
	fields := map[string]string{}
	typeName := tgType.Name
	if typeName == "InaccessibleMessage" {
		typeName = "Message"
	}
	snakeTypeNameId := titleToSnake(typeName) + "_id"

	for _, f := range tgMethod.Fields {
		if f.Name == snakeTypeNameId || f.Name == "id" {
			fields[snakeTypeNameId] = findIdField(tgType, f.Name)
		}
	}
	return fields
}

func findIdField(tgType TypeDescription, methodFieldName string) string {
	// And iterate over all the type fields, to see if any match the method field name.
	for _, f := range tgType.Fields {
		if methodFieldName == f.Name {
			return methodFieldName
		}
	}
	return "id" // we default to "id" if nothing else matches
}

var helperFuncTmpl = template.Must(template.New("helperFunc").Parse(helperFunc))

type helperFuncData struct {
	Receiver     string
	TypeName     string
	HelperName   string
	ReturnType   string
	FuncDefArgs  string
	Contents     string
	OptsName     string
	MethodName   string
	FuncCallArgs string
}

const helperFunc = `
func ({{.Receiver}} {{.TypeName}}) {{.HelperName}}({{.FuncDefArgs}}) ({{.ReturnType}}, error) {
	{{- if .Contents}}
		if opts == nil {
			opts = &{{.OptsName}}{}
		}
		{{.Contents}}

	{{end}}
	return b.{{.MethodName}}({{.FuncCallArgs}})
}
`
