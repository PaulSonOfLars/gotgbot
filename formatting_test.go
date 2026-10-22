package gotgbot

import (
	"reflect"
	"testing"
)

func Test_entitiesToRichBlocks(t *testing.T) {
	type args struct {
		text     string
		entities []MessageEntity
	}
	tests := []struct {
		name string
		args args
		want []RichBlock
	}{
		{
			name: "no formatting",
			args: args{
				text:     "hello, this is some content",
				entities: nil,
			},
			want: []RichBlock{
				RichBlockParagraph{Text: RichTextString("hello, this is some content")},
			},
		}, {
			name: "basic italic formatting",
			args: args{
				text: "hello there",
				entities: []MessageEntity{{
					Type:   "italic",
					Offset: 0,
					Length: 5,
				}},
			},
			want: []RichBlock{
				RichBlockParagraph{Text: RichTextArray{
					RichTextItalic{Text: RichTextString("hello")},
					RichTextString(" there")},
				},
			},
		}, {
			name: "basic username formatting",
			args: args{
				text: "@hello there",
				entities: []MessageEntity{{
					Type:   "mention",
					Offset: 0,
					Length: 6,
				}},
			},
			want: []RichBlock{
				RichBlockParagraph{Text: RichTextArray{
					RichTextMention{Text: RichTextString("@hello"), Username: "hello"},
					RichTextString(" there")},
				},
			},
		}, {
			name: "single pre block",
			args: args{
				text: "hello",
				entities: []MessageEntity{{
					Type:   "pre",
					Offset: 0,
					Length: 5,
				}},
			},
			want: []RichBlock{
				RichBlockPreformatted{Text: RichTextString("hello")},
			},
		}, {
			name: "combined pre and italic",
			args: args{
				text: "hello there",
				entities: []MessageEntity{{
					Type:   "pre",
					Offset: 0,
					Length: 5,
				}, {
					Type:   "italic",
					Offset: 6,
					Length: 5,
				}},
			},
			want: []RichBlock{
				RichBlockPreformatted{Text: RichTextString("hello")},
				RichBlockParagraph{Text: RichTextArray{RichTextString(" "), RichTextItalic{Text: RichTextString("there")}}},
			},
		}, {
			name: "nested blockquote",
			args: args{
				text: "hello there you",
				entities: []MessageEntity{{
					Type:   "blockquote",
					Offset: 0,
					Length: 15,
				}, {
					Type:   "italic",
					Offset: 6,
					Length: 5,
				}},
			},
			want: []RichBlock{
				RichBlockBlockQuotation{Blocks: []RichBlock{
					RichBlockParagraph{Text: RichTextArray{
						RichTextString("hello "),
						RichTextItalic{Text: RichTextString("there")},
						RichTextString(" you"),
					}},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entitiesToRichBlocks(tt.args.text, tt.args.entities)
			if len(got) != len(tt.want) {
				t.Errorf("incorrect number of entities got=%d, want=%d", len(got), len(tt.want))
			}

			for idx, want := range tt.want {
				if want.GetType() != got[idx].GetType() {
					t.Errorf("incorrect type got=%s, want=%s", got[idx].GetType(), want.GetType())
				}
				if !reflect.DeepEqual(got[idx], want) {
					t.Errorf("incorrect entity got=%+v, want=%+v", got[idx], want)
				}
			}
		})
	}
}

func Test_splitEdgeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		ent      MessageEntity
		wantPre  string
		wantCont string
		wantPost string
	}{
		{
			name:     "regular",
			text:     "  hello  ",
			ent:      MessageEntity{Type: "bold"},
			wantPre:  "  ",
			wantCont: "hello",
			wantPost: "  ",
		},
		{
			name:     "multiple trailing whitespace types",
			text:     "  hello \t\n ",
			ent:      MessageEntity{Type: "bold"},
			wantPre:  "  ",
			wantCont: "hello",
			wantPost: " \t\n ",
		},
		{
			name:     "all whitespace",
			text:     "   ",
			ent:      MessageEntity{Type: "bold"},
			wantPre:  "   ",
			wantCont: "",
			wantPost: "",
		},
		{
			name:     "pre with newlines",
			text:     " \ncode\n ",
			ent:      MessageEntity{Type: "pre"},
			wantPre:  " ",
			wantCont: "\ncode\n",
			wantPost: " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, cont, post := splitEdgeWhitespace(tt.text, tt.ent)
			if pre != tt.wantPre || cont != tt.wantCont || post != tt.wantPost {
				t.Errorf("splitEdgeWhitespace(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.text, pre, cont, post, tt.wantPre, tt.wantCont, tt.wantPost)
			}
		})
	}
}

func Test_OriginalMD_Pre(t *testing.T) {
	m := Message{
		Text: "print('hello')",
		Entities: []MessageEntity{{
			Type:   "pre",
			Offset: 0,
			Length: 14,
		}},
	}
	got := m.OriginalMD()
	want := "```print('hello')```"
	if got != want {
		t.Errorf("OriginalMD() = %q, want %q", got, want)
	}
}

func Test_wrapEntity_NilUser(t *testing.T) {
	// Should not panic on nil User
	ent := MessageEntity{
		Type: MessageEntityTypeTextMention,
	}
	res := wrapEntity(ent, RichTextString("name"))
	tm, ok := res.(RichTextTextMention)
	if !ok {
		t.Fatalf("expected RichTextTextMention, got %T", res)
	}
	if tm.User.Id != 0 {
		t.Errorf("expected 0 user id, got %d", tm.User.Id)
	}
}

func Test_parseEntity_Bounds(t *testing.T) {
	// Out of bounds offsets shouldn't panic
	ent := MessageEntity{
		Offset: 10,
		Length: 50,
	}
	pe := ParseEntity("short", ent)
	if pe.Text != "" {
		t.Errorf("expected empty text, got %q", pe.Text)
	}
}
