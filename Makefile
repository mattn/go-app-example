SRC = main.go main_nowasm.go main_wasm.go

ifeq ($(OS),Windows_NT) 
EXT = .exe
endif

all: go-app-example$(EXT) web/app.wasm

go-app-example$(EXT) : $(SRC)
	go build -o $@

web/app.wasm : $(SRC)
	GOARCH=wasm GOOS=js go build -o web/app.wasm
