package layang

import (
	"errors"
)

type Message struct {
	Subject    string
	Html       string
	Text       string
	Recipients Recipient
	From       Address
	Attachment []Attachment
}

type Recipient struct {
	Address []Address
}

type Address struct {
	Email string
	Name  string
}

type Attachment struct {
	Type     string
	Filename string
	Content  string
}

func (l *Layang) NewMessage(subject, text, html string, from Address, recipient Recipient) *Message {
	return &Message{
		Subject:    subject,
		Text:       text,
		Html:       html,
		Recipients: recipient,
		From:       from,
	}
}

func (l *Message) IsValid() error {
	if l == nil || (l == &Message{}) {
		return errors.New("message is empty")
	}
	if l.From.Email == "" {
		return errors.New("from is empty")
	}
	if l.Text == "" && l.Html == "" {
		return errors.New("text or html is empty")
	}
	if len(l.Recipients.Address) == 0 {
		return errors.New("recipient must not 0")
	}
	if !IsEmail(l.From.Email) {
		return errors.New(l.From.Email + " is not valid email address")
	}
	for _, to := range l.Recipients.Address {
		if !IsEmail(to.Email) {
			return errors.New(to.Email + " is not valid email address")
		}
	}

	return nil
}

func (l *Message) SetHtml(body string) error {
	if !IsHTML(body) {
		return errors.New("not valid html")
	}

	l.Html = body

	return nil
}

func (l *Message) SetAttachment(attachment []Attachment) {
	l.Attachment = attachment
}
