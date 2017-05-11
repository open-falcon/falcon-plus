package cron

type MailDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type SmsDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Content  string `json:"content"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type ChatDto struct {
	Priority int    `json:"priority"`
	Metric   string `json:"metric"`
	Content  string `json:"content"`
	IM       string `json:"im"`
	Status   string `json:"status"`
}
