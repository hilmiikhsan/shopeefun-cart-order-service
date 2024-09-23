-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL,
    payment_type_id UUID NOT NULL,
    order_number VARCHAR(100) NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    product_order JSONB, -- Added qty in product_order
    status VARCHAR(50) NOT NULL,
    is_paid BOOLEAN NOT NULL,
    ref_code VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders CASCADE;
-- +goose StatementEnd
