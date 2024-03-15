package onetp

import (
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateOTP(length int) string {
	characters := "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(characters))
		otp += string(characters[randomIndex])
	}

	return otp
}

func SendOtp(email, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ptkseries@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "SignUp verification code")
	m.SetBody("text/plain", "Your OTP for SignUp is: "+otp)

	d := gomail.NewDialer("smtp.gmail.com", 587, "ptkseries@gmail.com", "dcfgjdspvkwinhkj")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil

}
