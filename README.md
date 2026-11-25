# Mailtarget Go SDK

The official Go SDK for [Mailtarget](https://mailtarget.co). It covers two things:

- **Sending email** through Layang, Mailtarget's transactional sending API.
- **Managing your account** through the Mailtarget Open API: contacts, segments, campaigns,
  templates, forms, landing pages, the file manager, senders and domains, suppressions, webhooks,
  analytics, sub-accounts, API keys and usage.

Each capability uses its own credential, and you can use sending on its own.

## Contents

- [Installation](#installation)
- [Credentials](#credentials)
- [Quickstart](#quickstart)
- [Sending email](#sending-email)
- [Using the Open API](#using-the-open-api)
- [Guides](#guides)
- [Resource reference](#resource-reference)
- [Known limitations](#known-limitations)
- [License](#license)

## Installation

Requires Go 1.18 or later.

```sh
go get github.com/mailtarget/mailtarget-go-sdk
```

The root package is named `layang`, and the Open API lives in the `openapi` subpackage:

```go
import (
    "github.com/mailtarget/mailtarget-go-sdk"         // package layang
    "github.com/mailtarget/mailtarget-go-sdk/openapi"
)
```

## Credentials

| | Mailtarget API key | Open API secret key |
|---|---|---|
| Required | Yes | Only for Open API calls |
| Where to get it | Generate it yourself at [app.mailtarget.co/api-key](https://app.mailtarget.co/api-key) | Issued on request; contact your Mailtarget account manager |
| Used for | Sending email | Everything under `openapi` |

Keep both keys out of source control, for example by reading them from environment variables.

## Quickstart

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mailtarget/mailtarget-go-sdk"
)

func main() {
	client := layang.NewMailtargetClient(os.Getenv("MAILTARGET_API_KEY"))

	message := client.Layang.NewMessage(
		"Welcome aboard",                         // subject
		"Thanks for signing up.",                 // plain text body
		"<p>Thanks for <b>signing up</b>.</p>",   // HTML body
		layang.Address{Email: "hello@example.com", Name: "Example"},
		[]layang.Address{{Email: "recipient@example.com", Name: "Recipient"}},
	)

	success, apiErr, err := client.Send(message)
	if err != nil {
		if apiErr != nil {
			log.Fatalf("rejected: status=%d code=%s message=%s", apiErr.Status, apiErr.Code, apiErr.Message)
		}
		log.Fatal(err)
	}
	fmt.Println("transmission id:", success.TransmissionId)
}
```

## Sending email

`layang.NewMailtargetClient(apiKey)` is the entry point. It only needs the Mailtarget API key, and
sending never goes through the Open API, so it works whether or not you have a secret key.

If you only ever send email, `layang.NewLayang(apiKey)` gives you the sender on its own. It has the
same `NewMessage` and `Send` methods.

### Building a message

`NewMessage` takes the subject, text body, HTML body, sender and recipients. You can set the rest
with setters, or fill in the `layang.Message` struct yourself.

```go
message := client.Layang.NewMessage(subject, text, html, from, to)

message.SetCC([]layang.Address{{Email: "cc@example.com"}})
message.SetBCC([]layang.Address{{Email: "audit@example.com"}})
message.SetReplyTo([]layang.Address{{Email: "support@example.com", Name: "Support"}})
message.SetHeaders([]layang.Header{{Name: "X-Campaign", Value: "welcome"}})
message.SetMetadata(map[string]string{"userId": "42"})
message.SetOptionsAttributes(layang.OptionsAttributes{ClickTracking: true, OpenTracking: true})
```

### Attachments and inline images

`Value` is the file content encoded in base64. The example below references the attachment's
filename from the HTML body with `cid:` to show it inline:

```go
logo, err := os.ReadFile("logo.png")
if err != nil {
    log.Fatal(err)
}

message := client.Layang.NewMessage(
    "Your invoice",
    "Your invoice is attached.",
    `<p>Your invoice is attached.</p><img src="cid:logo.png">`,
    layang.Address{Email: "billing@example.com"},
    []layang.Address{{Email: "recipient@example.com"}},
)
message.SetAttachment([]layang.Attachment{{
    MimeType: "image/png",
    Filename: "logo.png",
    Value:    base64.StdEncoding.EncodeToString(logo),
}})
```

### Validation

`Send` checks the message before any request is made, and fails with an error if any of these is
not met:

- the client has an API key;
- `From.Email` is a valid address;
- at least one `To` address is given, and every `To`, `Cc`, `Bcc` and `ReplyTo` address is valid;
- a text body, an HTML body or both are set;
- an HTML body contains at least one HTML tag.

The subject is not checked locally.

### Handling the response

`Send` returns three values:

| Outcome | `*SuccessResponse` | `*ErrorResponse` | `error` |
|---|---|---|---|
| Accepted | has `TransmissionId` | `nil` | `nil` |
| Rejected by the API | `nil` | has `Status`, `Code`, `Message` | raw response body |
| Failed local validation | empty struct | empty struct | the reason |

Always check `err` first, and read the `*ErrorResponse` for details only when it is non-nil and
`Status` is set.

### Templates and substitution data

The Layang sender takes the message body directly. To send from a stored template with
substitution data, use the Open API's own send endpoint, `client.Transmissions.Send`
([see below](#transmissions)). It needs the Open API secret key as well.

## Using the Open API

### Enabling it

Pass the secret key when building the client. Open API resources are then available directly on
`client`:

```go
client := layang.NewMailtargetClient(
    os.Getenv("MAILTARGET_API_KEY"),
    layang.WithOpenAPISecretKey(os.Getenv("MAILTARGET_OPENAPI_SECRET")),
)

page, err := client.Contacts.List(nil)
```

Without the secret key every Open API call fails immediately with an `*openapi.ConfigError` and no
request is sent. Use `client.HasOpenAPIAccess()` to skip Open API features instead, or
`client.SetOpenAPISecretKey(key)` to supply the key later.

If you never send email, use the subpackage on its own:

```go
oa := openapi.New(openapi.WithSecretKey(os.Getenv("MAILTARGET_OPENAPI_SECRET")))
page, err := oa.Contacts.List(nil)
```

### Scopes

Each secret key is granted scopes. A call outside them is refused with HTTP 403.

| Resource | Read scope | Write scope |
|---|---|---|
| Contacts, contact imports, Segments | `contacts.read` | `contacts.write` |
| Forms | `forms.read` | `forms.write` |
| Pages | `pages.read` | `pages.write` |
| Files (file manager and its labels) | `files.read` | `files.write` |
| Campaigns | `campaigns.read` | `campaigns.write` |
| Templates | `campaigns.read` | — |
| Analytics | `analytics.read` | — |
| API Keys, Sub Accounts | `settings.read` | `settings.write` |
| Usage | `settings.read` | — |
| Suppressions | `suppressions.read` | `suppressions.write` |
| Webhooks, Webhook Events | `webhooks.read` | `webhooks.write` |
| Transmissions, `KeyInfo.CheckTransmissionKey` | — | `transmissions.write` |

The API documents no scope for Senders, Sending Domains, Labels and Settings. `KeyInfo.Get` needs
no scope.

### Errors

Open API calls return one of three kinds of error:

```go
_, err := client.Campaigns.Get("campaign-id")

var configErr *openapi.ConfigError
var apiErr *openapi.Error
switch {
case err == nil:
    // success
case errors.As(err, &configErr):
    // no secret key configured; nothing was sent
case errors.As(err, &apiErr):
    // the API answered with a non-2xx status
    log.Printf("status=%d code=%s message=%s", apiErr.StatusCode, apiErr.Code, apiErr.Message)
    // apiErr.Details holds the raw `details` JSON when the API sent one.
default:
    // a required field was missing (checked locally, nothing was sent) or a network error
    log.Print(err)
}
```

Calls whose response is `{"data": true|false}` (`APIKeys.Delete`, `Campaigns.Send`,
`Campaigns.SendTest`, `Senders.Delete`, `Forms.Delete`, `Forms.Submit`, `Pages.Delete`,
`Files.Delete`, `Files.DeleteLabel`, `Labels.Delete`, `Segments.Delete`, `Suppressions.Delete`, `Webhooks.Delete`) treat `data: false` as a failure: they
return an `*openapi.Error` with `Code` `openapi.ErrCodeOperationFailed` (`"OPERATION_FAILED"`) and
the HTTP status received, so a `nil` error always means the operation was carried out.

`apiErr.Details` is a `json.RawMessage` with the `details` member some errors carry, for example the
active imports of a 409 `IMPORT_IN_PROGRESS` or the report of a 422 `IMPORT_VALIDATION_FAILED`.
Decode it with `json.Unmarshal`; it is `nil` when absent.

Required fields are checked locally before any request. For example, `Contacts.Update` without
`Email`, `Firstname` and `Labels` returns an error straight away.

### Timeouts

Open API requests time out after 30 seconds. Change it with
`layang.WithOpenAPITimeout(60 * time.Second)`.

### Pagination

List methods return a `*openapi.Page[T]` holding `Items` and `Meta` (`Page`, `PerPage`, `Total`).
To walk every page:

```go
const perPage = 100
for p := 1; ; p++ {
    page, err := client.Contacts.List(&openapi.ListContactsParams{Page: p, PerPage: perPage})
    if err != nil {
        log.Fatal(err)
    }
    for _, contact := range page.Items {
        fmt.Println(contact.Email)
    }
    if len(page.Items) < perPage || p*perPage >= page.Meta.Total {
        break
    }
}
```

`SubAccounts.List` and `WebhookEvents.List` name the page size `Size` instead of `PerPage`.
`Senders.List` and `WebhookEvents.List` return a plain slice with no metadata.

## Guides

These cover the common end-to-end tasks, in the order the steps must happen.

### Build your audience from a CSV

Imports only accept labels that already exist, and map each CSV header to a contact field name.

```go
// 1. Make sure the label exists.
_, err := client.Labels.Create("newsletter")

// 2. Look up the field names the mapping can use.
fields, err := client.Contacts.Fields()
for _, f := range fields {
    fmt.Println(f.FieldName, f.DisplayType)
}

// 3. Upload the file.
f, err := os.Open("contacts.csv")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

result, err := client.Contacts.Import(&openapi.ImportContactsRequest{
    File:         f,
    Filename:     "contacts.csv",
    FieldMapping: map[string]string{"Email Address": "email", "First Name": "firstname"},
    Labels:       []string{"newsletter"},
    ImportMode:   openapi.ImportModeSkipExisting, // or ImportModeReplaceAll (default), ImportModeFillEmpty
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("rows expected to import:", result.Validation.EstimatedValid)

// 4. The import runs in the background; poll until it ends.
for {
    job, err := client.Contacts.GetImport(result.Import.ID)
    if err != nil {
        log.Fatal(err)
    }
    if job.State == "FINISH" || job.State == "INTERUPTED" {
        fmt.Println(job.State, job.NewContact, "new,", job.UpdatedContact, "updated", job.Message)
        break
    }
    time.Sleep(5 * time.Second)
}
```

Rules for the file and the import:

- UTF-8, `,` or `;` delimited, at most 5 MiB and 20,000 contacts, with an email column.
- Columns missing from `FieldMapping` are skipped, except `label` and `name`.
- Only one import can run per account, including imports started from the dashboard. Another one
  is refused with HTTP 409 `IMPORT_IN_PROGRESS`. Check `Contacts.CurrentImport()` first, or cancel
  a job with `Contacts.CancelImport(id)`.
- Problem rows are listed in `result.Validation`. With `Strict: true`, any problem row other than a
  `REPLACE_VALUE` replacement rejects the whole file with HTTP 422 `IMPORT_VALIDATION_FAILED`; the
  report is in `apiErr.Details`.
- A 409 `IMPORT_IN_PROGRESS` lists the active imports in `apiErr.Details`.

Once the contacts are in, group them with [segments](#segments).

### Verify a sender

A campaign can only be sent from a verified, permitted sender.

```go
sender, err := client.Senders.Create(&openapi.SenderRequest{Email: "news@example.com", Name: "Example News"})

status, err := client.Senders.CheckDomain("news@example.com")
fmt.Println("SPF", status.SPF, "DKIM", status.DKIM, "DMARC", status.DMARC, "permitted", status.Permitted)

// After publishing the DNS records, ask Mailtarget to check them again.
domains, err := client.SendingDomains.List(nil)
domain, err := client.SendingDomains.VerifyTXT(domains.Items[0].ID)
fmt.Println(sender.ID, domain.Status)
```

### Run a campaign

```go
// 1. Optionally start from a template. Its content and email type replace
//    HTMLContent and EmailType on the campaign.
templates, err := client.Templates.List(&openapi.ListTemplatesParams{Search: "Newsletter"})

// 2. Create the campaign. Recipients are chosen by label.
campaign, err := client.Campaigns.Create(&openapi.CampaignRequest{
    Subject:    "October newsletter",
    TemplateID: templates.Items[0].ID,
    Sender:     &openapi.CampaignSender{Email: "news@example.com", Name: "Example News"},
    Recipients: &openapi.CampaignRecipients{Labels: []string{"newsletter"}},
})

// 3. Send yourself a preview.
err = client.Campaigns.SendTest(campaign.ID, "qa@example.com")

// 4a. Send now.
err = client.Campaigns.Send(campaign.ID)

// 4b. Or schedule it: SetSchedule only saves the send time (the campaign stays DRAFT),
// Send then activates the schedule (READY, sent automatically at the due date).
// The due date is an RFC 3339 timestamp in the future.
campaign, err = client.Campaigns.SetSchedule(campaign.ID, time.Now().Add(24*time.Hour).Format(time.RFC3339))
err = client.Campaigns.Send(campaign.ID)

// 5. Follow the results.
stats, err := client.Campaigns.Analytics(campaign.ID)
opened, err := client.Campaigns.Recipients(campaign.ID, &openapi.ListCampaignRecipientsParams{Status: "opened"})
fmt.Println(stats.OpenRate, opened.Meta.Total)
```

### Keep your list clean

Suppressed addresses no longer receive email. Check before a one-off send, and subscribe a
webhook to learn about bounces and unsubscribes as they happen.

```go
found, err := client.Suppressions.Lookup(&openapi.LookupSuppressionsParams{Email: "recipient@example.com"})
if err == nil && len(found.Items) > 0 {
    fmt.Println("suppressed:", found.Items[0].Source)
}

// Event IDs come from the webhook event catalog.
events, err := client.WebhookEvents.List(nil)
ids := make([]int, 0, len(events))
for _, e := range events {
    ids = append(ids, e.ID)
}

hook, err := client.Webhooks.Create(&openapi.CreateWebhookRequest{
    Name:      "crm-sync",
    TargetURL: "https://crm.example.com/mailtarget-events",
    EventIDs:  ids,
})
fmt.Println(hook.ID)
```

### Provision a sub-account

Give a client or brand its own sub-account and sending key:

```go
sub, err := client.SubAccounts.Create(&openapi.CreateSubAccountRequest{Name: "Client A"})

// PermissionIDs is required; the Permissions of an existing key (APIKeys.Get)
// show which IDs are available.
key, err := client.APIKeys.Create(&openapi.CreateAPIKeyRequest{
    Name:          "client-a-sending",
    SubAccountID:  sub.ID,
    PermissionIDs: []int{1},
})

// The new key sends email for that sub-account.
subClient := layang.NewMailtargetClient(key.Key)
fmt.Println(subClient.HasOpenAPIAccess())
```

Check the sending volume before a large send with `client.Usage.Get()`. On plans without a fixed
monthly quota, the quota fields read 0; that does not block sending.

## Resource reference

Every method below is on `client` (or on a standalone `openapi.Client`).

### Contacts

```go
page, err := client.Contacts.List(&openapi.ListContactsParams{Page: 1, PerPage: 20, Search: "@example.com"})
contact, err := client.Contacts.Get("contact-id")
contact, err = client.Contacts.GetByEmail("recipient@example.com")

contact, err = client.Contacts.Create(&openapi.CreateContactRequest{
    Email:     "recipient@example.com",
    Firstname: "Sarah",
    Labels:    []string{"vip"},
    BirthDate: "1990-05-17", // YYYY-MM-DD
})

// Email, Firstname and Labels are required on update.
contact, err = client.Contacts.Update("contact-id", &openapi.UpdateContactRequest{
    Email:     "recipient@example.com",
    Firstname: "Sarah",
    Labels:    []string{"vip"},
    Note:      "renewed",
})

err = client.Contacts.Delete("contact-id")

// Count takes the same filters as segments; nil counts every contact.
total, err := client.Contacts.Count(map[string]any{
    "filters": []map[string]any{{"field": "labels", "type": "LIST", "operator": "CONTAIN_ANY", "values": []string{"vip"}}},
})

// Export returns one page as CSV. Paging past the last contact returns an
// *openapi.Error with StatusCode 404, which marks the end of the export.
csv, err := client.Contacts.Export(&openapi.ListContactsParams{Page: 1, PerPage: 100})

fields, err := client.Contacts.Fields()
```

Birth dates: send `BirthDate` (`YYYY-MM-DD`) on create and update. `DayOfBirth` on the request types
is deprecated: it is still accepted as `YYYY-MM-DD` for older clients, and `BirthDate` wins when both
are set. On a returned `Contact`, `BirthDate` is `YYYY-MM-DD` and `DayOfBirth` is a read-only `MM-DD`
value; do not copy it into a request.

Contact imports:

```go
result, err := client.Contacts.Import(&openapi.ImportContactsRequest{ /* see the guide above */ })
jobs, err := client.Contacts.Imports(&openapi.ListContactImportsParams{State: "RUNNING"})
queue, err := client.Contacts.CurrentImport() // queue.StillImporting
job, err := client.Contacts.GetImport("import-id")
cancelled, err := client.Contacts.CancelImport("import-id") // best effort while RUNNING
```

Import states are `READY`, `RUNNING`, `FINISH` and `INTERUPTED` (the API's spelling). A failed
import carries the reason in `Message`.

### Segments

Segments are saved, named contact filters. Each condition in `Filters` is
`{"field", "type", "operator", "values"}`:

- `type` is uppercase: `TEXT`, `NUMBER`, `DATE` or `LIST`.
- `operator` is uppercase: `MATCH`, `CONTAIN`, `NOT_CONTAIN`, `CONTAIN_ANY`, `CONTAIN_ALL`,
  `NOT_CONTAIN_ANY`, `IS_EMPTY`, `GREATER_THAN` or `LESS_THAN`.
- `values` is always a slice.

Conditions on different fields must all match.

```go
page, err := client.Segments.List(&openapi.ListSegmentsParams{Search: "vip", Sort: "-count"})
segment, err := client.Segments.Get("segment-id")

// Filters is required; pass an empty slice to match every contact.
segment, err = client.Segments.Create(&openapi.CreateSegmentRequest{
    Name: "VIP buyers",
    Filters: []map[string]any{
        {"field": "labels", "type": "LIST", "operator": "CONTAIN_ANY", "values": []string{"vip", "cold-leads"}},
        {"field": "firstname", "type": "TEXT", "operator": "CONTAIN", "values": []string{"a"}},
    },
})

// Omitted fields keep their current value.
segment, err = client.Segments.Update("segment-id", &openapi.UpdateSegmentRequest{Name: "Top buyers"})
err = client.Segments.Delete("segment-id")

// Live count of the active contacts a campaign to this segment would reach.
// Segment.Count is a cached value and may be stale.
count, err := client.Segments.RecipientsCount("segment-id")
```

### Labels

```go
page, err := client.Labels.List(&openapi.ListLabelsParams{Search: "vip"})
label, err := client.Labels.Create("newsletter")
label, err = client.Labels.Rename("newsletter", "monthly-newsletter")
err = client.Labels.Delete("vip", "cold-lead")
```

### Campaigns

```go
page, err := client.Campaigns.List(&openapi.ListCampaignsParams{Search: "Promo"})
campaign, err := client.Campaigns.Get("campaign-id")

campaign, err = client.Campaigns.Create(&openapi.CampaignRequest{
    Subject:    "Monthly newsletter",
    Sender:     &openapi.CampaignSender{Email: "news@example.com"},
    Recipients: &openapi.CampaignRecipients{Labels: []string{"vip"}},
})
campaign, err = client.Campaigns.Update("campaign-id", &openapi.CampaignRequest{Subject: "Updated"})
err = client.Campaigns.Delete("campaign-id")

err = client.Campaigns.SendTest("campaign-id", "qa@example.com")
err = client.Campaigns.Send("campaign-id") // sends now, or activates a saved schedule

// Scheduling takes two steps. SetSchedule alone sends nothing: the campaign stays DRAFT.
campaign, err = client.Campaigns.SetSchedule("campaign-id", "2026-10-01T09:00:00+07:00")
err = client.Campaigns.Send("campaign-id") // now READY, sent automatically at dueDate

// Only for an activated (READY) schedule: back to DRAFT, immediate send time, no dueDate.
// To move a schedule that is not activated yet, call SetSchedule again instead.
err = client.Campaigns.CancelSchedule("campaign-id")

stats, err := client.Campaigns.Analytics("campaign-id")
// Status is required: delivered, opened, clicked, bounced or complained.
recipients, err := client.Campaigns.Recipients("campaign-id", &openapi.ListCampaignRecipientsParams{
    Status: "opened", Page: 1,
})
```

### Templates

```go
// Type "team" (the default) lists your company's templates; "default" lists
// Mailtarget's published ones.
page, err := client.Templates.List(&openapi.ListTemplatesParams{Search: "Welcome", Type: "default"})
template, err := client.Templates.Get("template-id") // includes Content, Body and CSS
```

### Senders and sending domains

```go
senders, err := client.Senders.List() // not paginated
sender, err := client.Senders.Get("sender-id")
sender, err = client.Senders.Create(&openapi.SenderRequest{Email: "no-reply@example.com", Name: "No reply"})
sender, err = client.Senders.Update("sender-id", &openapi.SenderRequest{Name: "Renamed"})
err = client.Senders.Delete("sender-id")
status, err := client.Senders.CheckDomain("no-reply@example.com")

domains, err := client.SendingDomains.List(nil)
domain, err := client.SendingDomains.Get(5)
domain, err = client.SendingDomains.VerifyTXT(5)
```

### Suppressions

Email, Type and Source are required to create a suppression. Neither is checked against a fixed
enum: Type is a free-form label of your own choosing (max 255 characters), and so is Source,
*except* for the values below, which decide where the row also gets listed.

| `Source` | Also listed under |
|---|---|
| `"Bounce Rule"` | `Bounces` |
| `"List Unsubscribe"` | `Unsubscribes` |
| `"Link Unsubscribe"` | `Unsubscribes` |
| anything else, e.g. `"Manual"` | `List` only |

Source is matched exact and case-sensitive.

```go
page, err := client.Suppressions.List(&openapi.ListSuppressionsParams{
    Source: "List Unsubscribe,Link Unsubscribe", // comma separated to match any of several
})
bounces, err := client.Suppressions.Bounces(nil)
unsubscribes, err := client.Suppressions.Unsubscribes(nil)

// One address can be suppressed per sub-account plus account-wide, so this is a list.
found, err := client.Suppressions.Lookup(&openapi.LookupSuppressionsParams{Email: "recipient@example.com"})

// Leave SubAccountID zero to suppress account-wide.
sup, err := client.Suppressions.Create(&openapi.CreateSuppressionRequest{
    Email: "recipient@example.com", Type: "Manual block", Source: "Manual",
})
err = client.Suppressions.Delete(sup.ID)
```

### Webhooks

```go
events, err := client.WebhookEvents.List(nil) // the catalog of event types
event, err := client.WebhookEvents.Get(1)

page, err := client.Webhooks.List(&openapi.ListWebhooksParams{Search: "crm"})
hook, err := client.Webhooks.Get(5)

// TargetURL must be a public http(s) URL. With WebhookAuthBasic set Username
// and Password; with WebhookAuthOAuth set ClientID, ClientSecret and TokenURL.
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

### Forms

```go
page, err := client.Forms.List(&openapi.ListFormsParams{
    Status: []string{openapi.FormStatusPublished, openapi.FormStatusDraft},
    Labels: []string{"vip"},
})
form, err := client.Forms.Get("form-id") // form.Fields, form.PublicURL, form.EmbedHTML, form.UsedIn

copied, err := client.Forms.Duplicate("form-id") // unpublished "<name> (Duplicate)"
form, err = client.Forms.Publish("form-id")
form, err = client.Forms.Unpublish("form-id")
err = client.Forms.Delete("form-id") // unpublish first; forms used by an automation cannot be deleted

// From and To are Unix milliseconds and must be set together.
subs, err := client.Forms.Submissions("form-id", &openapi.ListFormSubmissionsParams{
    From: 1767225600000, To: 1769903999000,
})

// Keys are form.Fields[i].Name; a MultiValue field takes []string of option names.
err = client.Forms.Submit("form-id", &openapi.FormSubmissionRequest{Values: map[string]any{
    "email":     "sarah@example.com",
    "firstname": "Sarah",
    "interest":  []string{"Promo", "Newsletter"},
}})
```

`Forms.Submit` is server-side only: it needs the Open API secret key, which must never be shipped to
a browser or mobile app. To collect responses from visitors, show the form itself with `PublicURL`
or `EmbedHTML`. Submissions are validated against the last published version of the form, and forms
with captcha or file upload fields are not supported (HTTP 422).

### Pages

Landing pages. The models are `openapi.LandingPage`, `LandingPageDetail`, `LandingPageForm` and
`LandingPageFormRef`, so they do not clash with the generic `openapi.Page[T]` returned by list calls.

```go
page, err := client.Pages.List(&openapi.ListPagesParams{
    Status:   []string{openapi.PageStatusPublished, openapi.PageStatusDefault},
    Category: openapi.PageCategoryLanding, // or openapi.PageCategoryThankYou
    Sort:     openapi.PageSortUpdatedAtDesc, // default: openapi.PageSortCreatedAtDesc
})
lp, err := client.Pages.Get("page-id") // lp.URL, lp.Permalink, lp.IsDefault, lp.Forms

// Embedded forms; each ID works with client.Forms (forms.read / forms.write). Not paginated.
forms, err := client.Pages.Forms("page-id")

copied, err := client.Pages.Duplicate("page-id") // unpublished "<name> (copy)", shares the forms
lp, err = client.Pages.Publish("page-id")        // a draft counts against the published-page quota (422)
lp, err = client.Pages.SetDefault("page-id")     // must be published (409); served at the domain root
lp, err = client.Pages.Unpublish("page-id")      // not the default page (409)
err = client.Pages.Delete("page-id")             // not the default page (409); embedded forms are kept
```

`URL` is filled for drafts too but only serves once the page is published; the default page is
served at the domain root.

### Files

The file manager: images, documents and media served from a public CDN URL, organised with their
own labels (separate from contact labels).

```go
page, err := client.Files.List(&openapi.ListFilesParams{
    Search: "banner",                                  // whole words of the name, not prefixes
    Type:   []string{openapi.FileTypeImage, openapi.FileTypeDocument},
    Labels: []string{"Promo"},                         // files carrying any of them
    Sort:   openapi.FileSortSizeDesc,                  // default: openapi.FileSortCreatedAtDesc
})

// Upload from disk (named after the base name) or from any io.Reader.
file, err := client.Files.UploadPath("./banner.png", "Promo")
file, err = client.Files.Upload(&openapi.UploadFileRequest{
    File: bytes.NewReader(data), Filename: "report.pdf", Labels: []string{"Reports"},
})
fmt.Println(file.URL) // public CDN URL; file.ThumbnailURL is set for images only

file, err = client.Files.Get("file-id")

// Rename and/or replace labels; nil fields are kept. The extension cannot change.
name, labels := "hero-banner", []string{} // an empty slice clears the labels
file, err = client.Files.Update("file-id", &openapi.UpdateFileRequest{Name: &name, Labels: &labels})

err = client.Files.Delete("file-id") // its URL stops working everywhere it is used

stats, err := client.Files.Stats() // StorageLimitBytes / StorageRemainingBytes are nil when unlimited

labelsList, err := client.Files.Labels()           // name + file count
label, err := client.Files.CreateLabel("my banner") // stored as "My Banner"
label, err = client.Files.RenameLabel("my banner", "Hero")
err = client.Files.DeleteLabel("Hero") // removed from every file; the files are kept
```

Allowed extensions: jpg, jpeg, jfif, png, gif, webp, tiff, avif, pdf, doc, docx, xls, xlsx, ppt,
pptx, csv, txt, md, mp3, mp4, mov, avi. The name without extension may only contain letters, digits,
spaces, `-` and `_`, and the content must match the extension. Uploads count against the plan's
storage quota (HTTP 422 when it is used up).

### Key info

```go
info, err := client.KeyInfo.Get() // info.CompanyID, info.Scopes, ...

// true only when the Transmission API key is valid and belongs to the same company.
// An invalid key and another company's key both return false.
ok, err := client.KeyInfo.CheckTransmissionKey(apiKey)
```

### Analytics

Dates are `YYYY-MM-DD`, and both ends are required. `GroupBy` is used by `SummaryBreakdown` only;
`Summary` never sends it.

```go
summary, err := client.Analytics.Summary(&openapi.AnalyticsSummaryParams{From: "2026-08-01", To: "2026-08-31"})
breakdown, err := client.Analytics.SummaryBreakdown(&openapi.AnalyticsSummaryParams{
    From: "2026-08-01", To: "2026-08-31", GroupBy: "sender",
})
detail, err := client.Analytics.Transmission("transmission-id")
events, err := client.Analytics.TransmissionEvents("transmission-id")
```

### API keys, sub-accounts, settings and usage

API key calls authenticate with the Open API secret key; the keys they manage are sending keys.

```go
keys, err := client.APIKeys.List(&openapi.ListAPIKeysParams{Page: 1, PerPage: 10})
key, err := client.APIKeys.Get(7)
// PermissionIDs is required and must be non-empty on create and update.
key, err = client.APIKeys.Create(&openapi.CreateAPIKeyRequest{Name: "ci", PermissionIDs: []int{1, 2}})
key, err = client.APIKeys.Update(7, &openapi.UpdateAPIKeyRequest{Name: "renamed", PermissionIDs: []int{1, 2}})
err = client.APIKeys.Delete(7)

subs, err := client.SubAccounts.List(&openapi.ListSubAccountsParams{Page: 1, Size: 25})
sub, err := client.SubAccounts.Get(3) // includes its API keys and domains
sub, err = client.SubAccounts.Create(&openapi.CreateSubAccountRequest{Name: "sales"})
sub, err = client.SubAccounts.Update(3, &openapi.UpdateSubAccountRequest{Status: "Suspended"}) // or "Active"

company, err := client.Settings.Company()
profile, err := client.Settings.Profile()
usage, err := client.Usage.Get()
```

### Transmissions

The Open API's own send endpoint. It needs the Open API secret key as well as the API key, so
prefer `client.Send` unless you need a stored template or substitution data. The SDK fills in the
API key the client was built with.

```go
result, err := client.Transmissions.Send(&openapi.SendEmailRequest{
    From:             &openapi.Address{Email: "no-reply@example.com", Name: "No reply"},
    Subject:          "Your order",
    To:               []openapi.Address{{Email: "recipient@example.com"}},
    TemplateID:       "template-id",
    SubstitutionData: map[string]string{"orderId": "A-1001"},
})
fmt.Println(result.TransmissionID)
```

## Known limitations

- Label names containing a comma cannot be deleted with `Labels.Delete` or used in an import,
  because those endpoints take labels as a comma separated list.
- The `IS_NOT_EMPTY` filter operator currently fails with HTTP 502 on every field; avoid it.
- Campaign recipients can only be selected by label.
- `Contacts.Export` returns HTTP 404 past the last page, which is how you know the export is done.

## License

GNU General Public License v3.0. See [LICENSE](LICENSE).
