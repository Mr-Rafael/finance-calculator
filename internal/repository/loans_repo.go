package repository

import "github.com/Mr-Rafael/finance-calculator/internal/db"

type LoansRepo struct {
	queries *db.Queries
}

func NewLoansRepo(queries *db.Queries) *LoansRepo {
	return &LoansRepo{queries: queries}
}
