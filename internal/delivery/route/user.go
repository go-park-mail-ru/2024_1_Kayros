package route

import (
	"database/sql"

	dUser "2024_1_kayros/internal/delivery/user"

	rUser "2024_1_kayros/internal/repository/user"

	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
)

func AddUserRouter(db *sql.DB, mux *mux.Router) {
	repoUser := rUser.NewRepoLayer(db)
	usecaseUser := ucUser.NewUsecaseLayer(repoUser)
	deliveryUser := dUser.NewDeliveryLayer(usecaseUser)

	mux.HandleFunc("user", deliveryUser.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("user", deliveryUser.UpdateImage).Methods("PUT").Name("user_image")
}
