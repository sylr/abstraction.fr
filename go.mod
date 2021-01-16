module abstraction.fr

go 1.16

require (
	github.com/Masterminds/sprig/v3 v3.2.0
	github.com/google/uuid v1.1.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mileusna/useragent v1.0.2
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/prometheus/client_golang v1.9.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20201231184435-2d18734c6014 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	sylr.dev/libqd/config v0.0.0-20210116114649-aacc5175aef8
)

replace (
	github.com/go-kit/kit/log => github.com/sylr/go-kit/log v0.0.0-20210102181225-11d0b12e814f
	github.com/gorilla/handlers => github.com/sylr/gorilla-handlers v1.4.3-0.20200522195821-d4f92d62f121
	github.com/prometheus/common => github.com/sylr/prometheus-common v0.2.1-0.20210102181937-d132ea8268d2
)
