set -ex

if [ ! -d ./bin ]; then
    mkdir bin
fi

go build -o bin github.com/mohanson/dahlia/cmd/dahlia
