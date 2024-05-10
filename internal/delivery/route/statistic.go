package route

import (
	"database/sql"

	"2024_1_kayros/config"
	sessionproto "2024_1_kayros/gen/go/session"
	userproto "2024_1_kayros/gen/go/user"
	dQuiz "2024_1_kayros/internal/delivery/statistic"
	rQuiz "2024_1_kayros/internal/repository/statistic"
	ucSession "2024_1_kayros/internal/usecase/session"
	uQuiz "2024_1_kayros/internal/usecase/statistic"
	uUser "2024_1_kayros/internal/usecase/user"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func AddQuizRouter(db *sql.DB, sessionConn, userConn *grpc.ClientConn, mc *minio.Client, mux *mux.Router, logger *zap.Logger, cfg *config.Project) {
	repoQuiz := rQuiz.NewRepoLayer(db)
	// init grpc client interface 
	grpcSessionClient := sessionproto.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient)
	// init grpc user interface 
	grpcUserClient := userproto.NewUserManagerClient(userConn)
	usecaseUser := uUser.NewUsecaseLayer(grpcUserClient)

	usecaseQuiz := uQuiz.NewUsecaseLayer(repoQuiz)

	deliveryQuiz := dQuiz.NewDeliveryLayer(usecaseQuiz, usecaseUser, usecaseSession, logger, cfg)

	mux.HandleFunc("/quiz/stats", deliveryQuiz.GetStatistic).Methods("GET").Name("quiz-stats")
	mux.HandleFunc("/quiz/questions", deliveryQuiz.GetQuestions).Methods("GET").Name("get-questions")
	mux.HandleFunc("/quiz/question/rating", deliveryQuiz.AddAnswer).Methods("POST").Name("add-question-rating")
}
