package layang

import "github.com/google/uuid"

var privateAPIKey = "wnbj8XdJeeZ7Hraks0bU6tcK"
var sender = Address{
	Email: "sender@mtarget.co",
	Name:  "Sender",
}
var wrongSender = Address{
	Email: "wrongsender.com",
	Name:  "Sender",
}
var subject = "Fancy subject!"
var body = "Hello from Layang Go!"

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

var wrongTo = []Address{
	{
		Email: "wrongrecipient.com",
		Name:  "To",
	},
	{
		Email: "wrongrecipient2.com",
		Name:  "To2",
	},
}

var html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta content="text/plain; charset=utf-8" http-equiv="Content-Type"/>
    <meta content="IE=edge,chrome=1" http-equiv="X-UA-Compatible"/>
    <meta content="320" name="MobileOptimized"/>
    <meta content="true" name="HandHeldFriendly"/>
    <meta content="width=device-width, initial-scale=1.0" name="viewport"/>
    <style type="text/css">
        table {
        	border-collapse:collapse !important
        }
        /*MOBILE*/
        @media only screen and (max-width: 480px){
            .tableFull, .tableHAlf {
              width: 100% !important;
              max-width: 100% !important;
            }
            .btn-wrapper{
              width: 100% !important;
              max-width: 100% !important;
            }

          .btn-wrapper  a {
            display: block !important;
          }
        }
        /*DESKTOP*/
        @media only screen and (min-width: 481px){
        	.tableFull {
                width:600px !important;
            }
            .tableHAlf {
                width:250px !important;
            }
        }


    </style>
    <!--[if mso]>
    <style type="text/css">
.tableFull {
            width:600px !important;
      }
.tableHAlf {
            width:250px !important;
        }


    </style>
    <![endif]-->
    <title>Reset your MTARGET password</title>
</head>
<body align="center" bgcolor="#ffffff" leftmargin="0" marginheight="0" marginwidth="0"
      style="height:100%; width:100%; min-width:100%; margin:0; padding:0; background-color:#ffffff;" topmargin="0">
<div style="display:none;font-size:1px;color:#ffffff;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;mso-hide:all;">
    Tidak usah panik, ini hanya lupa password. Santai saja.
</div>
<table border="0" cellpadding="0" cellspacing="0" width="100%">
    <tr>
        <td align="center" bgcolor="#f4f4f4">
            <table align="center" border="0" cellpadding="0" cellspacing="0" class="tableFull" width="600">
                <tr>
                    <td align="center"
                        style="padding-right:50px; padding-left:50px;padding-top:52px;padding-bottom:52px">
                        <img alt="Mail Target" border="0" height="50"
                             src="https://files.mailtarget.co/assets/lg_mailtarget.png"
                             style="display: block; border: 0px;"/>
                    </td>
                </tr>
                <tr>
                    <td align="left"
                        bgcolor="#ffffff"
                        style="padding-right:40px; padding-left:40px;padding-top:22px;padding-bottom:22px;border-radius:3px;">

                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td align="left">
                                    <h1 style="font-size:21px; color:#4A4A4A;font-family: 'Roboto', sans-serif;">
                                        Halo, {fullname}
                                    </h1>
                                </td>
                            </tr>
                            <tr>
                                <td align="left">
                                    <p style="font-size:14px; color:#4A4A4A; line-height:1.4; font-family: 'Roboto', sans-serif;">
                                        Tidak usah panik, ini hanya lupa password. Santai saja.
                                        <br/>
                                        Buat password baru Anda, semudah mengklik tombol di bawah ini.
                                    </p>
                                </td>
                            </tr>
                            <tr>
                                <td align="center">
                                    <table border="0" cellpadding="0" cellspacing="0" class="btn-wrapper">
                                        <tr>
                                            <td style="padding-top:20px;padding-bottom:20px">
                                                <a class="btn-medium" href="{url}"
                                                   style="color:#FFFFFF; display:inline-block;background-color:#FF9700;font-family:Helvetica, Arial, sans-serif; font-size:16px; font-weight:bold; text-align:center;text-decoration:none;padding: 15px 30px;border-radius: 3px;"
                                                   target="_blank">Reset Password</a>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td align="left">
                                    <p style="font-size:14px; color:#BDBDBD; line-height:1.4; font-family: 'Roboto', sans-serif;">
                                        Abaikan pesan ini, jika Anda tidak melakukan permintaan reset password.
                                    </p>
                                </td>
                            </tr>
                        </table>

                    </td>
                </tr>
                <tr>
                    <td align="left" style="padding-top:27px;padding-bottom:27px">
                        <table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td align="center">
                                    <img alt="Mail Target" border="0" height="50"
                                         src="https://files.mailtarget.co/assets/logo.png"
                                         style="display: block; border: 0px;"/>
                                    <p style="font-size:14px; color:#4A4A4A; line-height:1.4; font-family: 'Roboto', sans-serif;">
                                        Send Email Marketing &
                                        <br/>
                                        Email Automation Easily.
                                    </p>
                                    <div style="Margin-top:24px">
                                        <a href="https://www.facebook.com/mtarget.co/"
                                           style="text-decoration:none;margin-right:5px">
                                            <img alt="" src="https://files.mailtarget.co/assets/fb.png"/>
                                        </a>
                                        <a href="https://twitter.com/mtargetco"
                                           style="text-decoration:none;margin-right:5px">
                                            <img alt="" src="https://files.mailtarget.co/assets/tw.png"/>
                                        </a>
                                        <a href="https://plus.google.com/u/0/102419518489858797401"
                                           style="text-decoration:none;margin-right:5px">
                                            <img alt="" src="https://files.mailtarget.co/assets/gp.png"/>
                                        </a>
                                        <a href="https://www.instagram.com/mtarget.co"
                                           style="text-decoration:none;margin-right:5px">
                                            <img alt="" src="https://files.mailtarget.co/assets/ig.png"/>
                                        </a>
                                    </div>
                                    <p style="font-size:10px; color:#4A4A4A; line-height:1.4; font-family: 'Roboto', sans-serif;">
                                        PT. TARGET SUKSES SINERGI
                                        <br/>
                                        88 Office Tower Lt. 10E
                                        <br/>
                                        JL. Casablanca Kav. 88 Jakarta Selatan 12870.
                                    </p>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
<img id="main-logo" style="display:inline-block;vertical-align:middle;max-width:100%;max-width:300px"; src=""></body>
</html>
`
var wrongHtml = `not html`
var successResponse = SuccessResponse{
	Id: uuid.NewString(),
}
var errorResponse = ErrorResponse{
	Error:  "Bad request",
	Status: 400,
}
var attachment = Attachment{
	MimeType:     "image/png",
	Filename: "AnImage.png",
	Value:  "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAWJQAAFiUBSVIk8AAAAXxJREFUOBFjvJVg84P5718WBjLAX2bmPyxMf/+xMDH8YyZDPwPDXwYGJkIaOXTNGdiUtHAqI2jA/18/GUQzGsg3gMfKg4FVQo6BiYcPqyF4XcChaczA4+DP8P//f4b/P3+SZgAzvxCDSGYjAyMjI8PvZw+AoYXdLuyiQLtE0uoZWAREwLb+fnKXQTipkngXcJu7MnACQx8G2FX1GHgs3bDGBlYX8HlFM/z9+JbhzewWhmf1CQyfti9j+PfzBwO/ZxTMTDiNmQKBfmZX1GB42V/K8P38YbDCX/dvMDAwMzPwuYbBNcIYmC4AhfjvXwx/376AqQHTf96+ZPj34xuKGIiDaQBQ8PPBTQwCoZkMjJzcYA3MgqIMAr7xDJ/3rAHzkQnGO7FWf5gZ/qLmBSZmBoHgNAZee1+Gf18/MzCyczJ83LyQ4fPetch6Gf4xMP3FbgBMGdAgJqAr/n37zABMTTBROA0ygAWUJUG5Civ4B8xwX78CpbD6FJiHmf4AAFicbTMTr5jAAAAAAElFTkSuQmCC",
}
var attachments = []Attachment{
	attachment,
}
var metadata = map[string]string{"key1": "value1", "key2": "value2"}
var optionsAttributes = OptionsAttributes{
	ClickTracking: true,
	OpenTracking:  true,
}

var headers = []Header{
	{
		Name:  "some header",
		Value: "some value",
	},
}