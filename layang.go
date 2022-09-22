package layang

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

const (
	APIBase = "https://smtpdevconf.mtarget.id"
)

type Layang struct {
	apiBase    string
	apiVersion string
	apiKey     string
	resty      *resty.Client
}

func NewLayang(apiKey string) *Layang {
	return &Layang{
		apiBase:    APIBase,
		apiKey:     apiKey,
		apiVersion: "v1",
		resty:      resty.New(),
	}
}

func (l *Layang) IsValid() error {
	if l.apiKey == "" {
		return errors.New("api key is empty")
	}
	return nil
}

func (l *Layang) Send(message *Message) (*SuccessResponse, *ErrorResponse, error) {

	if err := l.IsValid(); err != nil {
		return &SuccessResponse{}, &ErrorResponse{}, err
	}
	if err := message.IsValid(); err != nil {
		return &SuccessResponse{}, &ErrorResponse{}, err
	}

	return l.sendMessage(message)
}
