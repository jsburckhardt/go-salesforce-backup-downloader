# go-salesforce-backup-downloader

go-salesforce-backup-downloader is a Golang command line app to download Salesforce Org's backup files

Some of the features:

* Uses Cobra to read config and generate commands
* Use fan-out (worker) concurrency pattern
* Stores the downloaded files into orgs name
* Generates a csv with the download resutls
* Uses 5 threads as default, but can be customized

Example:
go-salesforce-backup-downloader.exe -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe --user sadmin@atyourcrazyorg --password mypasswordwithtoken
