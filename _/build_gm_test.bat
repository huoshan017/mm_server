go install mm_server_new/src/gm_test
if errorlevel 1 goto exit

go build -o ../bin/gm_test.exe mm_server_new/src/gm_test
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build gm_test failed!!!!!!!!!!!!!!!!!!!

:ok
echo build gm_test ok