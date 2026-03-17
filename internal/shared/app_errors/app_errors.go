package app_errors

type AppError interface {
	Code() 		string
	Message() string
	Status() 	int
}

// defining custom errors
// var (
// 	ErrAccountIDMissing = &AppError{Code: "ACCOUNT_ID_MISSING", Message: "AccountID is required", Status: 400 }
// )