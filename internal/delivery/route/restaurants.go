package route

import (
	"github.com/gorilla/mux"

	"go.uber.org/zap"

	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/gen/go/user"
	dComment "2024_1_kayros/internal/delivery/comment"
	"2024_1_kayros/internal/delivery/metrics"
	dRest "2024_1_kayros/internal/delivery/restaurants"
	dSearch "2024_1_kayros/internal/delivery/search"
	rFood "2024_1_kayros/internal/repository/food"
	rSearch "2024_1_kayros/internal/repository/search"
	ucComment "2024_1_kayros/internal/usecase/comment"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
	ucSearch "2024_1_kayros/internal/usecase/search"
	ucUser "2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"
)

func AddRestRouter(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, metrics *metrics.Metrics) {
	repoSearch := rSearch.NewRepoLayer(cluster.PsqlClient, metrics)
	repoFood := rFood.NewRepoLayer(cluster.PsqlClient, metrics)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood)
	usecaseSearch := ucSearch.NewUsecaseLayer(repoSearch)
	// init user grpc client
	grpcUser := user.NewUserManagerClient(clients.UserConn)
	usecaseUser := ucUser.NewUsecaseLayer(grpcUser, metrics)

	//init rest grpc client
	grpcRest := rest.NewRestWorkerClient(clients.RestConn)
	usecaseRest := ucRest.NewUsecaseLayer(grpcRest, metrics)

	// init comment grpc client
	grpcComment := comment.NewCommentWorkerClient(clients.CommentConn)
	usecaseComment := ucComment.NewUseCaseLayer(grpcComment, grpcUser, metrics)

	deliveryRest := dRest.NewRestaurantHandler(usecaseRest, usecaseFood, usecaseUser, logger)
	deliveryComment := dComment.NewDelivery(usecaseComment, logger)
	deliverySearch := dSearch.NewDelivery(usecaseSearch, logger)

	mux.HandleFunc("/api/v1/search", deliverySearch.Search).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/api/v1/restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/api/v1/restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
	mux.HandleFunc("/api/v1/restaurants/{id}/comment", deliveryComment.CreateComment).Methods("POST").Name("create-comment")
	mux.HandleFunc("/api/v1/restaurants/{id}/comment/{com_id}", deliveryComment.DeleteComment).Methods("DELETE").Name("delete-comment")
	mux.HandleFunc("/api/v1/restaurants/{id}/comments", deliveryComment.GetComments).Methods("GET").Name("comments-list")
	mux.HandleFunc("/api/v1/category", deliveryRest.CategoryList).Methods("GET").Name("category-list")
	mux.HandleFunc("/api/v1/recomendation", deliveryRest.Recomendation).Methods("GET").Name("recomendation")
}
