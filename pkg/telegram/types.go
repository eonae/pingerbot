package telegram

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
	Username  string `json:"username"`
}

type User struct {
	Id           int    `json:"id"`
	FirstName    string `json:"first_name"`
	IsBot        bool   `json:"is_bot"`
	LanguageCode string `json:"language_code"`
	Username     string `json:"username"`
}

type Entity struct {
	Length int    `json:"length"`
	Offset int    `json:"offset"`
	Type   string `json:"type"`
}

type Update struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		Chat      Chat
		Date      int `json:"date"`
		Entities  []Entity
		From      User
		MessageId int    `json:"message_id"`
		Text      string `json:"text"`
	}
}

type Me struct {
	Ok     bool `json:"ok"`
	Result struct {
		CanJoinGroups           bool   `json:"can_join_groups"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
		FirstName               string `json:"first_name"`
		Id                      int    `json:"id"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries"`
		Username                string `json:"username"`
	}
}
