package domain

type SenderMail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type SenderMailParams struct {
	From    string `json:"from" binding:"required,email"`
	To      string `json:"to" binding:"required,email,nefield=From"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}
