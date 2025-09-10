$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o mychat2-app main.go

chmod +x mychat2-app
nohup ./mychat2-app &