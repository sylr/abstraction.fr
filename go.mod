module abstraction.fr

go 1.16

require (
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mileusna/useragent v1.0.2
	github.com/mitchellh/copystructure v1.1.2 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.10.0
	github.com/rivo/uniseg v0.2.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	sylr.dev/libqd/config v0.0.0-20210116225432-33b27958b842
	sylr.dev/yaml/v3 v3.0.0-20210127132132-941109e4f08c // indirect
)

replace (
	github.com/go-kit/kit/log => github.com/sylr/go-kit/log v0.0.0-20210102181225-11d0b12e814f
	github.com/gorilla/handlers => github.com/sylr/gorilla-handlers v1.4.3-0.20200522195821-d4f92d62f121
	github.com/prometheus/common => github.com/sylr/prometheus-common v0.2.1-0.20210102181937-d132ea8268d2
)
