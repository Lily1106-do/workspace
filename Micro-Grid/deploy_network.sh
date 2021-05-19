#! /bin/bash

# 使用fabric部署平台搭建微电网区块链网络

url="http://192.168.49.2:30000"
# url="http://microgrid.com:443"

doFunc(){
    sleep 5
    curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic YWRtaW46cGFzcw==' -i "${url}/api/cluster/$1" --data "$2"
}

setConsensus(){
    echo "**************************************************************"
    echo "****************** set consensus ${consensusPara} ********************"
    echo "**************************************************************"
    consensusPara="{\"cluster_name\":\"cluster\",\"cluster_desc\":\"desc\",\"cluster_consensus\":\"$1\"}"
    doFunc "basic_info" $consensusPara
}

createOrgAndPeers(){
    orgName=${1^}
    orgMethod="create_org"
    peerMethod="create_node"
    ordererName="orderer"
    if [ $1 = "$ordererName" ];then
        sleep 5
        createOrg ${orgMethod} ${orgName} "orderer" 
        for k in 1 2 3 4 5;do
            createPeer $peerMethod $orgName "orderer" $1$k
        done
    else
        sleep 5
        createOrg ${orgMethod} ${orgName} "peer"
        for j in 1 2;do
            createPeer $peerMethod $orgName "peer" ${1}$j
        done
    fi
}

createOrg(){
    echo "**************************************************************"
    echo "****************** create org ${2} *********************"
    echo "**************************************************************"
    # 
    pare="{\"cluster_name\":\"cluster\",\"org_name\":\"$2\",\"org_type\":\"$3\",\"org_desc\":\"org_$2\"}"
    echo $pare
    doFunc $1 $pare
}

createPeer(){
    # peerName=$3$4
    echo "**************************************************************"
    echo "***************** create peer $4 ********************"
    echo "**************************************************************"
    pare="{\"cluster_name\":\"cluster\",\"org_name\":\"$2\",\"org_type\":\"$3\",\"node_name\":\"$4\",\"node_desc\":\"node_$4\"}"
    echo $pare
    doFunc $1 $pare
}

startCluster(){
    echo "**************************************************************"
    echo "********************* create cluster *************************"
    echo "**************************************************************"
    doFunc "start_cluster" "{\"cluster_name\":\"cluster\",\"org_names\":[\"${1^}\",\"${2^}\"]}"
}

channelFunc(){
    # echo "**************************************************************"
    # echo "******** create cluster for  ${1^} and ${2^} **********"
    # echo "**************************************************************"
    # doFunc "start_cluster" "{\"cluster_name\":\"cluster\",\"org_names\":[\"${1^}\",\"${2^}\"]}"
    # # create channel
    org1=$1
    org2=$2
    channelName="${org1:0:3}and${org2:0:3}"
    echo "**************************************************************"
    echo "************** create channel ${channelName} *****************"
    echo "**************************************************************"
    channelPara="{\"cluster_name\":\"cluster\",\"org_names\":[\"${1^}\",\"${2^}\"],\"channel_name\":\"$channelName\"}"
    doFunc "create_channel" $channelPara
    # join it
    for q in $1 $2;do
        for p in 1 2;do
        echo "**************************************************************"
        echo "************ $q${p} join channel ${channelName} **************"
        echo "**************************************************************"
        joinChannelPara="{\"cluster_name\":\"cluster\",\"org_name\":\"${q^}\",\"peer_name\":\"$q${p}\",\"channel_name\":\"$channelName\"}"
        doFunc "join_channel" $joinChannelPara
    done
    # for p in 1 2;do
    #     echo "**************************************************************"
    #     echo "************ $2${p} join channel ${channelName} **************"
    #     echo "**************************************************************"
    #     joinChannelPara="{\"cluster_name\":\"cluster\",\"org_name\":\"${2^}\",\"peer_name\":\"$2${p}\",\"channel_name\":\"$channelName\"}"
    #     doFunc "join_channel" $joinChannelPara
    done
    # update anchor peer
    for n in $1 $2;do
        echo "**************************************************************"
        echo "******* update anchor peer for ${n^} in ${channelName} *******"
        echo "**************************************************************"
        anchorPeerPara="{\"org_name\":\"${n^}\",\"cluster_name\":\"cluster\",\"channel_name\":\"$channelName\",\"peer_name\":\"${n}1\"}"
        doFunc "update_anchor_peer" $anchorPeerPara
    done
}

getCa(){
    echo "**************************************************************"
    echo "***************** get CA for ${caPeerPara} *******************"
    echo "**************************************************************"

    caPeerPara="{\"org_name\":\"${1^}\",\"cluster_name\":\"cluster\",\"user_name\":\"$11\"}"
    doFunc "apply_user" $caPeerPara
}

main(){
    # set consensus method
    consensusFunc="raft"
    setConsensus $consensusFunc
    # create org and peer
    name=(consumer producer orderer)
    for i in ${name[@]};do
        # create org and peer
        createOrgAndPeers ${i}
    done
    sleep 10
    # start cluster & create channel $ join it
    # startCluster ${name[0]} ${name[1]} ${name[2]}
    curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic YWRtaW46cGFzcw==' -i 'http://192.168.49.2:30000/api/cluster/create_channel' --data '{"cluster_name":"cluster","org_names":["Consumer","Producer"],"channel_name":"conandpro"}'
    # startCluster ${name[0]} ${name[1]}
    # startCluster ${name[0]} ${name[1]}
    channelFunc ${name[0]} ${name[1]}
    # channelFunc ${name[1]} ${name[2]}
    # channelFunc ${name[2]} ${name[0]}
    # curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic YWRtaW46cGFzcw==' -i "${url}/api/cluster/apply_user" --data '{"org_name":"Consumer","cluster_name":"cluster","user_name":"consumer1"}'
    # getCa "Consumer"
}
main

