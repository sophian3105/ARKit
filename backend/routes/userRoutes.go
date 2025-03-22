package routes

import (
	"aria/backend/database"
	"aria/backend/utility"
	"net/http"
)

type User struct {
	ID    *string `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func userFrom(u *database.User) User {
	return User{
		ID:    &u.ID,
		Name:  &u.Name,
		Email: &u.Email,
	}
}

func GetUser(ctx *utility.Context) {
	// TODO finish this method
	uid := ctx.UID
	u, err := ctx.GetUser(ctx.Context(), uid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Json(http.StatusOK, userFrom(&u))
}

// PostUser creates a new user
// Must be authenticated
func PostUser(ctx *utility.Context) {
	uid := ctx.UID
	email := ctx.GetEmail()
	if email == nil {
		ctx.AbortWithStatus(http.StatusBadRequest, "No email found")
		return
	}

	createUserParams := database.CreateUserParams{
		ID:    uid,
		Email: *email,
	}

	u, err := ctx.CreateUser(ctx.Context(), createUserParams)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Json(http.StatusCreated, userFrom(&u))
}
