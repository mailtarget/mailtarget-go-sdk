# Layang Golang SDK

The Layang Golang SDK enable Golang developer to work with Layang API efficiently.

## Getting Started

### Requirements
To run SDK, you will need go1.25+.

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
html := `<p>My fantastic HTML content.<br><br><b>MailTarget</b> <img src=\"cid:AnImage.png\"></p>`

// to
var to = []Address{
	{
		Email: "to@example.com",
		Name:  "To",
	},
	{
		Email: "to2@example.com",
		Name:  "To2",
	},
}

// attachment
attachment := Attachment{
    MimeType:     "image/png",
    Filename: "AnImage.png",
    Value:  "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAAXxJREFUOBFjvJVg84P5718WBjLAX2bmPyxMf/+xMDH8YyZDPwPDXwYGJkIaOXTNGdiUtHAqI2jA/18/GUQzGsg3gMfKg4FVQo6BiYcPqyF4XcChaczA4+DP8P//f4b/P3+SZgAzvxCDSGYjAyMjI8PvZw+AoYXdLuyiQLtE0uoZWAREwLb+fnKXQTipkngXcJu7MnACQx8G2FX1GHgs3bDGBlYX8HlFM/z9+JbhzewWhmf1CQyfti9j+PfzBwO/ZxTMTDiNmQKBfmZX1GB42V/K8P38YbDCX/dvMDAwMzPwuYbBNcIYmC4AhfjvXwx/376AqQHTf96+ZPj34xuKGIiDaQBQ8PPBTQwCoZkMjJzcYA3MgqIMAr7xDJ/3rAHzkQnGO7FWf5gZ/qLmBSZmBoHgNAZee1+Gf18/MzCyczJ83LyQ4fPetch6Gf4xMP3FbgBMGdAgJqAr/n37zABMTTBROA0ygAWUJUG5Civ4B8xwX78CpbD6FJiHmf4AAFicbTMTr5jAAAAAAElFTkSuQmCC",
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
message.SetOptionsAttributes(optionsAttributes)
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
	"github.com/mailtarget/mailtarget-go-sdk"
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
	html := `<p>My fantastic HTML content.<br><br><b>MailTarget</b> <img src=\"cid:AnImage.png\"></p>`
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
	message.SetOptionsAttributes(optionsAttributes)

	// Send the message
	successResponse, errorResponse, err := l.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Succcess Response: %+v Error Response: %+v\n", successResponse, errorResponse)
}
```

## Mailtarget Open API

Beyond sending email, the SDK can reach the Mailtarget Open API (contacts, analytics and more).
That API uses a **second credential**, the Open API secret key, which is issued on request through
your account manager — unlike the Mailtarget API key, you cannot generate it yourself from the
dashboard.

The two credentials are separate on purpose:

| | Mailtarget API key | Open API secret key |
|---|---|---|
| Required | yes | no, only for Open API calls |
| Obtained from | your dashboard | request to Mailtarget admin (CRM) |
| Used for | sending email | Open API resources |

Sending email always uses the Mailtarget API key and is never routed through the Open API, so a
client built without an Open API secret key sends email exactly as before.

Open API resources live in the `openapi` subpackage. `MailtargetClient` embeds it, so its resources
read the same way as in the Java and Python SDKs — directly off `client`:

```go
import (
    "github.com/mailtarget/mailtarget-go-sdk"
    "github.com/mailtarget/mailtarget-go-sdk/openapi"
)

client := layang.NewMailtargetClient(
    privateAPIKey,
    layang.WithOpenAPISecretKey(openAPISecretKey),
)

// Sending still works the same way, through the Transmission API.
successResponse, errorResponse, err := client.Send(message)

// Open API resources.
page, err := client.Contacts.List(nil)
```

Without the secret key, every Open API call fails immediately with an `*openapi.ConfigError` and no
HTTP request is made:

```go
client := layang.NewMailtargetClient(privateAPIKey) // no Open API secret key

_, err := client.Contacts.List(nil)

var configErr *openapi.ConfigError
if errors.As(err, &configErr) {
    // "open API secret key is required to call Contacts.List — pass ..."
}
```

Use `client.HasOpenAPIAccess()` to branch on availability instead of handling the error, or
`client.SetOpenAPISecretKey(...)` to supply the key later.

If you only need the Open API and never send email, use the subpackage on its own:

```go
oa := openapi.New(openapi.WithSecretKey(openAPISecretKey))
page, err := oa.Contacts.List(nil)
```

### Contacts

```go
// List returns one page plus its pagination metadata.
page, err := client.Contacts.List(&openapi.ListContactsParams{
    Page: 1, PerPage: 20, Search: "@mtarget.co",
})
fmt.Println(page.Items, page.Meta.Total)

contact, err := client.Contacts.Get("contact-id")
contact, err = client.Contacts.GetByEmail("recipient@example.com")

contact, err = client.Contacts.Create(&openapi.CreateContactRequest{
    Email:     "recipient@example.com",
    Firstname: "Sarah",
    Labels:    []string{"vip"},
})

contact, err = client.Contacts.Update("contact-id", &openapi.UpdateContactRequest{
    Note: "renewed",
})

err = client.Contacts.Delete("contact-id")

total, err := client.Contacts.Count(nil)
```

### Analytics

```go
summary, err := client.Analytics.Summary(&openapi.AnalyticsSummaryParams{
    From: "2026-08-01",
    To:   "2026-08-31",
})

breakdown, err := client.Analytics.SummaryBreakdown(&openapi.AnalyticsSummaryParams{
    From: "2026-08-01", To: "2026-08-31", GroupBy: "sender",
})

detail, err := client.Analytics.Transmission("transmission-id")
events, err := client.Analytics.TransmissionEvents("transmission-id")
```

### Errors and timeouts

Open API failures come back as `*openapi.Error`, carrying the HTTP status plus the API's own
`error` code and `message`:

```go
var apiErr *openapi.Error
if errors.As(err, &apiErr) {
    log.Printf("status=%d code=%s message=%s", apiErr.StatusCode, apiErr.Code, apiErr.Message)
}
```

Open API requests time out after 30 seconds by default; override with
`layang.WithOpenAPITimeout(d)`.

### Available resources

Currently implemented: **Contacts**, **Analytics**. The remaining Open API resources (API Keys,
Campaigns, Senders, Sending Domains, Labels, Settings, Sub Accounts and the Open API's own
transmissions endpoint) are being rolled out incrementally.

### Adding a new Open API resource

The layout is designed so a new resource stays a local change:

1. Add `openapi/<resource>.go` with its models and a `<Resource>Service` whose methods call the
   shared helpers (`object`, `objects`, `paged`, `bare`) or `c.do` directly.
2. Add one field for the service on `openapi.Client` and wire it in `New`.
3. Add `openapi/<resource>_test.go`.

Nothing in the root package changes, and the secret key guard applies automatically because every
call goes through `Client.do`.
