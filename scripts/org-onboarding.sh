#!/bin/bash
# Add a new organization to an existing Fabric channel.
# Usage: ./org-onboarding.sh <NEW_ORG_MSP_ID> <CHANNEL_NAME>
# Requires: peer binary, configtxlator, jq, and admin MSP credentials configured.

set -e

NEW_ORG_MSP=${1:?"Usage: $0 <NEW_ORG_MSP_ID> <CHANNEL_NAME>"}
CHANNEL_NAME=${2:?"Usage: $0 <NEW_ORG_MSP_ID> <CHANNEL_NAME>"}
ORDERER_CA=${ORDERER_CA:?"Set ORDERER_CA env var to the orderer TLS CA cert path"}
ORDERER_ADDRESS=${ORDERER_ADDRESS:-"orderer0.dxb.com:7050"}

echo "==> Generating org definition for ${NEW_ORG_MSP}..."
export FABRIC_CFG_PATH=$PWD
./bin/configtxgen -printOrg "${NEW_ORG_MSP}" > "./${NEW_ORG_MSP}.json"

echo "==> Fetching current channel config..."
peer channel fetch config config_block.pb \
  -o "${ORDERER_ADDRESS}" -c "${CHANNEL_NAME}" \
  --tls --cafile "${ORDERER_CA}"

echo "==> Decoding config block..."
configtxlator proto_decode --input config_block.pb --type common.Block \
  | jq .data.data[0].payload.data.config > config.json

echo "==> Injecting new org into config..."
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups":{"'"${NEW_ORG_MSP}"'":.[1]}}}}}' \
  config.json "./${NEW_ORG_MSP}.json" > modified_config.json

echo "==> Encoding configs to protobuf..."
configtxlator proto_encode --input config.json \
  --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json \
  --type common.Config --output modified_config.pb

echo "==> Computing config delta..."
configtxlator compute_update \
  --channel_id "${CHANNEL_NAME}" \
  --original config.pb \
  --updated modified_config.pb \
  --output "${NEW_ORG_MSP}_update.pb"

echo "==> Wrapping update in envelope..."
configtxlator proto_decode --input "${NEW_ORG_MSP}_update.pb" \
  --type common.ConfigUpdate | jq . > "${NEW_ORG_MSP}_update.json"

echo '{"payload":{"header":{"channel_header":{"channel_id":"'"${CHANNEL_NAME}"'","type":2}},"data":{"config_update":'$(cat "${NEW_ORG_MSP}_update.json)"'}}}' \
  | jq . > "${NEW_ORG_MSP}_update_envelope.json"

configtxlator proto_encode \
  --input "${NEW_ORG_MSP}_update_envelope.json" \
  --type common.Envelope \
  --output "${NEW_ORG_MSP}_update_envelope.pb"

echo "==> Sign and submit the config update with all required org admins, then run:"
echo "    peer channel update -f ${NEW_ORG_MSP}_update_envelope.pb -c ${CHANNEL_NAME} -o ${ORDERER_ADDRESS} --tls --cafile \$ORDERER_CA"
echo "==> Then join the new org's peer:"
echo "    peer channel fetch 0 ${CHANNEL_NAME}.block -o ${ORDERER_ADDRESS} -c ${CHANNEL_NAME} --tls --cafile \$ORDERER_CA"
echo "    peer channel join -b ${CHANNEL_NAME}.block"
