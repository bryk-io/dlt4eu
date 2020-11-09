.DEFAULT_GOAL:=help

## help: Prints this help message
help:
	@echo "Commands available"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | sort

## wasm: Build the WebAssembly module
# More information: https://github.com/golang/go/wiki/WebAssembly
wasm:
	GOOS=js GOARCH=wasm go build -v -o ./html/dlt4eu.wasm .
