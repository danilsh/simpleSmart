all: historycaptureservice

result = ../bin/historyCaptureService
source = main.go init.go dataAquisition.go

historycaptureservice: $(result)

$(result) : $(source)
	go mod tidy
	go build -o $(result) $(source)


