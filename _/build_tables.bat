go install mm_server_new/src/tables
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build tables failed!!!!!!!!!!!!!!!!!!!

:ok
echo build tables ok