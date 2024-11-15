package util

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func GetDefaultName(name string) string {
	if name == "" {
		return "Default Name"
	}
	return name
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
func IsValidDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	_, err := net.LookupMX(domain)
	return err == nil
}
func ValidateRequestRegister(request *model.RegisterRequest) error {
	if !IsValidEmail(request.Email) {
		return ErrInvalidEmail
	}
	if !IsValidDomain(request.Email) {
		return ErrInvalidDomain
	}
	return nil
}
func ValidateRequestLogin(request *model.LoginRequest) error {
	if !IsValidEmail(request.Email) {
		return ErrInvalidEmail
	}
	if !IsValidDomain(request.Email) {
		return ErrInvalidDomain
	}
	return nil
}

func GenerateSecretKey(email string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "digital-voter",
		AccountName: email,
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}
func Contain(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GenerateOTPLogin(secret string) (string, error) {
	otp, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period: 60,
	})
	if err != nil {
		return "", err
	}
	return otp, nil
}

func GenerateOTPRegister() (string, error) {
	max := big.NewInt(1000000)
	otp, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", otp.Int64()), nil
}
func HashOTP(otp string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyOTPRegister(hashedOTP, inputOTP string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedOTP), []byte(inputOTP))
	return err == nil
}
func SendOTPRegister(email, otp string, viperConfig *viper.Viper) error {
	smtpHost := viperConfig.GetString("SMTP_HOST")
	smtpPort := viperConfig.GetString("SMTP_PORT")
	smtpUser := viperConfig.GetString("SMTP_USER")
	smtpPass := viperConfig.GetString("SMTP_PASS")

	subject := "Email Registration Confirmation"
	body := fmt.Sprintf("Thanks for your registration. Please put your OTP for verification process:\n\nOTP: %s\n\nThis OTP is valid for 15 minutes.", otp)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)
	to := []string{email}

	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return err
	}
	defer client.Quit()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(smtpUser); err != nil {
		return err
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}