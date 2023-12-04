-- +goose Up

-- theme labels
ALTER TABLE theme ADD labels TEXT DEFAULT "" NOT NULL;

-- +goose Down

-- theme labels
ALTER TABLE theme DROP COLUMN labels;
