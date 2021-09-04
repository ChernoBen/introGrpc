#1
- Install protoc
    - pamac install protobuf || sudo pacman -S protobuf
    - more information about : https://grpc.io/docs/languages/go/quickstart/
    - protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.
