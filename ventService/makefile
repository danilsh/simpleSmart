all: ventservice

ventservice: ..\..\..\..\..\bin\ventService.exe

..\..\..\..\..\bin\ventService.exe : main.go proto\ventService.pb.micro.go
	go build -o ..\..\..\..\..\bin\ventService.exe main.go

proto\ventService.pb.micro.go proto\ventService.pb.go : proto\ventService.proto
	protoc -I. --micro_out=. --go_out=. proto\ventService.proto

