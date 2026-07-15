package gotgbot

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// RichTextMarkdown renders r as a Markdown string.
// Uses CommonMark conventions; for Telegram MarkdownV2 use RichTextMarkdownV2.
func RichTextMarkdown(rt RichText) string {
	var r renderCtx
	r.renderTextMarkdown(rt)
	return r.sb.String()
}

func (r *renderCtx) renderTextMarkdown(rt RichText) {
	if r == nil {
		return
	}
	switch v := rt.(type) {
	case RichTextString:
		r.sb.WriteString(mdEscape(string(v)))
	case RichTextArray:
		for _, child := range v {
			r.renderTextMarkdown(child)
		}
	case RichTextBold:
		r.sb.WriteString("**")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("**")
	case RichTextItalic:
		r.sb.WriteString("_")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("_")
	case RichTextUnderline:
		// No standard Markdown underline; use HTML fallback.
		r.sb.WriteString("<u>")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("</u>")
	case RichTextStrikethrough:
		r.sb.WriteString("~~")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("~~")
	case RichTextSpoiler:
		r.sb.WriteString(`||`)
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("||")
	case RichTextMarked:
		r.sb.WriteString("==")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("==")
	case RichTextSubscript:
		r.sb.WriteString("<sub>")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("</sub>")
	case RichTextSuperscript:
		r.sb.WriteString("<sup>")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("</sup>")
	case RichTextCode:
		r.sb.WriteString("`")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("`")
	case RichTextUrl:
		r.sb.WriteString("[")
		r.renderTextMarkdown(v.Text)
		r.sb.WriteString("](")
		r.sb.WriteString(v.Url)
		r.sb.WriteString(")")
	case RichTextEmailAddress:
		text := RichTextContent(v.Text)
		fmt.Fprintf(&r.sb, "[%s](mailto:%s)", mdEscape(text), v.EmailAddress)
	case RichTextPhoneNumber:
		text := RichTextContent(v.Text)
		fmt.Fprintf(&r.sb, "[%s](tel:%s)", mdEscape(text), v.PhoneNumber)
	case RichTextTextMention:
		r.sb.WriteString("[")
		r.renderTextMarkdown(v.Text)
		fmt.Fprintf(&r.sb, "](tg://user?id=%d)", v.User.Id)
	case RichTextMention:
		r.renderTextMarkdown(v.Text)
	case RichTextAnchor:
		fmt.Fprintf(&r.sb, "<a name=\"%s\"></a>\n", html.EscapeString(v.Name))
	case RichTextAnchorLink:
		r.sb.WriteString("[")
		r.renderTextMarkdown(v.Text)
		fmt.Fprintf(&r.sb, "](#%s)", v.AnchorName)
	case RichTextReference:
		fmt.Fprintf(&r.sb, "[^%s]: ", v.Name)
		r.renderTextMarkdown(v.Text)
	case RichTextReferenceLink:
		r.renderTextMarkdown(v.Text)
		fmt.Fprintf(&r.sb, "[^%s]", v.ReferenceName)
	case RichTextHashtag:
		r.renderTextMarkdown(v.Text)
	case RichTextCashtag:
		r.renderTextMarkdown(v.Text)
	case RichTextBotCommand:
		r.renderTextMarkdown(v.Text)
	case RichTextBankCardNumber:
		r.renderTextMarkdown(v.Text)
	case RichTextDateTime:
		r.sb.WriteString("![")
		r.renderTextMarkdown(v.Text)
		fmt.Fprintf(&r.sb, "](tg://time?unix=%dformat=%s)", v.UnixTime, v.DateTimeFormat)
	case RichTextCustomEmoji:
		fmt.Fprintf(&r.sb, "![%s](tg://emoji?id=%s)", v.AlternativeText, v.CustomEmojiId)
	case RichTextMathematicalExpression:
		r.sb.WriteString("$$")
		r.sb.WriteString(v.Expression)
		r.sb.WriteString("$$")
	}
}

// RichBlockMarkdown renders a RichBlock as a Markdown string.
func RichBlockMarkdown(b RichBlock) (string, []InputRichMessageMedia) {
	var r renderCtx
	r.renderBlockMarkdown(b, 0)
	return r.sb.String(), r.media
}

func (r *renderCtx) renderBlockMarkdown(b RichBlock, depth int) {
	indent := strings.Repeat("  ", depth)
	switch v := b.(type) {
	case RichBlockParagraph:
		r.sb.WriteString(RichTextMarkdown(v.Text))
		r.sb.WriteString("\n")
	case RichBlockSectionHeading:
		r.sb.WriteString(strings.Repeat("#", int(v.Size)) + " ")
		r.sb.WriteString(RichTextMarkdown(v.Text))
		r.sb.WriteString("\n")
	case RichBlockPreformatted:
		r.sb.WriteString("```" + v.Language + "\n")
		r.sb.WriteString(RichTextContent(v.Text)) // raw text, no escaping inside code blocks
		r.sb.WriteString("\n```\n")
	case RichBlockFooter:
		r.sb.WriteString("---\n")
		r.sb.WriteString(RichTextMarkdown(v.Text))
		r.sb.WriteString("\n")
	case RichBlockBlockQuotation:
		// Prefix each line with "> "
		for _, child := range v.Blocks {
			// create an internal render, only for the text
			inner := renderCtx{
				sb:      strings.Builder{},
				media:   r.media,
				counter: r.counter,
			}
			inner.renderBlockMarkdown(child, 0)

			for _, line := range strings.Split(strings.TrimRight(inner.sb.String(), "\n"), "\n") {
				r.sb.WriteString("> " + line + "\n")
			}

			r.media = inner.media
			r.counter = inner.counter
		}
		if v.Credit != nil {
			r.sb.WriteString("<aside>")
			r.sb.WriteString(RichTextMarkdown(v.Credit))
			r.sb.WriteString("</aside>")
		}
		r.sb.WriteString("\n")
	case RichBlockPullQuotation:
		r.sb.WriteString("<aside>")
		r.sb.WriteString(RichTextMarkdown(v.Text))
		if v.Credit != nil {
			r.sb.WriteString("<cite>")
			r.sb.WriteString(RichTextMarkdown(v.Credit))
			r.sb.WriteString("</cite>")
		}
		r.sb.WriteString("</aside>")
		r.sb.WriteString("\n")
	case RichBlockDetails:
		r.sb.WriteString("<details")
		if v.IsOpen {
			r.sb.WriteString(" open")
		}
		r.sb.WriteString(">\n<summary>")
		r.sb.WriteString(RichTextMarkdown(v.Summary))
		r.sb.WriteString("</summary>\n")
		for _, child := range v.Blocks {
			r.renderBlockMarkdown(child, depth)
		}
		r.sb.WriteString("</details>\n")
	case RichBlockList:
		ordered := len(v.Items) > 0 && v.Items[0].Value != 0
		for i, item := range v.Items {
			if ordered {
				valInt := item.Value
				if valInt == 0 {
					valInt = int64(i) + 1
				}

				val := strconv.FormatInt(valInt, 10)
				switch item.Type {
				case "A":
					val = toLetters(valInt)
				case "a":
					val = strings.ToLower(toLetters(valInt))
				case "I":
					val = toRoman(valInt)
				case "i":
					val = strings.ToLower(toRoman(valInt))
				}

				fmt.Fprintf(&r.sb, "%s%s. ", indent, val)

			} else {
				r.sb.WriteString(indent + "- ")
				if item.HasCheckbox {
					if item.IsChecked {
						r.sb.WriteString(indent + "[x] ")
					} else {
						r.sb.WriteString(indent + "[ ] ")
					}
				}
			}
			for _, child := range item.Blocks {
				r.renderBlockMarkdown(child, depth+1)
			}
		}
		r.sb.WriteString("\n")
	case RichBlockTable:
		// Emit as a GFM table.
		if len(v.Cells) == 0 {
			break
		}

		// simple attempt at detecting lossy table behaviour
		if row := v.Cells[0]; v.Caption != nil || len(row) == 0 || !row[0].IsHeader {
			r.renderBlockHTML(b)
			break
		}

		for i, row := range v.Cells {
			r.sb.WriteString("|")
			for _, cell := range row {
				text := ""
				if cell.Text != nil {
					text = RichTextMarkdown(cell.Text)
				}
				r.sb.WriteString(" " + text + " |")
			}
			r.sb.WriteString("\n")

			// NOTE: This builds the direction based on the first header row.
			// HTML supports per-row ordering, markdown does not; this is lossy behaviour.
			if i == 0 {
				r.sb.WriteString("|")
				for _, cell := range row {
					switch cell.Align {
					case "right":
						r.sb.WriteString("---:|")
					case "left":
						r.sb.WriteString(":---|")
					case "center":
						r.sb.WriteString(":---:|")
					default:
						r.sb.WriteString("---|")
					}
				}
				r.sb.WriteString("\n")
			}
		}
		r.sb.WriteString("\n")
	case RichBlockCollage:
		r.sb.WriteString("<tg-collage>\n")
		for _, child := range v.Blocks {
			r.renderBlockMarkdown(child, depth+1)
		}
		r.renderCaptionHTML(v.Caption)
		r.sb.WriteString("</tg-collage>\n")

	case RichBlockSlideshow:
		r.sb.WriteString("<tg-slideshow>\n")
		for _, child := range v.Blocks {
			r.renderBlockMarkdown(child, depth+1)
		}
		r.renderCaptionHTML(v.Caption)
		r.sb.WriteString("</tg-slideshow>\n")

	case RichBlockThinking:
		r.sb.WriteString("<tg-thinking>")
		r.sb.WriteString(RichTextMarkdown(v.Text))
		r.sb.WriteString("\n</tg-thinking>\n")
	case RichBlockMathematicalExpression:
		r.sb.WriteString("$$\n")
		r.sb.WriteString(v.Expression)
		r.sb.WriteString("\n$$\n")
	case RichBlockDivider:
		r.sb.WriteString("---\n")
	case RichBlockAnchor:
		fmt.Fprintf(&r.sb, "<a name=\"%s\"></a>\n", v.Name)
	case RichBlockPhoto:
		if v.Caption != nil {
			r.renderBlockHTML(v)
		} else {
			fmt.Fprintf(&r.sb, "![](%s)\n", r.addFile("photo", InputMediaPhoto{Media: InputFileByID(v.Photo[len(v.Photo)-1].FileId)}))
		}
	case RichBlockAnimation:
		if v.Caption != nil {
			r.renderBlockHTML(v)
		} else {
			fmt.Fprintf(&r.sb, "![](%s)\n", r.addFile("video", InputMediaAnimation{Media: InputFileByID(v.Animation.FileId)}))
		}
	case RichBlockVideo:
		if v.Caption != nil {
			r.renderBlockHTML(v)
		} else {
			fmt.Fprintf(&r.sb, "![](%s)\n", r.addFile("video", InputMediaVideo{Media: InputFileByID(v.Video.FileId)}))
		}
	case RichBlockAudio:
		if v.Caption != nil {
			r.renderBlockHTML(v)
		} else {
			fmt.Fprintf(&r.sb, "![](%s)\n", r.addFile("audio", InputMediaAudio{Media: InputFileByID(v.Audio.FileId)}))
		}
	case RichBlockVoiceNote:
		if v.Caption != nil {
			r.renderBlockHTML(v)
		} else {
			fmt.Fprintf(&r.sb, "![](%s)\n", r.addFile("audio", InputMediaVoiceNote{Media: InputFileByID(v.VoiceNote.FileId)}))
		}
	case RichBlockMap:
		// Maps aren't supported in markdown; require inline HTML.
		r.renderBlockHTML(v)
	}
}
