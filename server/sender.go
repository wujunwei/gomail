package server

type MailTask struct {
	from          string
	to            []string
	cc            []string
	bcc           []string
	subject       string
	LastMessageId string
	body          string
	contentType   string
	attachment    Attachment
}

type Attachment struct {
	name        string
	contentType string
	withFile    bool
}
