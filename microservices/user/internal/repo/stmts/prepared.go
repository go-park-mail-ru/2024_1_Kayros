package stmts

import "database/sql"

func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"getUser": `SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, ''), is_vk_user FROM "user" WHERE email = $1`,
		"deleteUser": `DELETE FROM "user" WHERE email = $1`,
		"createUser": `INSERT INTO "user" (name, email, phone, password, address, img_url, is_vk_user, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		"updateUser": `UPDATE "user" SET name = $1, email = $2, phone = $3, img_url = $4, password = $5, card_number = $6, address = $7, updated_at = $8 WHERE email = $9`,
		"getUnauthAddress": `SELECT address FROM unauth_address WHERE unauth_id = $1`,
		"updateUnauthAddress": `UPDATE unauth_address SET address = $1 WHERE unauth_id= $2`,
		"createUnauthAddress": `INSERT INTO unauth_address (unauth_id, address) VALUES ($1, $2)`,
	}
	preparedStatements := make(map[string]*sql.Stmt, len(statements))
	for key, value := range statements {
		stmt, err := db.Prepare(value)
		if err != nil {
			return nil, err
		}
		preparedStatements[key] = stmt
	}

	return preparedStatements, nil
}