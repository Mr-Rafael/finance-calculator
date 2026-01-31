package models

import "time"

type SavingsPlan struct {
	Plan []SavingsStatus `json:"plan"`
}

type SavingsStatus struct {
	Date         time.Time
	Interest     int `json:"interest"`
	Tax          int `json:"tax"`
	Contribution int `json:"contribution"`
	Increase     int `json:"increase"`
	Capital      int `json:"capital"`
}
