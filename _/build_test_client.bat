go install mm_server_new/src/test_client
if errorlevel 1 goto exit

go build -o ../bin/test_client.exe mm_server_new/src/test_client
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build test_client failed!!!!!!!!!!!!!!!!!!!

:ok
echo build test_client ok