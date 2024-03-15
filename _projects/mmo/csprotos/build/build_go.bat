protoc --go_out=.\go --go_opt=paths=source_relative --proto_path=..\protos cprotos.proto
protoc --go_out=.\go --go_opt=paths=source_relative --proto_path=..\protos data.proto
