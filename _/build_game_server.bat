call build_framework.bat
if errorlevel 1 goto exit

call build_tables.bat
if errorlevel 1 goto exit

go install mm_server_new/src/rpc_proto
go install mm_server_new/src/game_server
if errorlevel 1 goto exit

go build -i -o ../bin/game_server.exe mm_server_new/src/game_server
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build game_server failed!!!!!!!!!!!!!!!!!!!

:ok
echo build game_server ok