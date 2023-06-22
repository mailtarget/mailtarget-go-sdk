package layang

import (
	"errors"
)

type SuccessResponse struct {
	TransmissionId string `json:"transmissionId"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (l *Layang) sendMessage(message *Message) (*SuccessResponse, *ErrorResponse, error) {
	resp, err := l.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+l.apiKey).
		SetHeader("Accept", "application/json").
		SetBody(message).
		EnableTrace().
		SetResult(SuccessResponse{}).
		SetError(ErrorResponse{}).
		Post(l.apiBase + "/" + l.apiVersion + "/layang/transmissions")

	if resp.IsSuccess() {
		return resp.Result().(*SuccessResponse), nil, err
	}
	if resp.IsError() {
		errorResponse := resp.Error().(*ErrorResponse)
		errorResponse.Status = resp.StatusCode()
		return nil, errorResponse, errors.New(resp.String())
	}

	return nil, nil, err
}
