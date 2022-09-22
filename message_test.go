package layang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_IsValid_Success(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, recipient)
	err := message.IsValid()
	assert.NoErrorf(t, err, "")
}

func TestMessage_IsValid_ErrorMessageIsEmpty(t *testing.T) {
	message := &Message{}
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorMessageIsEmpty2(t *testing.T) {
	message := &Message{}
	message = nil
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorTextIsEmpty(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, "", "", sender, recipient)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorZeroRecipient(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, Recipient{Address: Address{}})
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorNotValidRecipientEmail(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, wrongRecipient)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_Error(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, wrongSender, wrongRecipient)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorHtmlNotValid(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, wrongHtml, sender, recipient)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_SetAttachment(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, wrongHtml, sender, wrongRecipient)
	message.SetAttachment([]Attachment{attachment})
}
