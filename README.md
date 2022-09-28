# Layang Golang SDK

The Layang Golang SDK enable Golang developer to work with Layang API efficiently.

## Getting Started

### Requirements
To run SDK, you will need go1.17.5+.

### Authentication
When you Sign Up, you can generate API Key in Layang Dashboard. To view your API Key in the Layang dashboard, click on Configuration on the left-hand navbar in the Layang dashboard and then API Key.

## Setup Client

### Client Configuration
Default client configuration :
```go
l := layang.NewLayang(privateAPIKey)
```

### Form Example
```go

// sender
sender := layang.Address{
    Email: "sender@example.com",
    Name:  "Sender",
}

// subject
subject := "Fancy subject!"

// text
text := "Hello from Layang Go!"

// html
html := `<p>My fantastic HTML content.<br><br><b>SparkPost</b> <img src=\"cid:AnImage.png\"></p>`

// recipient
recipient := layang.Recipient{Address: layang.Address{
    Email: "recipient@example.com",
    Name:  "Recipient",
}}

// attachment
attachment := Attachment{
    Type:     "image/png",
    Filename: "AnImage.png",
    Content:  "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAAXxJREFUOBFjvJVg84P5718WBjLAX2bmPyxMf/+xMDH8YyZDPwPDXwYGJkIaOXTNGdiUtHAqI2jA/18/GUQzGsg3gMfKg4FVQo6BiYcPqyF4XcChaczA4+DP8P//f4b/P3+SZgAzvxCDSGYjAyMjI8PvZw+AoYXdLuyiQLtE0uoZWAREwLb+fnKXQTipkngXcJu7MnACQx8G2FX1GHgs3bDGBlYX8HlFM/z9+JbhzewWhmf1CQyfti9j+PfzBwO/ZxTMTDiNmQKBfmZX1GB42V/K8P38YbDCX/dvMDAwMzPwuYbBNcIYmC4AhfjvXwx/376AqQHTf96+ZPj34xuKGIiDaQBQ8PPBTQwCoZkMjJzcYA3MgqIMAr7xDJ/3rAHzkQnGO7FWf5gZ/qLmBSZmBoHgNAZee1+Gf18/MzCyczJ83LyQ4fPetch6Gf4xMP3FbgBMGdAgJqAr/n37zABMTTBROA0ygAWUJUG5Civ4B8xwX78CpbD6FJiHmf4AAFicbTMTr5jAAAAAAElFTkSuQmCC",
}
var attachments := []Attachment{
    attachment,
}

// metadata
var metadata := map[string]string{"key1": "value1", "key2": "value2"}

// options attribute
var optionsAttributes := OptionsAttributes{
    ClickTracking: true,
    OpenTracking:  true,
}
```

### Creating Message
Message is the payload used to send email. There are some field that required like subject, text, html, sender and recipient.
```go
message := l.NewMessage(subject, text, html, sender, recipient)
```

### Set Attachment
```go
message.SetAttachment(attachments)
```

### Set Metadata
```go
message.SetMetadata(metadata)
```

### Set Options Attribute
```go
message.setOptionsAttributes(optionsAttributes)
```

### Sending Message
```go
successResponse, errorResponse, err := l.Send(message)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Succcess Response: %+v Error Response: %+v\n", successResponse, errorResponse)
```

### Full Example
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
	html := `<p>My fantastic HTML content.<br><br><b>SparkPost</b> <img src=\"cid:AnImage.png\"></p>`
	recipient := layang.Recipient{Address: layang.Address{
		Email: "recipient@example.com",
		Name:  "Recipient",
	}}
	attachment := layang.Attachment{
		Type:     "image/png",
		Filename: "AnImage.png",
		Content:  "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAAXxJREFUOBFjvJVg84P5718WBjLAX2bmPyxMf/+xMDH8YyZDPwPDXwYGJkIaOXTNGdiUtHAqI2jA/18/GUQzGsg3gMfKg4FVQo6BiYcPqyF4XcChaczA4+DP8P//f4b/P3+SZgAzvxCDSGYjAyMjI8PvZw+AoYXdLuyiQLtE0uoZWAREwLb+fnKXQTipkngXcJu7MnACQx8G2FX1GHgs3bDGBlYX8HlFM/z9+JbhzewWhmf1CQyfti9j+PfzBwO/ZxTMTDiNmQKBfmZX1GB42V/K8P38YbDCX/dvMDAwMzPwuYbBNcIYmC4AhfjvXwx/376AqQHTf96+ZPj34xuKGIiDaQBQ8PPBTQwCoZkMjJzcYA3MgqIMAr7xDJ/3rAHzkQnGO7FWf5gZ/qLmBSZmBoHgNAZee1+Gf18/MzCyczJ83LyQ4fPetch6Gf4xMP3FbgBMGdAgJqAr/n37zABMTTBROA0ygAWUJUG5Civ4B8xwX78CpbD6FJiHmf4AAFicbTMTr5jAAAAAAElFTkSuQmCC",
	}
	attachments := []layang.Attachment{
		attachment,
	}
	metadata := map[string]string{"key1": "value1", "key2": "value2"}
	optionsAttributes := layang.OptionsAttributes{
		ClickTracking: true,
		OpenTracking:  true,
	}

	// The message object allows you to add attachments meta and options attribute
	message := l.NewMessage(subject, body, html, sender, recipient)
	message.SetAttachment(attachments)
	message.SetMetadata(metadata)
	message.setOptionsAttributes(optionsAttributes)
	
	// Send the message
	successResponse, errorResponse, err := l.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Succcess Response: %+v Error Response: %+v\n", successResponse, errorResponse)
}
```