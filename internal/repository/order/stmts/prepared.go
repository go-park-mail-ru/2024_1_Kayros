package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"createOrder": `INSERT INTO "order" (user_id, created_at, updated_at, status) VALUES ($1, $2, $3, $4) RETURNING id`,
		"createNoUnauth": `INSERT INTO "order" (unauth_token, created_at, updated_at, status) VALUES ($1, $2, $3, $4) RETURNING id`,
		"getOrderByStatus": `SELECT id, user_id, order_created_at, status, address, extra_address, sum FROM "order" WHERE user_id= $1 AND status=$2`,
		"getBasketNoAuth": `SELECT id, created_at, updated_at, received_at, status, address, extra_address, sum FROM "order" WHERE unauth_token= $1 AND status=$2`,
		"getOrderById": `SELECT id, user_id, order_created_at, delivered_at, status, address, extra_address, sum, commented FROM "order" WHERE id= $1`,
		"": `SELECT id FROM "order" WHERE user_id= $1 AND status=$2`,
		"getBasketIdNoAuth": `SELECT id FROM "order" WHERE unauth_token=$1 AND status=$2`,
		"orderFood": `SELECT f.id, f.name, f.weight, f.price, fo.count, f.img_url, f.restaurant_id
		FROM food_order AS fo
		JOIN food AS f ON fo.food_id = f.id
		WHERE fo.order_id = $1`,
		"updateOrderAddress": `UPDATE "order" SET address=$1, extra_address=$2 WHERE id=$3 RETURNING id`,
		"updateOrderCreatedStatus": `UPDATE "order" SET status=$1, order_created_at=$2 WHERE id=$3 RETURNING id`,
		"udpateOrderDeliveredStatus": `UPDATE "order" SET status=$1, delivered_at=$2 WHERE id=$3 RETURNING id`,
		"updateOrderStatus": `UPDATE "order" SET status=$1 WHERE id=$2 RETURNING id`,
		"getOrderSum": `SELECT sum FROM "order" WHERE id=$1`,
		"getFoodPrice": `SELECT price FROM food WHERE id=$1`,
		"getFoodCountInOrder": `SELECT count FROM food_order WHERE order_id=$1 AND food_id=$2`,
		"updateSumOrder": `UPDATE "order" SET sum=$1 WHERE id=$2`,
		"addFoodToOrder": `INSERT INTO food_order (order_id, food_id, count,  created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		"updateFoodInOrder": `UPDATE food_order SET count=$1 WHERE order_id=$2 AND food_id=$3`,
		"deleteFoodFromOrder": `DELETE FROM food_order WHERE order_id=$1 AND food_id=$2`,
		"deleteOrder": `DELETE FROM "order" WHERE id=$1`,
		"cleanOrder": `DELETE FROM food_order WHERE order_id=$1`,
		"setOrderUser": `UPDATE "order" SET user_id=$1, unauth_token=NULL WHERE id=$2`,
		"getUserOrders": `SELECT count(*) FROM "order" WHERE user_id=$1 AND status=$2`,
		"getPromocode": `SELECT id, date, sale, type, restaurant_id, sum FROM promocode WHERE code=$1`,
		"wasPromocodeUsed": `SELECT count(*) FROM "order" WHERE user_id=$1 AND promocode_id=$2 AND status='delivered'`,
		"wasRestPromocodeUsed": `SELECT count(*) FROM "order" WHERE id=$1 AND promocode_id=$2`,
		"setPromocode": `UPDATE "order" SET promocode_id=$1 WHERE id=$2 RETURNING sum`,
		"deletePromocode": `UPDATE "order" SET promocode_id=NULL WHERE id=$1`,
		"getPromocodeId": `SELECT promocode_id FROM "order" WHERE id=$1`,
		"getPromocodeById": `SELECT id, code, date, sale, type, restaurant_id, sum FROM promocode WHERE id=$1`,
		"getActivePromocode": `SELECT id, code, date, sale, type, sum FROM promocode WHERE date > CURRENT_DATE`,
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
