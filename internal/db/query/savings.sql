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
    start_date
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetSavingsByUserID :many
SELECT * FROM savings
WHERE user_id = $1;

