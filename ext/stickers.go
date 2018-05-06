package ext

import (
	"gotgbot/types"
	"strconv"
	"log"
	"encoding/json"
	"net/url"
)

// TODO: inputfile version
// TODO: reply_markup version
func (b Bot) SendStickerStr(chatId int, sticker string) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker", sticker)

	r := Get(b, "sendSticker", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}

func (b Bot) ReplyStickerStr(chatId int, sticker string, replyToMessageId int) Message {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(chatId))
	v.Add("sticker", sticker)
	v.Add("reply_to_message_id", strconv.Itoa(replyToMessageId))

	r := Get(b, "sendSticker", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	return b.ParseMessage(r.Result)

}
func (b Bot) GetStickerSet(name string) types.StickerSet {
	v := url.Values{}
	v.Add("name", name)

	r := Get(b, "getStickerSet", v)
	if !r.Ok {
		log.Println("You done goofed")
		log.Println(r)
	}

	var ss types.StickerSet
	json.Unmarshal(r.Result, &ss)

	return ss
}

// TODO: input file stuff
//func (b Bot) UploadStickerFile(user_id int, png_sticker types.InputFile) types.File {
//
//}

// TODO: contains mask + mask position version
// TODO: InputFile version
// TODO: check return
func (b Bot) CreateNewStickerSetStr(userId int, name string, title string, pngSticker string, emojis string) bool {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("name", name)
	v.Add("title", title)
	v.Add("png_sticker", pngSticker)
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
func (b Bot) AddStickerToSetStr(userId int, name string, pngSticker string, emojis string) bool {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(userId))
	v.Add("name", name)
	v.Add("png_sticker", pngSticker)
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
