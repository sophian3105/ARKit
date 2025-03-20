package routes

import (
	"aria/backend/database"
	"aria/backend/utility"
	"encoding/json"
	"net/http"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func userFrom(u *database.User) User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func GetUser(ctx *utility.Context) {
	ctx.WriteHeader(http.StatusOK)
}

func PostUser(ctx *utility.Context) {
	var user database.User

	if err := json.NewDecoder(ctx.Body).Decode(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest, "Invalid request")
		return
	}
	u, err := ctx.CreateUser(ctx.Context(), database.CreateUserParams(user))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError, "Error in create user")
	}
	ctx.Json(http.StatusCreated, userFrom(&u))

}
