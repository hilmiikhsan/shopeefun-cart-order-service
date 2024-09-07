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

	queryAddCart = `
		INSERT INTO cart_items
		(
			user_id,
			product_id,
			qty,
			created_at
		) VALUES ($1, $2, $3, NOW()) RETURNING id
	`

	queryGetProductByUserIdAndProductId = `
		SELECT 
			product_id 
		FROM cart_items
		WHERE user_id = $1 AND product_id = ANY($2)
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
			qty = $1
		WHERE user_id = $2 AND product_id = $3
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
			deleted_at = NOW()
		WHERE user_id = $1 
		AND product_id = ANY($2::uuid[])
	`

	queryCheckUserExists = `
		SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1 AND deleted_at IS NULL)
	`

	queryCheckProductInCart = `
		SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1 AND product_id = ANY($2::uuid[]) AND deleted_at IS NULL)
	`
)
