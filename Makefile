.PHONY: build clean

export GODEBUG=netdns=cgo

build:
	mkdir -p dist
	cd dist && gox ../surfdap

clean:
	rm -rf dist
	rm -rf main
