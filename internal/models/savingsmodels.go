package models

type SavingsPlan struct {
	Plan []SavingsStatus `json:"plan"`
}

type SavingsStatus struct {
	Interest     int `json:"interest"`
	Contribution int `json:"contribution"`
	Increase     int `json:"increase"`
	Capital      int `json:"capital"`
}
