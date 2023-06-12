
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

// SAMPLE CONFIG .env, you should put the actual config details found on your project settings

VITE_FIREBASE_API_KEY=AIzaKJgkjhSdfSgkjhdkKJdkjowf
VITE_FIREBASE_AUTH_DOMAIN=yourauthdomin.firebaseapp.com
VITE_FIREBASE_DB_URL=https://yourdburl.firebaseio.com
VITE_FIREBASE_PROJECT_ID=yourproject-id
VITE_FIREBASE_STORAGE_BUCKET=yourstoragebucket.appspot.com
VITE_FIREBASE_MSG_SENDER_ID=43597918523958
VITE_FIREBASE_APP_ID=234598789798798fg3-034

2. Run the MySQL server

```
Can be used the Xampp, Apache or another things

```

3. In the root directory, please run 

```
go mod init and go mod tidy
```

2. Run 

```
go run main.go

```
