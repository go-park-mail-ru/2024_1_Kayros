package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/gen/go/user"
	dComment "2024_1_kayros/internal/delivery/comment"
	dRest "2024_1_kayros/internal/delivery/restaurants"
	dSearch "2024_1_kayros/internal/delivery/search"
	rFood "2024_1_kayros/internal/repository/food"
	rSearch "2024_1_kayros/internal/repository/search"
	ucComment "2024_1_kayros/internal/usecase/comment"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
	ucSearch "2024_1_kayros/internal/usecase/search"
)

func AddRestRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger, restConn, userConn, commentConn *grpc.ClientConn) {
	repoSearch := rSearch.NewRepoLayer(db)
	repoFood := rFood.NewRepoLayer(db)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood)
	usecaseSearch := ucSearch.NewUsecaseLayer(repoSearch)
	// init user grpc client
	grpcUser := user.NewUserManagerClient(userConn)

	//init rest grpc client
	grpcRest := rest.NewRestWorkerClient(restConn)
	usecaseRest := ucRest.NewUsecaseLayer(grpcRest)

	// init comment grpc client
	grpcComment := comment.NewCommentWorkerClient(commentConn)
	usecaseComment := ucComment.NewUseCaseLayer(grpcComment, grpcUser)

	deliveryRest := dRest.NewRestaurantHandler(usecaseRest, usecaseFood, logger)
	deliveryComment := dComment.NewDelivery(usecaseComment, logger)
	deliverySearch := dSearch.NewDelivery(usecaseSearch, logger)

	mux.HandleFunc("/search", deliverySearch.Search).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
	mux.HandleFunc("/restaurants/{id}/comment", deliveryComment.CreateComment).Methods("POST").Name("create-comment")
	mux.HandleFunc("/restaurants/{id}/comment/{com_id}", deliveryComment.DeleteComment).Methods("DELETE").Name("delete-comment")
	mux.HandleFunc("/restaurants/{id}/comments", deliveryComment.GetComments).Methods("GET").Name("comments-list")
	mux.HandleFunc("/category", deliveryRest.CategoryList).Methods("GET").Name("category-list")
}
