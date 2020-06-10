package controllers

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"os"
)



//    user_email: The email of the user. Needs permissions to access the Admin APIs.
func CreateGsuiteService(userEmail string) *admin.Service {

	ctx := context.Background()

	serviceAccString:=os.Getenv("SERVICE_ACCOUNT_JSON")
	serviceAccBytes := []byte(serviceAccString)
	config, err := google.JWTConfigFromJSON(serviceAccBytes, admin.AdminDirectoryUserScope)
	if err != nil {
		fmt.Println(err)
	}

	//jsonCredentials, err := ioutil.ReadFile(os.Getenv("SERVICE_ACCOUNT_JSON"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//config, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryUserScope)
	//if err != nil {
	//	fmt.Println(err)
	//}

	config.Subject = userEmail
	ts := config.TokenSource(ctx)
	client:= oauth2.NewClient(ctx, ts)

	srv, err := admin.New(client)
	if err != nil {
		fmt.Println("error creating new user")
		fmt.Println(err)
	}
	return srv
}
