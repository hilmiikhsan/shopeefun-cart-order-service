package order

const (
	queryCreateOrder = `
		INSERT INTO orders
		(
			user_id,
			payment_type_id,
			order_number,
			total_price,
			product_order,
			status,
			is_paid,
			ref_code,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) RETURNING id, ref_code
	`

	queryCreateOrderStatusLogs = `
		INSERT INTO order_status_logs
		(
			order_id,
			ref_code,
			from_status,
			to_status,
			notes,
			created_at
		) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING ref_code
	`
)
