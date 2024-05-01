package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	delivery "2024_1_kayros/internal/delivery/comment"
	repo "2024_1_kayros/internal/repository/comment"
	rUser "2024_1_kayros/internal/repository/user"
	uc "2024_1_kayros/internal/usecase/comment"
)

func AddCommentRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger) {
	repoComment := repo.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	ucComment := uc.NewUseCaseLayer(repoComment, repoUser)
	handler := delivery.NewCommentHandler(ucComment, logger)

	mux.HandleFunc("/comment", handler.CreateComment).Methods("POST")
	mux.HandleFunc("/comments/{rest_id}", handler.GetComments).Methods("GET")
	mux.HandleFunc("/comment", handler.DeleteComment).Methods("DELETE")
}
