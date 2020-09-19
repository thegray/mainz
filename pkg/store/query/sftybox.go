package query

const (
	GetSafetyBoxById = `SELECT * FROM safetybox WHERE id = ?`

	InsertSafetyBox = `
	INSERT INTO 
		safetybox (username, secret_pass, email, platform, details, secret_info, cred_id) 
	VALUES 
		(?, ?, ?, ?, ?, ?, ?)`

	InsertSBToCategory = `
	INSERT INTO 
		sb_categories (cat_id, sb_id) 
	VALUES 
		(?, ?)`

	UpdateSafetyBox = `
	UPDATE safetybox 
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
	UPDATE safetybox 
	SET 
		secret_pass = ?,
		secret_info = ?
	WHERE 
		id = ?
	`

	DeleteSafetyBox = `
	DELETE FROM safetybox 
	WHERE 
		id = ?
	`

	DeleteSBToCategory = `
	DELETE FROM sb_categories
	WHERE 
		sb_id = ?
	`
)
