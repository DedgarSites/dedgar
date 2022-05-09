package controllers

import (
	"strings"

	"github.com/dedgarsites/dedgar/datastores"

	"github.com/labstack/echo/v4"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	asession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"fmt"
	"net/http"
)

type filePath struct {
	Path string
}

// GET /all/*
func GetMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", datastores.PostMap)
}

// GET /login
func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// GET /about
func GetAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

// GET /contact
func GetContact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

// GET /login
func GetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// GET /privacy
func GetPrivacy(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy.html", nil)
}

// POST /post-contact
func PostContact(c echo.Context) error {

	if strings.Contains(c.FormValue("message"), "http") && strings.Contains(c.FormValue("message"), "dedgar.com/") == false {
		return c.String(http.StatusOK, "Form submitted")
	}

	TextBody := c.FormValue("name") + "\n" + c.FormValue("email") + "\n" + c.FormValue("message")

	sess, err := asession.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(datastores.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(datastores.CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(datastores.CharSet),
				Data:    aws.String(datastores.Subject),
			},
		},
		Source: aws.String(datastores.Sender),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

	}
	fmt.Println(c.FormValue("name"))
	fmt.Println(c.FormValue("email"))
	fmt.Println(c.FormValue("message"))
	fmt.Println("Email Sent to address: " + datastores.Recipient)
	fmt.Println(result)
	return c.String(http.StatusOK, "Form submitted")
}

// GET /post/:postname
func GetPost(c echo.Context) error {
	post := c.Param("postname")
	if _, ok := datastores.PostMap[post]; ok {
		return c.Render(http.StatusOK, post+".html", post)
	}
	return c.Render(http.StatusNotFound, "e04.html", "404 Post not found")
}

// GET /post
func GetPostView(c echo.Context) error {
	return c.Render(http.StatusOK, "post_view.html", datastores.PostMap)
}

// handle any error by attempting to render a custom page for it
func Custom404Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.Render(code, errorPage, code); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}
