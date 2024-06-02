package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"getAllRest": `SELECT id, name, short_description, address, img_url FROM restaurant ORDER BY rating DESC`,
		"getRestById": `SELECT id, name, long_description, address, img_url, rating, comment_count FROM restaurant WHERE id=$1`,
		"getRestListUsingFilter": `SELECT r.id, r.name, r.short_description, r.img_url FROM restaurant as r JOIN rest_categories AS rc ON r.id=rc.restaurant_id WHERE rc.category_id=$1`,
		"getRestsByCategory": `SELECT id, name FROM category WHERE type='rest'`,
		"getTopRests": `SELECT id, name, short_description, img_url FROM restaurant ORDER BY rating DESC LIMIT $1`,
		"getLastRests": `SELECT f.restaurant_id  FROM food AS f JOIN food_order AS fo ON f.id=fo.food_id 
		JOIN "order" AS o ON o.id=fo.order_id WHERE o.user_id=$1  GROUP BY f.restaurant_id
		ORDER BY MAX(o.delivered_at) DESC LIMIT $2`,
		"getShortRestById": `SELECT id, name, short_description, img_url FROM restaurant WHERE id=$1`,
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