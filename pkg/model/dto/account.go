package dto

import "time"

type Account struct {
	ID            int
	Username      string
	EmailPhonenum string
	Platform      string
	Details       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CredID        int
}

type AccountHistory struct {
	ID            int
	Username      string
	EmailPhonenum string
	Platform      string
	Details       string
	AddedAt       time.Time
	AccID         int
	Version       int
}
