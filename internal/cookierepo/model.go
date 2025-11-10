package cookierepo

type NoCookiesError struct {
}

func (e *NoCookiesError) Error() string {
	return "we have no cookies"
}
