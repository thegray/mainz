package query

const (
	CheckCredential = `SELECT EXISTS(SELECT id FROM credential WHERE user = ? AND pass = ?) AS found`
	GetCred         = `SELECT * FROM credential WHERE user = ?`
)
