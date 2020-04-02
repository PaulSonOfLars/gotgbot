package helpers

import (
	"strings"
)

var allMd = []string{"_", "*", "`", "[", "]", "(", ")", "\\"}
var mdRepl = strings.NewReplacer(func() (out []string) {
	for _, x := range allMd {
		out = append(out, x, "\\"+x)
	}
	return out
}()...)

func EscapeMarkdown(s string) string {
	return mdRepl.Replace(s)
}

var allMdV2 = []string{"_", "*", "`", "~", "[", "]", "(", ")", "\\"} // __ is not necessary because of _
var mdV2Repl = strings.NewReplacer(func() (out []string) {
	for _, x := range allMdV2 {
		out = append(out, x, "\\"+x)
	}
	return out
}()...)

func EscapeMarkdownV2(s string) string {
	return mdV2Repl.Replace(s)
}
