# Avanza Hyperledger Fabric Network

A production-grade Hyperledger Fabric **v2.5.5** blockchain network implementing a **multi-organization document management and digital signature workflow** system. Features a full PKI infrastructure, private data collections, and sequential/parallel multi-org approval flows.

---

## Overview

This project demonstrates end-to-end setup and operation of a permissioned blockchain network:

- Network configuration and Docker-based deployment
- Custom PKI with Root CA / Intermediate CA hierarchy using Fabric CA
- Go smart contract (chaincode) for document lifecycle management with digital signatures
- Private data collections for sensitive data isolation
- Sequential and parallel multi-org signature policies
- New organization onboarding to a live channel

---

## Architecture

### Network Topology

| Component       | Count | Details                                             |
|-----------------|-------|-----------------------------------------------------|
| Orderers        | 5     | Raft consensus — `orderer0`–`orderer4.dxb.com`      |
| Peer Orgs       | 2     | `org1.dxb.com`, `org2.dxb.com`                      |
| Peers           | 3     | 2× Org1 (ports 7051, 9051), 1× Org2 (port 8051)    |
| State Database  | CouchDB | 3 instances (`couchdb0`–`couchdb2`)               |
| Channel         | 1     | `avanzachannel` (Consortium: `AvanzaConsortium`)    |
| Fabric Version  | 2.5.5 |                                                     |

```
                    ┌──────────────────────────────────────────┐
                    │          avanzachannel                   │
                    │                                          │
  ┌─────────────┐   │  ┌──────────────┐  ┌──────────────┐    │
  │  Orderers   │   │  │    Org 1     │  │    Org 2     │    │
  │ (5× Raft)   │◄──┤  │  peer0:7051  │  │  peer0:8051  │    │
  │             │   │  │  peer1:9051  │  │              │    │
  └─────────────┘   │  │  CouchDB ×2  │  │  CouchDB ×1  │    │
                    │  └──────────────┘  └──────────────┘    │
                    └──────────────────────────────────────────┘
```

### PKI Architecture

```
Root CA (RCA) — prime256v1 EC key
└── Intermediate CA (ICA) — Fabric CA server
    ├── OrdererMSP  (orderer0–orderer4.dxb.com)
    ├── org1MSP     (peer0, peer1, users @ org1.dxb.com)
    └── org2MSP     (peer0, users @ org2.dxb.com)
```

---

## Chaincode — `smctest`

The core smart contract is written in **Go** and implements a document package management system with multi-org digital signature workflows. All arguments are passed via the transient map to prevent sensitive data from appearing in the transaction payload.

### Data Model

| Entity            | Private Collection  | Description                                              |
|-------------------|---------------------|----------------------------------------------------------|
| `Package`         | `package`           | Document bundle with status, progress, requested-by     |
| `Document`        | `document`          | Document with version history, signatures, hash         |
| `SignaturePolicy` | `SignaturePolicy`   | Ordered list of orgs with SLA and page coordinates      |
| `User`            | `user`              | Signatory with P12/public certificate data              |
| `Organization`    | `organization`      | Org-level stamp and public certificate                  |
| `DocumentType`    | `documentType`      | Document catalog (type, label, source, dispatch)        |
| `PackageType`     | `packageType`       | Package catalog with notification policies              |
| `Tasks`           | `tasks`             | Workflow tasks with SLA timers and stage tracking       |
| `OrgUser`         | `orgUsers`          | Org group membership for notifications                  |
| `DocumentHash`    | `documentHash`      | SHA-512 hash index per document                         |
| `DocumentMapping` | `documentMapping`   | Cross-package document reference                        |

### Chaincode Functions

| Function                | Description                                                    |
|-------------------------|----------------------------------------------------------------|
| `signDocument`          | Submit a digital signature/stamp for a document               |
| `upsertSignPolicy`      | Create or update a multi-org signature policy                  |
| `getSignPolicy`         | Retrieve a signature policy by key                             |
| `upsertUser` / `getUser`| Register or fetch a signatory's certificate data              |
| `upsertOrganization`    | Register or update an organization's stamp                     |
| `getOrganization`       | Fetch an organization's stamp configuration                    |
| `upsertDocumentType`    | Register or update a document type definition                  |
| `upsertPackageType`     | Register or update a package type with notification policy     |
| `getPackageType`        | Retrieve a package type definition                             |
| `upsertDocument`        | Create or update a document (stores hash, policy, signatures)  |
| `upsertPackage`         | Create or update a document package                            |
| `regeneratePkg`         | Regenerate a package (for corrections/amendments)             |

### Design Highlights

- **Transient data**: All arguments pass through the transient map (`PrivateArgs`) so sensitive data never enters the read/write set.
- **Private data collections**: 12 collections enforce strict data isolation per entity type.
- **MSP-based access control**: `orgType` (derived from the invoking MSP ID) gates which code paths a caller can execute.
- **Event system**: Every state mutation fires a `chainCodeEvent` with the affected keys and collections, enabling off-chain listeners.
- **SHA-512 hashing**: Document integrity is stored on-chain using `crypto/sha512`.
- **SLA tracking**: Tasks carry `slaTime` fields for service-level agreement enforcement.

---

## Repository Structure

```
avanza-fabric-network/
├── network/
│   ├── configtx.yaml               # Channel, orderer, and org MSP policies
│   ├── crypto-config.yaml          # MSP topology (orgs, peers, users)
│   ├── docker-compose-fabric.yaml  # Docker services for the full network
│   └── core.yaml                   # Peer node configuration
├── chaincode/
│   └── smctest/                    # Go chaincode source
│       ├── main.go                 # Entry point + Invoke function dispatch
│       ├── struct.go               # All data structures / types
│       ├── common.go               # CRUD helpers, MSP lookup, events
│       ├── helpers.go              # Sanitize, hash (SHA-512), event structs
│       ├── typedata.go             # TypeData CRUD operations
│       ├── testing.go              # Test payloads / development helpers
│       ├── go.mod                  # Module: fabric-chaincode-go dependency
│       └── go.sum
├── channel-artifacts/
│   └── collections_config.json     # Private data collection definitions
├── ca-config/                      # Fabric CA Docker configurations
│   ├── docker-compose-ca-org1.yaml
│   └── docker-compose-ca-orderer.yaml
└── scripts/
    ├── artifacts-setup.sh          # Genesis block + channel tx generation
    ├── rca-ica-setup.sh            # Root CA / Intermediate CA bootstrap
    └── org-onboarding.sh           # Add a new org to a running channel
```

---

## Prerequisites

| Dependency          | Version  | Notes                                  |
|---------------------|----------|----------------------------------------|
| Docker              | ≥ 20.x   | With Docker Compose                    |
| Fabric binaries     | 2.5.5    | `configtxgen`, `peer`, `orderer`, etc. |
| Fabric CA           | 1.5.6    | `fabric-ca-server`, `fabric-ca-client` |
| Go                  | ≥ 1.15   | Chaincode compilation                  |
| OpenSSL             | any      | PKI / RCA setup                        |
| `jq`                | any      | Used in onboarding script              |

Download Fabric binaries:
```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.5 1.5.6
```

---

## Getting Started

### 1. Generate PKI Materials

```bash
# Bootstrap Root CA and Intermediate CA for each org
bash scripts/rca-ica-setup.sh org1 dxb.com
bash scripts/rca-ica-setup.sh org2 dxb.com
bash scripts/rca-ica-setup.sh orderer dxb.com

# Start Fabric CA servers
docker-compose -f ca-config/docker-compose-ca-org1.yaml up -d
docker-compose -f ca-config/docker-compose-ca-orderer.yaml up -d
```

### 2. Generate Channel Artifacts

```bash
cd network/
bash ../scripts/artifacts-setup.sh
```

This creates `genesis.block`, `avanzachannel.tx`, and anchor peer transactions in `./channel-artifacts/`.

### 3. Start the Network

```bash
# Set CouchDB credentials (or create a .env file)
export COUCHDB_USER=admin
export COUCHDB_PASSWORD=<your-password>

cd network/
docker-compose -f docker-compose-fabric.yaml up -d
```

Verify all containers are running:
```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

### 4. Create Channel and Join Peers

```bash
export CHANNEL_NAME=avanzachannel
export ORDERER=orderer0.dxb.com:7050

# Create channel (from org1 admin)
peer channel create -o $ORDERER -c $CHANNEL_NAME \
  -f ./channel-artifacts/avanzachannel.tx \
  --tls --cafile $ORDERER_CA

# Join each peer
peer channel join -b ${CHANNEL_NAME}.block

# Update anchor peers
peer channel update -o $ORDERER -c $CHANNEL_NAME \
  -f ./channel-artifacts/org1MSPanchors.tx \
  --tls --cafile $ORDERER_CA
```

### 5. Deploy Chaincode

```bash
# Package
peer lifecycle chaincode package smctest.tar.gz \
  --path ./chaincode/smctest \
  --lang golang \
  --label smctest_1.0

# Install on all peers, then approve + commit with both orgs
peer lifecycle chaincode install smctest.tar.gz
peer lifecycle chaincode approveformyorg ...
peer lifecycle chaincode commit ...
```

### 6. Initialize and Invoke

```bash
# Initialize with MSP-to-orgType mapping
peer chaincode invoke \
  -o $ORDERER -C $CHANNEL_NAME -n smctest \
  --transient '{"PrivateArgs":"[{\"MSP\":\"org1MSP\",\"orgType\":\"Core\",\"ID\":\"org1\"}]"}' \
  -c '{"function":"Init","Args":["[{\"MSP\":\"org1MSP\",\"orgType\":\"Core\",\"ID\":\"org1\"}]"]}'

# Upsert a signature policy
peer chaincode invoke \
  -o $ORDERER -C $CHANNEL_NAME -n smctest \
  --transient '{"PrivateArgs":"<policy_json>|<collection>"}' \
  -c '{"function":"upsertSignPolicy","Args":[]}'
```

### 7. Onboard a New Organization

```bash
export ORDERER_CA=<path-to-orderer-tls-ca.crt>
export ORDERER_ADDRESS=orderer0.dxb.com:7050
bash scripts/org-onboarding.sh newOrgMSP avanzachannel
```

---

## Security Notes

> Before deploying to any environment:
> - Set `COUCHDB_USER` and `COUCHDB_PASSWORD` via environment variables or a `.env` file — never commit credentials.
> - Generate fresh cryptographic material for each deployment; never reuse keys.
> - Review `collections_config.json` policies for your specific org membership.
> - Enable TLS on all Fabric CA endpoints in production.

---

## License

[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)
