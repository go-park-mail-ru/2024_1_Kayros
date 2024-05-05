package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	dComment "2024_1_kayros/internal/delivery/comment"
	dRest "2024_1_kayros/internal/delivery/restaurants"
	dSearch "2024_1_kayros/internal/delivery/search"
	rFood "2024_1_kayros/internal/repository/food"
	rSearch "2024_1_kayros/internal/repository/search"
	rUser "2024_1_kayros/internal/repository/user"
	ucComment "2024_1_kayros/internal/usecase/comment"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
	ucSearch "2024_1_kayros/internal/usecase/search"
	comment "2024_1_kayros/microservices/comment/proto"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

func AddRestRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger, restConn *grpc.ClientConn, commentConn *grpc.ClientConn) {
	//repoRest := rRest.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	repoSearch := rSearch.NewRepoLayer(db)
	repoFood := rFood.NewRepoLayer(db)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood)
	usecaseSearch := ucSearch.NewUsecaseLayer(repoSearch)

	grpcRest := rest.NewRestWorkerClient(restConn)
	usecaseRest := ucRest.NewUsecaseLayer(grpcRest)

	grpcComment := comment.NewCommentWorkerClient(commentConn)
	usecaseComment := ucComment.NewUseCaseLayer(grpcComment, repoUser)

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
