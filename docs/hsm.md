# HSM
HSM stands for Hardware Security Module, a physical device that protects and manages cryptographic keys. HSMs are used to: 

- **Secure cryptographic processes:** HSMs generate, protect, and manage keys that are used to encrypt and decrypt data, and create digital signatures and certificates. 
- **Control access to digital security keys:** HSMs allow authorized users to access keys without giving them direct access. This helps reduce attack surface. 
- **Increase data security**: HSMs can help achieve higher levels of data security and trust. 

### Run HSM on your mac
```sh
# Install HSM using brew on Mac
brew install softhsm
```

Edit your softhsm2.conf file in `cmd/examples/hsm/softhsm2.conf`
```sh
export SOFTHSM2_CONF=$(pwd)/examples/hsm/softhsm2.conf
```

```sh
# List all available slots
softhsm2-util --show-slots

# Initialize 
softhsm2-util --init-token --slot 407959141 --label ForCryptkeeper --so-pin 1234 --pin 1234
```

### Run HSM on container (Work in progress)
For all purposes, you can use a Cloud HSM provider such as Thales Luna HSM or AWS CloudHSM. But for this project, we will use SoftHSM. SoftHSM is a software-based implementation of a Hardware Security Module (HSM) that allows users to explore PKCS#11 without the need for a physical HSM.

```sh
docker build -f Dockerfile-HSM --tag softhsm2:2.5.0 .
# Build with version
VERSION=2.6.1 && docker build --build-arg SOFTHSM2_VERSION=$VERSION --tag softhsm2:$VERSION .

# shell into the container
# docker run -ti --rm softhsm2:2.5.0 sh -l
```
```sh
docker run  -it --name softhsm \
  -v "$(pwd)/tokens:/var/lib/softhsm/tokens" \
  -v "$(pwd)/examples/hsm/softhsm2.conf:/etc/softhsm/softhsm2.conf" \
  -e SOFTHSM2_CONF=/etc/softhsm/softhsm2.conf \
  softhsm2:2.5.0
```
-v "$(pwd)/lib:/usr/local/lib" \
  --mount type=bind,source=/usr/local/lib/softhsm/libsofthsm2.so,target=$(pwd)/libsofthsm2.so \


- `-v "$(pwd)/tokens:/var/lib/softhsm/tokens"`: Maps a tokens directory on the host to the container. Update the softhsm2.conf with `directories.tokendir = /var/lib/softhsm/tokens` 


```sh
docker cp softhsm:/usr/local/lib/softhsm/libsofthsm2.so .

docker run -it --name softhsm \
  -v "$(pwd)/tokens:/var/lib/softhsm/tokens" \
  -e SOFTHSM2_CONF=/etc/softhsm/softhsm2.conf \
  softhsm2:2.5.0
```

## Initialize the token in SoftHSM
```sh
# Initialise a new token
softhsm2-util --init-token --slot 0  --label ForCryptkeeper --so-pin 1234 --pin 1234
softhsm2-util --show-slots

# Test the module
pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so -l -t

# RSA Key Pair
pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so -l --keypairgen --key-type rsa:2048 --id 100 --label MyKeyLabel

echo "Data to sign" > data.txt

pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so --id 100 -s -m RSA-PKCS --input-file data.txt --output-file data.sig

# extract pubkey
pkcs11-tool --module /usr/local/lib/softhsm/libsofthsm2.so -r --id 100 --type pubkey > pubkey.der
openssl rsa -inform DER -outform PEM -in pubkey.der -pubin > pubkey.pem

# Verify
openssl rsautl -verify -inkey pubkey.pem -in data.sig -pubin

```


You will see output like
```text
Available slots:
Slot 651578830
    Slot info:
        Description:      SoftHSM slot ID 0x26d64dce
        Manufacturer ID:  SoftHSM project
        Hardware version: 2.5
        Firmware version: 2.5
        Token present:    yes
    Token info:
        Manufacturer ID:  SoftHSM project
        Model:            SoftHSM v2
        Hardware version: 2.5
        Firmware version: 2.5
        Serial number:    f14141e0a6d64dce
        Initialized:      yes
        User PIN init.:   yes
        Label:            ForCryptkeeper
```



## Export required variables

```sh
# export PKCS11_LIB="./libsofthsm2.so"
export PKCS11_LIB=/usr/local/lib/softhsm/libsofthsm2.so
export PKCS11_LIB="/opt/homebrew/lib/softhsm/libsofthsm2.so"
export PKCS11_LABEL="ForCryptkeeper"
export PKCS11_PIN=1234
export PKCS11_SO_PIN=1234

```




