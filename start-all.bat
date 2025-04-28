@echo off
wt ^
new-tab -d .\frontend --title "Frontend" cmd /k "npm run dev" ; ^
new-tab -d .\backend\api-gateway --title "API Gateway" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-user --title "USER" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-email --title "EMAIL" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-redis --title "REDIS" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-media --title "MEDIA" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab --title "WSL" cmd /k "wsl"
