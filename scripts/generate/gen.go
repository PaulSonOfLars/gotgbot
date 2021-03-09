package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
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

func (d TypeDescription) receiverName() string {
	var rs []rune
	for _, r := range []rune(d.Name) {
		if unicode.IsUpper(r) {
			rs = append(rs, r)
		}
	}

	return strings.ToLower(string(rs))
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

const (
	// These are all base telegram types which make sense in other languages.
	tgTypeString  = "String"
	tgTypeBoolean = "Boolean"
	tgTypeFloat   = "Float"
	tgTypeInteger = "Integer"
	// These are all custom telegram types.
	tgTypeMessage              = "Message"
	tgTypeFile                 = "File"
	tgTypeInputFile            = "InputFile"
	tgTypeInputMedia           = "InputMedia"
	tgTypeInlineQueryResult    = "InlineQueryResult"
	tgTypeInputMessageContent  = "InputMessageContent"
	tgTypePassportElementError = "PassportElementError"
	tgTypeCallbackGame         = "CallbackGame"
	// This is actually a custom type.
	tgTypeReplyMarkup = "ReplyMarkup"
)

func generate(api io.ReadCloser) error {
	var d APIDescription
	if err := json.NewDecoder(api).Decode(&d); err != nil {
		return err
	}

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

func (f Field) getPreferredType() (string, error) {
	if len(f.Types) == 1 {
		return toGoType(f.Types[0]), nil
	}

	if len(f.Types) == 2 {
		if f.Types[0] == tgTypeInputFile && f.Types[1] == tgTypeString {
			return toGoType(f.Types[0]), nil
		} else if f.Types[0] == tgTypeInteger && f.Types[1] == tgTypeString {
			return toGoType(f.Types[0]), nil
		}
	}

	if f.Name == "media" {
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

	if f.Name == "reply_markup" && len(f.Types) == 4 {
		// Custom type used to handle the fact that reply_markup can take one of:
		// InlineKeyboardMarkup
		// ReplyKeyboardMarkup
		// ReplyKeyboardRemove
		// ForceReply
		return tgTypeReplyMarkup, nil
	}

	return "", fmt.Errorf("unable to choose one of %v for field %s", f.Types, f.Name)
}

func (m MethodDescription) GetReturnType(d APIDescription) (string, error) {
	prefRetVal := ""
	if len(m.Returns) == 1 {
		prefRetVal = m.Returns[0]
	} else if len(m.Returns) == 2 && m.Returns[0] == tgTypeMessage && m.Returns[1] == tgTypeBoolean {
		prefRetVal = m.Returns[0]
	} else {
		return "", fmt.Errorf("failed to determine return type for method %s", m.Name)
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
