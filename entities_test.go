package gotgbot

import "testing"

func Test_parseEntity_Bounds(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		entity   MessageEntity
		wantText string
	}{
		{
			name: "offset beyond text length",
			text: "short",
			entity: MessageEntity{
				Offset: 10,
				Length: 50,
			},
			wantText: "",
		},
		{
			name: "length beyond text length",
			text: "hello world",
			entity: MessageEntity{
				Offset: 6,
				Length: 100,
			},
			wantText: "world",
		},
		{
			name: "negative offset",
			text: "sample text",
			entity: MessageEntity{
				Offset: -5,
				Length: 6,
			},
			wantText: "s",
		},
		{
			name: "exact bounds",
			text: "hello",
			entity: MessageEntity{
				Offset: 0,
				Length: 5,
			},
			wantText: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := ParseEntity(tt.text, tt.entity)
			if pe.Text != tt.wantText {
				t.Errorf("ParseEntity(%q, %+v) got text %q, want %q", tt.text, tt.entity, pe.Text, tt.wantText)
			}
		})
	}
}
