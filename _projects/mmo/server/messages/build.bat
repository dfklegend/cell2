protoc -I="../../csprotos/protos" --go_out=. --go_opt=paths=source_relative --proto_path=. playerinfo.proto
protoc -I="../../../../modules/protoactor-go-0.2.0/actor" -I="../../csprotos/protos"  --go_out=. --go_opt=paths=source_relative --proto_path=. protos.proto
pause

