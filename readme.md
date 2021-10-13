#1
- Install protoc
    - pamac install protobuf || sudo pacman -S protobuf
    - more information about : https://grpc.io/docs/languages/go/quickstart/
    - protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.
    - Add these 2 important lines to your ~/.bash_profile:

    - export GO_PATH=~/go
    - export PATH=$PATH:/$GO_PATH/bin
    - chmod +x run.sh
    - chmod +x client.sh
    - chmod +x gen.sh
