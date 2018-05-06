package types

type Chat struct {
	Id                int
	Type              string
	Title             string
	Username          string
	First_name        string
	Last_name         string
	All_members_admin bool
	Photo             *ChatPhoto
	Description       string
	Invite_link       string
	Pinned_message    *Message
}

type ChatPhoto struct {
	Small_file_id string
	Big_file_id   string

}

type ChatMember struct {
	User                      *User
	Status                    string
	Until_date                int
	Can_be_edited             bool
	Can_change_info           bool
	Can_post_messages         bool
	Can_edit_messages         bool
	Can_delete_messages       bool
	Can_invite_users          bool
	Can_restrict_members      bool
	Can_pin_messages          bool
	Can_promote_members       bool
	Can_send_messages         bool
	Can_send_media_messages   bool
	Can_send_other_messages   bool
	Can_add_web_page_previews bool

}
