-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_logs (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id UUID NOT NULL,
    ref_code VARCHAR(100) NOT NULL,
    from_status VARCHAR(50) NOT NULL,
    to_status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT now(),

    FOREIGN KEY (order_id) REFERENCES orders(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_status_logs CASCADE;
-- +goose StatementEnd
