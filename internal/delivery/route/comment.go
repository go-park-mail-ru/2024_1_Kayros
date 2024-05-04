package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	delivery "2024_1_kayros/internal/delivery/comment"
	rUser "2024_1_kayros/internal/repository/user"
	uc "2024_1_kayros/internal/usecase/comment"
	comment "2024_1_kayros/microservices/comment/proto"
)

func AddCommentRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger, conn *grpc.ClientConn) {
	//repoComment := repo.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	grpcComment := comment.NewCommentWorkerClient(conn)
	ucComment := uc.NewUseCaseLayer(grpcComment, repoUser)
	handler := delivery.NewDelivery(ucComment, logger)

	mux.HandleFunc("/comment", handler.CreateComment).Methods("POST")
	mux.HandleFunc("/comments/{rest_id}", handler.GetComments).Methods("GET")
	mux.HandleFunc("/comment", handler.DeleteComment).Methods("DELETE")
}
