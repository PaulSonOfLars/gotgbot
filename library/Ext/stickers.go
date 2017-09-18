package Ext

import (
	"bot/library/Types"
	"strconv"
	"log"
	"encoding/json"
	"net/url"
)

// TODO: inputfile version
// TODO: reply_markup version
func (b Bot) SendStickerStr(chat_id int, sticker string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("sticker", sticker)

	r := Get(b, "sendSticker", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) ReplyStickerStr(chat_id int, sticker string, reply_to_message_id int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chat_id))
	v.Add("sticker", sticker)
	v.Add("reply_to_message_id", strconv.Itoa(reply_to_message_id))

	r := Get(b, "sendSticker", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) GetStickerSet(name string) Types.StickerSet {
	v := url.Values{}
	v.Add("name", name)

	r := Get(b, "getStickerSet", v)
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
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("name", name)
	v.Add("title", title)
	v.Add("png_sticker", png_sticker)
	v.Add("emojis", emojis)

	r := Get(b, "createNewStickerSet", v)
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
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(user_id))
	v.Add("name", name)
	v.Add("png_sticker", png_sticker)
	v.Add("emojis", emojis)

	r := Get(b, "addStickerToSet", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) SetStickerPositionInSet(sticker string, position int) bool {
	v := url.Values{}
	v.Add("sticker", sticker)
	v.Add("position", strconv.Itoa(position))

	r := Get(b, "setStickerPositionInSet", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}

func (b Bot) DeleteStickerFromSet(sticker string) bool {
	v := url.Values{}
	v.Add("sticker", sticker)

	r := Get(b, "deleteStickerFromSet", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var bb bool
	json.Unmarshal(r.Result, &bb)

	return bb
}
