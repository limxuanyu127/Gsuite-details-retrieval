package controllers

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"

	//"log"
	//"net/http"
	//"os"

	"golang.org/x/net/context"
)

func CreateGsheetService(userEmail string) *sheets.Service {

	ctx := context.Background()

	serviceAccString:=os.Getenv("SERVICE_ACCOUNT_JSON")
	serviceAccBytes := []byte(serviceAccString)
	config, err := google.JWTConfigFromJSON(serviceAccBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Println(err)
	}

	//jsonCredentials, err := ioutil.ReadFile(os.Getenv("SERVICE_ACCOUNT_JSON"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/spreadsheets")
	//if err != nil {
	//	fmt.Println(err)
	//}

	config.Subject = userEmail
	ts := config.TokenSource(ctx)
	client:= oauth2.NewClient(ctx, ts)

	srv, err := sheets.New(client)
	if err != nil {
		fmt.Println("error creating new user")
		fmt.Println(err)
	}
	return srv
}

func ClearSheet(gsheet_srv *sheets.Service, spreadsheetId string) {

	var batchclearValueReq sheets.BatchClearValuesRequest
	batchclearValueReq.Ranges = []string{"A2:AO"}
	_, err := gsheet_srv.Spreadsheets.Values.BatchClear(spreadsheetId, &batchclearValueReq).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func AddToSheet(gsheet_srv *sheets.Service, spreadsheetId string, batchData *sheets.ValueRange) {

	ctx := context.Background()
	rangeToCheck:= batchData.Range
	resp, err := gsheet_srv.Spreadsheets.Values.Append(spreadsheetId, rangeToCheck, batchData).Context(ctx).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.HTTPStatusCode)
	fmt.Println("Batch req sent")
}


// For normal user account
//func CreateGsheetService2() *sheets.Service{
//	b, err := ioutil.ReadFile("authorisation/gsheet_credentials.json")
//	if err != nil {
//		log.Fatalf("Unable to read client secret file: %v", err)
//	}
//
//	// If modifying these scopes, delete your previously saved token.json.
//	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
//	if err != nil {
//		log.Fatalf("Unable to parse client secret file to config: %v", err)
//	}
//	client := getClient(config)
//
//	srv, err := sheets.New(client)
//	if err != nil {
//		log.Fatalf("Unable to retrieve Sheets client: %v", err)
//	}
//
//	return srv
//}




// Retrieve a token, saves the token, then returns the generated client.
//func getClient(config *oauth2.Config) *http.Client {
//	// The file token.json stores the user's access and refresh tokens, and is
//	// created automatically when the authorization flow completes for the first
//	// time.
//	tokFile := "token.json"
//	tok, err := tokenFromFile(tokFile)
//	if err != nil {
//		print(err)
//	}
//	return config.Client(context.Background(), tok)
//}
// Retrieves a token from a local file.
//func tokenFromFile(file string) (*oauth2.Token, error) {
//	f, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//	tok := &oauth2.Token{}
//	err = json.NewDecoder(f).Decode(tok)
//	return tok, err
//}



