package services

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	"github.com/ihorlenko/weather_notifier/internal/types"
)

var _ interfaces.EmailService = (*EmailService)(nil)

type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	baseURL  string
}

func NewEmailService(cfg *config.Config) interfaces.EmailService {
	return &EmailService{
		from:     cfg.EmailConfig.From,
		password: cfg.EmailConfig.Password,
		smtpHost: cfg.EmailConfig.SMTPHost,
		smtpPort: cfg.EmailConfig.SMTPPort,
		baseURL:  cfg.AppConfig.BaseURL,
	}
}

func (s *EmailService) SendConfirmationEmail(email, city, token string) error {
	subject := "Confirm Your Weather Notifier Subscription"
	confirmURL := fmt.Sprintf("%s/api/confirm/%s", s.baseURL, token)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Confirm Your Subscription</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #5b86e5 0%%, #36d1dc 100%%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        .content {
            padding: 20px;
        }
        .button {
            display: inline-block;
            background: linear-gradient(135deg, #5b86e5 0%%, #36d1dc 100%%);
            color: white;
            text-decoration: none;
            padding: 12px 25px;
            border-radius: 4px;
            margin: 20px 0;
            font-weight: bold;
            text-align: center;
        }
        .footer {
            background-color: #f5f5f5;
            padding: 15px;
            text-align: center;
            font-size: 14px;
            color: #888888;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Weather Notifier</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>Thank you for subscribing to weather updates for <strong>%s</strong>.</p>
            <p>To confirm your subscription, please click the button below:</p>
            <div style="text-align: center;">
                <a href="%s" class="button">Confirm Subscription</a>
            </div>
            <p>If you did not subscribe to this service, please ignore this email.</p>
            <p>Best regards,<br>Weather Notifier Team</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Weather Notifier by `+
		`<a href="https://github.com/ihorlenko" style="color: #5b86e5;">ihorlenko</a></p>
        </div>
    </div>
</body>
</html>
`, city, confirmURL)

	return s.sendHTMLEmail(email, subject, htmlBody)
}

func (s *EmailService) SendWeatherUpdate(
	email,
	city string,
	weather *types.WeatherData,
	unsubscribeToken string,
) error {
	subject := fmt.Sprintf("Weather Update for %s", city)
	unsubscribeURL := fmt.Sprintf("%s/api/unsubscribe/%s", s.baseURL, unsubscribeToken)

	weatherEmoji := getWeatherEmoji(weather.Description)

	tempColor := "#3498db"
	if weather.Temperature > 20 {
		tempColor = "#e74c3c"
	} else if weather.Temperature > 10 {
		tempColor = "#f39c12"
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weather Update</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            max-width: 600px;
            margin: 0 auto;
            padding: 0;
            background-color: #f9f9f9;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            margin: 20px auto;
        }
        .header {
            background: linear-gradient(135deg, #5b86e5 0%%, #36d1dc 100%%);
            color: white;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 500;
        }
        .header p {
            margin: 10px 0 0;
            opacity: 0.9;
        }
        .content {
            padding: 30px 20px;
            text-align: center;
        }
        .weather-icon {
            font-size: 72px;
            margin: 20px 0;
            line-height: 1;
        }
        .temperature {
            font-size: 48px;
            font-weight: 700;
            margin: 5px 0;
            color: %s;
        }
        .description {
            font-size: 22px;
            margin: 10px 0 20px;
            font-weight: 500;
        }
        .details {
            background-color: #f9f9f9;
            border-radius: 8px;
            padding: 15px;
            margin: 20px auto;
            max-width: 80%%;
            text-align: left;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-row:last-child {
            border-bottom: none;
        }
        .detail-label {
            color: #666;
			padding-right: 10px;
        }
        .detail-value {
            font-weight: 500;
        }
        .unsubscribe {
            text-align: center;
            padding: 15px 20px;
            border-top: 1px solid #eee;
            color: #888;
            font-size: 14px;
        }
        .unsubscribe a {
            color: #5b86e5;
            text-decoration: none;
        }
        .unsubscribe a:hover {
            text-decoration: underline;
        }
        .footer {
            background-color: #f5f5f5;
            padding: 15px;
            text-align: center;
            font-size: 14px;
            color: #888888;
            border-top: 1px solid #eee;
        }
        .footer a {
            color: #5b86e5;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Weather Update for %s</h1>
            <p>Current conditions as of now</p>
        </div>
        
        <div class="content">
            <div class="weather-icon">%s</div>
            
            <div class="temperature">%.1f°C</div>
            
            <div class="description">%s</div>
            
            <div class="details">
                <div class="detail-row">
                    <span class="detail-label">Humidity</span>
                    <span class="detail-value">%.2f%%</span>
                </div>
            </div>
        </div>
        
        <div class="unsubscribe">
            <p>To unsubscribe from these weather updates, <a href="%s">click here</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; 2025 Weather Notifier by <a href="https://github.com/ihorlenko">ihorlenko</a></p>
        </div>
    </div>
</body>
</html>
`, tempColor, city, weatherEmoji, weather.Temperature, weather.Description, weather.Humidity, unsubscribeURL)

	return s.sendHTMLEmail(email, subject, htmlBody)
}

func (s *EmailService) sendHTMLEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message string
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}

func getWeatherEmoji(description string) string {
	description = strings.ToLower(description)

	switch {
	case strings.Contains(description, "clear") || strings.Contains(description, "sunny"):
		return "☀️"
	case strings.Contains(description, "partly cloudy"):
		return "⛅"
	case strings.Contains(description, "cloudy") || strings.Contains(description, "overcast"):
		return "☁️"
	case strings.Contains(description, "rain") && strings.Contains(description, "thunder"):
		return "⛈️"
	case strings.Contains(description, "thunder") || strings.Contains(description, "lightning"):
		return "🌩️"
	case strings.Contains(description, "drizzle") || strings.Contains(description, "light rain"):
		return "🌦️"
	case strings.Contains(description, "rain"):
		return "🌧️"
	case strings.Contains(description, "snow") && strings.Contains(description, "rain"):
		return "🌨️"
	case strings.Contains(description, "snow"):
		return "❄️"
	case strings.Contains(description, "sleet"):
		return "🌨️"
	case strings.Contains(description, "fog") || strings.Contains(description, "mist"):
		return "🌫️"
	case strings.Contains(description, "wind"):
		return "💨"
	default:
		return "🌡️"
	}
}
