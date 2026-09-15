package quoterepo

type NoQuotesError struct {
}

func (e *NoQuotesError) Error() string {
	return "no quotes available"
}
