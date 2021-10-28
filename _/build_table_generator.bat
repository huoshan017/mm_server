go install mm_server_new/src/table_generator
if errorlevel 1 goto exit

go build -o ../bin/table_generator.exe mm_server_new/src/table_generator
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build table_generator failed!!!!!!!!!!!!!!!!!!!

:ok
echo build table_generator ok