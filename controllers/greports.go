package controllers

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/reports/v1"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
)


func CreateReportService(userEmail string) *admin.Service {

	ctx := context.Background()

	//serviceAccString:=os.Getenv("SERVICE_ACCOUNT_JSON")
	serviceAccString:=os.Getenv("SERVICE_ACCOUNT_JSON")

	serviceAccBytes := []byte(serviceAccString)
	config, err := google.JWTConfigFromJSON(serviceAccBytes, "https://www.googleapis.com/auth/admin.reports.usage.readonly")
	if err != nil {
		fmt.Println(err)
	}

	//jsonCredentials, err := ioutil.ReadFile(os.Getenv("SERVICE_ACCOUNT_JSON"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//config, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/admin.reports.usage.readonly")
	//if err != nil {
	//	fmt.Println(err)
	//}

	config.Subject = userEmail
	ts := config.TokenSource(ctx)
	client:= oauth2.NewClient(ctx, ts)

	srv, err := admin.New(client)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	return srv
}


func GetNumLicense(report_srv *admin.Service, customerId string) (string, int, int){
	dt := time.Now().AddDate(0,0,-1)
	dateString := getDateString(dt)

	licenseResp, err:=report_srv.CustomerUsageReports.Get(dateString).CustomerId(customerId).Parameters("accounts:gsuite_unlimited_total_licenses,accounts:gsuite_unlimited_used_licenses").Do()
	if err != nil {
		log.Fatal(err)
	}
	numReports := len(licenseResp.UsageReports)
	for numReports <1{
		dt = dt.AddDate(0,0,-1)
		dateString = getDateString(dt)
		licenseResp, _=report_srv.CustomerUsageReports.Get(dateString).CustomerId(customerId).Parameters("accounts:gsuite_unlimited_total_licenses,accounts:gsuite_unlimited_used_licenses").Do()
		numReports= len(licenseResp.UsageReports)
	}

	// Once the DT with the latest report is found
	fmt.Println("DT of report is " + dateString)
	totalLicenses, usedLicenses := parseReport(licenseResp)
	return dateString, totalLicenses, usedLicenses
}


func getDateString(dt time.Time) string{
	year, month, day := dt.Date()
	yearString :=strconv.Itoa(year)
	monthString := strconv.Itoa(int(month))
	dayString := strconv.Itoa(day)
	if len(monthString) <2 {
		monthString = "0" + monthString
	}
	if len(dayString) <2 {
		dayString = "0" + dayString
	}

	dateString := yearString + "-" + monthString+ "-" + dayString
	return dateString
}


func parseReport(resp *admin.UsageReports) (int, int){

	var totalLicenses int
	var usedLicenses int

	usageReport := resp.UsageReports[0]
	reportParams := usageReport.Parameters
	fmt.Print("num of params is ")
	fmt.Println(len(reportParams))
	for i:=0; i<len(reportParams); i++{
		param := reportParams[i]
		if param.Name == "accounts:gsuite_unlimited_total_licenses"{
			totalLicenses = int(param.IntValue)
		}
		if param.Name == "accounts:gsuite_unlimited_used_licenses"{
			usedLicenses = int(param.IntValue)
		}
	}

	return totalLicenses, usedLicenses
}

func SendAlert(reportDate string, totalLicenses int, usedLicenses int, licenseThreshold int, sendgridApiKey string){

	if totalLicenses-usedLicenses < licenseThreshold{

		htmlContent := "<p>" + "Numbers accurate as of this date: " + reportDate + "<br>" +
			"Total number of licenses: " + strconv.Itoa(totalLicenses) + "<br>" +
			"Number of used licenses: " + strconv.Itoa(usedLicenses) + "<br>" +
			"<b>" + " Number of available licenses: " + strconv.Itoa(totalLicenses-usedLicenses) + "</b>" + "<br>" +
			"This email is generated because number of available licenses is less than " + strconv.Itoa(licenseThreshold) + "</p>"

		from := mail.NewEmail("Automated email from SendGrid", os.Getenv("EMAIL_SENDER"))
		subject := "Alert: Number of Gsuite Licenses running low (automated) "
		to := mail.NewEmail("Ninjavan User", os.Getenv("EMAIL_RECEIVER"))
		plainTextContent := htmlContent
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(sendgridApiKey)
		response, err := client.Send(message)
		if err != nil {
			fmt.Println("error sending email")
			log.Println(err)
		} else {
			fmt.Print("Email status: ")
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	}
}


