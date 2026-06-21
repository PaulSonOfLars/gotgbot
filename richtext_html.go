package gotgbot

import (
	"fmt"
	"html"
	"strings"
)

// RichTextHTML renders r as an HTML string, using standard HTML tags for each
// RichText node type. The output is suitable for embedding in an HTML document.
func RichTextHTML(r RichText) string {
	var sb strings.Builder
	renderTextHTML(r, &sb)
	return sb.String()
}

func renderTextHTML(r RichText, sb *strings.Builder) {
	if r == nil {
		return
	}
	switch v := r.(type) {
	case RichTextString:
		sb.WriteString(html.EscapeString(string(v)))
	case RichTextArray:
		for _, child := range v {
			renderTextHTML(child, sb)
		}

	// Formatting wrappers - recurse into .Text
	case RichTextBold:
		sb.WriteString("<b>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</b>")
	case RichTextItalic:
		sb.WriteString("<i>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</i>")
	case RichTextUnderline:
		sb.WriteString("<u>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</u>")
	case RichTextStrikethrough:
		sb.WriteString("<s>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</s>")
	case RichTextSpoiler:
		sb.WriteString(`<tg-spoiler>`)
		renderTextHTML(v.Text, sb)
		sb.WriteString("</tg-spoiler>")
	case RichTextMarked:
		sb.WriteString("<mark>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</mark>")
	case RichTextSubscript:
		sb.WriteString("<sub>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</sub>")
	case RichTextSuperscript:
		sb.WriteString("<sup>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</sup>")
	case RichTextCode:
		sb.WriteString("<code>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</code>")

	// Links and references
	case RichTextUrl:
		fmt.Fprintf(sb, `<a href="%s">`, html.EscapeString(v.Url))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")
	case RichTextEmailAddress:
		fmt.Fprintf(sb, `<a href="mailto:%s">`, html.EscapeString(v.EmailAddress))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")
	case RichTextPhoneNumber:
		fmt.Fprintf(sb, `<a href="tel:%s">`, html.EscapeString(v.PhoneNumber))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")
	case RichTextTextMention:
		fmt.Fprintf(sb, `<a href="tg://user?id=%d">`, v.User.Id)
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")
	case RichTextMention:
		renderTextHTML(v.Text, sb)
	case RichTextAnchor:
		fmt.Fprintf(sb, `<a name="%s"></a>`, html.EscapeString(v.Name))
	case RichTextAnchorLink:
		fmt.Fprintf(sb, `<a href="#%s">`, html.EscapeString(v.AnchorName))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")
	case RichTextReference:
		fmt.Fprintf(sb, `<tg-reference name="%s">`, html.EscapeString(v.Name))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</tg-reference>")
	case RichTextReferenceLink:
		fmt.Fprintf(sb, `<a href="#%s">`, html.EscapeString(v.ReferenceName))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</a>")

	// Entities with no wrapping text child
	case RichTextHashtag:
		renderTextHTML(v.Text, sb)
	case RichTextCashtag:
		renderTextHTML(v.Text, sb)
	case RichTextBotCommand:
		renderTextHTML(v.Text, sb)
	case RichTextBankCardNumber:
		renderTextHTML(v.Text, sb)
	case RichTextDateTime:
		fmt.Fprintf(sb, `<tg-time unix="%d" format="%s">`, v.UnixTime, html.EscapeString(v.DateTimeFormat))
		renderTextHTML(v.Text, sb)
		sb.WriteString("</tg-time>")
	case RichTextCustomEmoji:
		fmt.Fprintf(sb, `<tg-emoji emoji-id="%s">%s</tg-emoji>`, html.EscapeString(v.CustomEmojiId), html.EscapeString(v.AlternativeText))
	case RichTextMathematicalExpression:
		fmt.Fprintf(sb, `<tg-math>%s</tg-math>`, html.EscapeString(v.Expression))
	}
}

// RichBlockHTML renders a RichBlock and its descendants as an HTML fragment.
func RichBlockHTML(b RichBlock) string {
	var sb strings.Builder
	renderBlockHTML(b, &sb)
	return sb.String()
}

func renderBlockHTML(b RichBlock, sb *strings.Builder) {
	switch v := b.(type) {
	case RichBlockParagraph:
		sb.WriteString("<p>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</p>\n")
	case RichBlockSectionHeading:
		tag := fmt.Sprintf("h%d", v.Size)
		sb.WriteString("<" + tag + ">")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</" + tag + ">\n")
	case RichBlockPreformatted:
		if v.Language != "" {
			fmt.Fprintf(sb, `<pre><code class="language-%s">`, html.EscapeString(v.Language))
		} else {
			sb.WriteString("<pre><code>")
		}
		renderTextHTML(v.Text, sb)
		sb.WriteString("</code></pre>\n")
	case RichBlockFooter:
		sb.WriteString("<footer>")
		renderTextHTML(v.Text, sb)
		sb.WriteString("</footer>\n")
	case RichBlockBlockQuotation:
		sb.WriteString("<blockquote>\n")
		for _, child := range v.Blocks {
			renderBlockHTML(child, sb)
		}
		if v.Credit != nil {
			sb.WriteString("<cite>")
			renderTextHTML(v.Credit, sb)
			sb.WriteString("</cite>\n")
		}
		sb.WriteString("</blockquote>\n")
	case RichBlockPullQuotation:
		sb.WriteString(`<aside>`)
		renderTextHTML(v.Text, sb)
		if v.Credit != nil {
			sb.WriteString("<cite>")
			renderTextHTML(v.Credit, sb)
			sb.WriteString("</cite>")
		}
		sb.WriteString("</aside>\n")
	case RichBlockDetails:
		sb.WriteString("<details")
		if v.IsOpen {
			sb.WriteString(" open")
		}
		sb.WriteString(">\n<summary>")
		renderTextHTML(v.Summary, sb)
		sb.WriteString("</summary>\n")
		for _, child := range v.Blocks {
			renderBlockHTML(child, sb)
		}
		sb.WriteString("</details>\n")
	case RichBlockList:
		if len(v.Items) == 0 {
			return
		}

		// Detect ordered vs unordered from the first item.
		first := v.Items[0]
		tag := "ul"
		listAttrs := ""
		if first.Value > 0 {
			tag = "ol"
		}
		if first.Value > 1 {
			listAttrs += fmt.Sprintf(" start=\"%d\"", first.Value)
		}
		if first.Type != "" && first.Type != "1" {
			listAttrs += fmt.Sprintf(" type=\"%s\"", first.Type)
		}
		if isReversed(v.Items) {
			listAttrs += " reversed"
		}

		fmt.Fprintf(sb, "<%s%s>\n", tag, listAttrs)
		for _, item := range v.Items {
			itemAttrs := ""
			if tag == "ol" && item.Value > 1 {
				itemAttrs += fmt.Sprintf(" value=\"%d\"", item.Value)
			}
			if tag == "ol" && item.Type != "" && item.Type != "1" {
				itemAttrs += fmt.Sprintf(" type=\"%s\"", item.Type)
			}

			fmt.Fprintf(sb, "<li%s>", itemAttrs)
			if item.HasCheckbox {
				if item.IsChecked {
					sb.WriteString(`<input type="checkbox" checked> `)
				} else {
					sb.WriteString(`<input type="checkbox"> `)
				}
			}
			if len(item.Blocks) > 0 {
				sb.WriteString("\n")
				for _, child := range item.Blocks {
					renderBlockHTML(child, sb)
				}
			}
			sb.WriteString("</li>\n")
		}
		fmt.Fprintf(sb, "</%s>\n", tag)
	case RichBlockTable:
		attrs := ""
		if v.IsBordered {
			attrs += " bordered"
		}
		if v.IsStriped {
			attrs += " striped"
		}
		fmt.Fprintf(sb, "<table%s>\n", attrs)
		for _, row := range v.Cells {
			sb.WriteString("<tr>\n")
			for _, cell := range row {
				tag := "td"
				if cell.IsHeader {
					tag = "th"
				}

				attrs := ""
				if cell.Align != "" && cell.Align != "center" {
					attrs += fmt.Sprintf(" align=\"%s\"", cell.Align)
				}
				if cell.Valign != "" && cell.Valign != "middle" {
					attrs += fmt.Sprintf(" valign=\"%s\"", cell.Valign)
				}
				if cell.Colspan > 1 {
					attrs += fmt.Sprintf(" colspan=\"%d\"", cell.Colspan)
				}
				if cell.Rowspan > 1 {
					attrs += fmt.Sprintf(" rowspan=\"%d\"", cell.Rowspan)
				}
				fmt.Fprintf(sb, "<%s%s>", tag, attrs)
				if cell.Text != nil {
					renderTextHTML(cell.Text, sb)
				}
				fmt.Fprintf(sb, "</%s>\n", tag)
			}
			sb.WriteString("</tr>\n")
		}
		sb.WriteString("</table>\n")
	case RichBlockCollage:
		sb.WriteString("<tg-collage>\n")
		for _, child := range v.Blocks {
			renderBlockHTML(child, sb)
		}
		renderCaptionHTML(v.Caption, sb)
		sb.WriteString("</tg-collage>\n")
	case RichBlockSlideshow:
		sb.WriteString("<tg-slideshow>\n")
		for _, child := range v.Blocks {
			renderBlockHTML(child, sb)
		}
		renderCaptionHTML(v.Caption, sb)
		sb.WriteString("</tg-slideshow>\n")
	case RichBlockThinking:
		sb.WriteString(`<tg-thinking>`)
		renderTextHTML(v.Text, sb)
		sb.WriteString("</tg-thinking>\n")
	case RichBlockMathematicalExpression:
		fmt.Fprintf(sb, "<tg-math-block>%s</tg-math-block>\n", html.EscapeString(v.Expression))
	case RichBlockDivider:
		sb.WriteString("<hr/>\n")
	case RichBlockAnchor:
		fmt.Fprintf(sb, `<a name="%s"></a>`+"\n", html.EscapeString(v.Name))
	case RichBlockPhoto:
		attrs := ""
		if v.HasSpoiler {
			attrs += " tg-spoiler"
		}
		tgAnimation := fmt.Sprintf("<img src=\"fileId://%s\"%s></img>\n", v.Photo[len(v.Photo)-1].FileId, attrs)
		if v.Caption != nil {
			sb.WriteString("<figure>\n")
			sb.WriteString(tgAnimation)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>\n")
		} else {
			sb.WriteString(tgAnimation)
		}
	case RichBlockAnimation:
		attrs := ""
		if v.HasSpoiler {
			attrs += " tg-spoiler"
		}
		tgAnimation := fmt.Sprintf("<video src=\"fileId://%s\"%s></video>\n", v.Animation.FileId, attrs)
		if v.Caption != nil {
			sb.WriteString("<figure>\n")
			sb.WriteString(tgAnimation)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>\n")
		} else {
			sb.WriteString(tgAnimation)
		}
	case RichBlockVideo:
		attrs := ""
		if v.HasSpoiler {
			attrs += " tg-spoiler"
		}
		tgVideo := fmt.Sprintf("<video src=\"fileId://%s\"%s></video>\n", v.Video.FileId, attrs)
		if v.Caption != nil {
			sb.WriteString("<figure>\n")
			sb.WriteString(tgVideo)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>\n")
		} else {
			sb.WriteString(tgVideo)
		}
	case RichBlockAudio:
		tgAudio := fmt.Sprintf("<audio src=\"fileId://%s\"></audio>\n", v.Audio.FileId)
		if v.Caption != nil {
			sb.WriteString("<figure>\n")
			sb.WriteString(tgAudio)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>\n")
		} else {
			sb.WriteString(tgAudio)
		}
	case RichBlockVoiceNote:
		tgAudio := fmt.Sprintf("<audio src=\"fileId://%s\"></audio>\n", v.VoiceNote.FileId)
		if v.Caption != nil {
			sb.WriteString("<figure>\n")
			sb.WriteString(tgAudio)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>\n")
		} else {
			sb.WriteString(tgAudio)
		}
	case RichBlockMap:
		tgMap := fmt.Sprintf(
			`<tg-map lat="%g" long="%g" zoom="%d"/>`,
			v.Location.Latitude, v.Location.Longitude, v.Zoom)
		if v.Caption != nil {
			sb.WriteString("<figure>")
			sb.WriteString(tgMap)
			renderCaptionHTML(v.Caption, sb)
			sb.WriteString("</figure>")
		} else {
			sb.WriteString(tgMap)
		}
		sb.WriteString("\n")
	}
}

func isReversed(items RichBlockListItemArray) bool {
	if len(items) == 0 {
		return false
	}
	if len(items) == 1 {
		return items[0].Value > 1
	}
	return items[0].Value > items[1].Value
}
