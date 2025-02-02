package service

import (
	"fmt"
	"library-backend/config"
	"library-backend/internal/models"
	"log"
	"net/smtp"
)

func SendEmailConfirmation(scheduleInput models.PickupSchedule) {
	subject := "Pickup Confirmation"
	body := fmt.Sprintf(
		"Hello %s,\n\nYour book pickup is confirmed.\n\n"+
			"Details:\n"+
			"- Book: %s\n"+
			"- Date: %s\n"+
			"- Time: %s\n\nThank you!",
		scheduleInput.User.Name,
		scheduleInput.Book.Title,
		scheduleInput.PickupDate,
		scheduleInput.PickupTime,
	)

	err := SendEmail(scheduleInput.User.Email, subject, body)
	if err != nil {
		log.Printf("Error sending email to %s: %v\n", scheduleInput.User.Email, err)
		return
	}

	log.Printf("Email sent successfully to %s", scheduleInput.User.Email)
}

func SendEmail(to, subject, body string) error {
	cfg := config.Config

	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	msg := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"From: %s\r\n"+
			"To: %s\r\n"+
			"\r\n%s\r\n",
		subject,
		cfg.EmailFrom,
		to,
		body,
	))

	smtpAddr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(smtpAddr, auth, cfg.EmailFrom, []string{to}, msg)
}
