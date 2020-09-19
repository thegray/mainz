package query

const (
	InsertCategory = `
	INSERT INTO 
		category (id, name) 
	VALUES 
		(?, ?)`
)
