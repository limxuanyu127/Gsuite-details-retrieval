package controllers

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
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

func ClearSheet(sheetTabName string, gsheet_srv *sheets.Service, spreadsheetId string) {

	var batchclearValueReq sheets.BatchClearValuesRequest
	batchclearValueReq.Ranges = []string{sheetTabName+ "!" + "A2:AO"}
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



