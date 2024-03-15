set EXE=..\bin\csharp\protoc-3.6.1-win32\bin\protoc.exe
%EXE% --csharp_out=.\csharp --proto_path=..\protos cprotos.proto
%EXE% --csharp_out=.\csharp --proto_path=..\protos data.proto