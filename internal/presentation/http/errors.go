package http

type ErrorBag map[string][]string

// Create a new ErrorBag
func NewErrorBag() *ErrorBag {
	return &ErrorBag{}
}

// Add a new error message for a specific field
func (b *ErrorBag) Add(field, message string) {
	(*b)[field] = append((*b)[field], message)
}

// Returns true if there are any errors
func (b *ErrorBag) HasErrors() bool {
	return len(*b) > 0
}

// Returns true if there are no errors
func (b *ErrorBag) IsEmpty() bool {
	return len(*b) == 0
}
