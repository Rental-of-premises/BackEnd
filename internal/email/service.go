package email

import (
	"bytes"
    "fmt"
    "log"
    "html/template"
    "path/filepath"

    "rent/internal/config"

    "gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer   *gomail.Dialer
	from     string
	fromName string
	templateDir string
}

func NewEmailService(cfg *config.Config) *EmailService {
	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	return &EmailService{
		dialer:   dialer,
		from:     cfg.SMTPFromEmail,
		fromName: cfg.SMTPFromName,
		templateDir: "internal/email/patterns",
	}
}

func (s *EmailService) RenderTemplate(templateName string, data interface{}) (string, error) {
    basePath := filepath.Join(s.templateDir, "base.html")
    tmplPath := filepath.Join(s.templateDir, templateName)

    tmpl, err := template.ParseFiles(basePath, tmplPath)
    if err != nil {
        return "", fmt.Errorf("failed to parse templates: %w", err)
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, "base.html", data)
    if err != nil {
        return "", fmt.Errorf("failed to execute template: %w", err)
    }

    return buf.String(), nil
}

func (s *EmailService) SendWelcomeEmail(to, name string) error {
    subject := "Добро пожаловать в Rental Service!"

    data := struct {
        Name string
    }{
        Name: name,
    }

    body, err := s.RenderTemplate("welcome.html", data)
    if err != nil {
        return err
    }

    return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendEmail(to, subject, body string) error {
    log.Printf("📤 Отправка письма через SMTP на %s", to)
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)


	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}