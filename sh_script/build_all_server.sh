#!/bin/bash
export GOPATH=$(pwd)/../../go_projects
set -x
#go mod vendor 
go build -mod=mod -o ../bin/center_server ../src/center_server
go build -mod=mod -o ../bin/login_server ../src/login_server
go build -mod=mod -o ../bin/game_server ../src/game_server
go build -mod=mod -o ../bin/rpc_server ../src/rpc_server
go build -mod=mod -o ../bin/test_client ../src/test_client
