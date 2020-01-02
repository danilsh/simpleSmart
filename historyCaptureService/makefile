all: label historycaptureservice

result = ../bin/historyCaptureService
source = main.go init.go dataAquisition.go

label:
	@echo "***** Building historycaptureservice *****"

historycaptureservice: $(result)

$(result) : $(source)
	go mod tidy
	go build -o $(result) $(source)
