# Layang with Go
Go library for interacting with the Layang API.
## Usage

```go
package main

import (
	"fmt"
	"github.com/mailtarget/layang-sdk-go"
	"log"
)

var privateAPIKey = "your-private-key"

func main() {
	// Create an instance of the Layang Client
	l := layang.NewLayang(privateAPIKey)

	sender := layang.Address{
		Email: "sender@example.com",
		Name:  "Sender",
	}
	subject := "Fancy subject!"
	body := "Hello from Layang Go!"
	html := `
<p>My fantastic HTML content.<br><br><b>SparkPost</b> <img src=\"cid:AnImage.png\"></p>
`
	recipient := layang.Recipient{Address: []layang.Address{
		{
			Email: "recipient@example.com",
			Name:  "Recipient",
		},
		{
			Email: "recipient2@example.com",
			Name:  "Recipient 2",
		},
	}}

	// The message object allows you to add attachments and Bcc recipients
	message := l.NewMessage(subject, body, html, sender, recipient)

	// Send the message
	successResponse, errorResponse, err := l.Send(message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Succcess Response: %+v Error Response: %+v\n", successResponse, errorResponse)
}
```