package ext

import "github.com/PaulSonOfLars/gotgbot/types"

type Chat struct {
	types.Chat
	bot Bot
}

func (chat Chat) SendAction(action string) (bool, error) {
	return chat.bot.SendChatAction(chat.Id, action)
}

func (chat Chat) KickMember(userId int) (bool, error) {
	return chat.bot.KickChatMember(chat.Id, userId)
}

func (chat Chat) UnbanMember(userId int) (bool, error) {
	return chat.bot.UnbanChatMember(chat.Id, userId)
}

func (chat Chat) RestrictMember(userId int) (bool, error) {
	return chat.bot.RestrictChatMember(chat.Id, userId)
}

func (chat Chat) PromoteMember(userId int) (bool, error) {
	return chat.bot.PromoteChatMember(chat.Id, userId)
}

func (chat Chat) ExportInviteLink() (string, error) {
	return chat.bot.ExportChatInviteLink(chat.Id)
}

// TODO
//func (chat Chat) SetChatPhoto() (bool, error) {
//	return chat.bot.SetChatPhoto()
//}

func (chat Chat) DeletePhoto() (bool, error) {
	return chat.bot.DeleteChatPhoto(chat.Id)
}

func (chat Chat) SetTitle(title string) (bool, error) {
	return chat.bot.SetChatTitle(chat.Id, title)
}

func (chat Chat) SetDescription(description string) (bool, error) {
	return chat.bot.SetChatDescription(chat.Id, description)
}

func (chat Chat) PinMessage(messageId int) (bool, error) {
	return chat.bot.PinChatMessage(chat.Id, messageId)
}

func (chat Chat) UnpinMessage() (bool, error) {
	return chat.bot.UnpinChatMessage(chat.Id)
}

func (chat Chat) Leave(description string) (bool, error) {
	return chat.bot.LeaveChat(chat.Id)
}

func (chat Chat) Get() (*Chat, error) {
	return chat.bot.GetChat(chat.Id)
}

func (chat Chat) GetAdministrators() ([]types.ChatMember, error) {
	return chat.bot.GetChatAdministrators(chat.Id)
}

func (chat Chat) GetMembersCount() (int, error) {
	return chat.bot.GetChatMembersCount(chat.Id)
}

func (chat Chat) GetMember(userId int) (*types.ChatMember, error) {
	return chat.bot.GetChatMember(chat.Id, userId)
}

func (chat Chat) SetStickerSet(stickerSetName string) (bool, error) {
	return chat.bot.SetChatStickerSet(chat.Id, stickerSetName)
}

func (chat Chat) DeleteStickerSet() (bool, error) {
	return chat.bot.DeleteChatStickerSet(chat.Id)
}
