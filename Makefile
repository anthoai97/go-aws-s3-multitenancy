GOBIN = ./build/bin
GORUN = go

server:
	$(GORUN) build -o ${GOBIN}/server .
	${GOBIN}/server

