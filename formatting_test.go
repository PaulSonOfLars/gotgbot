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
