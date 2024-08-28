package types

type Email struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type EmailData struct {
	Email Email `json:"email"`
}
