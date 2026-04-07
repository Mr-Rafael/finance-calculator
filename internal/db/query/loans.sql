-- name: CreateLoan :one
INSERT INTO loans(user_id,
    name,
    starting_principal,
    yearly_interest_rate,
    monthly_payment,
    escrow_payment,
    start_date,
    monthly_interest_rate,
    duration_months,
    total_expenditure,
    total_paid,
    cost_of_credit
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetLoansByUserID :many
SELECT id, name, starting_principal FROM loans
WHERE user_id = $1;

-- name: GetLoan :one
SELECT * FROM loans
WHERE id = $1 AND user_id = $2;