package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"getPromocodeByCode": `SELECT id, date, sale, type, restaurant_id, sum FROM promocode WHERE code=$1`,
		"wasPromocodeUsed": `SELECT count(*) FROM "order" WHERE user_id=$1 AND promocode_id=$2 AND status='delivered'`,
		"wasRestPromocodeUsed": `SELECT count(*) FROM "order" WHERE id=$1 AND promocode_id=$2`,
		"setPromocode": `UPDATE "order" SET promocode_id=$1 WHERE id=$2 RETURNING sum`,
		"deletePromocode": `UPDATE "order" SET promocode_id=NULL WHERE id=$1`,
		"getPromocodeIdFromOrder": `SELECT promocode_id FROM "order" WHERE id=$1`,
		"getPromocodeById": `SELECT id, code, date, sale, type, restaurant_id, sum FROM promocode WHERE id=$1`,
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
