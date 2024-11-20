# PKI Engine
The PKI (Public Key Infrastructure) engine in this system provides functionalities to create and manage certificates, certificate authorities (CAs), and certificate policies. The engine supports use cases such as generating mTLS certificates, defining access policies, and managing multiple issuers.

## Features
- Root and Intermediate CA Management: Supports creating and managing root and intermediate certificate authorities.
- Certificate Management: Generate and manage certificates for mTLS authentication.
- Multiple Issuers: Configure and use multiple certificate issuers.
- Certificate Templates: Define reusable templates for certificate creation.
- Access Policies: Apply policies at both group and certificate levels to restrict access.
- Integration with Encryption/Tokenization: Works seamlessly with the encryption and tokenization engine for secured operations.



## Usage
You can create perform all operations on the PKI engine through the UI or using API.

1. Create a Root CA
Use the API or CLI to generate a root CA:

```bash
curl -X POST http://localhost:8080/api/pki/ca/root \
    -H "Authorization: Bearer <token>" \
    -H "Content-Type: application/json" \
    -d '{
        "common_name": "Root CA",
        "organization": "Your Organization",
        "country": "US",
        "validity_days": 3650
    }'
```

2. Create an Intermediate CA. This is equvivalent to creating a new path. The intermediate CA is signed by the root CA 

```bash
curl -X POST http://localhost:8080/api/pki/ca/intermediate \
    -H "Authorization: Bearer <token>" \
    -H "Content-Type: application/json" \
    -d '{
        "parent_ca_label": "RootCA",
        "common_name": "Intermediate CA",
        "organization": "Your Organization",
        "country": "US",
        "validity_days": 1825
    }'
```

3. Issue a Certificate. This is equivalent to creating a new secret on the path. This supports  secret operations like rotate (renew) and delete (revoke). To generate a certificate signed by an intermediate CA:

```bash
curl -X POST http://localhost:8080/api/pki/certificates \
    -H "Authorization: Bearer <token>" \
    -H "Content-Type: application/json" \
    -d '{
        "ca_label": "IntermediateCA",
        "common_name": "api.example.com",
        "organization": "Your Organization",
        "country": "US",
        "validity_days": 365,
        "key_usage": ["digitalSignature", "keyEncipherment"],
        "extended_key_usage": ["serverAuth", "clientAuth"]
    }'
```

4. Retrieve a Certificate
Fetch details of a certificate:

```bash
curl -X GET http://localhost:8080/api/pki/certificates/<cert-id> \
    -H "Authorization: Bearer <token>"
```

5. Define Certificate Templates
Templates simplify certificate creation by defining common attributes:

```bash
curl -X POST http://localhost:8080/api/pki/templates \
    -H "Authorization: Bearer <token>" \
    -H "Content-Type: application/json" \
    -d '{
        "template_name": "ServerCert",
        "key_usage": ["digitalSignature", "keyEncipherment"],
        "extended_key_usage": ["serverAuth"],
        "validity_days": 365
    }'
```

6. Apply Policies
Policies restrict access to groups or certificates:

```bash
curl -X POST http://localhost:8080/api/pki/policies \
    -H "Authorization: Bearer <token>" \
    -H "Content-Type: application/json" \
    -d '{
        "path": "/group-1/cert-1",
        "roles": {
            "editor": ["read", "write"],
            "viewer": ["read"]
        }
    }'
```

https://github.com/jsha/minica



git clone https://github.com/square/certstrap
cd certstrap
go build

./certstrap init \
     --organization "CryptKeeper" \
     --organizational-unit "Crypt Keeper" \
     --country "US" \
     --province "CA" \
     --locality "Fremont" \
     --common-name "CryptKeeper Root"

