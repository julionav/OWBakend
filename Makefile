GO_MODULE := github.com/julionav/OWBakend
GO_FLAGS := -ldflags '-s -w -extldflags "-static"'

PROTOS := $(wildcard protobuf/*.proto)
PBGO := $(patsubst protobuf/%.proto,server/%.pb.go,$(PROTOS))

server/%.pb.go: protobuf/%.proto
	protoc -I=protobuf/ --go_out=./ --go_opt=module=$(GO_MODULE) --experimental_allow_proto3_optional $?

all: $(PBGO)