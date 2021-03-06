call set_go_path.bat

cd../proto

md gen_go
cd gen_go

md rpc_message

cd ../../third_party/protobuf

move protoc.exe ../../proto
move protoc-gen-go.exe ../../proto

cd ../../proto
protoc.exe --go_out=./gen_go/rpc_message/ rpc_message.proto
cd ../_
if errorlevel 1 goto exit

cd ../proto
go install mm_server_new/proto/gen_go/rpc_message
cd ../_
if errorlevel 1 goto exit

cd ../proto
move protoc.exe ../third_party/protobuf
move protoc-gen-go.exe ../third_party/protobuf
cd ../_

goto ok

:exit
echo gen message failed!!!!!!!!!!!!!!!!!!!!!!!!!!!!

:ok
echo gen message ok