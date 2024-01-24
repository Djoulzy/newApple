
all:
	make main

clean:
	go clean -i -cache -modcache

main:
	go build -o newApple cmd/main/*

crtc: cmd/crtc/main.go crtc/crtc.go
	go build -o TestCrtc cmd/crtc/*