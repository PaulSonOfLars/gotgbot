package gotgbot

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRichBlockParsing(t *testing.T) {
	botToken := os.Getenv("BOT_TOKEN")
	chatIdStr := os.Getenv("CHAT_ID")

	var b *Bot
	var chatId int64
	if botToken == "" || chatIdStr == "" {
		t.Logf("BOT_TOKEN or CHAT_ID environment variable not set. Won't validate against live results.")
	} else {
		var err error
		b, err = NewBot(botToken, nil)
		if err != nil {
			t.Errorf("bot token provided, but failed to get bot: %s", err)
			t.FailNow()
			return
		}
		chatId, err = strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			t.Errorf("failed to parse chat id as integer: %s", err)
			t.FailNow()
			return
		}
		_, err = b.GetChat(chatId, nil)
		if err != nil {
			t.Errorf("failed to get chat: %s", err)
			t.FailNow()
			return
		}
	}

	tests := []struct {
		name         string
		b            RichBlock
		wantHTML     string
		wantMarkdown string
		wantText     string
		skipLiveTest bool
	}{
		{
			name:         "RichBlockParagraph",
			b:            RichBlockParagraph{Text: RichTextString("some text")},
			wantHTML:     `<p>some text</p>`,
			wantMarkdown: `some text`,
			wantText:     "some text",
		}, {
			name: "RichBlockSectionHeading",
			b: RichBlockSectionHeading{
				Text: RichTextString("header"),
				Size: 1,
			},
			wantHTML:     `<h1>header</h1>`,
			wantMarkdown: `# header`,
			wantText:     "header",
		}, {
			name: "RichBlockPreformatted",
			b: RichBlockPreformatted{
				Text:     RichTextString(`print("hi")`),
				Language: "python",
			},
			wantHTML:     "<pre><code class=\"language-python\">print(&#34;hi&#34;)</code></pre>",
			wantMarkdown: "```python\nprint(\"hi\")\n```",
			wantText:     "print(\"hi\")",
		}, {
			name: "RichBlockFooter",
			b: RichBlockFooter{
				Text: RichTextString("hello"),
			},
			wantHTML:     "<footer>hello</footer>",
			wantMarkdown: "---\nhello",
			wantText:     "hello",
		}, {
			name:         "RichBlockDivider",
			b:            RichBlockDivider{},
			wantHTML:     `<hr>`,
			wantMarkdown: `---`,
			wantText:     "",
			skipLiveTest: true, // empty!
		}, {
			name: "RichBlockMathematicalExpression",
			b: RichBlockMathematicalExpression{
				Expression: "1+2",
			},
			wantHTML:     `<tg-math-block>1+2</tg-math-block>`,
			wantMarkdown: "$$\n1+2\n$$",
			wantText:     "1+2",
		}, {
			name:         "RichBlockAnchor",
			b:            RichBlockAnchor{Name: "hello"},
			wantHTML:     `<a name="hello"></a>`,
			wantMarkdown: `<a name="hello"></a>`,
			wantText:     "",
			skipLiveTest: true, // empty
		}, {
			name: "RichBlockList",
			b: RichBlockList{
				Items: RichBlockListItemArray{
					{
						Label:       "",
						Blocks:      RichBlockArray{RichBlockParagraph{Text: RichTextString("some text")}},
						HasCheckbox: false,
						IsChecked:   false,
						Value:       0,
						Type:        "",
					}, {
						Label:       "",
						Blocks:      RichBlockArray{RichBlockParagraph{Text: RichTextString("some more text")}},
						HasCheckbox: false,
						IsChecked:   false,
						Value:       0,
						Type:        "",
					},
				},
			},
			wantHTML: "<ul>\n" +
				"<li>\n<p>some text</p>\n</li>\n" +
				"<li>\n<p>some more text</p>\n</li>\n" +
				"</ul>",
			wantMarkdown: "- some text\n- some more text",
			wantText:     "some text\nsome more text",
		}, {
			name: "RichBlockBlockQuotation",
			b: RichBlockBlockQuotation{
				Blocks: RichBlockArray{RichBlockParagraph{Text: RichTextString("some text")}},
				Credit: nil,
			},
			wantHTML:     "<blockquote>\n<p>some text</p>\n</blockquote>",
			wantMarkdown: `> some text`,
			wantText:     "some text",
		}, {
			name: "RichBlockPullQuotation",
			b: RichBlockPullQuotation{
				Text:   RichTextString("some text"),
				Credit: nil,
			},
			wantHTML:     `<aside>some text</aside>`,
			wantMarkdown: `<aside>some text</aside>`,
			wantText:     "some text",
		}, {
			name: "RichBlockCollage",
			b: RichBlockCollage{
				Blocks: RichBlockArray{
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
				},
				Caption: nil,
			},
			wantHTML:     "<tg-collage>\n<img src=\"fileId://1234\"></img>\n<img src=\"fileId://1234\"></img>\n</tg-collage>",
			wantMarkdown: "<tg-collage>\n![](fileId://1234)\n![](fileId://1234)\n</tg-collage>",
			wantText:     "",
		}, {
			name: "RichBlockSlideshow",
			b: RichBlockSlideshow{
				Blocks: RichBlockArray{
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
				},
				Caption: nil,
			},
			wantHTML:     "<tg-slideshow>\n<img src=\"fileId://1234\"></img>\n<img src=\"fileId://1234\"></img>\n</tg-slideshow>",
			wantMarkdown: "<tg-slideshow>\n![](fileId://1234)\n![](fileId://1234)\n</tg-slideshow>",
			wantText:     "",
		}, {
			name: "RichBlockTable",
			b: RichBlockTable{
				Cells: [][]RichBlockTableCell{
					{
						{Text: RichTextString("first column"), IsHeader: true, Colspan: 0, Rowspan: 0, Align: "", Valign: ""},
						{Text: RichTextString("second column"), IsHeader: true, Colspan: 0, Rowspan: 0, Align: "", Valign: ""},
					},
					{
						{Text: RichTextString("value")},
						{Text: RichTextString("value two")},
					}},
				IsBordered: false,
				IsStriped:  false,
				Caption:    nil,
			},
			wantHTML: "<table>\n" +
				"<tr>\n" +
				"<th align=\"center\" valign=\"middle\">first column</th>\n" +
				"<th align=\"center\" valign=\"middle\">second column</th>\n" +
				"</tr>\n" +
				"<tr>\n" +
				"<td align=\"center\" valign=\"middle\">value</td>\n" +
				"<td align=\"center\" valign=\"middle\">value two</td>\n" +
				"</tr>\n" +
				"</table>",
			wantMarkdown: `| first column | second column |
|:---:|:---:|
| value | value two |`,
			wantText: "first column\tsecond column\nvalue\tvalue two",
		}, {
			name: "RichBlockDetails",
			b: RichBlockDetails{
				Summary: RichTextString("some text"),
				Blocks:  RichBlockArray{RichBlockParagraph{Text: RichTextString("some more text")}},
				IsOpen:  false,
			},
			wantHTML:     "<details>\n<summary>some text</summary>\n<p>some more text</p>\n</details>",
			wantMarkdown: "<details>\n<summary>some text</summary>\nsome more text\n</details>",
			wantText:     "some text\nsome more text",
		}, {
			name: "RichBlockMap",
			b: RichBlockMap{
				Location: Location{
					Latitude:  41.9,
					Longitude: 12.5,
				},
				Zoom:    14,
				Caption: nil,
			},
			wantHTML:     `<tg-map lat="41.9" long="12.5" zoom="14"/>`,
			wantMarkdown: `<tg-map lat="41.9" long="12.5" zoom="14"/>`,
			wantText:     "",
			skipLiveTest: true, // Floating point errors
		}, {
			name: "RichBlockAnimation",
			b: RichBlockAnimation{
				Animation: Animation{
					FileId: "1234",
				},
				HasSpoiler: false,
				Caption:    nil,
			},
			wantHTML:     `<video src="fileId://1234"></video>`,
			wantMarkdown: `![](fileId://1234)`,
			wantText:     "",
		}, {
			name: "RichBlockAudio",
			b: RichBlockAudio{
				Audio: Audio{
					FileId: "1234",
				},
				Caption: nil,
			},
			wantHTML:     `<audio src="fileId://1234"></audio>`,
			wantMarkdown: `![](fileId://1234)`,
			wantText:     "",
		}, {
			name: "RichBlockPhoto",
			b: RichBlockPhoto{
				Photo: []PhotoSize{{
					FileId: "1234",
				}},
				HasSpoiler: false,
				Caption:    nil,
			},
			wantHTML:     `<img src="fileId://1234"></img>`,
			wantMarkdown: `![](fileId://1234)`,
			wantText:     "",
		}, {
			name: "RichBlockVideo",
			b: RichBlockVideo{
				Video: Video{
					FileId: "1234",
				},
				HasSpoiler: false,
				Caption:    nil,
			},
			wantHTML:     `<video src="fileId://1234"></video>`,
			wantMarkdown: `![](fileId://1234)`,
			wantText:     "",
		}, {
			name: "RichBlockVoiceNote",
			b: RichBlockVoiceNote{
				VoiceNote: Voice{
					FileId: "1234",
				},
				Caption: nil,
			},
			wantHTML:     `<audio src="fileId://1234"></audio>`,
			wantMarkdown: `![](fileId://1234)`,
			wantText:     "",
			skipLiveTest: true, // doesn't play nice with PMs
		}, {
			name:         "RichBlockThinking",
			b:            RichBlockThinking{Text: RichTextString("Thinking...")},
			wantHTML:     `<tg-thinking>Thinking...</tg-thinking>`,
			wantMarkdown: "<tg-thinking>Thinking\\.\\.\\.\n</tg-thinking>",
			wantText:     "Thinking...",
			skipLiveTest: true, // only for sendRichMessageDraft
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("HTML", func(t *testing.T) {
				if got := strings.TrimSpace(RichBlockHTML(tt.b)); got != tt.wantHTML {
					t.Errorf("RichBlockHTML() = %v, want %v", got, tt.wantHTML)
				}
			})

			t.Run("Markdown", func(t *testing.T) {
				if got := strings.TrimSpace(RichBlockMarkdown(tt.b)); got != tt.wantMarkdown {
					t.Errorf("RichBlockMarkdown() = %v, want %v", got, tt.wantMarkdown)
				}
			})

			t.Run("Text", func(t *testing.T) {
				if got := strings.TrimSpace(RichBlockContent(tt.b)); got != tt.wantText {
					t.Errorf("RichBlockContent() = %v, want %v", got, tt.wantText)
				}
			})

			t.Run("Bot API", func(t *testing.T) {
				if tt.skipLiveTest {
					return
				}

				if b == nil {
					t.Skip("set BOT_TOKEN and CHAT_ID to test against bot API")
					t.SkipNow()
					return
				}

				time.Sleep(time.Second) // sleep one second to avoid rate limits

				m, err := b.SendRichMessage(chatId, InputRichMessage{
					Html: strings.Replace(tt.wantHTML, "fileId://1234", getFile(t, tt.name), -1),
				}, nil)
				if err != nil {
					t.Fatal(err)
				}

				if got := replaceLiveId(t, strings.TrimSpace(m.RichMessage.HTML())); got != tt.wantHTML {
					t.Errorf("HTML() = %v, want %v", got, tt.wantHTML)
				}

				if got := replaceLiveId(t, strings.TrimSpace(m.RichMessage.Markdown())); got != tt.wantMarkdown {
					t.Errorf("Markdown() = %v, want %v", got, tt.wantMarkdown)
				}

				// no ids to change
				if got := strings.TrimSpace(m.RichMessage.PlainText()); got != tt.wantText {
					t.Errorf("Text() = %v, want %v", got, tt.wantText)
				}
			})
		})
	}
}

var fileRe = regexp.MustCompile("fileId://[A-Za-z0-9_-]+")

func replaceLiveId(t *testing.T, text string) string {
	t.Helper()
	return fileRe.ReplaceAllString(text, `fileId://1234`)
}

func getFile(t *testing.T, kind string) string {
	t.Helper()
	switch kind {
	case "RichBlockAudio":
		return "https://telegram.org/example/audio.mp3"
	case "RichBlockAnimation":
		return "https://telegram.org/example/animation.gif"
	case "RichBlockPhoto", "RichBlockCollage", "RichBlockSlideshow":
		return "https://telegram.org/example/photo.jpg"
	case "RichBlockVideo":
		return "https://telegram.org/example/video.mp4"
	case "RichBlockVoiceNote":
		return "https://telegram.org/example/audio.ogg"
	default:
		// assume ok!
		return ""
	}
}

func TestRichTextMarkdown(t *testing.T) {
	tests := []struct {
		name string
		r    RichText
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RichTextMarkdown(tt.r); got != tt.want {
				t.Errorf("RichTextMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
