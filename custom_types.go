package gotgbot

type RichTextString string

func (r RichTextString) GetType() string {
	return "string"
}

func (r RichTextString) richText() {}

func (v RichTextString) Children() []RichText { return nil }

type RichTextArray []RichText

func (r RichTextArray) GetType() string {
	return "array"
}

func (r RichTextArray) richText() {}

func (v RichTextArray) Children() []RichText { return []RichText(v) }

type RichBlockArray []RichBlock

type RichBlockListItemArray []RichBlockListItem
