package layang

import (
	"errors"
	"fmt"
)

type Message struct {
	Subject           string            `json:"subject"`
	From              Address           `json:"from"`
	ReplyTo           []Address         `json:"replyTo"`
	To                []Address         `json:"to"`
	Cc                []Address         `json:"cc"`
	Bcc               []Address         `json:"bcc"`
	BodyText          string            `json:"bodyText"`
	BodyHtml          string            `json:"bodyHtml"`
	Headers           []Header          `json:"headers"`
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
	MimeType string `json:"mimeType"`
	Filename string `json:"filename"`
	Value    string `json:"value"`
}

type OptionsAttributes struct {
	ClickTracking bool `json:"clickTracking"`
	OpenTracking  bool `json:"openTracking"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (l *Layang) NewMessage(subject, text, html string, from Address, to []Address) *Message {
	return &Message{
		Subject:  subject,
		BodyText: text,
		BodyHtml: html,
		To:       to,
		From:     from,
	}
}

func (l *Message) IsValid() error {
	if l == nil || (l == &Message{}) {
		return errors.New("message is empty")
	}
	if l.From.Email == "" {
		return errors.New("from is empty")
	}
	if l.BodyText == "" && l.BodyHtml == "" {
		return errors.New("text or html is empty")
	}
	if !IsEmail(l.From.Email) {
		return errors.New(l.From.Email + " is not valid email address")
	}
	if len(l.To) == 0 {
		return errors.New("to is empty")
	}
	for i := range l.To {
		if !IsEmail(l.To[i].Email) {
			return errors.New(l.To[i].Email + " is not valid email address")
		}
	}

	for i := range l.Cc {
		if !IsEmail(l.Cc[i].Email) {
			return errors.New(l.Cc[i].Email + " is not valid email address")
		}
	}

	for i := range l.Bcc {
		if !IsEmail(l.Bcc[i].Email) {
			return errors.New(l.Bcc[i].Email + " is not valid email address")
		}
	}

	for i := range l.ReplyTo {
		if !IsEmail(l.ReplyTo[i].Email) {
			return errors.New(l.ReplyTo[i].Email + " is not valid email address")
		}
	}

	if !IsHTML(l.BodyHtml) {
		fmt.Println("not valid")
		return errors.New("not valid html")
	} else {
		fmt.Println("valid")
	}
	return nil
}

func (l *Message) SetReplyTo(replyTo []Address){
	l.ReplyTo = replyTo
}

func (l *Message) SetCC(cc []Address)  {
	l.Cc = cc
}

func (l *Message) SetBCC(bcc []Address)  {
	l.Bcc = bcc
}

func (l *Message) SetHeaders(headers []Header)  {
	l.Headers = headers
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
