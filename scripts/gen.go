package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
)

type APIDescription struct {
	Types   map[string]TypeDescription   `json:"types"`
	Methods map[string]MethodDescription `json:"methods"`
}

type TypeDescription struct {
	Name        string   `json:"name"`
	Description []string `json:"description"`
	Fields      []Field  `json:"fields"`
	Href        string   `json:"href"`
	Subtypes    []string `json:"subtypes"`
	SubtypeOf   []string `json:"subtype_of"`
}

type MethodDescription struct {
	Name        string   `json:"name"`
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

func isTgType(d APIDescription, goType string) bool {
	_, ok := d.Types[goType]
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

func (f Field) getPreferredType() (string, error) {
	if len(f.Types) == 1 {
		return f.Types[0], nil
	}
	if len(f.Types) == 2 {
		if f.Types[0] == "InputFile" && f.Types[1] == "String" {
			return f.Types[0], nil
		} else if f.Types[0] == "Integer" && f.Types[1] == "String" {
			return f.Types[0], nil
		}
	}
	if f.Name == "media" {
		var arrayType bool
		// TODO: check against API description type
		for _, t := range f.Types {
			arrayType = arrayType || isTgArray(t)

			if !strings.Contains(t, "InputMedia") {
				return "", fmt.Errorf("mediatype %s is not of kind InputMedia for field %s", t, f.Name)
			}
		}
		if arrayType {
			return "Array of InputMedia", nil
		}
		return "InputMedia", nil
	}

	if f.Name == "reply_markup" && len(f.Types) == 4 {
		// Custom type used to handle the fact that reply_markup can take one of:
		// InlineKeyboardMarkup
		// ReplyKeyboardMarkup
		// ReplyKeyboardRemove
		// ForceReply
		return "ReplyMarkup", nil
	}

	return f.Types[0], fmt.Errorf("unable to choose one of %v for field %s", f.Types, f.Name)
}

func (m MethodDescription) GetReturnType(d APIDescription) (string, error) {
	prefRetVal := ""
	switch len(m.Returns) {
	case 1:
		prefRetVal = m.Returns[0]
	case 2:
		if m.Returns[0] == "Message" && m.Returns[1] == "Boolean" {
			prefRetVal = m.Returns[0]
		} else {
			return "", fmt.Errorf("failed to determine return type for method from %v", m.Returns)
		}
	default:
		return "", fmt.Errorf("failed to determine return type for method")
	}

	retType := toGoType(prefRetVal)
	if isTgType(d, retType) {
		retType = "*" + retType
	}
	return retType, nil
}

func (m MethodDescription) optsName() string {
	return snakeToTitle(m.Name) + "Opts"
}
