package cart

const (
	queryGetAllCart = `
		SELECT
			id,
			user_id,
			product_id,
			qty,
			created_at,
			updated_at,
			deleted_at
		FROM cart_items
	`

	queryInsertOrUpdateCart = `
		INSERT INTO cart_items
		(
			user_id,
			product_id,
			qty,
			created_at
		) VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, product_id) 
		DO UPDATE SET
			qty = EXCLUDED.qty,
			updated_at = NOW()
		RETURNING id, user_id, product_id, qty, created_at, updated_at
	`

	queryGetProductByUserIdAndProductId = `
		SELECT 
			product_id,
			deleted_at 
		FROM cart_items
		WHERE user_id = $1 AND product_id = $2
	`

	queryLockUpdateQty = `
		SELECT 
			1
		FROM cart_items
		WHERE user_id = $1
		FOR UPDATE
	`

	queryUpdateQty = `
		UPDATE cart_items
		SET 
			qty = $1,
			updated_at = NOW()
		WHERE user_id = $2 AND product_id = $3
		RETURNING id, user_id, product_id, qty, created_at, updated_at, deleted_at
	`

	queryLockSoftDeleteProduct = `
		SELECT 
			1
		FROM cart_items	
		WHERE user_id = $1
		FOR UPDATE
	`

	queryUpdateDeletedAt = `
		UPDATE cart_items
		SET 
			deleted_at = NOW(),
			qty = 0
		WHERE user_id = $1 
		AND product_id = $2
	`

	queryCheckUserExists = `
		SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1)
	`

	queryCheckProductInCart = `
		SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1 AND product_id = $2)
	`

	queryRestoreDeletedProduct = `
		UPDATE cart_items
		SET 
			deleted_at = NULL, 
			qty = $3
		WHERE user_id = $1 AND product_id = $2 AND deleted_at IS NOT NULL
	`
)
