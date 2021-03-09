package main

import (
	"regexp"
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
