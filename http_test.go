package layang

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendMessage_Success(t *testing.T) {

	//setup layang
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, recipient)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:3000/v1/send",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, successResponse)
		})

	//send
	successResponseActual, _, _ := layang.Send(message)

	//check
	assert.Equal(t, &successResponse, successResponseActual)
}

func TestSendMessage_Error(t *testing.T) {

	//setup layang
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, recipient)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:3000/v1/send",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	//send
	_, errorResponseActual, _ := layang.Send(message)

	//check
	assert.Equal(t, &errorResponse, errorResponseActual)
}

func TestSendMessage_ErrorResty(t *testing.T) {

	//setup layang
	layang := NewLayang(privateAPIKey)
	message := layang.NewMessage(subject, body, html, sender, recipient)

	//mock http
	httpmock.ActivateNonDefault(layang.resty.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:3000/v1/send/cus",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	//send
	_, _, err := layang.Send(message)

	//check
	assert.Error(t, err)
}
