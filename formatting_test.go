package gotgbot

import (
	"testing"
)

func TestSplitEdgeWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		entType string
		wantPre string
		wantCnt string
		wantPos string
	}{
		{
			name:    "no whitespace",
			text:    "hello",
			entType: "bold",
			wantPre: "",
			wantCnt: "hello",
			wantPos: "",
		},
		{
			name:    "leading and trailing",
			text:    "  hello  ",
			entType: "bold",
			wantPre: "  ",
			wantCnt: "hello",
			wantPos: "  ",
		},
		{
			name:    "trailing tabs and spaces mixed",
			text:    "hello \t",
			entType: "bold",
			wantPre: "",
			wantCnt: "hello",
			wantPos: " \t",
		},
		{
			name:    "only spaces",
			text:    "   ",
			entType: "bold",
			wantPre: "   ",
			wantCnt: "",
			wantPos: "",
		},
		{
			name:    "newlines in pre",
			text:    "\n hello \n",
			entType: "pre",
			wantPre: "",
			wantCnt: "\n hello",
			wantPos: " \n",
		},
		{
			name:    "newlines in bold",
			text:    "\n hello \n",
			entType: "bold",
			wantPre: "\n ",
			wantCnt: "hello",
			wantPos: " \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, cnt, post := splitEdgeWhitespace(tt.text, MessageEntity{Type: tt.entType})
			if pre != tt.wantPre || cnt != tt.wantCnt || post != tt.wantPos {
				t.Errorf("splitEdgeWhitespace(%q, %q) = (%q, %q, %q); want (%q, %q, %q)",
					tt.text, tt.entType, pre, cnt, post, tt.wantPre, tt.wantCnt, tt.wantPos)
			}
		})
	}
}
