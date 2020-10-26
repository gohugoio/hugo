brew install upx
wget https://golang.org/dl/go1.15.3.darwin-amd64.pkg
sudo installer -pkg go1.15.3.darwin-amd64.pkg -target /
go version
go install --tags extended
