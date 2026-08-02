package gotgbot

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRichBlockParsing(t *testing.T) {
	b, chatId := tryLocalBot(t)

	tests := []struct {
		name         string
		b            RichBlock
		wantHTML     string
		wantMarkdown string
		wantText     string
		wantMedia    []InputRichMessageMedia
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
			wantMarkdown: "<footer>hello</footer>", // HTML syntax
			wantText:     "hello",
		}, {
			name:         "RichBlockDivider",
			b:            RichBlockDivider{},
			wantHTML:     `<hr/>`,
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
			wantHTML: "<tg-collage>\n<img src=\"tg://photo?id=1\"></img>\n<img src=\"tg://photo?id=2\"></img>\n</tg-collage>",
			// collage unreliable in markdown
			wantMarkdown: "<tg-collage>\n<img src=\"tg://photo?id=1\"></img>\n<img src=\"tg://photo?id=2\"></img>\n</tg-collage>",
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{Media: InputFileByID("1234")}},
				{Id: "2", Media: InputMediaPhoto{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockSlideshow",
			b: RichBlockSlideshow{
				Blocks: RichBlockArray{
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
					RichBlockPhoto{Photo: []PhotoSize{{FileId: "1234"}}},
				},
				Caption: nil,
			},
			wantHTML: "<tg-slideshow>\n<img src=\"tg://photo?id=1\"></img>\n<img src=\"tg://photo?id=2\"></img>\n</tg-slideshow>",
			// collage unreliable in markdown
			wantMarkdown: "<tg-slideshow>\n<img src=\"tg://photo?id=1\"></img>\n<img src=\"tg://photo?id=2\"></img>\n</tg-slideshow>",
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{Media: InputFileByID("1234")}},
				{Id: "2", Media: InputMediaPhoto{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockTable",
			b: RichBlockTable{
				Cells: [][]RichBlockTableCell{
					{
						{Text: RichTextString("first column"), IsHeader: true, Colspan: 0, Rowspan: 0, Align: "left"},
						{Text: RichTextString("second column"), IsHeader: true, Colspan: 0, Rowspan: 0, Align: "left"},
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
				"<th align=\"left\">first column</th>\n" +
				"<th align=\"left\">second column</th>\n" +
				"</tr>\n" +
				"<tr>\n" +
				"<td>value</td>\n" +
				"<td>value two</td>\n" +
				"</tr>\n" +
				"</table>",
			wantMarkdown: `| first column | second column |
|:---|:---|
| value | value two |`,
			wantText: "first column\tsecond column\nvalue\tvalue two",
		}, {
			name: "RichBlockDetails",
			b: RichBlockDetails{
				Summary: RichTextString("some text"),
				Blocks:  RichBlockArray{RichBlockParagraph{Text: RichTextString("some more text")}},
				IsOpen:  false,
			},
			wantHTML: "<details>\n<summary>some text</summary>\n<p>some more text</p>\n</details>",
			// identical; details can't be displayed in markdown.
			wantMarkdown: "<details>\n<summary>some text</summary>\n<p>some more text</p>\n</details>",
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
			wantHTML:     `<video src="tg://video?id=1"></video>`,
			wantMarkdown: `![](tg://video?id=1)`,
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaAnimation{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockAudio",
			b: RichBlockAudio{
				Audio: Audio{
					FileId: "1234",
				},
				Caption: nil,
			},
			wantHTML:     `<audio src="tg://audio?id=1"></audio>`,
			wantMarkdown: `![](tg://audio?id=1)`,
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaAudio{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockPhoto",
			b: RichBlockPhoto{
				Photo: []PhotoSize{{
					FileId: "1234",
				}},
				HasSpoiler: false,
				Caption:    nil,
			},
			wantHTML:     `<img src="tg://photo?id=1"></img>`,
			wantMarkdown: `![](tg://photo?id=1)`,
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockVideo",
			b: RichBlockVideo{
				Video: Video{
					FileId: "1234",
				},
				HasSpoiler: false,
				Caption:    nil,
			},
			wantHTML:     `<video src="tg://video?id=1"></video>`,
			wantMarkdown: `![](tg://video?id=1)`,
			wantText:     "",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaVideo{Media: InputFileByID("1234")}},
			},
		}, {
			name: "RichBlockVoiceNote",
			b: RichBlockVoiceNote{
				VoiceNote: Voice{
					FileId: "1234",
				},
				Caption: nil,
			},
			wantHTML:     `<audio src="tg://audio?id=1"></audio>`,
			wantMarkdown: `![](tg://audio?id=1)`,
			wantText:     "",
			skipLiveTest: true, // doesn't play nice with PMs
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaVoiceNote{Type: "voice_note", Media: InputFileByID("1234")}},
			},
		}, {
			name:         "RichBlockThinking",
			b:            RichBlockThinking{Text: RichTextString("Thinking...")},
			wantHTML:     `<tg-thinking>Thinking...</tg-thinking>`,
			wantMarkdown: "<tg-thinking>Thinking...</tg-thinking>",
			wantText:     "Thinking...",
			skipLiveTest: true, // only for sendRichMessageDraft
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("HTML", func(t *testing.T) {
				html, media := RichBlockHTML(tt.b)
				if got := strings.TrimSpace(html); got != tt.wantHTML {
					t.Errorf("HTML() mismatch:\n got: %q\nwant: %q", got, tt.wantHTML)
				}
				_ = media
			})

			t.Run("Markdown", func(t *testing.T) {
				markdown, media := RichBlockMarkdown(tt.b)
				if got := strings.TrimSpace(markdown); got != tt.wantMarkdown {
					t.Errorf("Markdown() mismatch:\n got: %q\nwant: %q", got, tt.wantMarkdown)
				}
				_ = media
			})

			t.Run("Text", func(t *testing.T) {
				if got := strings.TrimSpace(RichBlockContent(tt.b)); got != tt.wantText {
					t.Errorf("Content() mismatch:\n got: %q\nwant: %q", got, tt.wantText)
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
					Html: tt.wantHTML,
					Media: func() []InputRichMessageMedia {
						var ms []InputRichMessageMedia
						for idx := range tt.wantMedia {
							ms = append(ms, InputRichMessageMedia{Id: strconv.Itoa(idx + 1), Media: getFile(t, tt.name)})
						}
						return ms
					}(),
				}, nil)
				if err != nil {
					t.Fatal(err)
				}

				html, htmlMedia := m.RichMessage.HTML()
				if got := strings.TrimSpace(html); got != tt.wantHTML {
					t.Errorf("HTML() mismatch:\n got: %q\nwant: %q", got, tt.wantHTML)
				}
				equalMedias(t, htmlMedia, tt.wantMedia)

				markdown, mdMedia := m.RichMessage.Markdown()
				if got := strings.TrimSpace(markdown); got != tt.wantMarkdown {
					t.Errorf("Markdown() mismatch:\n got: %q\nwant: %q", got, tt.wantMarkdown)
				}
				equalMedias(t, mdMedia, tt.wantMedia)

				// no ids to change
				if got := strings.TrimSpace(m.RichMessage.PlainText()); got != tt.wantText {
					t.Errorf("PlainText() mismatch:\n got: %q\nwant: %q", got, tt.wantText)
				}
			})
		})
	}
}

func equalMedias(t *testing.T, got, want []InputRichMessageMedia) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("len(a) != len(b): got %d, want %d", len(got), len(want))
	}
	for idx, want := range want {
		if want.Id != got[idx].Id {
			t.Errorf("incorrect id got=%s, want=%s", got[idx].Id, want.Id)
		}
		if want.Media.GetType() != got[idx].Media.GetType() {
			t.Errorf("incorrect type got=%s, want=%s", got[idx].Media.GetType(), want.Media.GetType())
		}
		// We don't need the media to match more than this; IDs will change.
	}
}

func tryLocalBot(t *testing.T) (*Bot, int64) {
	t.Helper()

	botToken := os.Getenv("BOT_TOKEN")
	chatIdStr := os.Getenv("CHAT_ID")

	if botToken == "" || chatIdStr == "" {
		t.Logf("BOT_TOKEN or CHAT_ID environment variable not set. Won't validate against live results.")
	} else {
		b, err := NewBot(botToken, nil)
		if err != nil {
			t.Errorf("bot token provided, but failed to get bot: %s", err)
			t.FailNow()
		}

		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			t.Errorf("failed to parse chat id as integer: %s", err)
			t.FailNow()
		}

		_, err = b.GetChat(chatId, nil)
		if err != nil {
			t.Errorf("failed to get chat: %s", err)
			t.FailNow()
		}

		return b, chatId
	}

	return nil, 0
}

func getFile(t *testing.T, kind string) InputMedia {
	t.Helper()
	switch kind {
	case "RichBlockAudio":
		return InputMediaAudio{Media: InputFileByURL("https://telegram.org/example/audio.mp3")}
	case "RichBlockAnimation":
		return InputMediaAnimation{Media: InputFileByURL("https://telegram.org/example/animation.gif")}
	case "RichBlockPhoto", "RichBlockCollage", "RichBlockSlideshow":
		return InputMediaPhoto{Media: InputFileByURL("https://telegram.org/example/photo.jpg")}
	case "RichBlockVideo":
		return InputMediaVideo{Media: InputFileByURL("https://telegram.org/example/video.mp4")}
	case "RichBlockVoiceNote":
		return InputMediaVoiceNote{Type: "voice_note", Media: InputFileByURL("https://telegram.org/example/audio.ogg")}
	default:
		// assume nothing?
		return nil
	}
}

func TestRichTextParsing(t *testing.T) {
	tests := []struct {
		name         string
		r            RichText
		wantHTML     string
		wantMarkdown string
		wantText     string
	}{
		{
			name:         "RichTextBold",
			r:            RichTextBold{Text: RichTextString("hi")},
			wantHTML:     "<b>hi</b>",
			wantMarkdown: "**hi**",
			wantText:     "hi",
		}, {
			name:         "RichTextItalic",
			r:            RichTextItalic{Text: RichTextString("hi")},
			wantHTML:     "<i>hi</i>",
			wantMarkdown: "_hi_",
			wantText:     "hi",
		}, {
			name:         "RichTextUnderline",
			r:            RichTextUnderline{Text: RichTextString("hi")},
			wantHTML:     "<u>hi</u>",
			wantMarkdown: "<u>hi</u>",
			wantText:     "hi",
		}, {
			name:         "RichTextStrikethrough",
			r:            RichTextStrikethrough{Text: RichTextString("hi")},
			wantHTML:     "<s>hi</s>",
			wantMarkdown: "~~hi~~",
			wantText:     "hi",
		}, {
			name:         "RichTextSpoiler",
			r:            RichTextSpoiler{Text: RichTextString("hi")},
			wantHTML:     "<tg-spoiler>hi</tg-spoiler>",
			wantMarkdown: "||hi||",
			wantText:     "hi",
		}, {
			name: "RichTextDateTime",
			r: RichTextDateTime{
				Text:           RichTextString("22:45 tomorrow"),
				UnixTime:       1647531900,
				DateTimeFormat: "wDT",
			},
			wantHTML:     "<tg-time unix=\"1647531900\" format=\"wDT\">22:45 tomorrow</tg-time>",
			wantMarkdown: "![22:45 tomorrow](tg://time?unix=1647531900format=wDT)",
			wantText:     "22:45 tomorrow",
		}, {
			name: "RichTextTextMention",
			r: RichTextTextMention{
				Text: RichTextString("some text"),
				User: User{Id: 1234},
			},
			wantHTML:     "<a href=\"tg://user?id=1234\">some text</a>",
			wantMarkdown: "[some text](tg://user?id=1234)",
			wantText:     "some text",
		}, {
			name:         "RichTextSubscript",
			r:            RichTextSubscript{Text: RichTextString("some text")},
			wantHTML:     "<sub>some text</sub>",
			wantMarkdown: "<sub>some text</sub>",
			wantText:     "some text",
		}, {
			name:         "RichTextSuperscript",
			r:            RichTextSuperscript{Text: RichTextString("some text")},
			wantHTML:     "<sup>some text</sup>",
			wantMarkdown: "<sup>some text</sup>",
			wantText:     "some text",
		}, {
			name:         "RichTextMarked",
			r:            RichTextMarked{Text: RichTextString("some text")},
			wantHTML:     "<mark>some text</mark>",
			wantMarkdown: "==some text==",
			wantText:     "some text",
		}, {
			name:         "RichTextCode",
			r:            RichTextCode{Text: RichTextString("some text")},
			wantHTML:     "<code>some text</code>",
			wantMarkdown: "`some text`",
			wantText:     "some text",
		}, {
			name: "RichTextCustomEmoji",
			r: RichTextCustomEmoji{
				CustomEmojiId:   "1234",
				AlternativeText: "hi",
			},
			wantHTML:     "<tg-emoji emoji-id=\"1234\">hi</tg-emoji>",
			wantMarkdown: "![hi](tg://emoji?id=1234)",
			wantText:     "hi",
		}, {
			name: "RichTextMathematicalExpression",
			r: RichTextMathematicalExpression{
				Expression: "1+2",
			},
			wantHTML:     "<tg-math>1+2</tg-math>",
			wantMarkdown: "$$1+2$$",
			wantText:     "1+2",
		}, {
			name: "RichTextUrl",
			r: RichTextUrl{
				Text: RichTextString("hello"),
				Url:  "example.com",
			},
			wantHTML:     "<a href=\"example.com\">hello</a>",
			wantMarkdown: "[hello](example.com)",
			wantText:     "hello",
		}, {
			name: "RichTextEmailAddress",
			r: RichTextEmailAddress{
				Text:         RichTextString("some text"),
				EmailAddress: "some@email.com",
			},
			wantHTML:     "<a href=\"mailto:some@email.com\">some text</a>",
			wantMarkdown: "[some text](mailto:some@email.com)",
			wantText:     "some text",
		}, {
			name: "RichTextPhoneNumber",
			r: RichTextPhoneNumber{
				Text:        RichTextString("some text"),
				PhoneNumber: "1234",
			},
			wantHTML:     "<a href=\"tel:1234\">some text</a>",
			wantMarkdown: "[some text](tel:1234)",
			wantText:     "some text",
		}, {
			name: "RichTextBankCardNumber",
			r: RichTextBankCardNumber{
				Text:           RichTextString("4242 4242 4242"),
				BankCardNumber: "4242 4242 4242",
			},
			wantHTML:     "4242 4242 4242",
			wantMarkdown: "4242 4242 4242",
			wantText:     "4242 4242 4242",
		}, {
			name: "RichTextMention",
			r: RichTextMention{
				Text:     RichTextString("@someuser"),
				Username: "someuser",
			},
			wantHTML:     "@someuser",
			wantMarkdown: "@someuser",
			wantText:     "@someuser",
		}, {
			name: "RichTextHashtag",
			r: RichTextHashtag{
				Text:    RichTextString("#hash"),
				Hashtag: "hash",
			},
			wantHTML:     "#hash",
			wantMarkdown: "\\#hash",
			wantText:     "#hash",
		}, {
			name: "RichTextCashtag",
			r: RichTextCashtag{
				Text:    RichTextString("$CASH"),
				Cashtag: "CASH",
			},
			wantHTML:     "$CASH",
			wantMarkdown: "$CASH",
			wantText:     "$CASH",
		}, {
			name: "RichTextBotCommand",
			r: RichTextBotCommand{
				Text:       RichTextString("/command"),
				BotCommand: "command",
			},
			wantHTML:     "/command",
			wantMarkdown: "/command",
			wantText:     "/command",
		}, {
			name: "RichTextAnchor",
			r: RichTextAnchor{
				Name: "test",
			},
			wantHTML:     "<a name=\"test\"></a>",
			wantMarkdown: "<a name=\"test\"></a>",
			wantText:     "",
		}, {
			name: "RichTextAnchorLink",
			r: RichTextAnchorLink{
				Text:       RichTextString("some text"),
				AnchorName: "test",
			},
			wantHTML:     "<a href=\"#test\">some text</a>",
			wantMarkdown: "[some text](#test)",
			wantText:     "some text",
		}, {
			name: "RichTextReference",
			r: RichTextReference{
				Text: RichTextString("some text"),
				Name: "test",
			},
			wantHTML:     "<tg-reference name=\"test\">some text</tg-reference>",
			wantMarkdown: "[^test]: some text",
			wantText:     "some text",
		}, {
			name: "RichTextReferenceLink",
			r: RichTextReferenceLink{
				Text:          RichTextString("some text"),
				ReferenceName: "test",
			},
			wantHTML:     "<a href=\"#test\">some text</a>",
			wantMarkdown: "some text[^test]",
			wantText:     "some text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("HTML", func(t *testing.T) {
				if got := strings.TrimSpace(RichTextHTML(tt.r)); got != tt.wantHTML {
					t.Errorf("RichTextHTML() = %v, want %v", got, tt.wantHTML)
				}
			})

			t.Run("Markdown", func(t *testing.T) {
				if got := strings.TrimSpace(RichTextMarkdown(tt.r)); got != tt.wantMarkdown {
					t.Errorf("RichTextMarkdown() = %v, want %v", got, tt.wantMarkdown)
				}
			})

			t.Run("Text", func(t *testing.T) {
				if got := strings.TrimSpace(RichTextContent(tt.r)); got != tt.wantText {
					t.Errorf("RichTextContent() = %v, want %v", got, tt.wantText)
				}
			})
		})
	}
}

func TestRichMessageSending(t *testing.T) {
	b, chatId := tryLocalBot(t)

	if b == nil {
		t.Skip("Skiping message sending - no bot configured.")
		return
	}

	tests := []struct {
		name         string
		inputHTML    string
		wantHTML     string
		wantMarkdown string
		wantText     string
		wantMedia    []InputRichMessageMedia
	}{
		{
			name: "basic 1",
			inputHTML: `<a name="chapter-0"></a>
				<b>bold text</b>, <strong>bold text</strong>
				<i>italic text</i>, <em>italic text</em>
				<u>underlined text</u>, <ins>underlined text</ins>
				<s>strikethrough text</s>, <strike>strikethrough text</strike>, <del>strikethrough text</del>`,
			wantHTML: "<p><a name=\"chapter-0\"></a> " +
				"<b>bold text</b>, <b>bold text</b> " +
				"<i>italic text</i>, <i>italic text</i> " +
				"<u>underlined text</u>, <u>underlined text</u> " +
				"<s>strikethrough text</s>, <s>strikethrough text</s>, <s>strikethrough text</s></p>",
			wantMarkdown: "<a name=\"chapter-0\"></a> " +
				"**bold text**, **bold text** " +
				"_italic text_, _italic text_ " +
				"<u>underlined text</u>, <u>underlined text</u> " +
				"~~strikethrough text~~, ~~strikethrough text~~, ~~strikethrough text~~",
			wantText: "bold text, bold text italic text, italic text underlined text, underlined text strikethrough text, strikethrough text, strikethrough text",
		}, {
			name: "basic 2",
			inputHTML: `<code>inline fixed-width code</code>
				<mark>marked text</mark>
				<sub>subscript text</sub>
				<sup>superscript text</sup>
				<tg-spoiler>spoiler</tg-spoiler>`,
			wantHTML: "<p><code>inline fixed-width code</code> " +
				"<mark>marked text</mark> " +
				"<sub>subscript text</sub> " +
				"<sup>superscript text</sup>" +
				" <tg-spoiler>spoiler</tg-spoiler></p>",
			wantMarkdown: "`inline fixed\\-width code` " +
				"==marked text== " +
				"<sub>subscript text</sub> " +
				"<sup>superscript text</sup> " +
				"||spoiler||",
			wantText: "inline fixed-width code marked text subscript text superscript text spoiler",
		}, {
			name: "anchors",
			inputHTML: `<a href="#note-1">Reference</a>
				<a href="https://t.me/">inline URL</a>
				<a href="mailto:user@example.com">inline e-mail</a>
				<a href="tel:+123456789">inline phone number</a>
				<a href="tg://user?id=123456789">inline mention of a user</a>
				<a href="#chapter-1">in-document link</a>
				<a name="chapter-1"></a>`,
			wantHTML: "<p><a href=\"#note-1\">Reference</a> " +
				"<a href=\"https://t.me/\">inline URL</a> " +
				"<a href=\"mailto:user@example.com\">inline e-mail</a> " +
				"<a href=\"tel:+123456789\">inline phone number</a> " +
				"inline mention of a user " + // never seen this user; can't mention!
				"<a href=\"#chapter-1\">in-document link</a> " +
				"<a name=\"chapter-1\"></a></p>",
			wantMarkdown: "[Reference](#note-1) " +
				"[inline URL](https://t.me/) " +
				"[inline e\\-mail](mailto:user@example.com) " +
				"[inline phone number](tel:+123456789) " +
				"inline mention of a user " +
				"[in\\-document link](#chapter-1) " +
				"<a name=\"chapter-1\"></a>",
			wantText: "Reference " +
				"inline URL " +
				"inline e-mail " +
				"inline phone number " +
				"inline mention of a user " +
				"in-document link",
		}, {
			name: "tg-types",
			inputHTML: `<tg-reference name="note-1">Referenced text</tg-reference>
				<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>
				<img src="tg://emoji?id=5368324170671202286" alt="👍"/>
				<tg-time unix="1647531900" format="wDT">22:45 tomorrow</tg-time>
				<tg-math>x^2 + y^2</tg-math>`,
			wantHTML: "<p><tg-reference name=\"note-1\">Referenced text</tg-reference> " +
				"<tg-emoji emoji-id=\"5368324170671202286\">👍</tg-emoji> " +
				"<tg-emoji emoji-id=\"5368324170671202286\">👍</tg-emoji> " +
				"<tg-time unix=\"1647531900\" format=\"wDT\">22:45 tomorrow</tg-time> " +
				"<tg-math>x^2 + y^2</tg-math></p>",
			wantMarkdown: "[^note-1]: Referenced text " +
				"![👍](tg://emoji?id=5368324170671202286) " +
				"![👍](tg://emoji?id=5368324170671202286) " +
				"![22:45 tomorrow](tg://time?unix=1647531900format=wDT) " +
				"$$x^2 + y^2$$",
			wantText: `Referenced text 👍 👍 22:45 tomorrow x^2 + y^2`,
		}, {
			name: "tg-types strings",
			inputHTML: `
				#hashtag $USD +12345678901, card: 4242 4242 4242 4242, https://t.me t.me a@t.me /command @username
			
				all the text above was on the same line`,
			wantHTML: "<p>#hashtag $USD " +
				"<a href=\"tel:+12345678901\">+12345678901</a>, " +
				"card: 4242 4242 4242 4242, " +
				"<a href=\"https://t.me\">https://t.me</a> " +
				"<a href=\"t.me\">t.me</a> " +
				"<a href=\"mailto:a@t.me\">a@t.me</a> /command " +
				"@username " +
				"all the text above was on the same line</p>",
			wantMarkdown: "\\#hashtag $USD " +
				"[\\+12345678901](tel:+12345678901), " +
				"card: 4242 4242 4242 4242, " +
				"[https://t.me](https://t.me) " +
				"[t.me](t.me) " +
				"[a@t.me](mailto:a@t.me) " +
				"/command " +
				"@username " +
				"all the text above was on the same line",
			wantText: "#hashtag $USD +12345678901, " +
				"card: 4242 4242 4242 4242, " +
				"https://t.me t.me a@t.me " +
				"/command @username " +
				"all the text above was on the same line",
		}, {
			name: "headings",
			inputHTML: `<h1>Heading 1</h1>
				<h2>Heading 2</h2>
				<h3>Heading 3</h3>
				<h4>Heading 4</h4>
				<h5>Heading 5</h5>
				<h6>Heading 6</h6>`,
			wantHTML:     "<h1>Heading 1</h1>\n<h2>Heading 2</h2>\n<h3>Heading 3</h3>\n<h4>Heading 4</h4>\n<h5>Heading 5</h5>\n<h6>Heading 6</h6>",
			wantMarkdown: "# Heading 1\n## Heading 2\n### Heading 3\n#### Heading 4\n##### Heading 5\n###### Heading 6",
			wantText:     "Heading 1\nHeading 2\nHeading 3\nHeading 4\nHeading 5\nHeading 6",
		}, {
			name: "html basics",
			inputHTML: `<p>Paragraph text</p>
				<pre>pre-formatted fixed-width code block</pre>
				<pre><code class="language-python">  print('pre-formatted fixed-width code block written in the Python programming language')</code></pre>
				<footer>Footer text</footer>
				<hr/>
				<ul><li>unordered list item</li></ul>
				<ol><li>ordered list item</li></ol>
				<ol start="3" type="a" reversed><li>ordered list item</li></ol>
				<ol><li value="7" type="i">ordered list item with explicit number</li></ol>
				<ul>
				<li><input type="checkbox" checked>Checked checkbox</li>
				<li><input type="checkbox">Unchecked checkbox</li>
				</ul>`,
			wantHTML: "<p>Paragraph text</p>\n" +
				"<pre><code>pre-formatted fixed-width code block</code></pre>\n" +
				"<pre><code class=\"language-python\">  " +
				"print(&#39;pre-formatted fixed-width code block written in the Python programming language&#39;)" +
				"</code></pre>\n" +
				"<footer>Footer text</footer>\n" +
				"<hr/>\n" +
				"<ul>\n<li>\n<p>unordered list item</p>\n</li>\n</ul>\n" +
				"<ol>\n<li>\n<p>ordered list item</p>\n</li>\n</ol>\n" +
				"<ol start=\"3\" type=\"a\" reversed>\n<li value=\"3\" type=\"a\">\n<p>ordered list item</p>\n</li>\n</ol>\n" +
				"<ol start=\"7\" type=\"i\" reversed>\n<li value=\"7\" type=\"i\">\n<p>ordered list item with explicit number</p>\n</li>\n</ol>\n" +
				"<ul>\n<li><input type=\"checkbox\" checked> \n<p>Checked checkbox</p>\n</li>\n" +
				"<li><input type=\"checkbox\"> \n<p>Unchecked checkbox</p>\n</li>\n" +
				"</ul>",
			wantMarkdown: "Paragraph text\n" +
				"```\npre-formatted fixed-width code block\n```\n" +
				"```python\n  print('pre-formatted fixed-width code block written in the Python programming language')\n```" +
				"\n<footer>Footer text</footer>\n" +
				"\n---\n" +
				"- unordered list item\n\n" +
				"1. ordered list item\n\n" +
				"c. ordered list item\n\n" +
				"vii. ordered list item with explicit number\n\n" +
				"- [x] Checked checkbox\n" +
				"- [ ] Unchecked checkbox",
			wantText: "Paragraph text\n" +
				"pre-formatted fixed-width code block\n" +
				"  print('pre-formatted fixed-width code block written in the Python programming language')\n" +
				"Footer text\n\n" +
				"unordered list item\n\n" +
				"ordered list item\n\n" +
				"ordered list item\n\n" +
				"ordered list item with explicit number\n\n" +
				"Checked checkbox\n" +
				"Unchecked checkbox",
		}, {
			name: "quotes",
			inputHTML: `<blockquote>Block quotation started<br>Block quotation continued<br>The last line of the block quotation<cite>The Author</cite></blockquote>
				<aside>Pull quote<cite>The Author</cite></aside>`,
			wantHTML: "<blockquote>\n" +
				"<p>Block quotation started<br/>" +
				"Block quotation continued<br/>" +
				"The last line of the block quotation</p>\n" +
				"<cite>The Author</cite>\n" +
				"</blockquote>\n" +
				"<aside>Pull quote<cite>The Author</cite></aside>",
			wantMarkdown: "<blockquote>\n" +
				"<p>Block quotation started<br/>" +
				"Block quotation continued<br/>" +
				"The last line of the block quotation</p>\n" +
				"<cite>The Author</cite>\n" +
				"</blockquote>\n\n<aside>Pull quote<cite>The Author</cite></aside>",
			wantText: "Block quotation started\nBlock quotation continued\nThe last line of the block quotation\nThe Author\nPull quote\nThe Author",
		}, {
			name: "files",
			// Skip <audio src="https://telegram.org/example/audio.ogg"></audio>
			// Can't send voice files in PM
			inputHTML: `<img src="https://telegram.org/example/photo.jpg"/>
				<video src="https://telegram.org/example/video.mp4"></video>
				<audio src="https://telegram.org/example/audio.mp3"></audio>
				<video src="https://telegram.org/example/animation.gif"></video>`,
			wantHTML: "<img src=\"tg://photo?id=1\"></img>\n" +
				"<video src=\"tg://video?id=2\"></video>\n" +
				"<audio src=\"tg://audio?id=3\"></audio>\n" +
				"<video src=\"tg://video?id=4\"></video>",
			wantMarkdown: "![](tg://photo?id=1)\n" +
				"![](tg://video?id=2)\n" +
				"![](tg://audio?id=3)\n" +
				"![](tg://video?id=4)",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{}},
				{Id: "2", Media: InputMediaVideo{}},
				{Id: "3", Media: InputMediaAudio{}},
				{Id: "4", Media: InputMediaAnimation{}},
			},
		}, {
			name: "captioned files",
			// skipping <figure><audio src="https://telegram.org/example/audio.ogg"></audio><figcaption>Voice note caption</figcaption></figure>
			inputHTML: `<figure><img src="https://telegram.org/example/photo.jpg" tg-spoiler/><figcaption>Photo caption<cite>Photo credit</cite></figcaption></figure>
				<figure><video src="https://telegram.org/example/video.mp4" tg-spoiler></video><figcaption>Video caption</figcaption></figure>
				<figure><audio src="https://telegram.org/example/audio.mp3"></audio><figcaption>Audio caption</figcaption></figure>
				<figure><video src="https://telegram.org/example/animation.gif" tg-spoiler></video><figcaption>Animation caption</figcaption></figure>`,
			wantHTML: "<figure>\n<img src=\"tg://photo?id=1\" tg-spoiler></img>\n" +
				"<figcaption>Photo caption<cite>Photo credit</cite></figcaption>\n</figure>\n" +
				"<figure>\n<video src=\"tg://video?id=2\" tg-spoiler></video>\n" +
				"<figcaption>Video caption</figcaption>\n</figure>\n" +
				"<figure>\n<audio src=\"tg://audio?id=3\"></audio>\n" +
				"<figcaption>Audio caption</figcaption>\n</figure>\n" +
				"<figure>\n<video src=\"tg://video?id=4\" tg-spoiler></video>\n" +
				"<figcaption>Animation caption</figcaption>\n</figure>",
			wantMarkdown: "<figure>\n<img src=\"tg://photo?id=1\" tg-spoiler></img>\n" +
				"<figcaption>Photo caption<cite>Photo credit</cite></figcaption>\n</figure>\n" +
				"<figure>\n<video src=\"tg://video?id=2\" tg-spoiler></video>\n" +
				"<figcaption>Video caption</figcaption>\n</figure>\n" +
				"<figure>\n<audio src=\"tg://audio?id=3\"></audio>\n" +
				"<figcaption>Audio caption</figcaption>\n</figure>\n" +
				"<figure>\n<video src=\"tg://video?id=4\" tg-spoiler></video>\n" +
				"<figcaption>Animation caption</figcaption>\n</figure>",
			wantText: "Photo caption\nPhoto credit\n" +
				"Video caption\nAudio caption\nAnimation caption",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{}},
				{Id: "2", Media: InputMediaVideo{}},
				{Id: "3", Media: InputMediaAudio{}},
				{Id: "4", Media: InputMediaAnimation{}},
			},
		}, {
			name: "maps",
			inputHTML: `<tg-map lat="41.9" long="12.5" zoom="14"/>
				<figure><tg-map lat="41.9" long="12.5" zoom="14"/><figcaption>Map caption</figcaption></figure>`,
			wantHTML: "<tg-map lat=\"41.900006\" long=\"12.500002\" zoom=\"14\"/>\n" +
				"<figure><tg-map lat=\"41.900006\" long=\"12.500002\" zoom=\"14\"/><figcaption>Map caption</figcaption>\n" +
				"</figure>",
			wantMarkdown: "<tg-map lat=\"41.900006\" long=\"12.500002\" zoom=\"14\"/>\n" +
				"<figure><tg-map lat=\"41.900006\" long=\"12.500002\" zoom=\"14\"/><figcaption>Map caption</figcaption>\n" +
				"</figure>",
			wantText: "Map caption",
		}, {
			name: "collage",
			inputHTML: `<tg-collage><img src="https://telegram.org/example/photo.jpg"/><video src="https://telegram.org/example/video.mp4"/></tg-collage>
				<tg-collage><video src="https://telegram.org/example/video.mp4"/><img src="https://telegram.org/example/photo.jpg"/><figcaption>Collage caption</figcaption></tg-collage>`,
			wantHTML: "<tg-collage>\n" +
				"<img src=\"tg://photo?id=1\"></img>\n" +
				"<video src=\"tg://video?id=2\"></video>\n" +
				"</tg-collage>\n" +
				"<tg-collage>\n" +
				"<video src=\"tg://video?id=3\"></video>\n" +
				"<img src=\"tg://photo?id=4\"></img>\n" +
				"<figcaption>Collage caption</figcaption>\n" +
				"</tg-collage>",
			wantMarkdown: "<tg-collage>\n" +
				"<img src=\"tg://photo?id=1\"></img>\n" +
				"<video src=\"tg://video?id=2\"></video>\n" +
				"</tg-collage>\n" +
				"<tg-collage>\n" +
				"<video src=\"tg://video?id=3\"></video>\n" +
				"<img src=\"tg://photo?id=4\"></img>\n" +
				"<figcaption>Collage caption</figcaption>\n" +
				"</tg-collage>",
			wantText: "Collage caption",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{}},
				{Id: "2", Media: InputMediaVideo{}},
				{Id: "3", Media: InputMediaVideo{}},
				{Id: "4", Media: InputMediaPhoto{}},
			},
		}, {
			name: "slideshow",
			inputHTML: `<tg-slideshow><img src="https://telegram.org/example/photo.jpg"/><video src="https://telegram.org/example/video.mp4"/></tg-slideshow>
				<tg-slideshow><video src="https://telegram.org/example/video.mp4"/><img src="https://telegram.org/example/photo.jpg"/><figcaption>Slideshow caption</figcaption></tg-slideshow>`,
			wantHTML: "<tg-slideshow>\n" +
				"<img src=\"tg://photo?id=1\"></img>\n" +
				"<video src=\"tg://video?id=2\"></video>\n" +
				"</tg-slideshow>\n" +
				"<tg-slideshow>\n" +
				"<video src=\"tg://video?id=3\"></video>\n" +
				"<img src=\"tg://photo?id=4\"></img>\n" +
				"<figcaption>Slideshow caption</figcaption>\n" +
				"</tg-slideshow>",
			wantMarkdown: "<tg-slideshow>\n" +
				"<img src=\"tg://photo?id=1\"></img>\n" +
				"<video src=\"tg://video?id=2\"></video>\n" +
				"</tg-slideshow>\n" +
				"<tg-slideshow>\n" +
				"<video src=\"tg://video?id=3\"></video>\n" +
				"<img src=\"tg://photo?id=4\"></img>\n" +
				"<figcaption>Slideshow caption</figcaption>\n" +
				"</tg-slideshow>",
			wantText: "Slideshow caption",
			wantMedia: []InputRichMessageMedia{
				{Id: "1", Media: InputMediaPhoto{}},
				{Id: "2", Media: InputMediaVideo{}},
				{Id: "3", Media: InputMediaVideo{}},
				{Id: "4", Media: InputMediaPhoto{}},
			},
		}, {
			name: "table",
			inputHTML: `<table><tr><th>Header 1</th><th>Header 2</th></tr>
				<tr><td>Value 1</td><td>Value 2</td></tr></table>
				<table bordered striped><caption>Table caption</caption>
				<tr><td colspan="2" rowspan="2" align="left">Value</td><td align="center">Value2</td><td align="right">Value3</td></tr>
				<tr><td valign="top">Value4</td><td valign="middle">Value5</td><td valign="bottom">Value6</td></tr>
				<tr><td>Value7</td></tr></table>`,
			wantHTML: "<table>\n<tr>\n" +
				"<th>Header 1</th>\n" +
				"<th>Header 2</th>\n" +
				"</tr>\n<tr>\n" +
				"<td>Value 1</td>\n" +
				"<td>Value 2</td>\n" +
				"</tr>\n</table>\n" +
				"<table bordered striped>\n" +
				"<caption>Table caption</caption>\n" +
				"<tr>\n" +
				"<td colspan=\"2\" rowspan=\"2\">Value</td>\n" +
				"<td align=\"center\">Value2</td>\n" +
				"<td align=\"right\">Value3</td>\n" +
				"<td valign=\"top\"></td>\n" +
				"</tr>\n<tr>\n" +
				"<td valign=\"top\">Value4</td>\n" +
				"<td>Value5</td>\n" +
				"<td valign=\"bottom\">Value6</td>\n" +
				"</tr>\n<tr>\n" +
				"<td>Value7</td>\n" +
				"<td valign=\"top\" colspan=\"4\"></td>\n" +
				"</tr>\n</table>",
			wantMarkdown: "| Header 1 | Header 2 |\n" +
				"|:---:|:---:|\n" +
				"| Value 1 | Value 2 |\n\n" +
				"<table bordered striped>\n<caption>Table caption</caption>\n<tr>\n" +
				"<td colspan=\"2\" rowspan=\"2\">Value</td>\n<td align=\"center\">Value2</td>\n<td align=\"right\">Value3</td>\n<td valign=\"top\"></td>\n</tr>\n" +
				"<tr>\n<td valign=\"top\">Value4</td>\n<td>Value5</td>\n<td valign=\"bottom\">Value6</td>\n</tr>\n" +
				"<tr>\n<td>Value7</td>\n<td valign=\"top\" colspan=\"4\"></td>\n</tr>\n" +
				"</table>",
			wantText: "Header 1\tHeader 2\n" +
				"Value 1\tValue 2\n\n" +
				"Value\tValue2\tValue3\t\n" +
				"Value4\tValue5\tValue6\n" +
				"Value7",
		}, {
			name: "details",
			inputHTML: `<details><summary>Title</summary>Content</details>
				<details open><summary>Title</summary>Content</details>
				<tg-math-block>E = mc^2</tg-math-block>`,
			wantHTML: "<details>\n" +
				"<summary>Title</summary>\n" +
				"<p>Content</p>\n</details>\n" +
				"<details open>\n<summary>Title</summary>\n" +
				"<p>Content</p>\n</details>\n" +
				"<tg-math-block>E = mc^2</tg-math-block>",
			wantMarkdown: "<details>\n" +
				"<summary>Title</summary>\n" +
				"<p>Content</p>\n</details>\n" +
				"<details open>\n<summary>Title</summary>\n" +
				"<p>Content</p>\n</details>\n" +
				"$$\nE = mc^2\n$$",
			wantText: "Title\nContent\nTitle\nContent\nE = mc^2",
		}, {
			name:         "simple newlines",
			inputHTML:    "<p>test<br/>me</p>",
			wantHTML:     "<p>test<br/>me</p>",
			wantMarkdown: "test<br/>me", // Double the newlines for markdown to keep them.
			wantText:     "test\nme",
			wantMedia:    nil,
		}, {
			name:         "double newlines",
			inputHTML:    "<p>test<br/><br/>me</p>",
			wantHTML:     "<p>test<br/><br/>me</p>",
			wantMarkdown: "test<br/><br/>me", // Double the newlines for markdown to keep them.
			wantText:     "test\n\nme",
			wantMedia:    nil,
		}, {
			name:         "divider",
			inputHTML:    "<h1>title</h1><p>paragraph</p><hr/>",
			wantHTML:     "<h1>title</h1>\n<p>paragraph</p>\n<hr/>",
			wantMarkdown: "# title\nparagraph\n\n---", // Need 2 newlines to show as a divider
			wantText:     "title\nparagraph",
			wantMedia:    nil,
		}, {
			name:         "footer",
			inputHTML:    "<h1>title</h1><p>paragraph</p><footer>Footer</footer>",
			wantHTML:     "<h1>title</h1>\n<p>paragraph</p>\n<footer>Footer</footer>",
			wantMarkdown: "# title\nparagraph\n<footer>Footer</footer>",
			wantText:     "title\nparagraph\nFooter",
			wantMedia:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(time.Second)
			m, err := b.SendRichMessage(chatId, InputRichMessage{Html: tt.inputHTML}, nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			html, htmlMedia := m.RichMessage.HTML()
			if got := strings.TrimSpace(html); got != tt.wantHTML {
				t.Errorf("HTML() mismatch:\n got: %q\nwant: %q", got, tt.wantHTML)
			}
			equalMedias(t, htmlMedia, tt.wantMedia)

			markdown, mdMedia := m.RichMessage.Markdown()
			if got := strings.TrimSpace(markdown); got != tt.wantMarkdown {
				t.Errorf("Markdown() mismatch:\n got: %q\nwant: %q", got, tt.wantMarkdown)
			}
			equalMedias(t, mdMedia, tt.wantMedia)

			// no ids to change
			if got := strings.TrimSpace(m.RichMessage.PlainText()); got != tt.wantText {
				t.Errorf("Text() mismatch:\n got: %q\nwant: %q", got, tt.wantText)
			}
		})
	}
}
