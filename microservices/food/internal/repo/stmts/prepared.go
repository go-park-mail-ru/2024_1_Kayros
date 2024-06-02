package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"getByRestId": 		`SELECT c.name, f.id, f.name, restaurant_id, weight, price, img_url FROM food as f
		JOIN category as c ON f.category_id=c.id WHERE restaurant_id = $1 ORDER BY category_id`,
		"getById": 		`SELECT id, name, restaurant_id, category_id, weight, price, img_url
		FROM food WHERE id=$1`,
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