package dto

import "time"

type SafetyBox struct {
	ID           int
	SecretPass   string
	SecretInfo   string
	LastModified time.Time
	AccID        int
}
