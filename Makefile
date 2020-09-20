export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on

OBJ = tcpdup-server

default: $(OBJ)

$(OBJ):
	go build -gcflags "-N -l" -o $@ ./src

clean:
	rm -fr $(OBJ)

-include .deps
