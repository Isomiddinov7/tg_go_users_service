#!/bin/bash
CURRENT_DIR=$1
IFS=$'\n'
for folder in $(find "$CURRENT_DIR/tg_protos"/* -type d)
do
     echo ">>>>${folder}"
     sudo protoc --plugin="protoc-gen-go=${GOPATH}/bin/protoc-gen-go" --plugin="protoc-gen-go-grpc=${GOPATH}/bin/protoc-gen-go-grpc" -I=${folder} -I="${CURRENT_DIR}/protos" -I /usr/local/include --go_out="${CURRENT_DIR}" \
   --go-grpc_out="${CURRENT_DIR}" ${folder}/*.proto
done
IFS="$OIFS"