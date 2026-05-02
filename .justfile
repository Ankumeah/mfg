build:
  go build -ldflags="-s -w" -o ./bin/mfg .

run *args:
  just build
  ./bin/mfg {{ args }}
