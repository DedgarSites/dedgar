package models

type AppSecrets struct {
	CookieSecret    string
	GoogleAuthID    string
	GoogleAuthKey   string
	PsqlPassword    string
	PsqlUser        string
	PsqlServicePort string
	PsqlDatabase    string
	PsqlServiceHost string
	Subject         string
	CharSet         string
	Sender          string
	Recipient       string
	AuthMap         map[string]bool
}
