package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"strings"
)

type APIDescription struct {
	Types   map[string]TypeDescription   `json:"types"`
	Methods map[string]MethodDescription `json:"methods"`
}

type TypeDescription struct {
	Description []string `json:"description"`
	Fields      []Field  `json:"fields"`
	Href        string   `json:"href"`
	Subtypes    []string `json:"subtypes"`
	SubtypeOf   []string `json:"subtype_of"`
}

type MethodDescription struct {
	Fields      []Field  `json:"fields"`
	Returns     []string `json:"returns"`
	Description []string `json:"description"`
	Href        string   `json:"href"`
}

type Field struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
}

func main() {
	api, err := os.Open("api.json")
	if err != nil {
		panic(err)
	}

	var d APIDescription
	if err = json.NewDecoder(api).Decode(&d); err != nil {
		panic(err)
	}

	// TODO: Use golang templates instead of string builders
	err = generateTypes(d)
	if err != nil {
		panic(err)
	}
	err = generateMethods(d)
	if err != nil {
		panic(err)
	}
}

func writeGenToFile(file strings.Builder, filename string) error {
	write, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	bs := []byte(file.String())

	_, err = write.WriteAt(bs, 0)
	if err != nil {
		return err
	}

	fmted, err := format.Source(bs)
	if err != nil {
		return err
	}

	_, err = write.WriteAt(fmted, 0)
	if err != nil {
		return err
	}
	return nil
}

func isTgType(tgTypes map[string]TypeDescription, goType string) bool {
	_, ok := tgTypes[goType]
	return ok
}

func snakeToTitle(s string) string {
	bd := strings.Builder{}
	for _, s := range strings.Split(s, "_") {
		bd.WriteString(strings.Title(s))
	}
	return bd.String()
}

func snakeToCamel(s string) string {
	title := snakeToTitle(s)
	return strings.ToLower(title[:1]) + title[1:]
}

func toGoType(s string) string {
	pref := ""
	for isTgArray(s) {
		pref += "[]"
		s = strings.TrimPrefix(s, "Array of ")
	}

	switch s {
	case "Integer":
		return pref + "int64"
	case "Float":
		return pref + "float64"
	case "Boolean":
		return pref + "bool"
	case "String":
		return pref + "string"
	}
	return pref + s
}

func isTgArray(s string) bool {
	return strings.HasPrefix(s, "Array of ")
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

func getPreferredType(f Field) string {
	if len(f.Types) == 1 {
		return f.Types[0]
	}
	if len(f.Types) == 2 {
		if f.Types[0] == "InputFile" && f.Types[1] == "String" {
			return f.Types[0]
		} else if f.Types[0] == "Integer" && f.Types[1] == "String" {
			return f.Types[0]
		}
	}
	if f.Name == "media" {
		var arrayType bool
		// TODO: check against API description type
		for _, t := range f.Types {
			arrayType = arrayType || isTgArray(t)

			if !strings.Contains(t, "InputMedia") {
				fmt.Printf("%s: mediatype %s is not of kind InputMedia\n", f.Name, t)
				return f.Types[0]
			}
		}
		if arrayType {
			return "Array of InputMedia"
		}
		return "InputMedia"
	}

	// skip reply_markup for now to avoid logspam
	if f.Name == "reply_markup" {
		// TODO: Handle reply_markup results
		return f.Types[0]
	}

	fmt.Printf("%s: unable to choose one of %v\n", f.Name, f.Types)
	return f.Types[0]
}
