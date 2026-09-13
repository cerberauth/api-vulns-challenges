module github.com/cerberauth/api-vulns-challenges/challenges/proxy-http2-authority-spoofing

go 1.26.0

require golang.org/x/net v0.59.0

require (
	github.com/spf13/cobra v1.10.2 // indirect
	golang.org/x/text v0.42.0 // indirect
)

require (
	github.com/cerberauth/api-vulns-challenges/common v0.0.0-00010101000000-000000000000
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
)

replace github.com/cerberauth/api-vulns-challenges/common => ../../common
