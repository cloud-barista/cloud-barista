regenerate:
	# gogo proto compile
	protoc \
			./gogoproto/gogo.proto \
			-I . \
			-I $(GOPATH)/src/github.com/gogo/protobuf/protobuf \
			-I $(GOPATH)/src \
			--gofast_out=plugins=grpc,paths=source_relative,\
	Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:\
	.

	# cbmcks proto compile
	protoc \
			./cbmcks/cbmcks.proto \
			-I . \
			-I $(GOPATH)/src/github.com/gogo/protobuf/protobuf \
			-I $(GOPATH)/src/github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf \
			--gofast_out=plugins=grpc:\
	.	
