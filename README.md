
# Introduction to gRPC with Golang!

## Changes due to protoc update:

### Option in the proto file:
There should be a structure of package/file[pb]
```option go_package="greet/greetpb";```

### Command structure:
```protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.```

