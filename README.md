# cv-manager-server-extension
Server for extension queries

# Dependences.
In order to run this project you need to execute the commands to obtain the required dependencies.

```bash
go get cloud.google.com/go/cloudsqlconn
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/stdlib
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option
```

# To run the server.

```bash
go run server.go
```