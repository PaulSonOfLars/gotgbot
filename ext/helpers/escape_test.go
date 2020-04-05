package helpers

import (
	"testing"
)

var sampleText = "Hey this is _markdown_ text with _items _ contained_ `everywhere`, [how](does) that ~look~ to *you*?"

func BenchmarkEscapeMarkdown(b *testing.B) {
	var v string
	for i := 0; i < b.N; i++ {
		v = EscapeMarkdown(sampleText)
	}
	_ = v
}

func BenchmarkEscapeMarkdownV2(b *testing.B) {
	var v string
	for i := 0; i < b.N; i++ {
		v = EscapeMarkdownV2(sampleText)
	}
	_ = v
}
