CC_SRC_PATH=github.com/vote
CHANNEL_NAME=mychannel
CCNAME=vote
VERSION=1.0
PEER1=peer0.org1.example.com:7051
PEER2=peer0.org2.example.com:7051
PEER3=peer0.org3.example.com:7051

#chaincode install
docker exec cli peer chaincode install -n $CCNAME -v $VERSION -p ${CC_SRC_PATH}
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" \
-e "CORE_PEER_ADDRESS=$PEER2" cli peer chaincode install -n $CCNAME -v $VERSION -p ${CC_SRC_PATH}
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp" \
-e "CORE_PEER_ADDRESS=$PEER3" cli peer chaincode install -n $CCNAME -v $VERSION -p ${CC_SRC_PATH}

#chaincode instantiate
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C $CHANNEL_NAME -n $CCNAME -v $VERSION -c '{"Args":[]}'
sleep 3

#chaincode test
docker exec cli peer chaincode invoke -C $CHANNEL_NAME -n $CCNAME -c '{"Args":["initLedger"]}' 
sleep 3
docker exec cli peer chaincode query -C $CHANNEL_NAME -n $CCNAME -c '{"Args":["queryAllCars"]}' 
