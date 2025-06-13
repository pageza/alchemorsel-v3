package errors

// AppError represents a structured application error.
type AppError struct {
	Code    string
	Message string
	Status  int
	Details map[string]any
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// New creates a new AppError without details.
func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

// NewWithDetails creates a new AppError with additional detail data.
func NewWithDetails(code, message string, status int, details map[string]any) *AppError {
	return &AppError{Code: code, Message: message, Status: status, Details: details}
}

var (
	// ErrUserNotFound is returned when a user cannot be located.
	ErrUserNotFound = New("user_not_found", "user not found", 404)
	// ErrInvalidInput indicates invalid client supplied data.
	ErrInvalidInput = New("invalid_input", "invalid input", 400)
)
