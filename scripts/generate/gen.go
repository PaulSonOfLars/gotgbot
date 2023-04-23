package main

import (
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
	"unicode"
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

func (td TypeDescription) receiverName() string {
	var rs []rune
	for _, r := range td.Name {
		if unicode.IsUpper(r) {
			rs = append(rs, r)
		}
	}

	return strings.ToLower(string(rs))
}

func (td TypeDescription) sentByAPI(d APIDescription) bool {
	checked := map[string]bool{}

	for _, m := range d.Methods {
		for _, r := range m.Returns {
			if r == td.Name {
				return true
			}

			child, ok := d.Types[r]
			if !ok || checked[r] {
				continue
			}

			if CheckChildTypes(d, child, td.Name, []string{td.Name}) {
				return true
			}
			checked[r] = true
		}
	}
	return false
}

func CheckChildTypes(d APIDescription, tgType TypeDescription, typeName string, skip []string) bool {
	for _, f := range tgType.Fields {
		for _, t := range f.Types {
			if t == typeName {
				return true
			}

			if contains(t, skip) {
				continue
			}

			if child, ok := d.Types[t]; ok && t != tgType.Name {
				if CheckChildTypes(d, child, typeName, append(skip, tgType.Name)) {
					return true
				}
			}
		}
	}
	return false
}

func (td TypeDescription) getTypeNameFromParent(parentType string) string {
	// Telegram inconsistencies
	if td.Name == "ChatMemberOwner" {
		return "creator"
	} else if td.Name == "ChatMemberBanned" {
		return "kicked"
	}

	typeName := strings.TrimPrefix(td.Name, parentType)
	typeName = strings.TrimPrefix(typeName, "Cached") // some of them are "Cached"
	typeName = strings.TrimSuffix(typeName, "Field")  // some of them are "Field"
	return titleToSnake(typeName)
}

func (td TypeDescription) getConstantFieldFromParent(d APIDescription) (string, error) {
	if len(td.Subtypes) == 0 {
		return "", fmt.Errorf("expected %s to be a parent", td.Name)
	}

	subTypes, err := getTypesByName(d, td.Subtypes)
	if err != nil {
		return "", fmt.Errorf("failed to get parent type %s: %w", td.Name, err)
	}

	common := getCommonFields(subTypes)
	if len(common) == 0 {
		return "", fmt.Errorf("no common fields for parenttype %s", td.Name)
	}
	return common[0].Name, nil
}

func (m MethodDescription) docs() string {
	return docs(strings.Title(m.Name), m.Href, m.Description)
}

func (td TypeDescription) docs() string {
	return docs(td.Name, td.Href, td.Description)
}

func docs(name string, href string, ds []string) string {
	bd := strings.Builder{}
	// Start by mentioning name, with docs link.
	bd.WriteString(fmt.Sprintf("\n// %s (%s)", name, href))
	// add a newline for nice paragraphing and readability
	bd.WriteString("\n//")
	for _, text := range ds {
		// if this is a list item, we should add extra indentation to make sure it renders right in the docs pages.
		if strings.HasPrefix(text, "- ") {
			text = " " + text
		}

		bd.WriteString("\n// " + text)
	}
	return bd.String()
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

func (f Field) isConstantField(d APIDescription, tgType TypeDescription) bool {
	for _, parent := range tgType.SubtypeOf {
		constantField, err := d.Types[parent].getConstantFieldFromParent(d)
		if err != nil {
			continue
		}
		if constantField == f.Name {
			return true
		}
	}
	return false
}

const (
	// These are all base telegram types which make sense in other languages.
	tgTypeString  = "String"
	tgTypeBoolean = "Boolean"
	tgTypeFloat   = "Float"
	tgTypeInteger = "Integer"
	// These are all custom telegram types.
	tgTypeMessage    = "Message"
	tgTypeFile       = "File"
	tgTypeInputFile  = "InputFile"
	tgTypeInputMedia = "InputMedia"
	// This is actually a custom type.
	tgTypeReplyMarkup = "ReplyMarkup"
)

func generate(d APIDescription) error {
	// TODO: Use golang templates instead of string builders
	if err := generateTypes(d); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	if err := generateMethods(d); err != nil {
		return fmt.Errorf("failed to generate helpers: %w", err)
	}

	if err := generateHelpers(d); err != nil {
		return fmt.Errorf("failed to generate helpers: %w", err)
	}

	if err := generateConsts(d); err != nil {
		return fmt.Errorf("failed to generate consts: %w", err)
	}

	return nil
}

func writeGenToFile(file strings.Builder, filename string) error {
	write, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	bs := []byte(file.String())

	_, err = write.WriteAt(bs, 0)
	if err != nil {
		return fmt.Errorf("failed to write unformatted file %s: %w", filename, err)
	}

	fmted, err := format.Source(bs)
	if err != nil {
		return fmt.Errorf("failed to format file %s: %w", filename, err)
	}

	err = write.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file %s: %w", filename, err)
	}

	_, err = write.WriteAt(fmted, 0)
	if err != nil {
		return fmt.Errorf("failed to write final file %s: %w", filename, err)
	}

	return nil
}

func orderedTgTypes(d APIDescription) []string {
	types := make([]string, 0, len(d.Types))
	for k := range d.Types {
		types = append(types, k)
	}

	sort.Strings(types)

	return types
}

func orderedMethods(d APIDescription) []string {
	methods := make([]string, 0, len(d.Methods))
	for k := range d.Methods {
		methods = append(methods, k)
	}

	sort.Strings(methods)

	return methods
}

func isTgType(d APIDescription, goType string) bool {
	_, ok := d.Types[goType]
	return ok
}

func (f Field) getPreferredType(d APIDescription) (string, error) {
	if f.Name == "media" {
		if len(f.Types) == 1 && f.Types[0] == "String" {
			return tgTypeInputFile, nil
		}
		var arrayType bool
		// TODO: check against API description type
		for _, t := range f.Types {
			arrayType = arrayType || isTgArray(t)

			if !strings.Contains(t, tgTypeInputMedia) {
				return "", fmt.Errorf("mediatype %s is not of kind InputMedia for field %s", t, f.Name)
			}
		}

		if arrayType {
			return "[]" + tgTypeInputMedia, nil
		}

		return tgTypeInputMedia, nil
	}

	if f.Name == "reply_markup" {
		if len(f.Types) == 4 {
			// TODO: this should check if the values are replymarkup children, instead of checking length
			// Custom type used to handle the fact that reply_markup can take one of:
			// InlineKeyboardMarkup
			// ReplyKeyboardMarkup
			// ReplyKeyboardRemove
			// ForceReply
			return tgTypeReplyMarkup, nil

		} else if len(f.Types) == 1 {
			return toGoType(f.Types[0]), nil

		} else {
			return "", fmt.Errorf("unable to handle reply_markup field with types: %v", f.Types)
		}
	}

	if len(f.Types) == 1 {
		goType := toGoType(f.Types[0])

		// Optional TG types should be pointers, unless they're already an interface type.
		if !f.Required && isTgType(d, f.Types[0]) && !isArray(goType) && goType != tgTypeInputFile {
			rawType, err := getTypeByName(d, f.Types[0])
			if err != nil {
				return "", fmt.Errorf("failed to get parent for %s: %w", f.Types[0], err)
			}

			if len(rawType.Subtypes) == 0 {
				return "*" + goType, nil
			}
		}

		// Some fields are marked as "May be empty", in which case the empty values are still meaningful.
		// These should be handled as pointers, so we can differentiate the empty case.
		if (strings.Contains(f.Description, "May be empty") || strings.Contains(f.Description, "If not specified")) && !(isPointer(goType) || isArray(goType)) {
			// If "not specified" has the same outcome as "empty", then we dont need to use a pointer.
			if strings.Contains(f.Description, "If not specified or empty") {
				return goType, nil
			}

			// Note: we override the AnswerCallbackQuery 'text' field to not be a pointer, because "" and nil have the same outcome here.
			if f.Name == "text" && strings.Contains(f.Description, "nothing will be shown") {
				return goType, nil
			}
			return "*" + goType, nil
		}

		return goType, nil
	}

	if len(f.Types) == 2 {
		if f.Types[0] == tgTypeInputFile && f.Types[1] == tgTypeString {
			return toGoType(f.Types[0]), nil
		} else if f.Types[0] == tgTypeInteger && f.Types[1] == tgTypeString {
			return toGoType(f.Types[0]), nil
		}
	}

	return "", fmt.Errorf("unable to choose one of %v for field %s", f.Types, f.Name)
}

func (m MethodDescription) GetReturnTypes(d APIDescription) ([]string, error) {
	// We currently only support dual returns for msg+bool
	if len(m.Returns) >= 2 && !(len(m.Returns) == 2 && m.Returns[0] == tgTypeMessage && m.Returns[1] == tgTypeBoolean) {
		return nil, fmt.Errorf("no support for multiple returns for types %s", m.Name)
	}

	var retTypes []string
	for _, prefRetVal := range m.Returns {
		goRetType := toGoType(prefRetVal)
		if isTgType(d, goRetType) && len(d.Types[prefRetVal].Subtypes) == 0 {
			goRetType = "*" + goRetType
		}

		retTypes = append(retTypes, goRetType)
	}

	return retTypes, nil
}

func (m MethodDescription) optsName() string {
	return snakeToTitle(m.Name) + "Opts"
}
