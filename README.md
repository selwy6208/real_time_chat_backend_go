
## Getting Started

1. Create an `.env` file and add the following variables.

DB_HOST=127.0.0.1
DB_DRIVER=mysql
DB_USER=root
DB_PASSWORD=
DB_NAME=real-chat-backend-db
DB_PORT=3306
TOKEN_HOUR_LIFESPAN=1
API_SECRET=My_API_Secret_Key

2. Run the MySQL server

```
Can be used the Xampp, navicat or another things

```

3. In the root directory, please run 

```
go mod init and go mod tidy
```

2. Run 

```
go run main.go

```
