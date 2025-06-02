@echo off
wt ^
new-tab --title "WSL" cmd /k "wsl" ; ^
new-tab -d .\backend\service-user --title "USER" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-notification --title "NOTIFICATION" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-thread --title "THREAD" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-media --title "MEDI" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\util-email --title "EMAIL" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\util-redis --title "REDIS" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\api-gateway --title "GATEWAY" cmd /k "go build -o main.exe && main.exe" 
new-tab -d .\frontend --title "Frontend" cmd /k "npm run dev" ; ^
