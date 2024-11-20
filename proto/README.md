brew  install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

git clone https://github.com/googleapis/googleapis.git


export PATH="$PATH:$(go env GOPATH)/bin"

protoc -I . -I ./googleapis \
  --go_out . --go-grpc_out . \
  --grpc-gateway_out . \
  transit.proto


Explanation of Flags
-I .: Includes the current directory (where your .proto file is located).
-I ./googleapis: Includes the googleapis directory where annotations.proto and other proto files are located.
--go_out .: Generates the basic Go structures for your .proto file.
--go-grpc_out .: Generates the gRPC-specific code (server and client).
--grpc-gateway_out .: Generates the gRPC Gateway code for HTTP mapping.


Verifying the Generated Files
After running this command, you should see the following files (depending on your .proto definitions):

transit.pb.go: Contains Go structures and types.
transit_grpc.pb.go: Contains gRPC server and client interfaces.
transit.pb.gw.go: Contains the gRPC Gateway code for HTTP mappings.
