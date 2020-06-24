package main

import (
	"GsuiteRetrieval/controllers"
	"encoding/json"
	"fmt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"

)

func main() {

	switch job:= os.Getenv("JOB"); job{

	case "email":
		fmt.Println("Starting Email job")
		licenseThreshold,_ := strconv.Atoi(os.Getenv("LICENSE_THRESHOLD"))
		customerId := os.Getenv("CUSTOMER_ID")
		sendgridApiKey := os.Getenv("SENDGRID_API_KEY")



		report_srv := controllers.CreateReportService(os.Getenv("ADMIN_EMAIL"))
		reportDate, totalLicenses, usedLicenses := controllers.GetNumLicense(report_srv, customerId)
		if totalLicenses-usedLicenses < licenseThreshold{
			controllers.SendAlert(reportDate, totalLicenses, usedLicenses, licenseThreshold, sendgridApiKey)
		}

	case "sheets":
		fmt.Println("Starting Gsheets job")
		gsheet_srv := controllers.CreateGsheetService(os.Getenv("ADMIN_EMAIL"))
		gsuite_srv := controllers.CreateGsuiteService(os.Getenv("ADMIN_EMAIL"))
		//report_srv := controllers.CreateReportService(os.Getenv("ADMIN_EMAIL"))

		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		//customerId := os.Getenv("CUSTOMER_ID")
		domainName:= os.Getenv("DOMAIN_NAME")
		batchSize := os.Getenv("BATCH_SIZE")
		sheetTabName := os.Getenv("TAB_NAME")
		batchSizeInt,_ := strconv.Atoi(batchSize)

		// Variables for loop below
		var vr sheets.ValueRange
		var rowFormula string
		rowCounter := 2 // row at which u start appending
		batchCounter := 0 //counter to keep track of # of values in vr before appending as a batch


		// Clear Sheet
		controllers.ClearSheet(sheetTabName,  gsheet_srv, spreadsheetId)

		// Get first 500 users
		call := gsuite_srv.Users.List().Domain(domainName).MaxResults(500)
		allUsers, err := call.Do()
		nextPageToken := allUsers.NextPageToken
		if err != nil {
			log.Fatal(err)
		}

		for true{

			fmt.Println("fetching")
			for i := 0; i < len(allUsers.Users); i++ {

				rowFormula = sheetTabName + "!" + "A" + strconv.Itoa(rowCounter)
				target := allUsers.Users[i]
				fmt.Println("Building user")

				// Get all fields of target user
				primaryEmail := target.PrimaryEmail
				firstName := target.Name.GivenName
				lastName := target.Name.FamilyName
				isAdmin := target.IsAdmin
				isDelegatedAdmin := target.IsDelegatedAdmin
				isEnrolledIn2Sv := target.IsEnrolledIn2Sv
				isEnforcedIn2Sv := target.IsEnforcedIn2Sv
				orgUnitPath := target.OrgUnitPath
				suspended := target.Suspended
				suspensionReason := target.SuspensionReason

				// Fields that require datetime converino
				creationTime := target.CreationTime
				creationTime = convertToDt(creationTime)
				lastLoginTime := target.LastLoginTime
				lastLoginTime = convertToDt(lastLoginTime)

				// Fields that when extracted from admin.User, are interfacs{} , ie. probably returning array
				phoneIntf := target.Phones
				phoneValue, phonePrimary, phoneType := getPhoneInfo(phoneIntf)
				relationsIntf := target.Relations
				rel0value, rel0type := getRelInfo(relationsIntf)
				organizaionIntf := target.Organizations
				org0map, org1map := getOrgInfo(organizaionIntf)
				org0name := org0map["org0name"]
				org0title := org0map["org0title"]
				org0primary := org0map["org0primary"]
				org0custom := org0map["org0custom"]
				org0dept := org0map["org0dept"]
				org0desc := org0map["org0desc"]
				org0location := org0map["org0location"]
				org0symbol := org0map["org0symbol"]
				org0domain := org0map["org0domain"]
				org0costCenter := org0map["org0costCenter"]
				org1name := org1map["org1name"]
				org1title := org1map["org1title"]
				//org1primary := org1map["primary"]
				org1custom := org1map["org1custom"]
				org1dept := org1map["org1dept"]
				org1desc := org1map["org1desc"]
				org1location := org1map["org1location"]
				org1symbol := org1map["org1symbol"]
				org1domain := org1map["org1domain"]
				org1costCenter := org1map["org1costCenter"]

				// Adding a user field to vr
				var myval []interface{}
				myval = []interface{}{primaryEmail, firstName, lastName, nil, nil, isAdmin, isDelegatedAdmin, isEnrolledIn2Sv, isEnforcedIn2Sv, creationTime, lastLoginTime, orgUnitPath, suspended, suspensionReason, nil, org0title, org0primary, org0custom, org0dept, org0desc, org0costCenter, rel0value, rel0type, org0name, phoneValue, phonePrimary, phoneType, org0location, org0symbol, org0domain, org1name, org1title, org1custom, org1dept, org1symbol, org1location, org1desc, org1domain, org1costCenter}
				//if i==0{
				//	// Add license data to first row
				//	myval = append(myval, totalLicenses, usedLicenses)
				//}
				vr.Values = append(vr.Values, myval)
				vr.MajorDimension = "ROWS"
				if vr.Range == ""{
					vr.Range = rowFormula
				}

				// When batch size is reached, add to gsheets and update counters(reset vr and batchcounter)
				if batchCounter == batchSizeInt {
					controllers.AddToSheet(gsheet_srv, spreadsheetId, &vr)
					batchCounter = 0
					vr.Values= vr.Values[:0]
					vr.Range=""
				}
				// Increment row and batch counters (building up vr, before adding to sheet)
				rowCounter += 1
				batchCounter += 1
			}
			// Get next page of users
			allUsers, nextPageToken = LoopGetUsers(gsuite_srv, domainName, nextPageToken)
			if allUsers == nil{
				// Add remaining values that weren't enough to cause a batch insert
				if len(vr.Values) > 0 {
					fmt.Println("Final add to sheet")
					controllers.AddToSheet(gsheet_srv, spreadsheetId, &vr)
				}
				break
			}
		}


	}
}


func LoopGetUsers(gsuite_srv *admin.Service, domainName string, nextPageToken string ) (*admin.Users, string){

	if nextPageToken == ""{
		return nil, ""
	}else{
		call := gsuite_srv.Users.List().Domain(domainName).MaxResults(500).PageToken(nextPageToken)
		allUsers, err := call.Do()
		nextPageToken = allUsers.NextPageToken
		if err != nil {
			log.Fatal(err)
		}
		return allUsers, nextPageToken
	}
}

func getPhoneInfo(phoneIntf interface{}) (interface{}, interface{}, interface{}) {
	var phoneValue interface{}
	var phonePrimary interface{}
	var phoneType interface{}
	if phoneIntf != nil {
		//var phoneMaps []map[string]interface{}
		var phoneMap []map[string]interface{}
		m, _ := json.Marshal(phoneIntf.([]interface{}))
		err := json.Unmarshal(m, &phoneMap)
		if err != nil {
			log.Fatal(err)
		}
		// Only take first phone
		phone := phoneMap[0]
		phoneValue = phone["value"]
		phonePrimary = phone["primary"]
		phoneType = phone["type"]
	}

	return phoneValue, phonePrimary, phoneType
}

func getRelInfo(relationsIntf interface{}) (interface{}, interface{}) {
	var rel0value interface{}
	var rel0type interface{}

	if relationsIntf != nil {
		var relationsMap []map[string]interface{}
		r, _ := json.Marshal(relationsIntf.([]interface{}))
		err := json.Unmarshal(r, &relationsMap)
		if err != nil {
			log.Fatal(err)
		}
		rel0 := relationsMap[0]
		rel0value = rel0["value"]
		rel0type = rel0["type"]
	}

	return rel0value, rel0type
}

func getOrgInfo(organizaionIntf interface{}) (map[string]interface{}, map[string]interface{}) {

	org0map := map[string]interface{}{
		"org0name":       "",
		"org0title":      "",
		"org0primary":    "",
		"org0custom":     "",
		"org0dept":       "",
		"org0desc":       "",
		"org0location":   "",
		"org0symbol":     "",
		"org0domain":     "",
		"org0costCenter": "",
	}

	org1map := map[string]interface{}{
		"org1name":       "",
		"org1title":      "",
		"org1primary":    "",
		"org1custom":     "",
		"org1dept":       "",
		"org1desc":       "",
		"org1location":   "",
		"org1symbol":     "",
		"org1domain":     "",
		"org1costCenter": "",
	}

	if organizaionIntf != nil {
		var organizationMap []map[string]interface{}
		m, _ := json.Marshal(organizaionIntf.([]interface{}))
		err := json.Unmarshal(m, &organizationMap)
		if err != nil {
			log.Fatal(err)
		}

		var org0 map[string]interface{}
		var org1 map[string]interface{}
		if len(organizationMap) >= 1 {
			org0 = organizationMap[0]

			org0map["org0name"] = org0["name"]
			org0map["org0title"] = org0["title"]
			org0map["org0primary"] = org0["primary"]
			org0map["org0custom"] = org0["custom"]
			org0map["org0dept"] = org0["department"]
			org0map["org0desc"] = org0["description"]
			org0map["org0location"] = org0["location"]
			org0map["org0symbol"] = org0["symbol"]
			org0map["org0domain"] = org0["domain"]
			org0map["org0costCenter"] = org0["costcenter"]

		}
		if len(organizationMap) >= 2 {
			org1 = organizationMap[1]

			org1map["org1name"] = org1["name"]
			org1map["org1title"] = org1["title"]
			org1map["org1primary"] = org1["primary"]
			org1map["org1custom"] = org1["custom"]
			org1map["org1dept"] = org1["department"]
			org1map["org1desc"] = org1["description"]
			org1map["org1location"] = org1["location"]
			org1map["org1symbol"] = org1["symbol"]
			org1map["org1domain"] = org1["domain"]
			org1map["org1costCenter"] = org1["costcenter"]
		}
	}
	return org0map, org1map
}

func convertToDt( s string) string{
	dateString := s[:10]
	timeString := s[12:19]
	return dateString + " " + timeString

}
