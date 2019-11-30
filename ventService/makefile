all: ventservice

ventservice: ..\..\bin\ventService.exe

..\..\bin\ventService.exe : main.go
	go build -o ..\..\bin\ventService.exe main.go

proto\ventService.pb.go : ventService.proto
	cd protoc
	protoc -I. --go_out=plugins=micro:. ventService.proto
	cd ..

