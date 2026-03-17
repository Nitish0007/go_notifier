package dto



type SMTPConfiguration struct {
	Host     string `json:"host" validate:"required,min=1"`
	Port     int    `json:"port" validate:"required,gt=0,lte=65535"`
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
	From     string `json:"from" validate:"required,email"`
}

// type WebAppConfiguration struct {
// 	WebAppURL    string `json:"web_app_url" validate:"required,url"`
// 	WebAppSecret string `json:"web_app_secret" validate:"required,min=1"`
// 	WebAppToken  string `json:"web_app_token" validate:"required,min=1"`
// }