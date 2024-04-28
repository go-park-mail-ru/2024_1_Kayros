package route

import (
	"database/sql"

	dQuiz "2024_1_kayros/internal/delivery/statistic"
	"2024_1_kayros/internal/repository/minios3"
	rQuiz "2024_1_kayros/internal/repository/statistic"
	rUser "2024_1_kayros/internal/repository/user"
	uQuiz "2024_1_kayros/internal/usecase/statistic"
	uUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func AddQuizRouter(db *sql.DB, mc *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoQuiz := rQuiz.NewRepoLayer(db, logger)
	repoUser := rUser.NewRepoLayer(db)
	repoMinio := minios3.NewRepoLayer(mc)
	usecaseQuiz := uQuiz.NewUsecaseLayer(repoQuiz)
	usecaseUser := uUser.NewUsecaseLayer(repoUser, repoMinio)
	deliveryQuiz := dQuiz.NewDeliveryLayer(usecaseQuiz, usecaseUser, logger)

	mux.HandleFunc("/quiz/stats", deliveryQuiz.GetStatistic).Methods("GET").Name("quiz-stats")
	mux.HandleFunc("/quiz/questions", deliveryQuiz.GetQuestions).Methods("GET").Name("get-questions")
	mux.HandleFunc("/quiz/question/rating", deliveryQuiz.AddAnswer).Methods("POST").Name("add-question-rating")
}
