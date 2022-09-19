package layang

import (
	"fmt"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func (l *Layang) sendMessage(message *Message) (*SuccessResponse, *ErrorResponse, error) {
	resp, err := l.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Key", l.apiKey).
		SetBody(message).
		SetResult(SuccessResponse{}).
		SetError(ErrorResponse{}).
		Post(l.apiBase + "/" + l.apiVersion + "/send")
	if resp.IsSuccess() {
		println("success")
		fmt.Printf("%+v", resp.Result())
		return resp.Result().(*SuccessResponse), nil, err
	}
	if resp.IsError() {
		println("error")
		errorResponse := resp.Error().(*ErrorResponse)
		errorResponse.Status = resp.StatusCode()
		return nil, errorResponse, err
	}

	return nil, nil, err
}
