[linux, macos]
build:
  go build -ldflags="-s -w" -o ./bin/mfg main.go

[linux, macos]
run *args:
  just build
  ./bin/mfg {{ args }}
