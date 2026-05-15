#!/bin/bash
# Generate channel artifacts (genesis block, channel tx, anchor peer tx)
# Run from the network/ directory after placing Fabric binaries in ./bin/

set -e

export FABRIC_CFG_PATH=$PWD

echo "==> Generating genesis block..."
./bin/configtxgen -profile AvanzaEtcdRaft \
  -channelID orderer-system-channel \
  -outputBlock ./channel-artifacts/genesis.block

echo "==> Generating channel transaction..."
./bin/configtxgen -profile AvanzaChannel \
  -outputCreateChannelTx ./channel-artifacts/avanzachannel.tx \
  -channelID avanzachannel

echo "==> Generating anchor peer transactions..."
./bin/configtxgen -profile AvanzaChannel \
  -outputAnchorPeersUpdate ./channel-artifacts/org1MSPanchors.tx \
  -channelID avanzachannel -asOrg org1MSP

./bin/configtxgen -profile AvanzaChannel \
  -outputAnchorPeersUpdate ./channel-artifacts/org2MSPanchors.tx \
  -channelID avanzachannel -asOrg org2MSP

echo "==> Done. Artifacts written to ./channel-artifacts/"
