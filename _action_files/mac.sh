# ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" < /dev/null 2> /dev/null
brew install upx
wget https://golang.org/dl/go1.15.3.darwin-amd64.pkg
sudo installer -pkg go1.15.3.darwin-amd64.pkg -target /
go version
go install --tags extended
