ifeq ($(OS),Windows_NT) 
EXT = .exe
endif

all: go-app-example$(EXT) web/app.wasm

go-app-example$(EXT) : main.go
	go build -o $@

web/app.wasm : main.go
	GOARCH=wasm GOOS=js go build -o web/app.wasm
