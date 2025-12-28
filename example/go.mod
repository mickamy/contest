module github.com/mickamy/contest/example

go 1.25.3

replace github.com/mickamy/contest => ../

require (
	connectrpc.com/connect v1.19.1
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/mickamy/contest v0.0.0
	github.com/mickamy/gokitx v0.0.2
	github.com/rs/cors v1.11.1
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
	golang.org/x/net v0.46.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool go.uber.org/mock/mockgen
