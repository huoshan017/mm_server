call build_table_config.bat

go install mm_server_new/src/rpc_proto
go install mm_server_new/src/rpc_server
if errorlevel 1 goto exit

go build -i -o ../bin/rpc_server.exe mm_server_new/src/rpc_server
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build rpc_server failed!!!!!!!!!!!!!!!!!!!

:ok
echo build rpc_server ok