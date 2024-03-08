package helper

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)
func GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
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
