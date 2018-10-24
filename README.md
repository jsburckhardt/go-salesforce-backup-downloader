# go-salesforce-backup-downloader

go-salesforce-backup-downloader is a command line app written in GO to download backup files from Salesforce Org's

> Some of the features:

- Uses Cobra to read config and generate commands
- Use fan-out (worker) concurrency pattern
- Stores the downloaded files into orgs name
- Generates a csv with the download resutls
- Uses 5 threads as default, but can be customized


[![Build status](https://juanburckhardt.visualstudio.com/go-salesforce-backup-downloader/_apis/build/status/go-salesforce-backup-downloader-Go%20(preview)-CI)](https://juanburckhardt.visualstudio.com/go-salesforce-backup-downloader/_build/latest?definitionId=1) [![Github Issues](http://githubbadges.herokuapp.com/jsburckhardt/go-salesforce-backup-downloader/issues.svg?style=flat-square)](https://github.com/jsburckhardt/go-salesforce-backup-downloader/issues) [![Pending Pull-Requests](http://githubbadges.herokuapp.com/jsburckhardt/go-salesforce-backup-downloader/pulls.svg?style=flat-square)](https://github.com/jsburckhardt/go-salesforce-backup-downloader/pulls) [![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org) 
[![Badges](http://img.shields.io/:badges-9/9-ff6799.svg?style=flat-square)](https://github.com/jsburckhardt/go-salesforce-backup-downloader)

---

## Table of Contents

- [Commands](#Commands)
- [Examples](#Examples)
- [License](#License)


## Commands:

- download
- numberOfFiles
- testCredentials
- help
- License

## Examples:
> download
```shell
> go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken 
> go-salesforce-backup-downloader.exe download --user sadmin@atyourcrazyorg --password mypasswordwithtoken
> go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken -m 5
> go-salesforce-backup-downloader.exe download -u sadmin@atyourcrazyorg -p mypasswordwithtoken --maxworkers 5
```
> numberOfFiles
```shell
> go-salesforce-backup-downloader.exe numberOfFiles -u sadmin@atyourcrazyorg -p mypasswordwithtoken
> go-salesforce-backup-downloader.exe numberOfFiles --user sadmin@atyourcrazyorg --password mypasswordwithtoken
```

> testCredentials
```shell
> go-salesforce-backup-downloader.exe numberOfFiles -u sadmin@atyourcrazyorg -p mypasswordwithtoken
> go-salesforce-backup-downloader.exe numberOfFiles --user sadmin@atyourcrazyorg --password mypasswordwithtoken
```

---
## License

[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

- **[MIT license](http://opensource.org/licenses/mit-license.php)**
