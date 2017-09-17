package Ext

import (
	"bot/library/Types"
	"strconv"
	"log"
	"encoding/json"
)

// TODO: inputfile version
// TODO: reply version
// TODO: reply_markup version
func (b Bot) SendStickerStr(chat_id int, sticker string) Message {
	m := make(map[string]string)
	m["chat_id"] = strconv.Itoa(chat_id)
	m["sticker"] = sticker

	r := Get(b, "sendSticker", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) GetStickerSet(name string) Types.StickerSet {
	m := make(map[string]string)
	m["name"] = name

	r := Get(b, "getStickerSet", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var ss Types.StickerSet
	json.Unmarshal(r.Result, &ss)

	return ss
}

// TODO: input file stuff
//func (b Bot) UploadStickerFile(user_id int, png_sticker Types.InputFile) Types.File {
//
//}

// TODO: contains mask + mask position version
// TODO: InputFile version
// TODO: check return
func (b Bot) CreateNewStickerSetStr(user_id int, name string, title string, png_sticker string, emojis string) bool {
	m := make(map[string]string)
	m["user_id"] = strconv.Itoa(user_id)
	m["name"] = name
	m["title"] = title
	m["png_sticker"] = png_sticker
	m["emojis"] = emojis

	r := Get(b, "createNewStickerSet", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

// TODO: InputFile version
// TODO: mask position version
func (b Bot) AddStickerToSetStr(user_id int, name string, png_sticker string, emojis string) bool {
	m := make(map[string]string)
	m["user_id"] = strconv.Itoa(user_id)
	m["name"] = name
	m["png_sticker"] = png_sticker
	m["emojis"] = emojis

	r := Get(b, "addStickerToSet", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) SetStickerPositionInSet(sticker string, position int) bool {
	m := make(map[string]string)
	m["sticker"] = sticker
	m["position"] = strconv.Itoa(position)

	r := Get(b, "setStickerPositionInSet", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) DeleteStickerFromSet(sticker string) bool {
	m := make(map[string]string)
	m["sticker"] = sticker

	r := Get(b, "deleteStickerFromSet", m)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}
