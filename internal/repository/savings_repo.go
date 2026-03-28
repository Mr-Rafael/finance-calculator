package repository

import (
	"github.com/Mr-Rafael/finance-calculator/internal/db"
)

type SavingsRepo struct {
	queries *db.Queries
}

func NewSavingsRepo(queries *db.Queries) *SavingsRepo {
	return &SavingsRepo{queries: queries}
}
