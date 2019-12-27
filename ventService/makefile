all: ventservice

result = ../bin/ventService
rpc_interface = ./proto/ventService.pb.go proto/ventService.pb.micro.go
proto = ./proto/ventService.proto
source = main.go init.go ventServiceImpl.go dataAquisition.go dataCalc.go

ventservice: $(result)

$(result) : $(source) $(rpc_interface)
	go mod tidy
	go build -o $(result) $(source)

$(rpc_interface) : $(proto)
	protoc -I. --micro_out=. --go_out=. $(proto)

