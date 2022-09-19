package layang

import "github.com/google/uuid"

var privateAPIKey = "j7dlt0ojjczewvvkll70"
var sender = Address{
	Email: "sender@example.com",
	Name:  "Sender",
}
var wrongSender = Address{
	Email: "wrongsender.com",
	Name:  "Sender",
}
var subject = "Fancy subject!"
var body = "Hello from Layang Go!"
var recipient = Recipient{Address: []Address{
	{
		Email: "recipient@example.com",
		Name:  "Recipient",
	},
	{
		Email: "recipient2@example.com",
		Name:  "Recipient 2",
	},
}}
var wrongRecipient = Recipient{Address: []Address{
	{
		Email: "wrongrecipient.com",
		Name:  "Recipient",
	},
	{
		Email: "recipient2@example.com",
		Name:  "Recipient 2",
	},
}}
var html = `
<p>My fantastic HTML content.<br><br><b>SparkPost</b> <img src=\"cid:AnImage.png\"></p>
`
var wrongHtml = `not html`
var successResponse = SuccessResponse{
	Message: "Queued. Thank you.",
	Id:      uuid.NewString(),
}
var errorResponse = ErrorResponse{
	Error:  "Bad request",
	Status: 400,
}
var attachment = Attachment{
	Type:     "image/png",
	Filename: "AnImage.png",
	Content:  "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAAXxJREFUOBFjvJVg84P5718WBjLAX2bmPyxMf/+xMDH8YyZDPwPDXwYGJkIaOXTNGdiUtHAqI2jA/18/GUQzGsg3gMfKg4FVQo6BiYcPqyF4XcChaczA4+DP8P//f4b/P3+SZgAzvxCDSGYjAyMjI8PvZw+AoYXdLuyiQLtE0uoZWAREwLb+fnKXQTipkngXcJu7MnACQx8G2FX1GHgs3bDGBlYX8HlFM/z9+JbhzewWhmf1CQyfti9j+PfzBwO/ZxTMTDiNmQKBfmZX1GB42V/K8P38YbDCX/dvMDAwMzPwuYbBNcIYmC4AhfjvXwx/376AqQHTf96+ZPj34xuKGIiDaQBQ8PPBTQwCoZkMjJzcYA3MgqIMAr7xDJ/3rAHzkQnGO7FWf5gZ/qLmBSZmBoHgNAZee1+Gf18/MzCyczJ83LyQ4fPetch6Gf4xMP3FbgBMGdAgJqAr/n37zABMTTBROA0ygAWUJUG5Civ4B8xwX78CpbD6FJiHmf4AAFicbTMTr5jAAAAAAElFTkSuQmCC",
}
var attachments = []Attachment{
	attachment,
}
