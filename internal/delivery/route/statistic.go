package route

import (
	"database/sql"

	dQuiz "2024_1_kayros/internal/delivery/statistic"
	"2024_1_kayros/internal/repository/minios3"
	rSession "2024_1_kayros/internal/repository/session"
	rQuiz "2024_1_kayros/internal/repository/statistic"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	uQuiz "2024_1_kayros/internal/usecase/statistic"
	uUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AddQuizRouter(db *sql.DB, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, mc *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoQuiz := rQuiz.NewRepoLayer(db, logger)
	repoUser := rUser.NewRepoLayer(db)
	repoMinio := minios3.NewRepoLayer(mc)
	repoSession := rSession.NewRepoLayer(clientRedisSession)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf)

	usecaseQuiz := uQuiz.NewUsecaseLayer(repoQuiz)
	usecaseUser := uUser.NewUsecaseLayer(repoUser, repoMinio)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrf, logger)

	deliveryQuiz := dQuiz.NewDeliveryLayer(usecaseQuiz, usecaseUser, usecaseCsrf, usecaseSession, logger)

	mux.HandleFunc("/quiz/stats", deliveryQuiz.GetStatistic).Methods("GET").Name("quiz-stats")
	mux.HandleFunc("/quiz/questions", deliveryQuiz.GetQuestions).Methods("GET").Name("get-questions")
	mux.HandleFunc("/quiz/question/rating", deliveryQuiz.AddAnswer).Methods("POST").Name("add-question-rating")
}
