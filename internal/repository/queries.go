package repository

const (

	// Balance queries
	queryGetAccrued = `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE login = $1
	`

	queryGetWithdrawn = `
		SELECT COALESCE(SUM(withdrawn), 0)
		FROM withdrawals
		WHERE login = $1
	`

	queryInsertWithdrawal = `
		INSERT INTO withdrawals (login, order_number, withdrawn, processed_at)
		VALUES ($1, $2, $3, $4)
	`

	queryGetWithdrawalsByUser = `
		SELECT order_number, withdrawn, processed_at
		FROM withdrawals
		WHERE login = $1
		ORDER BY processed_at DESC
	`

	// Authorization queries
	queryInsertUser = `
		INSERT INTO users(login, password, token)
		VALUES($1, $2, $3)
	`

	queryGetUserByLogin = `
		SELECT login, password, token
		FROM users
		WHERE login = $1
	`

	// Order queries
	queryGetOrder = `
		SELECT number, login, uploaded_at, status, accrual
		FROM orders
		WHERE number = $1
	`

	queryInsertOrder = `
		INSERT INTO orders (number, login, uploaded_at, status)
		VALUES ($1, $2, $3, $4)
	`

	queryUpdateOrder = `
		UPDATE orders
		SET status = $2,
		    accrual = $3
		WHERE number = $1
	`

	queryGetOrdersByUser = `
		SELECT number, status, accrual, uploaded_at
		FROM orders
		WHERE login = $1
		ORDER BY uploaded_at DESC
	`
)
