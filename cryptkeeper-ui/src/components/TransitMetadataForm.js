import React, { useState } from 'react';
import { Form } from 'react-bootstrap';

const keyTypes = [
    'fpe', 'aes128-gcm96', 'aes256-gcm96', 'aes256S-iv', 'chacha20-poly1305', 'ed25519',
    'ecdsa-p256', 'ecdsa-p384', 'ecdsa-p521', 'rsa-2048', 'rsa-3072', 'rsa-4096', 'hmac'
];

const keyTypeOperations = {
    'fpe': ['Encrypt', 'Decrypt'],
    'aes128-gcm96': ['Encrypt', 'Decrypt', 'Key Derivation', 'Convergent Encryption'],
    'aes256-gcm96': ['Encrypt', 'Decrypt', 'Key Derivation', 'Convergent Encryption'],
    'aes256S-iv': ['Encrypt', 'Decrypt'],
    'chacha20-poly1305': ['Encrypt', 'Decrypt', 'Key Derivation', 'Convergent Encryption'],
    'ed25519': ['Sign', 'Verify', 'Key Derivation'],
    'ecdsa-p256': ['Sign', 'Verify'],
    'ecdsa-p384': ['Sign', 'Verify'],
    'ecdsa-p521': ['Sign', 'Verify'],
    'rsa-2048': ['Encrypt', 'Decrypt', 'Sign', 'Verify'],
    'rsa-3072': ['Encrypt', 'Decrypt', 'Sign', 'Verify'],
    'rsa-4096': ['Encrypt', 'Decrypt', 'Sign', 'Verify'],
    'hmac': ['HMAC Generation', 'HMAC Verification']
};

const TransitMetadataForm = ({ metadata, setMetadata }) => {
    const [selectedKey, setSelectedKey] = useState("");

    const selectKey = (keyType) => {
        setSelectedKey(keyType);
        setMetadata({
            ...metadata,
            key_type: keyType,
            key_operations: keyTypeOperations[keyType] || []
        });
    };

    const rows = [
        { key: 'fpe', method: 'FPE-FF3', description: 'Format-preserving encryption', benefit: 'Maintains the structure of sensitive data (e.g., credit card numbers) while encrypting.' },
        { key: 'aes128-gcm96', method: 'aes128-gcm96', description: 'AES-GCM with a 128-bit AES key and a 96-bit nonce; supports encryption, decryption, key derivation, and convergent encryption', benefit: 'Provides high-performance encryption suitable for resource-constrained systems.' },
        { key: 'aes256-gcm96', method: 'aes256-gcm96', description: 'AES-GCM with a 256-bit AES key and a 96-bit nonce; supports encryption, decryption, key derivation, and convergent encryption (default)', benefit: 'Ensures strong security with a longer key length, ideal for sensitive data.' },
        { key: 'chacha20-poly1305', method: 'chacha20-poly1305', description: 'ChaCha20-Poly1305 with a 256-bit key; supports encryption, decryption, key derivation, and convergent encryption', benefit: 'Offers high-speed encryption and resistance to timing attacks, suitable for modern CPUs.' },
        { key: 'ed25519', method: 'ed25519', description: 'Ed25519; supports signing, signature verification, and key derivation', benefit: 'Provides high-performance digital signatures with robust security.' },
        { key: 'ecdsa-p256', method: 'ecdsa-p256', description: 'ECDSA using curve P-256; supports signing and signature verification', benefit: 'Efficient signing and verification with widely used cryptographic standards.' },
        { key: 'ecdsa-p384', method: 'ecdsa-p384', description: 'ECDSA using curve P-384; supports signing and signature verification', benefit: 'Stronger security with moderate performance trade-offs for signing and verification.' },
        { key: 'ecdsa-p521', method: 'ecdsa-p521', description: 'ECDSA using curve P-521; supports signing and signature verification', benefit: 'Highest security level for signing and verification, suitable for long-term data protection.' },
        { key: 'rsa-2048', method: 'rsa-2048', description: '2048-bit RSA key; supports encryption, decryption, signing, and signature verification', benefit: 'Standard RSA strength for general-purpose encryption and signatures.' },
        { key: 'rsa-3072', method: 'rsa-3072', description: '3072-bit RSA key; supports encryption, decryption, signing, and signature verification', benefit: 'Enhanced security over RSA-2048 for sensitive data without significant performance loss.' },
        { key: 'rsa-4096', method: 'rsa-4096', description: '4096-bit RSA key; supports encryption, decryption, signing, and signature verification', benefit: 'Maximum security for RSA encryption and signing, suitable for long-term archival data.' },
        { key: 'hmac', method: 'hmac', description: 'HMAC; supporting HMAC generation and verification', benefit: 'Ensures message integrity and authenticity using efficient key-based hashing.' },
    ];

    const handleKeyTypeChange = (e) => {
        const keyType = e.target.value;
        setMetadata({
            ...metadata,
            key_type: keyType,
            key_operations: keyTypeOperations[keyType] || []
        });
    };

    const handleKeyOperationsChange = (e) => {
        const options = e.target.options;
        const selectedOperations = [];
        for (const option of options) {
            if (option.selected) {
                selectedOperations.push(option.value);
            }
        }
        setMetadata({ ...metadata, key_operations: selectedOperations });
    };

    return (
        <div className='mt-3 '>

            {/* <Form.Group className="form-floating mt-1">
                <Form.Control
                    as="select"
                    placeholder="Key Type"
                    value={metadata.key_type || ''}
                    onChange={handleKeyTypeChange}
                >
                    <option value="" disabled>Select Key Type</option>
                    {keyTypes.map((keyType) => (
                        <option key={keyType} value={keyType}>
                            {keyType}
                        </option>
                    ))}
                </Form.Control>
                <Form.Label>Key Type</Form.Label>
            </Form.Group> */}

            <b>Select Key Type:</b>
            <table className='table table-sm table-hover table-striped'>
                <thead>
                    <tr>
                        <th>Method</th>
                        <th width="35%">Description</th>
                        <th width="50%">Benefit</th>
                    </tr>
                </thead>
                <tbody>
                    {rows.map((row) => (
                        <tr
                            style={{cursor: "pointer"}}
                            key={row.key}
                            className={selectedKey === row.key ? 'active' : ''}
                            onClick={() => selectKey(row.key)}
                        >
                            <td ><b>{row.method}</b></td>
                            <td>{row.description}</td>
                            <td>{row.benefit}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default TransitMetadataForm;
