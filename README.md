# go-salesforce-backup-downloader

go-salesforce-backup-downloader is a Golang command line app to download Salesforce Org's backup files

> Some of the features:

- Uses Cobra to read config and generate commands
- Use fan-out (worker) concurrency pattern
- Stores the downloaded files into orgs name
- Generates a csv with the download resutls
- Uses 5 threads as default, but can be customized

## Commands:

> Examples:

- [download](#download)
- [numberOfFiles](#numberOfFiles)
- [testCredentials](#testCredentials)
- help

## download
go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken 
go-salesforce-backup-downloader.exe download --user sadmin@atyourcrazyorg --password mypasswordwithtoken
go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken -m 5
go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken --maxworkers 5

## numberOfFiles
go-salesforce-backup-downloader.exe numberOfFiles -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe numberOfFiles --user sadmin@atyourcrazyorg --password mypasswordwithtoken

## testCredentials
go-salesforce-backup-downloader.exe numberOfFiles -u sadmin@atyourcrazyorg -p mypasswordwithtoken
go-salesforce-backup-downloader.exe numberOfFiles --user sadmin@atyourcrazyorg --password mypasswordwithtoken

