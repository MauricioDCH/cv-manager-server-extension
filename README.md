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
go get github.com/joho/godotenv
```

If there are some error with dependences, please execute.
```bash
go mod tidy
```

# To run the server.


Before running the server make sure the database is working.

Then, change the .env_example file name to .env and introduce the credentials.

Then, execute the command below.
```bash
go run server.go
```

The server listens on port 5522.

If you whant to send some POST, you can execute in POSTMAN.

```bash
http://localhost:5522/submit-data
```

In the body of the post, use this example.
```json
{
    "inputs": [
        {
            "name": "email",
            "type": "text"
        },
        {
            "name": "password",
            "type": "pass"
        }
    ],
    "email": "ana.gomez@email.com"
}

```