1. Run Powershell in backend folder, insert command: go run main.go
2. Keep it running
3. Run another Powershell in frontend folder, insert command: "npm install", then "npm run serve"
4. Should be able to run on port 8081

PostgreSQL is needed for database.
if database can't be created through backend, proceed to go into pgAdmin and execute "CREATE DATABASE wira_db" and let the .go code do the rest.

