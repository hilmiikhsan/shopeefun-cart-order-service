-- +goose Up
-- +goose StatementBegin
INSERT INTO cart_items (user_id, product_id, qty)
VALUES
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440001', 2),
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440002', 1),
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440003', 3),
    ('c396f23e-a097-476d-aae5-cfc9973634f3', '550e8400-e29b-41d4-a716-446655440005', 4),
    ('c396f23e-a097-476d-aae5-cfc9973634f3', '550e8400-e29b-41d4-a716-446655440006', 5);

-- Insert data into orders
INSERT INTO orders (
    user_id, payment_type_id, order_number, total_price, 
    product_order, status, is_paid, ref_code, created_at, updated_at, deleted_at
) VALUES 
(
    'b3d1a4f6-1f12-4f1a-87d8-1e82671e3f1a', '4d5a2c3f-2b19-4e8a-8f4b-1e932c2d3a4b', 'ORD001', 110.50,
    '[{"product_id": "prod1", "qty": 2, "price": 50.25}]', 'pending', FALSE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    'a2c3e4b5-3d11-4f2a-87d9-1e93762e4f2b', '5e6b3d4f-3c2a-4f9b-9f5c-2f032d3e4a5c', 'ORD002', 220.00,
    '[{"product_id": "prod2", "qty": 4, "price": 50.00}]', 'completed', TRUE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    'd55e1169-3c3f-423d-bb26-e5b47345a931', '6f7c4e5f-4d3b-5f1a-af5d-3f142e4f5b6d', 'ORD003', 330.00,
    '[{"product_id": "prod3", "qty": 6, "price": 50.00}]', 'shipped', TRUE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    'd55e1169-3c3f-423d-bb26-e5b47345a931', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD004', 440.00,
    '[{"product_id": "prod4", "qty": 8, "price": 50.00}]', 'cancelled', FALSE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD005', 550.00,
    '[{"product_id": "prod5", "qty": 10, "price": 50.00}]', 'pending', FALSE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD006', 660.00,
    '[{"product_id": "prod6", "qty": 12, "price": 50.00}]', 'completed', TRUE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD007', 770.00,
    '[{"product_id": "prod7", "qty": 14, "price": 50.00}]', 'shipped', TRUE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD008', 880.00,
    '[{"product_id": "prod8", "qty": 16, "price": 50.00}]', 'cancelled', FALSE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD009', 990.00,
    '[{"product_id": "prod9", "qty": 18, "price": 50.00}]', 'pending', FALSE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
),
(
    '346df8a0-1cb0-4fd5-bf39-9739dc848f16', '346df8a0-1cb0-4fd5-bf39-9739dc848f16', 'ORD010', 1100.00,
    '[{"product_id": "prod10", "qty": 20, "price": 50.00}]', 'completed', TRUE, 'REF' || EXTRACT(EPOCH FROM NOW()), now(), NULL, NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM cart_items 
WHERE (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440001')
   OR (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440002')
   OR (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440003')
   OR (user_id = 'c396f23e-a097-476d-aae5-cfc9973634f3' AND product_id = '550e8400-e29b-41d4-a716-446655440005')
   OR (user_id = 'c396f23e-a097-476d-aae5-cfc9973634f3' AND product_id = '550e8400-e29b-41d4-a716-446655440006');

DELETE FROM orders 
WHERE order_number IN ('ORD001', 'ORD002', 'ORD003', 'ORD004', 'ORD005', 'ORD006', 'ORD007', 'ORD008', 'ORD009', 'ORD010');
-- +goose StatementEnd
