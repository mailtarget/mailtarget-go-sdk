# Layang Golang SDK

The Layang Golang SDK enable Golang developer to work with Layang API efficiently.

## Getting Started

### Requirements
To run SDK, you will need go1.18+.

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

// Email, Firstname and Labels are required on update.
contact, err = client.Contacts.Update("contact-id", &openapi.UpdateContactRequest{
    Email:     "recipient@example.com",
    Firstname: "Sarah",
    Labels:    []string{"vip"},
    Note:      "renewed",
})

err = client.Contacts.Delete("contact-id")

total, err := client.Contacts.Count(nil)

// Export returns one page as CSV. Paging past the last contact returns an
// *openapi.Error with StatusCode 404, which marks the end of the export.
csv, err := client.Contacts.Export(&openapi.ListContactsParams{Page: 1, PerPage: 100})

// Every field a contact has, built in and custom; use FieldName in an import mapping.
fields, err := client.Contacts.Fields()
```

#### Importing contacts from CSV

`Import` uploads the file as multipart form data. The file is checked before it
is queued, and the result carries that validation report next to the job.
`File`, `Filename`, `FieldMapping` (CSV header → contact field) and at least one
label are required.

```go
f, err := os.Open("contacts.csv")
defer f.Close()

result, err := client.Contacts.Import(&openapi.ImportContactsRequest{
    File:         f,
    Filename:     "contacts.csv",
    FieldMapping: map[string]string{"Email Address": "email", "Name": "firstname"},
    Labels:       []string{"newsletter"},
    ImportMode:   openapi.ImportModeSkipExisting, // or ImportModeReplaceAll, ImportModeFillEmpty
})
fmt.Println(result.Import.ID, result.Import.State, result.Validation.EstimatedValid)

jobs, err := client.Contacts.Imports(&openapi.ListContactImportsParams{State: "RUNNING"})
queue, err := client.Contacts.CurrentImport()      // queue.StillImporting
job, err := client.Contacts.GetImport(result.Import.ID)
cancelled, err := client.Contacts.CancelImport(result.Import.ID) // best effort while RUNNING
```

### Segments

Segments are saved, named contact filters. `Filters` uses the same structure as the filters body of
`Contacts.Count`: each condition is `{"field", "type", "operator", "values"}`, with uppercase `type`
(`TEXT`, `NUMBER`, `DATE`, `LIST`) and uppercase `operator` (`MATCH`, `CONTAIN`, `NOT_CONTAIN`,
`CONTAIN_ANY`, `CONTAIN_ALL`, `NOT_CONTAIN_ANY`, `IS_EMPTY`, `GREATER_THAN`, `LESS_THAN`); `values` is
always a slice. Conditions on different fields are ANDed together. The `IS_NOT_EMPTY` operator
currently 502s on every field — avoid it until it's fixed upstream.

```go
page, err := client.Segments.List(&openapi.ListSegmentsParams{Search: "vip", Sort: "-count"})
segment, err := client.Segments.Get("segment-id")

// Filters is required; pass an empty slice to match every contact.
segment, err = client.Segments.Create(&openapi.CreateSegmentRequest{
    Name: "VIP buyers",
    Filters: []map[string]any{
        {"field": "labels", "type": "LIST", "operator": "CONTAIN_ANY", "values": []string{"vip", "cold-leads"}},
        {"field": "firstname", "type": "TEXT", "operator": "CONTAIN", "values": []string{""}},
    },
})

// Omitted fields keep their current value.
segment, err = client.Segments.Update("segment-id", &openapi.UpdateSegmentRequest{Name: "Top buyers"})
err = client.Segments.Delete("segment-id")

// Live count of the active contacts a campaign to this segment would reach.
count, err := client.Segments.RecipientsCount("segment-id")
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

### API Keys

Manages the Mailtarget API keys used for sending. These calls authenticate with the Open API secret
key, while the keys they return are the sending credentials.

```go
page, err := client.APIKeys.List(&openapi.ListAPIKeysParams{Page: 1, PerPage: 10})
key, err := client.APIKeys.Get(7)
key, err = client.APIKeys.Create(&openapi.CreateAPIKeyRequest{
    Name: "ci", PermissionIDs: []int{1, 2},
})
// PermissionIDs is required and must be non-empty on create and update.
key, err = client.APIKeys.Update(7, &openapi.UpdateAPIKeyRequest{
    Name: "renamed", PermissionIDs: []int{1, 2},
})
err = client.APIKeys.Delete(7)
```

### Campaigns

```go
page, err := client.Campaigns.List(&openapi.ListCampaignsParams{Search: "Promo"})
campaign, err := client.Campaigns.Get("campaign-id")

campaign, err = client.Campaigns.Create(&openapi.CampaignRequest{
    Subject:    "Monthly newsletter",
    Sender:     &openapi.CampaignSender{Email: "no-reply@example.com"},
    Recipients: &openapi.CampaignRecipients{Labels: []string{"vip"}},
})
campaign, err = client.Campaigns.Update("campaign-id", &openapi.CampaignRequest{Subject: "Updated"})
err = client.Campaigns.Delete("campaign-id")

stats, err := client.Campaigns.Analytics("campaign-id")

// Status is required: delivered, opened, clicked, bounced or complained.
recipients, err := client.Campaigns.Recipients("campaign-id", &openapi.ListCampaignRecipientsParams{
    Status: "opened", Page: 1,
})

err = client.Campaigns.Send("campaign-id")
err = client.Campaigns.SendTest("campaign-id", "qa@example.com")

// Schedule instead of sending now; the due date is required.
campaign, err = client.Campaigns.SetSchedule("campaign-id", "2026-10-01 09:00")
err = client.Campaigns.CancelSchedule("campaign-id")
```

### Templates

```go
page, err := client.Templates.List(&openapi.ListTemplatesParams{Search: "Welcome", EmailType: "regular"})
template, err := client.Templates.Get("template-id") // includes Content, Body and CSS
```

### Senders and sending domains

```go
// Senders is not paginated; it returns every identity.
senders, err := client.Senders.List()
sender, err := client.Senders.Get("sender-id")
sender, err = client.Senders.Create(&openapi.SenderRequest{Email: "no-reply@example.com", Name: "No reply"})
sender, err = client.Senders.Update("sender-id", &openapi.SenderRequest{Name: "Renamed"})
err = client.Senders.Delete("sender-id")

status, err := client.Senders.CheckDomain("no-reply@example.com")

domains, err := client.SendingDomains.List(nil)
domain, err := client.SendingDomains.Get(5)
domain, err = client.SendingDomains.VerifyTXT(5)
```

### Labels

```go
page, err := client.Labels.List(&openapi.ListLabelsParams{Search: "vip"})
label, err := client.Labels.Create("newsletter")
label, err = client.Labels.Rename("newsletter", "monthly-newsletter")

// The endpoint takes the names as one comma separated segment.
err = client.Labels.Delete("vip", "cold-lead")
```

### Suppressions

Suppressed addresses no longer receive email.

```go
page, err := client.Suppressions.List(&openapi.ListSuppressionsParams{
    Source: "List Unsubscribe,Link Unsubscribe", // comma separated to match any of several
})
bounces, err := client.Suppressions.Bounces(nil)
unsubscribes, err := client.Suppressions.Unsubscribes(nil)

// One address can be suppressed per sub-account plus account-wide, so this is a list.
found, err := client.Suppressions.Lookup(&openapi.LookupSuppressionsParams{Email: "recipient@example.com"})

// Email, Type and Source are required; Type and Source are free-form (max 255
// characters). Only Source changes where the row shows up: "Bounce Rule" puts it
// under Bounces, "List Unsubscribe" or "Link Unsubscribe" under Unsubscribes.
// Leave SubAccountID zero to suppress account-wide.
sup, err := client.Suppressions.Create(&openapi.CreateSuppressionRequest{
    Email: "recipient@example.com", Type: "Non-transactional", Source: "Manual",
})
err = client.Suppressions.Delete(sup.ID)
```

### Webhooks

```go
// The catalog of event types; their IDs go into EventIDs.
events, err := client.WebhookEvents.List(nil)
event, err := client.WebhookEvents.Get(1)

page, err := client.Webhooks.List(&openapi.ListWebhooksParams{Search: "crm"})
hook, err := client.Webhooks.Get(5)

hook, err = client.Webhooks.Create(&openapi.CreateWebhookRequest{
    Name:             "crm",
    TargetURL:        "https://crm.example.com/hook",
    EventIDs:         []int{1, 2},
    AuthenticationID: openapi.WebhookAuthBasic,
    Username:         "bot",
    Password:         "secret",
})

// Update replaces the whole configuration. Leave Password or IsActive nil to keep them.
active := false
hook, err = client.Webhooks.Update(5, &openapi.UpdateWebhookRequest{
    Name: "crm", TargetURL: "https://crm.example.com/hook", EventIDs: []int{1}, IsActive: &active,
})
err = client.Webhooks.Delete(5)
```

### Settings, usage and sub-accounts

```go
company, err := client.Settings.Company()
profile, err := client.Settings.Profile()

// Quota fields are 0 on plans without a fixed monthly quota; that does not block sending.
usage, err := client.Usage.Get()

// Note: this endpoint spells the page size "size", not "perPage".
page, err := client.SubAccounts.List(&openapi.ListSubAccountsParams{Page: 1, Size: 25})
sub, err := client.SubAccounts.Get(3)
sub, err = client.SubAccounts.Create(&openapi.CreateSubAccountRequest{Name: "sales"})
sub, err = client.SubAccounts.Update(3, &openapi.UpdateSubAccountRequest{Status: "Suspended"})
```

### Transmissions (Open API)

The Open API has its own send endpoint. **Prefer `client.Send` / `Layang.Send`** for sending email:
it needs only the Mailtarget API key, while this one also requires the Open API secret key. This is
here so the SDK covers the whole Open API surface.

The endpoint carries the Mailtarget API key in its body rather than the header; the SDK fills that
in from the key the client was built with, so you never pass it twice.

```go
result, err := client.Transmissions.Send(&openapi.SendEmailRequest{
    From:    &openapi.Address{Email: "no-reply@example.com", Name: "No reply"},
    Subject: "Hello",
    To:      []openapi.Address{{Email: "recipient@example.com"}},
    BodyHTML: "<p>Hello</p>",
})
fmt.Println(result.TransmissionID)
```

### Available resources

All Open API resources are implemented: **Contacts**, **Segments**, **Analytics**, **API Keys**,
**Campaigns**, **Senders**, **Sending Domains**, **Labels**, **Suppressions**, **Webhooks**,
**Webhook Events**, **Settings**, **Usage**, **Sub Accounts**, **Templates** and **Transmissions**.

### Adding a new Open API resource

The layout is designed so a new resource stays a local change:

1. Add `openapi/<resource>.go` with its models and a `<Resource>Service` whose methods call the
   shared helpers (`object`, `objects`, `paged`, `bare`) or `c.do` directly.
2. Add one field for the service on `openapi.Client` and wire it in `New`.
3. Add `openapi/<resource>_test.go`.

Nothing in the root package changes, and the secret key guard applies automatically because every
call goes through `Client.do`.
