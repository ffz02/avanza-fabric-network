#!/bin/bash
# PKI Bootstrap: Generate Root CA and Intermediate CA materials using OpenSSL
# Produces keys and certificates for the Fabric CA servers.
# Requires: openssl, a ca.cnf file in the working directory.

set -e

ORG=${1:-"org1"}
DOMAIN=${2:-"dxb.com"}
SUBJ="/C=US/ST=State/L=City/O=${ORG}/CN=${ORG}.${DOMAIN}"

echo "==> Step 1: Generate Root CA private key"
openssl ecparam -name prime256v1 -genkey -noout -out rca.key

echo "==> Step 2: Generate self-signed Root CA certificate (10 years)"
openssl req -config ca.cnf -new -x509 -sha256 -extensions v3_ca \
  -key rca.key -out rca.cert -days 3650 -subj "$SUBJ"

openssl x509 -in rca.cert -out rca.pem

echo "==> Step 3: Generate Intermediate CA private key"
openssl ecparam -name prime256v1 -genkey -noout -out ica.key

echo "==> Step 4: Generate Intermediate CA CSR"
openssl req -new -sha256 -key ica.key -out ica.csr -subj "$SUBJ"

echo "==> Step 5: RCA issues certificate to ICA"
touch index.txt serial
echo 1000 > serial
echo 1000 > crlnumber
openssl ca -batch -config ca.cnf -extensions v3_intermediate_ca \
  -days 2920 -notext -md sha256 -in ica.csr -out ica.cert

echo "==> Step 6: Build certificate chain"
cat ica.cert rca.cert > chain.cert

echo "==> Done. Files: rca.key, rca.cert, rca.pem, ica.key, ica.cert, chain.cert"
echo "    Mount chain.cert and ica.key into the Fabric CA server container."
