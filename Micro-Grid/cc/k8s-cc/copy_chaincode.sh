# !bin/bash
set -x
namespace=mychannel-cli
pod=$(kubectl get pods --namespace=$namespace | grep "mychannel-cli" | cut -d" " -f 1)
kubectl cp ../chaincode -n mychannel-cli ${pod}:/opt/gopath/src/github.com/hyperledger/fabric/chaincode/
kubectl cp ../scripts/excute_chaincode.sh -n mychannel-cli ${pod}:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/
