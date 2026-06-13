package gotgbot

import "strings"

type RichTextString string

func (r RichTextString) GetType() string {
	return "string"
}

func (r RichTextString) GetText() string {
	return string(r)
}

func (r RichTextString) richText() {}

type RichTextArray []RichText

func (r RichTextArray) GetType() string {
	return "array"
}

func (r RichTextArray) GetText() string {
	bd := strings.Builder{}
	for _, v := range r {
		bd.WriteString(v.GetText())
	}
	return bd.String()
}

func (r RichTextArray) richText() {}

type RichBlockArray []RichBlock

func (r RichBlockArray) GetText() string {
	bd := strings.Builder{}
	for _, b := range r {
		bd.WriteString(b.GetText() + "\n")
	}
	return bd.String()
}

type RichBlockListItemArray []RichBlockListItem

func (r RichBlockListItemArray) GetText() string {
	bd := strings.Builder{}
	for _, b := range r {
		bd.WriteString(b.GetText())
	}
	return bd.String()
}

func (v RichBlockDivider) GetText() string {
	return "\n"
}

func (v RichBlockMathematicalExpression) GetText() string {
	return v.Expression
}

func (v RichBlockAnchor) GetText() string {
	return v.Name
}

func (v RichTextCustomEmoji) GetText() string {
	return v.AlternativeText
}

func (v RichTextMathematicalExpression) GetText() string {
	return v.Expression
}

func (v RichTextAnchor) GetText() string {
	return v.Name
}
