package gotgbot

type RichTextString string

func (r RichTextString) GetType() string {
	return "string"
}

func (r RichTextString) richText() {}

type RichTextArray []RichText

func (r RichTextArray) GetType() string {
	return "array"
}

func (r RichTextArray) richText() {}
