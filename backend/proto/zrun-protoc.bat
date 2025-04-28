@echo off
setlocal enabledelayedexpansion

set PROTO_DIR=.
set API_GATEWAY_DIR=..\api-gateway\proto
set SERVICE_DIR=..\service



for %%F in (*.proto) do (
    set "FILENAME=%%~nF"
    set "PROTO_FILE=%%F"

    echo Processing !FILENAME!.proto...

    protoc --go_out=. --go-grpc_out=. --proto_path=. !PROTO_FILE!
    
    if not exist "%API_GATEWAY_DIR%\!FILENAME!\" (
        mkdir "%API_GATEWAY_DIR%\!FILENAME!\"
    )

    copy /Y "!FILENAME!\!FILENAME!.pb.go" "%API_GATEWAY_DIR%\!FILENAME!\" 
    copy /Y "!FILENAME!\!FILENAME!_grpc.pb.go" "%API_GATEWAY_DIR%\!FILENAME!\" 

    set "SERVICE_PROTO_DIR=%SERVICE_DIR%-!FILENAME!\proto\!FILENAME!"

    if not exist "!SERVICE_PROTO_DIR!" (
        mkdir "!SERVICE_PROTO_DIR!"
    )

    copy /Y "!FILENAME!\!FILENAME!.pb.go" "!SERVICE_PROTO_DIR!\" 
    copy /Y "!FILENAME!\!FILENAME!_grpc.pb.go" "!SERVICE_PROTO_DIR!\" 
)

echo.
echo Successfully generated and copied gRPC files.
