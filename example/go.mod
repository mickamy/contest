module github.com/mickamy/connecttest-example

go 1.25.3

replace github.com/mickamy/connecttest => ../

require (
	connectrpc.com/connect v1.19.1
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/google/uuid v1.6.0
	github.com/mickamy/connecttest v0.0.0
	github.com/mickamy/gokitx v0.0.2
	github.com/rs/cors v1.11.1
	golang.org/x/net v0.46.0
	google.golang.org/protobuf v1.36.10
)

require golang.org/x/text v0.30.0 // indirect
