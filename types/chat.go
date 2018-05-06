package types

type Chat struct {
	Id              int        `json:"id"`
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	AllMembersAdmin bool       `json:"all_members_admin"`
	Photo           *ChatPhoto `json:"photo"`
	Description     string     `json:"description"`
	InviteLink      string     `json:"invite_link"`
	PinnedMessage   *Message   `json:"pinned_message"`
}

type ChatPhoto struct {
	SmallFileId string `json:"small_file_id"`
	BigFileId   string `json:"big_file_id"`
}

type ChatMember struct {
	User                  *User  `json:"user"`
	Status                string `json:"status"`
	UntilDate             int    `json:"until_date"`
	CanBeEdited           bool   `json:"can_be_edited"`
	CanChangeInfo         bool   `json:"can_change_info"`
	CanPostMessages       bool   `json:"can_post_messages"`
	CanEditMessages       bool   `json:"can_edit_messages"`
	CanDeleteMessages     bool   `json:"can_delete_messages"`
	CanInviteUsers        bool   `json:"can_invite_users"`
	CanRestrictMembers    bool   `json:"can_restrict_members"`
	CanPinMessages        bool   `json:"can_pin_messages"`
	CanPromoteMembers     bool   `json:"can_promote_members"`
	CanSendMessages       bool   `json:"can_send_messages"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews"`
}
