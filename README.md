
 <img src="./docs/images/logo.png" width="100" height="100">



# CryptKeeper

CryptKeeper is an advanced secrets management platform designed to securely store and manage sensitive information. It supports various secret engines like Key-Value (KV), Transit, PKI, and Database. The platform also integrates with external systems such as AWS Secrets Manager, Azure Key Vault, and HashiCorp Vault.

<img src="./docs/images/diagram1.png" >

> **CryptKeeper is still under active development and not production ready**
> 


## Table of Contents

- [CryptKeeper](#cryptkeeper)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
    - [Secrets Management:](#secrets-management)
    - [Transit Encryption:](#transit-encryption)
    - [Internal PKI:](#internal-pki)
  - [Architecture Overview](#architecture-overview)
    - [Components](#components)
      - [Tech Stack](#tech-stack)
      - [Optional depedencies](#optional-depedencies)
    - [Cryptography](#cryptography)
    - [Authentication](#authentication)
    - [Policy Engine](#policy-engine)
    - [Secrets Engine](#secrets-engine)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Running CryptKeeper locally](#running-cryptkeeper-locally)
    - [Running CryptKeeper from Source](#running-cryptkeeper-from-source)
  - [Usage](#usage)
    - [Web Interface](#web-interface)
    - [CLI](#cli)
    - [API Documentation](#api-documentation)
    - [Agent and k8s operator](#agent-and-k8s-operator)
  - [Optional Components Setup](#optional-components-setup)
    - [TLS Server](#tls-server)
    - [KeyCloak SSO](#keycloak-sso)
    - [Kafka](#kafka)
    - [SpiceDB](#spicedb)
    - [External Secrets Sync](#external-secrets-sync)
      - [Vault Setup](#vault-setup)
  - [TODO](#todo)
  - [Known Bugs](#known-bugs)
  - [License](#license)

## Features

<img src="./docs/images/1.png" >


### Secrets Management:
- **Web Interface**: Intuitive web interface for managing secrets and policies.
- **Integration**: Sync secrets with HashiCorp Vault, AWS SSM and Azure KV.
- **Secret versioning**: Keep track of every secret and changes; Revert back if needed.
- **Secret Rotation**: Rotate secrets at regular intervals for services like PostgreSQL, MySQL, AWS IAM, and more.
- **Dynamic Secrets**: Generate ephemeral secrets for services like PostgreSQL and more.
- **Secret Detection:** Prevent secrets from leaking.
- **Database Engine**: Generate and manage database credentials dynamically.
- **Audit Logs**: Detailed audit logs for all actions performed within the system.
- **CLI**: Command-line interface for automation and integration.
- **Kubernetes Operator**: Deliver secrets to your Kubernetes workloads and automatically reload deployments.
- **Agent:** Inject secrets into applications without modifying any code logic.

### Transit Encryption:
- Cryptograhic Keys: Centrally manage keys across projects through a user-friendly interface or via the API.
- Encrypt and Decrypt Data: Use symmetric keys to encrypt and decrypt data.
- Sign and Verify: Ensures message integrity and authenticity using efficient key-based hashing.

### Internal PKI:
- Private Certificate Authority: Create CA hierarchies, configure certificate templates for policy enforcement, and start issuing X.509 certificates.
- Certificate Management: Manage the certificate lifecycle from issuance to revocation with support for CRL.
- Alerting: Configure alerting for expiring CA and end-entity certificates.
- PKI Issuer for Kubernetes: Deliver TLS certificates to your Kubernetes workloads with automatic renewal.

## Architecture Overview

The architecture of CryptKeeper is designed to ensure security, scalability, and flexibility. Below is a high-level overview of the components and their interactions:

<img src="./docs/images/diagram.png" >

### Components
The core components or CryptKeeper is the frontend, API server and database. The code is structures in a way that kafka, HSM, Zanzibar etc are all optional components. 

#### Tech Stack
- **React** - UI
- **GO** - API Server
- **Postgres DB**: Central database for storing secrets, paths, policies, and audit logs.

#### Optional depedencies
- **HSM (Hardware Security Module)**: Used for generating and encrypting path keys securely.
  - **Google Tink**: Used for encryption and cryptographic operations.
- **Kafka**: Message queue for publishing and subscribing to policy updates.
- **Zanzibar**: Access control system for enforcing policies.
- **Sync Services**: External synchronization with services like AWS, Azure, and Vault.
- **External Root CA**: External Certificate Authority for signing root certificates.

### Cryptography
Refer to the [Cryptograpy Documentation](docs/cryptography.md) for details on cryptography and how it ensures the security of stored data/secrets.

### Authentication
Refer to the [Authentication Documentation](docs/authentication.md) for details on authentication within CryptKeeper.

### Policy Engine
Refer to the [Policy Engine Documentation](docs/policy.md) for details on how to write policies and how they work within CryptKeeper.

### Secrets Engine
Refer to the [Secrets Engine Documentation](docs/secrets_engine.md) for details on how KV, PKI, transit, and Database secrets are generated, stored and used with CryptKeeper. For more detailed information, you can refer to induvidual engines and how they operate.

- [KV Engine](docs/kv_engine.md)
- [Transit Engine](docs/transit_engine.md)
- [PKI Engine](docs/pki_engine.md)
- [Database Engine](docs/database_engine.md)

## Installation
Refer to the [k8s Documentation](deploy/README.md) for details on deploying this to Kubernetes.

### Prerequisites
- Docker and Docker Compose
- Go
- PostgreSQL
- NPM


### Running CryptKeeper locally 
To set up and run CryptKeeper locally, make sure you have prerequisites installed. Then run the command to get things going
```sh
git clone https://github.com/cryptkeeperhq/cryptkeeper.git
cd cryptkeeper
docker-compose up -d
```
If after changes, you want to restart a specific container you can do:

```sh
docker-compose up -d --no-deps --build app
```

If you want to start optional components, you can use various docker-compose files in the repo.

```sh
docker-compose -f docker-compose.kakfa.yml -f docker-compose.spicedb.yml -f docker-compose.yml up -d
```

### Running CryptKeeper from Source 
**Start Database**
```sh
# Pull and start postgres container
docker pull postgres
docker run --name mypostgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres

# Exec into the container
docker exec -it mypostgres psql -U postgres

# Create Database Cryptkeeper
create database cryptkeeper;
\c cryptkeeper;
```

**Verify Configuration**: 

Configuration is done via environment variables or a `config.yaml` file. 

**Run API Server**

Run `go run cmd/cryptkeeper/main.go` to start the API server. 


**Get UI Up and Running**

Run UI `cd cryptkeeper-ui && npm install` and `npm start`. 
Point your browser to http://localhost:3000/

## Usage
CryptKeeper can be consumed through the CLI, API, or web interface. Ensure that you have the necessary permissions to access the secrets you need.

### Web Interface
Access the web interface at `http://localhost:8080`. The web interface provides a user-friendly way to manage secrets, policies, and audit logs.

### CLI

The CLI allows you to interact with CryptKeeper from the command line.
Refer to the [CLI Documentation](docs/cli.md) for details on how to use various CLI commands to manage your secrets.


### API Documentation

CryptKeeper provides a RESTful API for integration. 
Refer to the [API Documentation](docs/api.md) for more details.

### Agent and k8s operator
TODO


## Optional Components Setup

### TLS Server
Enable TLS on CryptKeeper UI

### KeyCloak SSO
KeyCloak can be used instead of local username/password for authenticataion. You can download the keycloak container and run it as below. Don't forget to change your `auth` configuration section in your config.yaml file and restart.

```sh
docker run -p 8080:9999 -e KC_BOOTSTRAP_ADMIN_USERNAME=admin -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:26.0.2 start-dev
# http://localhost:9999/realms/myrealm
# http://localhost:9999/realms/myrealm/account
```

### Kafka

### SpiceDB
SpiceDB is the most mature open source project inspired by Google's internal authorization system - Zanzibar. You can download the SpiceDB container and run it as below. Don't forget to change your `zanzibar` configuration section in your config.yaml file and restart.

```sh
docker exec -it spicedb spicedb migrate head --datastore-engine postgres --datastore-conn-uri postgres://postgres:mysecretpassword@postgres:5432/cryptkeeper?sslmode=disable
```

### External Secrets Sync
#### Vault Setup
You can follow documentation directly on Hashicorp's website https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install to get vault installed and running locally.

```sh
brew tap hashicorp/tap
brew install hashicorp/tap/vault
```

Run the Vault server:

```sh
vault server -dev
```

```sh
# Test out your vault installation
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN=''
vault login $VAULT_TOKEN

# Enable KV engine
vault secrets enable -path=cryptkeeper kv

# Test things out
vault kv put secret/my-secret value="my-secret-value"
vault kv get secret/my-secret
```

Don't forget to change your `sync.vault` configuration section in your config.yaml file and restart.

```yaml
sync:
  platform: "vault"
  vault:
    address: "http://127.0.0.1:8200"
    token: ""
```


## TODO
- [ ] Implement group-based access permissions in Zanzibar
- [ ] Improve unit test coverage
- [ ] Add support for additional secret engines
- [ ] Enhance the web interface with more features
- [ ] Revocation: Revoke not only single secrets, but a tree of secrets, for example, all secrets read by a specific user, or all secrets of a particular type. Revocation assists in key rolling as well as locking down systems in the case of an intrusion.


## Known Bugs
- [x] Unique constratints on path name and path+key combination
- [ ] SpiceDB permission lookup fails if the key has `.` in the name. Ex: localhost.com for the PKI path engine type
- [ ] gRPC service for transit encryption
- [ ] mTLS certificates for authentication


## License

[MIT License](LICENSE)


