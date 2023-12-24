deadcode -test ./... | grep -v go_templ
# staticcheck ./... | grep -v go_templ | grep  U1000 