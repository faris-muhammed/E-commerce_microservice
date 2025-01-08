package helper

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// GenerateOtp creates a 6-digit random OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SendOtp sends the generated OTP to the user's email
func SendOTP(email, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("APPEMAIL"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code for Signup")
	m.SetBody("text", "Your OTP for signup is: "+otp)
	fmt.Println("APPEMAIL:", os.Getenv("APPEMAIL"))
	fmt.Println("APPPASSWORD:", os.Getenv("APPPASSWORD"))

	// Use the email and password from environment variables
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("APPEMAIL"), os.Getenv("APPPASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}
func HashPassword(password string) (string, error) {
	// Generate a hashed password with bcrypt's default cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
