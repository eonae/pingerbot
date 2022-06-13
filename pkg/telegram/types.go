package telegram

type Chat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	Title     string `json:"title"`
}

type User struct {
	Id           int64   `json:"id"`
	FirstName    string  `json:"first_name"`
	Username     string  `json:"username"`
	IsBot        bool    `json:"is_bot"`
	LanguageCode *string `json:"language_code"`
}

type Entity struct {
	Length int    `json:"length"`
	Offset int    `json:"offset"`
	Type   string `json:"type"`
}

type IncomingMessage struct {
	Chat            Chat
	Date            int `json:"date"`
	Entities        []Entity
	From            User
	Id              int    `json:"message_id"`
	Text            string `json:"text"`
	NewMember       *User  `json:"new_chat_member,omitempty"`
	NewParticipant  *User  `json:"new_chat_participant,omitempty"`
	LeftMember      *User  `json:"left_chat_participant,omitempty"`
	LeftParticipant *User  `json:"left_chat_member,omitempty"`
}

type JoinLeave struct {
	Chat          Chat
	Date          int `json:"date"`
	From          User
	OldChatMember struct {
		Status string `json:"status"`
		User   User   `json:"user"`
	} `json:"old_chat_member,omitempty"`
	NewChatMember struct {
		Status string `json:"status"`
		User   User   `json:"user"`
	} `json:"new_chat_member,omitempty"`
}

type Update struct {
	UpdateId      int64            `json:"update_id"`
	Message       *IncomingMessage `json:"message,omitempty"`        // Optional
	EditedMessage *IncomingMessage `json:"edited_message,omitempty"` // Optional
	MyChatMember  *JoinLeave       `json:"my_chat_member,omitempty"` // Optional
}

type Me struct {
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	FirstName               string `json:"first_name"`
	Id                      int64  `json:"id"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
	Username                string `json:"username"`
}

type MsgContent struct {
	Text           string      `json:"text"`
	ParseMode      string      `json:"parse_mode,omitempty"`
	Entities       []Entity    `json:"entities,omitempty"`
	ProtectContent bool        `json:"protect_content,omitempty"`
	ReplyMarkup    interface{} `json:"reply_markup,omitempty"`
}

type OutgoingMessage struct {
	ChatId  int64 `json:"chat_id"`
	ReplyTo int   `json:"reply_to_message_id,omitempty"`
	MsgContent
}
