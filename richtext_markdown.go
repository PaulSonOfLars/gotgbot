package gotgbot

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// RichTextMarkdown renders r as a Markdown string.
// Uses CommonMark conventions; for Telegram MarkdownV2 use RichTextMarkdownV2.
func RichTextMarkdown(r RichText) string {
	var sb strings.Builder
	renderTextMarkdown(r, &sb)
	return sb.String()
}

func renderTextMarkdown(r RichText, sb *strings.Builder) {
	if r == nil {
		return
	}
	switch v := r.(type) {
	case RichTextString:
		sb.WriteString(mdEscape(string(v)))
	case RichTextArray:
		for _, child := range v {
			renderTextMarkdown(child, sb)
		}
	case RichTextBold:
		sb.WriteString("**")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("**")
	case RichTextItalic:
		sb.WriteString("_")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("_")
	case RichTextUnderline:
		// No standard Markdown underline; use HTML fallback.
		sb.WriteString("<u>")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("</u>")
	case RichTextStrikethrough:
		sb.WriteString("~~")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("~~")
	case RichTextSpoiler:
		sb.WriteString(`||`)
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("||")
	case RichTextMarked:
		sb.WriteString("==")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("==")
	case RichTextSubscript:
		sb.WriteString("<sub>")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("</sub>")
	case RichTextSuperscript:
		sb.WriteString("<sup>")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("</sup>")
	case RichTextCode:
		sb.WriteString("`")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("`")
	case RichTextUrl:
		sb.WriteString("[")
		renderTextMarkdown(v.Text, sb)
		sb.WriteString("](")
		sb.WriteString(v.Url)
		sb.WriteString(")")
	case RichTextEmailAddress:
		text := RichTextContent(v.Text)
		fmt.Fprintf(sb, "[%s](mailto:%s)", mdEscape(text), v.EmailAddress)
	case RichTextPhoneNumber:
		text := RichTextContent(v.Text)
		fmt.Fprintf(sb, "[%s](tel:%s)", mdEscape(text), v.PhoneNumber)
	case RichTextTextMention:
		sb.WriteString("[")
		renderTextMarkdown(v.Text, sb)
		fmt.Fprintf(sb, "](tg://user?id=%d)", v.User.Id)
	case RichTextMention:
		renderTextMarkdown(v.Text, sb)
	case RichTextAnchor:
		fmt.Fprintf(sb, "<a name=\"%s\"></a>\n", html.EscapeString(v.Name))
	case RichTextAnchorLink:
		sb.WriteString("[")
		renderTextMarkdown(v.Text, sb)
		fmt.Fprintf(sb, "](#%s)", v.AnchorName)
	case RichTextReference:
		fmt.Fprintf(sb, "[^%s]: ", v.Name)
		renderTextMarkdown(v.Text, sb)
	case RichTextReferenceLink:
		renderTextMarkdown(v.Text, sb)
		fmt.Fprintf(sb, "[^%s]", v.ReferenceName)
	case RichTextHashtag:
		renderTextMarkdown(v.Text, sb)
	case RichTextCashtag:
		renderTextMarkdown(v.Text, sb)
	case RichTextBotCommand:
		renderTextMarkdown(v.Text, sb)
	case RichTextBankCardNumber:
		renderTextMarkdown(v.Text, sb)
	case RichTextDateTime:
		sb.WriteString("![")
		renderTextMarkdown(v.Text, sb)
		fmt.Fprintf(sb, "](tg://time?unix=%dformat=%s)", v.UnixTime, v.DateTimeFormat)
	case RichTextCustomEmoji:
		fmt.Fprintf(sb, "![%s](tg://emoji?id=%s)", v.AlternativeText, v.CustomEmojiId)
	case RichTextMathematicalExpression:
		sb.WriteString("$$")
		sb.WriteString(v.Expression)
		sb.WriteString("$$")
	}
}

// RichBlockMarkdown renders a RichBlock as a Markdown string.
func RichBlockMarkdown(b RichBlock) string {
	var sb strings.Builder
	renderBlockMarkdown(b, &sb, 0)
	return sb.String()
}

func renderBlockMarkdown(b RichBlock, sb *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)
	switch v := b.(type) {
	case RichBlockParagraph:
		sb.WriteString(RichTextMarkdown(v.Text))
		sb.WriteString("\n")
	case RichBlockSectionHeading:
		sb.WriteString(strings.Repeat("#", int(v.Size)) + " ")
		sb.WriteString(RichTextMarkdown(v.Text))
		sb.WriteString("\n")
	case RichBlockPreformatted:
		sb.WriteString("```" + v.Language + "\n")
		sb.WriteString(RichTextContent(v.Text)) // raw text, no escaping inside code blocks
		sb.WriteString("\n```\n")
	case RichBlockFooter:
		sb.WriteString("---\n")
		sb.WriteString(RichTextMarkdown(v.Text))
		sb.WriteString("\n")
	case RichBlockBlockQuotation:
		for _, child := range v.Blocks {
			// Prefix each line with "> "
			inner := strings.Builder{}
			renderBlockMarkdown(child, &inner, 0)
			for _, line := range strings.Split(strings.TrimRight(inner.String(), "\n"), "\n") {
				sb.WriteString("> " + line + "\n")
			}
		}
		if v.Credit != nil {
			sb.WriteString("<aside>")
			sb.WriteString(RichTextMarkdown(v.Credit))
			sb.WriteString("</aside>")
		}
		sb.WriteString("\n")
	case RichBlockPullQuotation:
		sb.WriteString("<aside>")
		sb.WriteString(RichTextMarkdown(v.Text))
		if v.Credit != nil {
			sb.WriteString("<cite>")
			sb.WriteString(RichTextMarkdown(v.Credit))
			sb.WriteString("</cite>")
		}
		sb.WriteString("</aside>")
		sb.WriteString("\n")
	case RichBlockDetails:
		sb.WriteString("<details")
		if v.IsOpen {
			sb.WriteString(" open")
		}
		sb.WriteString(">\n<summary>")
		sb.WriteString(RichTextMarkdown(v.Summary))
		sb.WriteString("</summary>\n")
		for _, child := range v.Blocks {
			renderBlockMarkdown(child, sb, depth)
		}
		sb.WriteString("</details>\n")
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

				fmt.Fprintf(sb, "%s%s. ", indent, val)

			} else {
				sb.WriteString(indent + "- ")
				if item.HasCheckbox {
					if item.IsChecked {
						sb.WriteString(indent + "[x] ")
					} else {
						sb.WriteString(indent + "[ ] ")
					}
				}
			}
			for _, child := range item.Blocks {
				renderBlockMarkdown(child, sb, depth+1)
			}
		}
		sb.WriteString("\n")
	case RichBlockTable:
		// Emit as a GFM table.
		if len(v.Cells) == 0 {
			break
		}
		for i, row := range v.Cells {
			sb.WriteString("|")
			for _, cell := range row {
				text := ""
				if cell.Text != nil {
					text = RichTextMarkdown(cell.Text)
				}
				sb.WriteString(" " + text + " |")
			}
			sb.WriteString("\n")
			if i == 0 {
				sb.WriteString("|")
				for _, cell := range row {
					switch cell.Align {
					case "right":
						sb.WriteString("---:|")
					case "left":
						sb.WriteString(":---|")
					case "center", "":
						sb.WriteString(":---:|")
					default:
						sb.WriteString("---|")
					}
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n")
	case RichBlockCollage:
		sb.WriteString("<tg-collage>\n")
		for _, child := range v.Blocks {
			renderBlockMarkdown(child, sb, depth+1)
		}
		renderCaptionHTML(v.Caption, sb)
		sb.WriteString("</tg-collage>\n")

	case RichBlockSlideshow:
		sb.WriteString("<tg-slideshow>\n")
		for _, child := range v.Blocks {
			renderBlockMarkdown(child, sb, depth+1)
		}
		renderCaptionHTML(v.Caption, sb)
		sb.WriteString("</tg-slideshow>\n")

	case RichBlockThinking:
		sb.WriteString("<tg-thinking>")
		sb.WriteString(RichTextMarkdown(v.Text))
		sb.WriteString("\n</tg-thinking>\n")
	case RichBlockMathematicalExpression:
		sb.WriteString("$$\n")
		sb.WriteString(v.Expression)
		sb.WriteString("\n$$\n")
	case RichBlockDivider:
		sb.WriteString("---\n")
	case RichBlockAnchor:
		fmt.Fprintf(sb, "<a name=\"%s\"></a>\n", v.Name)
	case RichBlockPhoto:
		if v.Caption != nil {
			renderBlockHTML(v, sb)
		} else {
			sb.WriteString(fileLink(v.Photo[len(v.Photo)-1].FileId))
		}
	case RichBlockAnimation:
		if v.Caption != nil {
			renderBlockHTML(v, sb)
		} else {
			sb.WriteString(fileLink(v.Animation.FileId))
		}
	case RichBlockVideo:
		if v.Caption != nil {
			renderBlockHTML(v, sb)
		} else {
			sb.WriteString(fileLink(v.Video.FileId))
		}
	case RichBlockAudio:
		if v.Caption != nil {
			renderBlockHTML(v, sb)
		} else {
			sb.WriteString(fileLink(v.Audio.FileId))
		}
	case RichBlockVoiceNote:
		if v.Caption != nil {
			renderBlockHTML(v, sb)
		} else {
			sb.WriteString(fileLink(v.VoiceNote.FileId))
		}
	case RichBlockMap:
		// Maps aren't supported in markdown; require inline HTML.
		renderBlockHTML(v, sb)
	}
}

func fileLink(fileId string) string {
	return fmt.Sprintf("![](fileId://%s)\n", fileId)
}
