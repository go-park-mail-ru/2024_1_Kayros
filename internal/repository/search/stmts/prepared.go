package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"selectRestBySearch": `SELECT id, name, img_url FROM restaurant 
		WHERE LOWER(name) LIKE LOWER('%' || $1 || '%')`,
		"getRestsByCategory": `SELECT DISTINCT r.id, r.name, r.img_url FROM restaurant AS r
		JOIN rest_categories AS rc ON r.id=rc.restaurant_id JOIN category AS c
		ON rc.category_id=c.id WHERE LOWER(c.name) LIKE LOWER('%' || $1 || '%')`,
		"selectRests": `SELECT id, name FROM category AS c
		JOIN rest_categories AS rc ON c.id=rc.category_id WHERE rc.restaurant_id=$1`,
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
