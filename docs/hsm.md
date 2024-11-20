# HSM
HSM stands for Hardware Security Module, a physical device that protects and manages cryptographic keys. HSMs are used to: 

- **Secure cryptographic processes:** HSMs generate, protect, and manage keys that are used to encrypt and decrypt data, and create digital signatures and certificates. 
- **Control access to digital security keys:** HSMs allow authorized users to access keys without giving them direct access. This helps reduce attack surface. 
- **Increase data security**: HSMs can help achieve higher levels of data security and trust. 


### Install Soft HSM on your Mac
For all purposes, you can use a Cloud HSM provider such as Thales Luna HSM or AWS CloudHSM. But for this project, we will use SoftHSM. SoftHSM is a software-based implementation of a Hardware Security Module (HSM) that allows users to explore PKCS#11 without the need for a physical HSM.

```sh
# Install HSM using brew on Mac
brew install softhsm

# List all available slots
softhsm2-util --show-slots

# Initialize 
softhsm2-util --init-token --slot 407959141 --label ForCryptkeeper --so-pin 1234 --pin 1234
# The token has been initialized on slot 407959141

# mkdir ~/softhsm-tokens
# directories.tokendir = /Users/yourusername/softhsm-tokens
# softhsm2-util --init-token --slot 0 --label "MyToken" --pin 1234 --so-pin 0000

# --slot 7 specifies the slot where the token will be initialized. If you're initializing the first token, slot 0 is typically used.
# --label "ForCryptkeeper" sets a label for the token, which you'll use to reference it.
# --pin 98765432 sets the user PIN for the token, used for user operations. 
# --so-pin 1234 sets the Security Officer (SO) PIN, used for administrative operations. Change 0000 to a secure PIN as well.

# Export environment variables
export PKCS11_LIB="/opt/homebrew/lib/softhsm/libsofthsm2.so"
export PKCS11_LABEL="ForCryptkeeper"
export PKCS11_PIN=1234
export PKCS11_SO_PIN=1234
```


