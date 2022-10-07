
all:
	make main

clean:
	go clean -i -cache -modcache

main:
	export PKG_CONFIG_PATH=/Library/Frameworks/Python.framework/Versions/3.7/lib/pkgconfig; go build -o newApple cmd/main/*

crtc: cmd/crtc/main.go crtc/crtc.go
	go build -o TestCrtc cmd/crtc/*