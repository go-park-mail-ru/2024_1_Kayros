package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"createComment": `INSERT INTO "comment" (user_id, restaurant_id, text, rating) VALUES ($1, $2, $3, $4) RETURNING id`,
		"getRestRating": `SELECT rating, comment_count FROM restaurant WHERE id=$1`,
		"updateRestRating": `UPDATE restaurant SET rating=$1, comment_count=$2 WHERE id=$3`,
		"updateOrderCommentedStatus": `UPDATE "order" SET commented=true WHERE id=$1`,
		"getRestComments": `SELECT c.id, u.name, u.img_url, c.text, c.rating FROM "comment" AS c JOIN "user" AS u ON c.user_id = u.id WHERE restaurant_id=$1 AND c.text IS NOT NULL AND c.text !=''`,
		"deleteComment": `DELETE FROM "comment" WHERE id=$1`,
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