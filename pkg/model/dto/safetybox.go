package dto

import "time"

type SafetyBox struct {
	ID         int
	Username   string
	SecretPass string
	Email      string
	Platform   string
	Details    string
	SecretInfo string
	DateAdd    time.Time
	DateModif  time.Time
	CredID     int
}
