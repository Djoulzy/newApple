
all:
	make main

main:
	go build -o newApple cmd/newApple/*

crtc: cmd/crtc/main.go crtc/crtc.go
	go build -o TestCrtc cmd/crtc/*