package layang

import (
	"errors"
	"fmt"
)

type Message struct {
	Subject           string            `json:"subject"`
	Html              string            `json:"html"`
	Text              string            `json:"text"`
	Recipients        Recipient         `json:"recipients"`
	From              Address           `json:"from"`
	Attachment        []Attachment      `json:"attachment"`
	Metadata          map[string]string `json:"metadata"`
	OptionsAttributes OptionsAttributes `json:"optionsAttributes"`
}

type Recipient struct {
	Address Address `json:"address"`
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Attachment struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type OptionsAttributes struct {
	ClickTracking bool `json:"clickTracking"`
	OpenTracking  bool `json:"openTracking"`
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
	if !IsEmail(l.From.Email) {
		return errors.New(l.From.Email + " is not valid email address")
	}
	if !IsEmail(l.Recipients.Address.Email) {
		return errors.New(l.Recipients.Address.Email + " is not valid email address")
	}
	if !IsHTML(l.Html) {
		fmt.Println("not valid")
		return errors.New("not valid html")
	} else {
		fmt.Println("valid")
	}
	return nil
}

func (l *Message) SetAttachment(attachment []Attachment) {
	l.Attachment = attachment
}

func (l *Message) SetMetadata(metadata map[string]string) {
	l.Metadata = metadata
}

func (l *Message) setOptionsAttributes(attributes OptionsAttributes) {
	l.OptionsAttributes = attributes
}
