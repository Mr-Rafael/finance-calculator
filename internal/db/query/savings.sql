-- name: CreateSavings :one
INSERT INTO savings (user_id,
    name,
    starting_capital,
    yearly_interest_rate,
    interest_rate_type,
    monthly_contribution,
    duration_years,
    tax_rate,
    yearly_inflation_rate,
    start_date,
    monthly_interest_rate,
    total_interest_earnings,
    rate_of_return,
    inflation_adjusted_ror
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetSavingsByUserID :many
SELECT id, name, starting_capital FROM savings
WHERE user_id = $1;

-- name: GetSavings :one
SELECT * FROM savings
WHERE id = $1;