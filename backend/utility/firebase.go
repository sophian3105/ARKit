package utility

import (
	"context"
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
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

	if err != nil { return err }
	
	ctx := context.Background()
	opt := option.WithCredentialsFile(home + "")
	config := &firebase.Config{DatabaseURL: ""}
	app, err := firebase.NewApp(ctx, config, opt)

	if err != nil { return fmt.Errorf("error initializing app %v", err) }

	dbClient, dbErr := app.Database(ctx)
	authClient, authErr := app.Auth(ctx)

	if dbErr != nil || authErr != nil {
		return fmt.Errorf("error initializing clients: %v %v", dbErr, authErr)
	}

	fc.Database = dbClient
	fc.Auth = authClient
	return nil
}

func FBClient() *FirebaseClient {
	// Synchronize creation across requests
	mu.Lock()
	if fireDB == nil {
		tmp := &FirebaseClient{}
		err := tmp.connect() 
		
		if err == nil {
			fireDB = tmp
		}
	}
	mu.Unlock()

	return fireDB
}
