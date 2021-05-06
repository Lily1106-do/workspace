#!/bin/bash
set -x
INGRESS=true
YARDMAN_NAME=microgrid
CLUSTER_NAME=cluster
DOMAIN=${CLUSTER_NAME}.${YARDMAN_NAME}
CHANNEL_NAME=conandpro
CC_NAME=usercc
CC_VERSION=1.0

if [ ${INGRESS} = 'false' ]; then
    ORDERER_PORT=7050
    PEER_PORT=7051
else
    ORDERER_PORT=443
    PEER_PORT=443
fi

ORDERER_URL=orderer1.${DOMAIN}:${ORDERER_PORT}
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config/ordererOrganizations/${DOMAIN}/orderers/orderer1.${DOMAIN}/tls/ca.crt

function env_set {
    org_name=$1
    peer_name=$2
    org_name_lower=$(echo ${org_name} | tr '[A-Z]' '[a-z]')
    export CORE_PEER_LOCALMSPID=${org_name}MSP
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config/peerOrganizations/${org_name_lower}.${DOMAIN}/users/Admin@${org_name_lower}.${DOMAIN}/msp
    export CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config/peerOrganizations/${org_name_lower}.${DOMAIN}/users/Admin@${org_name_lower}.${DOMAIN}/tls/client.crt
    export CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config/peerOrganizations/${org_name_lower}.${DOMAIN}/users/Admin@${org_name_lower}.${DOMAIN}/tls/client.key
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/crypto-config/peerOrganizations/${org_name_lower}.${DOMAIN}/users/Admin@${org_name_lower}.${DOMAIN}/tls/ca.crt
    export CORE_PEER_ADDRESS=${peer_name}.${org_name_lower}.${DOMAIN}:${PEER_PORT}
}

function create_channel {
    echo "create channel................................................."
    peer channel create -o ${ORDERER_URL} -c ${CHANNEL_NAME} -f channel.tx
    echo "create channel.............................................done"
}

function join_channel {
    echo "join channel................................................."
    peer channel join -b ${CHANNEL_NAME}.block
    echo "join channel.............................................done"
}

function update_anchor {
    echo "update anchor................................................."
    org_name=$1
    peer channel update -o ${ORDERER_URL} -c ${CHANNEL_NAME} -f ./${org_name}MSPanchors.tx
    echo "update anchor.............................................done"
}

function install_chaincode {
    echo "install chaincode................................................."
    peer chaincode install -n ${CC_NAME} -v ${CC_VERSION} -p github.com/hyperledger/fabric/chaincode/chaincode/consumer
    echo "install chaincode.............................................done"
}

function instantiate_chaincode {
    echo "instantiate chaincode................................................."
    org1_name=$1
    org2_name=$2
    peer chaincode instantiate -o ${ORDERER_URL} -C ${CHANNEL_NAME} -n ${CC_NAME} -v ${CC_VERSION} -c '{"Args":["init","a","100","b","200"]}' -P "OR ('${org1_name}MSP.peer','${org2_name}MSP.peer')"  --tls --cafile ${ORDERER_CA}
    echo "instantiate chaincode.............................................done"
}

function invoke_chaincode {
    echo "invoke chaincode................................................."
    echo peer chaincode invoke -o ${ORDERER_URL}  -C ${CHANNEL_NAME} -n ${CC_NAME} -c '{"Args":["invoke","a","b", "10"]}' --tls --cafile ${ORDERER_CA}
    peer chaincode invoke -o ${ORDERER_URL}  -C ${CHANNEL_NAME} -n ${CC_NAME} -c '{"Args":["invoke","a","b", "10"]}' --tls --cafile ${ORDERER_CA}
    echo "invoke chaincode.............................................done"
}

function query_chaincode {
    echo "query chaincode................................................."
    peer chaincode query -C ${CHANNEL_NAME} -n ${CC_NAME} -c '{"Args":["query","a"]}'
    echo "query chaincode.............................................done"
}




# start from here
org_name=Consumer
env_set Consumer consumer1
install_chaincode ${cc_name} ${cc_version}

# env_set ${org_name} peer1
# install_chaincode ${cc_name} ${cc_version}
sleep 5

instantiate_chaincode Consumer Producer
sleep 10
# query_chaincode
# sleep 10
# invoke_chaincode
# sleep 10
# query_chaincode

# Org2
org_name=Producer
env_set Producer producer1
install_chaincode ${cc_name} ${cc_version}
# sleep 10
# invoke_chaincode
# sleep 10
# query_chaincode
