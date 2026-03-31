-- +goose Up
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    starting_principal INT NOT NULL,
    yearly_interest_rate TEXT NOT NULL,
    monthly_payment INT NOT NULL,
    escrow_payment INT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    duration_months INT NOT NULL,
    total_expenditure INT NOT NULL,
    total_paid INT NOT NULL,
    cost_of_credit TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE loan_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID REFERENCES loans(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    payment INT NOT NULL,
    interest INT NOT NULL,
    other_payments INT NOT NULL,
    paydown INT NOT NULL,
    principal INT NOT NULL
);

-- +goose Down
DROP TABLE loan_state;
DROP TABLE loans;