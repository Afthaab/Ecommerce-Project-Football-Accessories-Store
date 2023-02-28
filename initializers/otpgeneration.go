package initializers

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	gomail "gopkg.in/mail.v2"
)

func GenCaptchaCode() (string, error) {
	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + (codes[i] % 10))
	}

	return string(codes), nil
}

func Otpgeneration(emails string) string {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("EMAIL"))

	// Set E-Mail receivers
	m.SetHeader("To", emails)

	// Set E-Mail subject
	m.SetHeader("Subject", "OTP to verify your Gmail")

	//otp generation
	onetimepassword, err := GenCaptchaCode()
	if err != nil {
		panic(err)
	}

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", onetimepassword+" is your OTP to register to our site. Thank you registering to our site. Happy Shopping :)")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("OTP has been sent successfully")
	}
	return onetimepassword

}

func Twilio(phone string) string {
	fmt.Println("hi")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("ACCOUNT_SID"),
		Password: os.Getenv("AUTH_TOKEN"),
	})
	onetimepassword, err := GenCaptchaCode()
	if err != nil {
		panic(err)
	}
	params := &twilioApi.CreateMessageParams{}
	params.SetTo("+919741968809")
	params.SetFrom("+12706389186")
	params.SetBody(onetimepassword + "is your OTP to register to our site. Thank you registering to our site. Happy Shopping :) " + phone)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return onetimepassword
}
