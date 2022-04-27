
all:
	make main

main:
	go build -o newApple cmd/newApple/*

crtc: cmd/crtc/main.go crtc/crtc.go
	go build -o TestCrtc cmd/crtc/*

mkasset: cmd/mkasset/main.go graphic/makeAsset.go
	go build -o mkasset cmd/mkasset/*