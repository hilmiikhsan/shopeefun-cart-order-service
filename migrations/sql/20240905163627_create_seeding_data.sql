-- +goose Up
-- +goose StatementBegin
INSERT INTO cart_items (user_id, product_id, qty)
VALUES
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440001', 2),
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440002', 1),
    ('a4b7a3f1-751a-4a10-b506-99202581427b', '550e8400-e29b-41d4-a716-446655440003', 3),
    ('c396f23e-a097-476d-aae5-cfc9973634f3', '550e8400-e29b-41d4-a716-446655440005', 4),
    ('c396f23e-a097-476d-aae5-cfc9973634f3', '550e8400-e29b-41d4-a716-446655440006', 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM cart_items 
WHERE (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440001')
   OR (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440002')
   OR (user_id = 'a4b7a3f1-751a-4a10-b506-99202581427b' AND product_id = '550e8400-e29b-41d4-a716-446655440003')
   OR (user_id = 'c396f23e-a097-476d-aae5-cfc9973634f3' AND product_id = '550e8400-e29b-41d4-a716-446655440005')
   OR (user_id = 'c396f23e-a097-476d-aae5-cfc9973634f3' AND product_id = '550e8400-e29b-41d4-a716-446655440006');
-- +goose StatementEnd
