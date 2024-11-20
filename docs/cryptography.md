# Cryptography
CryptKeeper uses a layered encryption approach to ensure the security of stored secrets. The system leverages three types of keys: the Master Key, Path Keys, and Secret Keys, to avoid single point of compromise. You can either use HSM or Software encryption for encryption. 

`Currently HSM is under development`. 
You can refer to [HSM Documentation](docs/hsm.md) to learn more about setting up local HSM. 

## Encryption Keys in CryptKeeper

1. Master Key:
	- Used to encrypt and decrypt the Path Key using Google Tink. 
2. Path Key:
   - A unique key generated for each path. 
   - Encrypted with the Master Key and stored securely in database
3. Secret Key:
   - A unique key generated for each secret.
   - Encrypted with the Path Key and stored securely in database.

Here's a brief explanation of each:
#### 1. Master Key

**Purpose:**  
The Master Key is the root key for the entire CryptKeeper system. It is used to encrypt and decrypt Path Keys, adding a layer of security and ensuring that even if a Path Key is compromised, other paths remain secure.

**Setup:**  
- The Master Key is generated and securely stored when CryptKeeper is initialized.
- It is essential to protect the Master Key since it is crucial for unsealing the vault and accessing Path Keys.
- The Master Key can be encrypted using techniques such as Shamir's Secret Sharing to split the key into multiple shares. These shares are then distributed to multiple administrators to enhance security.

```sh
go run cmd/gen-master-key/main.go
```
This will generate master.key file. This is the most secret attribute which you need to protect. 
This will generate 5 shared using SSS. 
SSS (Shamir's Secret Sharing) is a cryptographic algorithm used to divide a secret (in this case, the master key) into multiple parts, called shares. These shares can then be distributed to multiple parties. A minimum number of shares (a threshold) is required to reconstruct the secret.


#### 2. Path Keys

**Purpose:**  
Path Keys are unique keys generated for each path within CryptKeeper. They are used to encrypt and decrypt the secrets stored within their respective paths.

**Setup:**  
- When a new path is created, a 32-byte AES key is generated for that path using Google TINK primitive. 
- The Path Key is then encrypted using the Master Key before being stored in the database.
- During secret operations, the Master Key is used to decrypt the Path Key, which in turn is used to encrypt and decrypt the secrets within that path.


**Example:**
- Path `/kvstore/attribute` might have its unique Path Key.
- This key is used to encrypt the DEKs for secrets like `/kvstore/attribute/foo` and `/kvstore/attribute/bar`.

#### 3. Secret Key / Data Encryption Keys (DEKs)

**Purpose:**  
DEKs are randomly generated AES keys used to encrypt the actual secret values. Each secret has its own unique DEK, providing an additional layer of security.

**Setup:**  
- When a secret is created, a new 32-byte DEK is generated.
- The secret value is encrypted using the DEK.
- The DEK is then encrypted using the Path Key.
- Both the encrypted secret value and the encrypted DEK are stored in the database.

**Encryption Process:**

1. **Generate DEK:** A new 32-byte DEK is generated for the secret.
2. **Encrypt Secret Value:** The secret value is encrypted using the DEK.
3. **Encrypt DEK:** The DEK is encrypted using the Path Key.
4. **Store Encrypted Values:** Both the encrypted secret value and the encrypted DEK are stored.



### Summary

- **Master Key:** The root key for the entire system, used to encrypt and decrypt Path Keys. It is vital for unsealing the vault.
- **Path Keys:** Unique keys for each path, used to encrypt and decrypt the DEKs for the secrets within their respective paths. They provide an additional layer of security by isolating different paths.
- **Data Encryption Keys (DEKs):** Randomly generated AES keys used to encrypt the actual secret values. Each secret has its own unique DEK, which is then encrypted using the Path Key.

This layered approach ensures that the compromise of a single secret or path does not lead to the exposure of other secrets, providing a robust security framework for managing sensitive information.