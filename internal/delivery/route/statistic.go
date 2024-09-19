package route

import (
	sessionproto "2024_1_kayros/gen/go/session"
	userproto "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	dQuiz "2024_1_kayros/internal/delivery/statistic"
	rQuiz "2024_1_kayros/internal/repository/statistic"
	ucSession "2024_1_kayros/internal/usecase/session"
	uQuiz "2024_1_kayros/internal/usecase/statistic"
	uUser "2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func AddQuizRouter(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, metrics *metrics.Metrics) {
	repoQuiz := rQuiz.NewRepoLayer(cluster.PsqlClient, metrics)
	// init grpc client interface
	grpcSessionClient := sessionproto.NewSessionManagerClient(clients.SessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient, metrics)
	// init grpc user interface
	grpcUserClient := userproto.NewUserManagerClient(clients.UserConn)
	usecaseUser := uUser.NewUsecaseLayer(grpcUserClient, metrics)

	usecaseQuiz := uQuiz.NewUsecaseLayer(repoQuiz)

	deliveryQuiz := dQuiz.NewDeliveryLayer(usecaseQuiz, usecaseUser, usecaseSession, logger, metrics)

	mux.HandleFunc("/api/v1/quiz/stats", deliveryQuiz.GetStatistic).Methods("GET").Name("quiz-stats")
	mux.HandleFunc("/api/v1/quiz/questions", deliveryQuiz.GetQuestions).Methods("GET").Name("get-questions")
	mux.HandleFunc("/api/v1/quiz/question/rating", deliveryQuiz.AddAnswer).Methods("POST").Name("add-question-rating")
}
