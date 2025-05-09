build:
	GOOS=windows GOARCH=amd64 go build -o extract.exe
	osslsigncode sign \
		-pkcs12 moncertificat.pfx \
		-pass motdepasse \
		-n "Extract images" \
		-i https://www.roumet.com \
		-t http://timestamp.sectigo.com \
		-in extract.exe \
		-out extract-secure.exe