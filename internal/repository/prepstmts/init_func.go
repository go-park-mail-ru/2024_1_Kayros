package prepstmts

import (
	stmtsFood "2024_1_kayros/internal/repository/food/stmts"
	stmtsOrder "2024_1_kayros/internal/repository/order/stmts"
	stmtsPromocode "2024_1_kayros/internal/repository/promocode/stmts"
	stmtsSearch "2024_1_kayros/internal/repository/search/stmts"
	stmtsStatistic "2024_1_kayros/internal/repository/statistic/stmts"
	"database/sql"

	"go.uber.org/zap"
)

func InitMonolithPreparedStatements(db *sql.DB, logger *zap.Logger) map[string]map[string]*sql.Stmt{
	foodStatements, err := stmtsFood.InitPrepareStatements(db)
	if err != nil {
		logger.Fatal("Cant' initialize foodStatements")
	}
	orderStatements, err := stmtsOrder.InitPrepareStatements(db)
	if err != nil {
		logger.Fatal("Cant' initialize orderStatements")
	}
	promocodeStatements, err := stmtsPromocode.InitPrepareStatements(db)
	if err != nil {
		logger.Fatal("Cant' initialize promocodeStatements")
	}
	searchStatements, err := stmtsSearch.InitPrepareStatements(db)
	if err != nil {
		logger.Fatal("Cant' initialize searchStatements")
	}
	statisticStatements, err := stmtsStatistic.InitPrepareStatements(db)
	if err != nil {
		logger.Fatal("Cant' initialize statisticStatements")
	}
	resultMap := make(map[string]map[string]*sql.Stmt, 5)

	resultMap["food"] = foodStatements
	resultMap["order"] = orderStatements
	resultMap["promocode"] = promocodeStatements
	resultMap["search"] = searchStatements
	resultMap["statistic"] = statisticStatements

	return resultMap
}