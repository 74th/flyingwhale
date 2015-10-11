build:
	go build flyingwhale.go
build4linux:
	GOOS=linux GOARCH=amd64 go build -o ./dist/flyingwhale_linux flyingwhale.go
build4macos:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/flyingwhale_darwin flyingwhale.go
build4winwow:
	# not work
	GOOS=windows GOARCH=amd64 go build -o ./dist/flyingwhale_win.exe flyingwhale.go
build4release: build4linux build4macos
