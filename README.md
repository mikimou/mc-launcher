# Custom sipmle minecraft launcher

![logo](https://i.imgur.com/ZTLvKhH.png)

## Installation
### Windows
*Download latest `launcher-windows.zip` from releases and run executable or build executable from code.
**If you are building from code your need to get minecraft files manually to be aable to run client!***

### Other OS
Other oses are not compatible yet.

## Configuration
Change nickname in username.txt file.

## Build 
```
go build launcher/launcher.go
```
```
go-winres patch launcher.exe
```

### Get minecraft files

```
python -m portablemc --work-dir . --main-dir . start forge:<version> --dry
```


## Screenshots
<img src="https://i.imgur.com/8nJu9Sj.png" width="500">
Screenshot from MacOS but launcher is still compatible only with Windows!

***

Michal Hicz
