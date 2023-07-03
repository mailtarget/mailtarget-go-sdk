package layang

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIsValid_Success(t *testing.T) {
	layang := NewLayang(privateAPIKey)
	err := layang.IsValid()
	assert.NoErrorf(t, err, "")
}

func TestIsValid_Error(t *testing.T) {
	layang := NewLayang("")
	err := layang.IsValid()
	assert.Error(t, err)
}

func TestSend_Success(t *testing.T) {

	//setup layang
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, to)
	message.SetAttachment(attachments)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://apiconfig.mailtarget.co/v1/layang/transmissions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, successResponse)
		})

	//send
	successResponseActual, _, _ := layang.Send(message)

	//check
	assert.Equal(t, &successResponse, successResponseActual)
}

func TestSend_ErrorLayangInvalid(t *testing.T) {

	//setup layang
	layang := NewLayang("")
	message := layang.NewMessage(subject, body, html, sender, to)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:3000/send",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, successResponse)
		})

	//send
	_, _, err := layang.Send(message)

	//check
	assert.Error(t, err)
}

func TestSend_ErrorMessageInvalid(t *testing.T) {

	//setup layang
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, wrongSender, to)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:3000/send",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, successResponse)
		})

	//send
	_, _, err := layang.Send(message)

	//check
	assert.Error(t, err)
}

func TestLayang_Send(t *testing.T) {
	var subject = "Hello from Mailtarget"
	var body = "Congratulation, you just sent email with Mailtarget. You are truly awesome!"
	var html = "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF - 8\"><title>Hello from Mailtarget</title></head><body><p>Congratulation, you just sent email with Mailtarget. You are truly awesome!</p></body></html>"
	var sender = Address{
		Email: "default@sandbox.mailtarget.co",
		Name:  "MailTarget Sandbox",
	}
	var to = []Address{
		{
			Email: "penguinoflostatlantica@gmail.com",
			Name:  "",
		},
	}

	//setup layang
	layang := NewLayang("8U1UyCCo5GlG6KJAIR5Ebzsj")
	message := layang.NewMessage(subject, body, html, sender, to)

	//send
	successResponseActual, errorResponseActual, err := layang.Send(message)
	if err != nil {
		println(err.Error())
	}
	fmt.Println("Success Response ", successResponseActual)
	fmt.Println("Error Response ", errorResponseActual)
}
