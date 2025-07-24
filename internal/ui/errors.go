package ui

// CancelledError represents when the user cancels the selection (ESC key or Ctrl+C)
type CancelledError struct {
    Message string
}

func (e *CancelledError) Error() string {
    return e.Message
}

// IsCancelled checks if the error represents a user cancellation
func IsCancelled(err error) bool {
    _, ok := err.(*CancelledError)
    return ok
}

// NewCancelledError creates a new CancelledError
func NewCancelledError(message string) *CancelledError {
    return &CancelledError{Message: message}
}
