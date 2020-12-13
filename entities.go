package gotgbot

import "unicode/utf16"

type ParsedMessageEntity struct {
	MessageEntity
	Text string `json:"text"`
}

func (m Message) ParseEntities() (out []ParsedMessageEntity) {
	return m.ParseEntityTypes(nil)

}

func (m Message) ParseCaptionEntities() (out []ParsedMessageEntity) {
	return m.ParseCaptionEntityTypes(nil)
}

func (m Message) ParseEntityTypes(accepted map[string]struct{}) (out []ParsedMessageEntity) {
	utf16Text := utf16.Encode([]rune(m.Text))
	for _, ent := range m.Entities {
		if _, ok := accepted[ent.Type]; ok || accepted == nil {
			out = append(out, parseEntity(ent, utf16Text))
		}
	}
	return out
}

func (m Message) ParseCaptionEntityTypes(accepted map[string]struct{}) (out []ParsedMessageEntity) {
	utf16Caption := utf16.Encode([]rune(m.Caption))
	for _, ent := range m.CaptionEntities {
		if _, ok := accepted[ent.Type]; ok || accepted == nil {
			out = append(out, parseEntity(ent, utf16Caption))
		}
	}
	return out
}

func (m Message) ParseEntity(entity MessageEntity) ParsedMessageEntity {
	return parseEntity(entity, utf16.Encode([]rune(m.Text)))
}

func (m Message) ParseCaptionEntity(entity MessageEntity) ParsedMessageEntity {
	return parseEntity(entity, utf16.Encode([]rune(m.Caption)))
}

func parseEntity(entity MessageEntity, utf16Text []uint16) ParsedMessageEntity {
	text := string(utf16.Decode(utf16Text[entity.Offset : entity.Offset+entity.Length]))

	if entity.Type == "url" {
		entity.Url = text
	}
	entity.Offset = int64(len(string(utf16.Decode(utf16Text[:entity.Offset]))))
	return ParsedMessageEntity{
		MessageEntity: entity,
		Text:          text,
	}
}
