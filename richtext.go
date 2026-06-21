package gotgbot

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// GetRichTextTypes gets the list of all richtext types which match the "want" function.
// Eg: m.GetRichTextTypes(IsRichTextType[RichTextString], IsRichTextType[RichTextUrl]).
func (m *Message) GetRichTextTypes(want ...func(RichText) bool) []RichText {
	if m.RichMessage == nil {
		return nil
	}

	var out []RichText
	m.RichMessage.Walk(nil, func(node RichText) bool {
		for _, w := range want {
			if w(node) {
				out = append(out, node)
				break
			}
		}
		return true
	})
	return out
}

// Checks whether an input type is a specific kind of RichText.
// Use like this: gotgbot.IsRichTextType[RichTextString](x).
func IsRichTextType[T RichText](n RichText) bool {
	_, ok := n.(T)
	return ok
}

// GetRichTexts allows for specifying a RichText type, walking over all the blocks, and extracting the result.
func GetRichTexts[T RichText](m *RichMessage) []T {
	var out []T
	m.Walk(nil, func(text RichText) bool {
		if v, ok := text.(T); ok {
			out = append(out, v)
		}
		return true
	})
	return out
}

// GetRichBlocks allows for specifying a RichBlock type, walking over all the blocks, and extracting the result.
func GetRichBlocks[T RichBlock](m *RichMessage) []T {
	var out []T
	m.Walk(func(block RichBlock) bool {
		if v, ok := block.(T); ok {
			out = append(out, v)
		}
		return true
	}, nil)
	return out
}

// WalkRichText recursively visits r and every descendant, calling fn on each node.
// Returning false from fn skips that node's children (but not its siblings).
func WalkRichText(r RichText, fn func(RichText) bool) {
	if r == nil {
		return
	}
	if !fn(r) {
		return
	}
	for _, child := range r.Children() {
		WalkRichText(child, fn)
	}
}

// WalkRichBlock recursively visits b and every descendant block, calling fnBlock
// on each. For each block, all RichText fields are also walked via fnText.
// Returning false from fnBlock skips that block's children.
func WalkRichBlock(b RichBlock, fnBlock func(RichBlock) bool, fnText func(RichText) bool) {
	if b == nil {
		return
	}
	if !fnBlock(b) {
		return
	}
	for _, child := range b.RichTextChildren() {
		WalkRichText(child, fnText)
	}
	for _, child := range b.RichBlockChildren() {
		WalkRichBlock(child, fnBlock, fnText)
	}
}

// Walk visits every RichBlock and RichText node in a RichMessage.
// Pass nil for either callback to skip that class of node.
func (m RichMessage) Walk(fnBlock func(RichBlock) bool, fnText func(RichText) bool) {
	if fnBlock == nil {
		fnBlock = func(RichBlock) bool { return true }
	}
	if fnText == nil {
		fnText = func(RichText) bool { return true }
	}
	for _, block := range m.Blocks {
		WalkRichBlock(block, fnBlock, fnText)
	}
}

// =============================================================================
// RichTextContent - flat text extraction
// =============================================================================

// RichTextContent walks r and concatenates all RichTextString leaf nodes,
// producing the plain-text content with no formatting. This is equivalent to
// the generated GetText() methods but works on any RichText value including
// RichTextArray and the interface itself.
func RichTextContent(r RichText) string {
	var sb strings.Builder
	WalkRichText(r, func(node RichText) bool {
		switch v := node.(type) {
		case RichTextString:
			sb.WriteString(string(v))
		case RichTextMathematicalExpression:
			sb.WriteString(v.Expression)
		case RichTextCustomEmoji:
			sb.WriteString(v.AlternativeText)
		}
		return true
	})
	return sb.String()
}

// RichBlockContent walks b and returns its full plain-text content, descending
// into nested blocks and text nodes. Block-level elements are separated by
// newlines appropriate to their type.
func RichBlockContent(b RichBlock) string {
	var sb strings.Builder
	renderBlockText(b, &sb)
	return sb.String()
}

func renderBlockText(b RichBlock, sb *strings.Builder) {
	switch v := b.(type) {
	case RichBlockParagraph:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
	case RichBlockSectionHeading:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
	case RichBlockPreformatted:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
	case RichBlockFooter:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
	case RichBlockBlockQuotation:
		for _, child := range v.Blocks {
			renderBlockText(child, sb)
		}
		if v.Credit != nil {
			sb.WriteString(RichTextContent(v.Credit))
			sb.WriteString("\n")
		}
	case RichBlockPullQuotation:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
		if v.Credit != nil {
			sb.WriteString(RichTextContent(v.Credit))
			sb.WriteString("\n")
		}
	case RichBlockDetails:
		sb.WriteString(RichTextContent(v.Summary))
		sb.WriteString("\n")
		for _, child := range v.Blocks {
			renderBlockText(child, sb)
		}
	case RichBlockList:
		for _, item := range v.Items {
			for _, child := range item.Blocks {
				renderBlockText(child, sb)
			}
		}
		sb.WriteString("\n")
	case RichBlockTable:
		for _, row := range v.Cells {
			for j, cell := range row {
				if j > 0 {
					sb.WriteString("\t")
				}
				if cell.Text != nil {
					sb.WriteString(RichTextContent(cell.Text))
				}
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	case RichBlockCollage:
		for _, child := range v.Blocks {
			sb.WriteString(RichBlockContent(child))
		}
		captionContent(v.Caption, sb)

	case RichBlockSlideshow:
		for _, child := range v.Blocks {
			sb.WriteString(RichBlockContent(child))
		}
		captionContent(v.Caption, sb)

	case RichBlockThinking:
		sb.WriteString(RichTextContent(v.Text))
		sb.WriteString("\n")
	case RichBlockMathematicalExpression:
		sb.WriteString(v.Expression)
		sb.WriteString("\n")
	case RichBlockDivider:
		sb.WriteString("\n")
	case RichBlockAnchor:
		// No text content.
	case RichBlockPhoto:
		captionContent(v.Caption, sb)
	case RichBlockAnimation:
		captionContent(v.Caption, sb)
	case RichBlockVideo:
		captionContent(v.Caption, sb)
	case RichBlockAudio:
		captionContent(v.Caption, sb)
	case RichBlockVoiceNote:
		captionContent(v.Caption, sb)
	case RichBlockMap:
		captionContent(v.Caption, sb)
	}
}

func captionContent(caption *RichBlockCaption, sb *strings.Builder) {
	if caption == nil {
		return
	}

	sb.WriteString(RichTextContent(caption.Text))
	if caption.Credit != nil {
		sb.WriteString("\n")
		sb.WriteString(RichTextContent(caption.Credit))
	}
	sb.WriteString("\n")
}

// =============================================================================
// RichTextHTML - HTML rendering
// =============================================================================

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

// =============================================================================
// RichTextMarkdown - Markdown rendering
// =============================================================================

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

func fileLink(fileId string) string {
	return fmt.Sprintf("![](fileId://%s)\n", fileId)
}

// =============================================================================
// Shared helpers
// =============================================================================

func renderCaptionHTML(cap *RichBlockCaption, sb *strings.Builder) {
	if cap == nil {
		return
	}
	sb.WriteString("<figcaption>")
	renderTextHTML(cap.Text, sb)
	if cap.Credit != nil {
		sb.WriteString("<cite>")
		renderTextHTML(cap.Credit, sb)
		sb.WriteString("</cite>")
	}
	sb.WriteString("</figcaption>\n")
}

// mdEscape escapes characters that have special meaning in CommonMark.
var mdSpecial = strings.NewReplacer(
	`\`, `\\`,
	"`", "\\`",
	"*", `\*`,
	"_", `\_`,
	"{", `\{`,
	"}", `\}`,
	"[", `\[`,
	"]", `\]`,
	"(", `\(`,
	")", `\)`,
	"#", `\#`,
	"+", `\+`,
	"-", `\-`,
	".", `\.`,
	"!", `\!`,
	"|", `\|`,
)

func mdEscape(s string) string {
	return mdSpecial.Replace(s)
}

func toRoman(num int64) string {
	values := []int64{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	result := ""
	for i := 0; i < len(values); i++ {
		for num >= values[i] {
			result += symbols[i]
			num -= values[i]
		}
	}
	return result
}

func toLetters(num int64) string {
	if num <= 0 {
		return ""
	}

	result := ""
	for num > 0 {
		num-- // shift into 0-25 range for this digit
		letter := rune('a' + (num % 26))
		result = string(letter) + result
		num /= 26
	}
	return result
}
