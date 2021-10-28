export GOPATH=$(pwd)/../../..

go build -o ../tools/code_generator github.com/huoshan017/mysql-go/code_generator
go build -o ../tools/db_proxy_server github.com/huoshan017/mysql-go/proxy/server