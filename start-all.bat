@echo off
wt ^
new-tab -d .\backend\service-user --title "USER" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-notification --title "NOTIFICATION" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-community --title "COMMUNITY" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-thread --title "THREAD" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-media --title "MEDIA" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-message --title "MESSAGE" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\util-email --title "EMAIL" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\util-redis --title "REDIS" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\backend\service-ai --title "AI" cmd /k "python app.py" ; ^
new-tab -d .\backend\api-gateway --title "GATEWAY" cmd /k "go build -o main.exe && main.exe" ; ^
new-tab -d .\frontend --title "Frontend" cmd /k "npm run dev"
