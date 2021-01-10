package query

const (
	GetAccountById = `SELECT * FROM Account WHERE id = ?`

	InsertAccount = `
	INSERT INTO 
		account (username, email_phonenum, platform, details, cred_id) 
	VALUES 
		(?, ?, ?, ?, ?)`

	InsertAccountHistory = `
	INSERT INTO 
		account_history (username, email_phonenum, platform, details, acc_id, version) 
	VALUES 
		(?, ?, ?, ?, ?, ?)`

	InsertAccToCategories = `
	INSERT INTO 
		acc_categories (cat_id, acc_id) 
	VALUES 
		(?, ?)`

	UpdateAccount = `
	UPDATE Account 
	SET 
		username = ?,
		email = ?,
		platform = ?,
		details = ?,
		cred_id = ?
	WHERE 
		id = ?
	`

	UpdateSecret = `
	UPDATE Account 
	SET 
		secret_pass = ?,
		secret_info = ?
	WHERE 
		id = ?
	`

	DeleteAccount = `
	DELETE FROM Account 
	WHERE 
		id = ?
	`

	DeleteSBToCategory = `
	DELETE FROM sb_categories
	WHERE 
		sb_id = ?
	`
)
