#!/bin/bash
export baseDir=$(pwd)
echo "baseDir $baseDir"
echo "Building web..."
rm -rf $baseDir/server/www/dist
cd $baseDir/web
npm run build
echo "Building server..."
sudo rm -rf $baseDir/dist
rm -rf $baseDir/server/data/workdir
rm -rf $baseDir/server/tmp
rm -rf $baseDir/server/logs
cd $baseDir/server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $baseDir/dist/minas_amd64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $baseDir/dist/minas_arm64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $baseDir/dist/minas.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o $baseDir/dist/minas_arm64.exe main.go
echo "Done!"
