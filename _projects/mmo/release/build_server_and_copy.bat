cd ..\server
go build .\main
cd ..\release
xcopy ..\server\main.exe .\server  /S/Q/Y
pause