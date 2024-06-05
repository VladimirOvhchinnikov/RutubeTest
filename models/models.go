package models

type ShortUserInfo struct {
	ID        int
	IDTG      int
	FirstName string
	LastName  string
	BirthDate string
}

type UserInfo struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name,omitempty"`
			Username     string `json:"username,omitempty"`
			LanguageCode string `json:"language_code,omitempty"`
			IsPremium    bool   `json:"is_premium"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Username  string `json:"username,omitempty"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text,omitempty"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities,omitempty"`
	} `json:"message"`
}

type AllMessage struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
			IsPremium    bool   `json:"is_premium"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}
