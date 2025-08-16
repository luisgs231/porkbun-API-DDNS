# Porkbun API DDNS Updater

Small tool to update DNS records on Porkbun using the API.

## How to use
- Fill in your API key and secret in the `config.yaml` file.
- Add some domains and existing records to update.
- (Optional) change the IP source if needed.
- Change the report endpoint.
- After each update, the program will send a GET request to the report endpoint with the new IP as the value of the `msg` parameter.  
  - For example, you can use [Uptime Kuma](https://github.com/louislam/uptime-kuma) with its **push monitor** feature as the report endpoint. This way, each successful update will be visible in Uptime Kuma, and you’ll also see the current IP as part of the status message.

**Note:** The `config.yaml` file must be placed in the same directory as the source file or the compiled binary.

## How to run
```bash
go run porkbunDDNS.go
```

## How to build
```bash
go build -o porkbunDDNS porkbunDDNS.go
```

On Windows:
```bash
go build -o porkbunDDNS.exe porkbunDDNS.go
```

You can also run:
```bash
build4All.bat
```
to automatically build binaries for multiple OS/ARCH combinations.
