package internal

type Parent struct {
	ID           string
	Email        string
	Country      string
	Lang         string
	PasswordHash string
}

type Consent struct {
	ID       string
	ParentID string
	Method   string
	Verified bool
}
