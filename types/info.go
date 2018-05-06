package types

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
}
