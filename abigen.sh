#! /bin/bash

search_dir=./out

mkdir -p ./tmp

forge build --extra-output-files abi bin --force --skip script

echo Running abigen for: Counter

cat ./out/Counter.sol/Counter.json | jq .abi > ./tmp/Counter.abi
cat ./out/Counter.sol/Counter.json | jq -r .bytecode.object | cut -c3- > ./tmp/Counter.bin

  abigen \
    --bin=./tmp/Counter.bin \
    --abi=./tmp/Counter.abi \
    --type=Counter \
    --pkg=main \
    --out=./go/Counter.go

rm -rf ./tmp

