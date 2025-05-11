build:
	goreleaser release --snapshot --clean
	osslsigncode sign \
		-pkcs12 moncertificat.pfx \
		-pass motdepasse \
		-n "Extract images" \
		-i https://www.roumet.com \
		-t http://timestamp.sectigo.com \
		-in extract.exe \
		-out extract-secure.exe

run:
	# export CFLAGS="-I/opt/homebrew/include"
	# export LDFLAGS="-L/opt/homebrew/lib"
	# export PKG_CONFIG_PATH="/opt/homebrew/lib/pkgconfig:$PKG_CONFIG_PATH"
	# go env -w CGO_CXXFLAGS="-O2 -g -ID:/opt/homebrew/include"
	# CGO_CFLAGS="-I/opt/homebrew/include/leptonica" CGO_LDFLAGS="-L/opt/homebrew/lib" go run main.go	
	CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib -llept -ltesseract" go run main.go
	
test:
	export CGO_ENABLED=1
	CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib" go run main.go