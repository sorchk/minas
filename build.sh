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
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $baseDir/dist/minas_linux_amd64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $baseDir/dist/minas_linux_arm64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $baseDir/dist/minas_darwin_amd64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $baseDir/dist/minas_darwin_arm64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $baseDir/dist/minas_windows_amd64.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o $baseDir/dist/minas_windows_arm64.exe main.go
cd $baseDir/dist && tar -czf minas_linux_amd64.tar.gz minas_linux_amd64
cd $baseDir/dist && tar -czf minas_linux_arm64.tar.gz minas_linux_arm64
cd $baseDir/dist && tar -czf minas_darwin_amd64.tar.gz minas_darwin_amd64
cd $baseDir/dist && tar -czf minas_darwin_arm64.tar.gz minas_darwin_arm64
cd $baseDir/dist && zip -j minas_windows_amd64.zip minas_windows_amd64.exe
cd $baseDir/dist && zip -j minas_windows_arm64.zip minas_windows_arm64.exe

echo "Done!"
