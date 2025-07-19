## 直播
#protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./live/*.proto

## 站点
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./site/*.proto



##协议生成后，请务必执行脚本
go run remove_omitempty.go
