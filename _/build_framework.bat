call gen_server_message.bat
if errorlevel 1 goto exit

call gen_client_message.bat
if errorlevel 1 goto exit

go install mm_server_new/libs/log
if errorlevel 1 goto exit

go install mm_server_new/libs/timer
if errorlevel 1 goto exit

go install mm_server_new/libs/socket
if errorlevel 1 goto exit

go install mm_server_new/libs/server_conn
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build framework failed!!!!!!!!!!!!!!!!!!

:ok
echo build framework ok