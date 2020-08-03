# Retrieval of employee's Gsuite account details

The objective of this script is to retrieve all the details of all Gsuite users in a given domain. This requested information will be updated on a spreadsheet.

The second function of this script is to check if there are sufficient Gsuite licenses in the domain.
<hr style="border:1px solid gray">

### Controllers Directory
 
Each file in this directory is responsible for making API request to google platforms. A service account key (json) is required to get authentication. 

The 3 services that we will make API requests to are the Reports API (under Gsuite SDK), Directory API (under Gsuite SDK), Gsheet API
<hr style="border:1px solid gray">

### main.go

**Function 1: Retrieve information**
- Fetch all Gsuite users using Gsuite SDK Api (directory)
- Consolidate a user's info in an array of  _interface{}_
- Consolidate all users info in an array of these arrays 
- Update the Gsheet with these info, batch by batch

**Function 2: Check number of licenses**
- Use Reports API to fetch number of license (can be normal or business account)
- Use SendGrid API to send an email to an admin if the license count is below a threshold