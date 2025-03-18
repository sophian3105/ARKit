package utility

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	Database *db.Client
	Auth *auth.Client
}

var mu sync.Mutex
var fireDB *FirebaseClient = nil
 
func (fc *FirebaseClient) connect() error {
	home, err := os.Getwd()
	dbUrl := os.Getenv("DB_URL")

	if err != nil { return err }
	if dbUrl == "" { return fmt.Errorf("empty DB_URL string") }
	
	ctx := context.Background()
	opt := option.WithCredentialsFile(home + "/firebase-adminsdk.json")
	config := &firebase.Config{DatabaseURL: dbUrl}
	app, err := firebase.NewApp(ctx, config, opt)

	if err != nil { return fmt.Errorf("error initializing app %v", err) }

	dbClient, dbErr := app.Database(ctx)
	authClient, authErr := app.Auth(ctx)

	if dbErr != nil || authErr != nil {
		return fmt.Errorf("error initializing clients: %v %v", authErr, dbErr)
	}

	fc.Database = dbClient
	fc.Auth = authClient
	return nil
}

// Gets a singleton instance of FirebaseClient
func FBClient() *FirebaseClient {
	mu.Lock()
	if fireDB == nil {
		tmp := &FirebaseClient{}
		err := tmp.connect() 
		
		if err == nil {
			// Successful so assign
			fireDB = tmp
		} else {
			// Should always connect
			log.Fatal(err)
		}
	}
	mu.Unlock()

	return fireDB
}

func AuthMiddleware() *Middleware {
	return NewMiddleware(func(ctx *Context) {
		client := FBClient()

		authHeader := ctx.Request.Header.Get("Authorization")
		authTokenString, found := strings.CutPrefix(authHeader, "Bearer ")

		if !found {
			ctx.AbortWithStatus(http.StatusUnauthorized, "No auth token found")
			return
		}

		authToken, err := client.Auth.VerifyIDToken(ctx.Context(), authTokenString)
		
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Token = authToken
	})
}
