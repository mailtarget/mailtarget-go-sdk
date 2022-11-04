package layang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_IsValid_Success(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, to)
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
	message := layang.NewMessage(subject, "", "", sender, to)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorZeroRecipient(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, []Address{})
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorNotValidRecipientEmail(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, wrongTo)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorNotValidCC(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, to)
	message.SetCC(wrongTo)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorNotValidBCC(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, to)
	message.SetBCC(wrongTo)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorNotValidReplyTo(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, to)
	message.SetReplyTo(wrongTo)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_Error(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, wrongSender, wrongTo)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_IsValid_ErrorHtmlNotValid(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, wrongHtml, sender, to)
	err := message.IsValid()
	assert.Error(t, err)
}

func TestMessage_SetAttachment(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, wrongHtml, sender, wrongTo)
	message.SetAttachment([]Attachment{attachment})
}

func TestMessage_ShouldSuccessSetHeader(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, wrongHtml, sender, wrongTo)
	message.SetHeaders(headers)
}
