import React, { useState, useEffect, useRef } from 'react';
import { Card, InputGroup, Container, Form, Button, Alert, Col, Row, Tab, Tabs } from 'react-bootstrap';
import Title from './common/Title';
import { TransitEncryptionHelp } from './help/Help';
import { useApi } from '../api/api'
import CLIUsage from './CLIUsage';


const TransitEncryption = ({ keyId, setTitle, setHelp }) => {
    const { get, post, put, del } = useApi();
    const [keys, setKeys] = useState([]);
    const [selectedKey, setSelectedKey] = useState('');
    const [selectedKeyType, setSelectedKeyType] = useState('');
    const [plaintext, setPlaintext] = useState('');
    const [ciphertext, setCiphertext] = useState('');
    const [encryptedText, setEncryptedText] = useState('');
    const [decodedText, setDecodedText] = useState('');
    const [signature, setSignature] = useState('');
    const [message, setMessage] = useState('');
    const [signVerifyResponseText, setSignVerifyResponseText] = useState('');
    const [decryptedText, setDecryptedText] = useState('');
    const [error, setError] = useState('');

    const hasFetchKeys = useRef(false);
    useEffect(() => {
        if (!hasFetchKeys.current) {
            hasFetchKeys.current = true;
            fetchKeys()
        }

        
    }, []);

    useEffect(() => {
        if (keyId != null) {
            console.log(keyId)
            setKey(keyId)
            
        }
    }, [keys]);

    useEffect(() => {
        renderOperations()
    }, [selectedKeyType]);

    useEffect(() => {
        // setTitle({ heading: "Transit Encryption", subheading: "Test out the transit encryption functionality" })


        // setHelp(
        //     <>
        //         <TransitEncryptionHelp />
        //     </>
        // )

    }, []);


    // Encoding to Base64
    const encodeToBase64 = (str) => {
        return btoa(str);
    };

    // Decoding from Base64
    const decodeFromBase64 = (base64) => {
        return atob(base64);
    };



    const fetchKeys = async () => {
        try {
            const data = await get(`/transit/keys`);
            setKeys(data || []);
        } catch (error) {
            console.log(error.message);
        }

    };


    const handleKeyChange = (e) => {
        const keyId = e.target.value;
        setKey(keyId)

    };

    const setKey = (keyId) => {
        setSelectedKey(keyId);
        const key = keys.find(k => k.id === keyId);
        setSelectedKeyType(key ? key.key_type : '');
    }




    const hmacText = async () => {
        try {
            let encodedText = encodeToBase64(plaintext)
            const data = await post(`/transit/hmac`, {
                key_id: selectedKey,
                message: encodedText
            });
            setEncryptedText(data.hmac);
            setDecodedText(data.hmac)
        } catch (error) {
            console.log(error.message);
        }

    };

    const hmacTextVerify = async () => {
        try {
            let encodedText = encodeToBase64(plaintext)

            const data = await post(`/transit/hmac/verify`, {
                key_id: selectedKey,
                message: encodedText,
                hmac: encryptedText
            });
            setSignVerifyResponseText(data.verified ? "Verification Successful" : "Failed Verification");
            // setDecodedText(atob(data.ciphertext))
        } catch (error) {
            setSignVerifyResponseText(error.message);
        }

    };

    const signText = async () => {

        try {
            let encodedText = encodeToBase64(message)
            const data = await post(`/transit/sign`, {
                key_id: selectedKey,
                message: encodedText
            });
            setSignVerifyResponseText(data.signature);
            setSignature(data.signature)
        } catch (error) {
            setSignVerifyResponseText(error.message);
        }

    };

    const verifyText = async () => {

        try {
            let encodedText = encodeToBase64(message)
            const data = await post(`/transit/verify`, {
                key_id: selectedKey,
                message: encodedText,
                signature: signature
            });
            setSignVerifyResponseText(data.verified ? "Verified Message and Signature" : "Not Verified");
        } catch (error) {
            setSignVerifyResponseText(error.message);
        }

    };



    const encryptText = async () => {
        try {
            setError("")
            setEncryptedText("")
            setCiphertext("")
            setDecodedText("")
            let encodedText = encodeToBase64(plaintext)

            let url = "/transit/encrypt"
            const data = await post(url, {
                key_id: selectedKey,
                plaintext: encodedText
            });
            setEncryptedText(data.ciphertext);
            setCiphertext(data.ciphertext)

            let dText = decodeFromBase64(data.ciphertext.split(':')[1])
            setDecodedText(dText)
        } catch (error) {
            setEncryptedText(error.message);
        }
    };

    const decryptText = async () => {
        try {
            setError("")
            setDecryptedText("")
            let url = "/transit/decrypt"

            const data = await post(url, {
                key_id: selectedKey,
                ciphertext
            });
            let decodedText = decodeFromBase64(data.plaintext)
            setDecryptedText(decodedText);
        } catch (error) {
            setDecryptedText("")
            setError(error.message);
        }

    };


    const renderOperations = () => {
        switch (selectedKeyType) {
            case 'aes128-gcm96':
            case 'aes256-gcm96':
            case 'aes256S-iv':
            case 'chacha20-poly1305':
            case 'fpe':
                return (
                    <>
                        <Tabs
                            defaultActiveKey="encrypt"
                            id="mgt"
                            variant='underline'
                            className="mb-3 ps-3"
                        >
                            <Tab eventKey="encrypt" title="Encrypt">
                                <Card>
                                    <Card.Header>Encrypt Operation</Card.Header>
                                    <Card.Body>
                                        <InputGroup>
                                            <Form.Group className='form-floating' controlId="formPlaintext">
                                                <Form.Control
                                                    type="text"
                                                    placeholder="Enter plaintext to encrypt"
                                                    value={plaintext}
                                                    onChange={(e) => setPlaintext(e.target.value)}
                                                />
                                                <Form.Label>Plaintext</Form.Label>
                                            </Form.Group>
                                            <Button variant="primary" onClick={encryptText}>
                                                Encrypt
                                            </Button>
                                        </InputGroup>
                                    </Card.Body>
                                    <Card.Footer>
                                        <CLIUsage cmd="transit encrypt [key_id] $(echo -n VALUE | base64)" />
                                        {/* <pre>
                                            +------------+------------+<br />
                                            | KEY        | VALUE      |<br />
                                            +------------+------------+<br />
                                            | CipherText | 0:ajAwWGY= |<br />
                                            +------------+------------+
                                        </pre> */}
                                    </Card.Footer>
                                    <Card.Footer>

                                        {encryptedText && (
                                            <>
                                                <ul className='mt-0 small ms-3 me-3'>
                                                    <li>The returned ciphertext starts with <code>version:encrypted_value</code>. The version indicates the key version which was used to encrypt the plaintext; therefore, when you rotate keys, CryptKeeper knows which version to use for decryption. The rest is a base64 ciphertext value.</li>
                                                    <li>Note that CryptKeeper does not store any of this data. The caller is responsible for storing the encrypted ciphertext. When the caller wants the plaintext, it must provide the ciphertext back to CryptKeeper to decrypt the value.</li>
                                                </ul>

                                                <Alert variant='success'>
                                                    <strong>Encrypted (and encoded) Text:</strong><br /><code>{encryptedText}</code>
                                                    {selectedKeyType == "fpe" && <div>
                                                        <br />
                                                        <strong>Decoded Text:</strong><br /><code>{decodedText}</code>
                                                    </div>}
                                                </Alert>

                                            </>

                                        )}


                                    </Card.Footer>
                                </Card>
                            </Tab>
                            <Tab eventKey="decrypt" title="Decrypt">

                                <Card className='mt-3'>
                                    <Card.Header>Decrypt Operation</Card.Header>
                                    <Card.Body>
                                        <InputGroup className=''>
                                            <Form.Group className='form-floating ' controlId="formCiphertext" >
                                                <Form.Control
                                                    type="text"
                                                    placeholder="Enter ciphertext to decrypt"
                                                    value={ciphertext}
                                                    onChange={(e) => setCiphertext(e.target.value)}
                                                />
                                                <Form.Label>Ciphertext</Form.Label>
                                            </Form.Group>
                                            <Button variant="primary" onClick={decryptText}>
                                                Decrypt
                                            </Button>
                                        </InputGroup>

                                    </Card.Body>
                                    <Card.Footer>
                                        <CLIUsage cmd="transit decrypt [key_id] [ciphertext]" />

                                        {/* <pre>
                                            +-----------+----------+<br />
                                            | KEY       | VALUE    |<br />
                                            +-----------+----------+<br />
                                            | PlainText | VkFMVUU= |<br />
                                            +-----------+----------+<br />
                                        </pre>
                                        <pre>echo "VkFMVUU=" | base64 --decode</pre> */}
                                    </Card.Footer>
                                    <Card.Footer>
                                        {decryptedText && (
                                            <Alert variant='success'>
                                                <strong>Decrypted Text:</strong><br /><code>{decryptedText}</code>
                                            </Alert>
                                        )}


                                    </Card.Footer>
                                </Card>
                            </Tab>
                        </Tabs>


                    </>
                );
            case 'ed25519':
            case 'ecdsa-p256':
            case 'ecdsa-p384':
            case 'ecdsa-p521':
            case 'rsa-2048':
            case 'rsa-3072':
            case 'rsa-4096':
                return (
                    <>

                        <Tabs
                            defaultActiveKey="sign"
                            id="mgt-sign-verify"
                            variant='underline'
                            className="mb-3 ps-3"
                        >
                            <Tab eventKey="sign" title="Sign">
                                <Card>
                                    <Card.Header>Sign Operation</Card.Header>
                                    <Card.Body>
                                        <InputGroup className=''>
                                            <Form.Group className='form-floating' controlId="formMessage">
                                                <Form.Control
                                                    type="text"
                                                    placeholder="Enter message to sign"
                                                    value={message}
                                                    onChange={(e) => setMessage(e.target.value)}
                                                />
                                                <Form.Label>Message</Form.Label>
                                            </Form.Group>
                                            <Button variant="primary" onClick={signText}>
                                                Sign
                                            </Button>
                                        </InputGroup>

                                        {signature && (
                                            <div className="bg-success-soft text-green border-0 mt-2 mb-2 rounded-2 p-3">
                                                <code>{signature}</code>
                                            </div>
                                        )}

                                    </Card.Body>
                                    <Card.Footer>
                                        <CLIUsage cmd="transit sign [key_id] [plaintext]" />
                                    </Card.Footer>
                                </Card>

                            </Tab>
                            <Tab eventKey="verify" title="Verify">
                                <Card className='mt-3'>
                                    <Card.Header>Verify Operation</Card.Header>
                                    <Card.Body>
                                        {/* <InputGroup className='mt-2'> */}

                                        <Form.Group className='form-floating' controlId="formMessage">
                                            <Form.Control
                                                type="text"
                                                placeholder="Enter message to sign"
                                                value={message}
                                                onChange={(e) => setMessage(e.target.value)}
                                            />
                                            <Form.Label>Message</Form.Label>
                                        </Form.Group>

                                        <Form.Group className='form-floating mt-2' controlId="formSignature">
                                            <Form.Control
                                                type="text"
                                                placeholder="Enter signature to verify"
                                                value={signature}
                                                onChange={(e) => setSignature(e.target.value)}
                                            />
                                            <Form.Label>Signature</Form.Label>
                                        </Form.Group>
                                        <Button variant="primary" className='w-100 mt-2' onClick={verifyText}>
                                            Verify
                                        </Button>
                                        {/* </InputGroup> */}
                                        {signVerifyResponseText && (
                                            <div className="bg-success-soft text-green border-0 mt-2 mb-2 rounded-2 p-3">
                                                <code>{signVerifyResponseText}</code>
                                            </div>
                                        )}
                                    </Card.Body>
                                    <Card.Footer>
                                        <CLIUsage cmd="transit verify [key_id] [plaintext] [signature]" />
                                    </Card.Footer>
                                </Card>
                            </Tab>
                        </Tabs>



                    </>
                );
            case 'hmac':
                return (
                    <>

                        <Tabs
                            defaultActiveKey="hmac"
                            id="hmac-sign-verify"
                            variant='underline'
                            className="mb-3 ps-3"
                        >
                            <Tab eventKey="hmac" title="Sign">
                                <Card>
                                    <Card.Header>HMAC Operation</Card.Header>
                                    <Card.Body>
                                        <Form.Group className='form-floating' controlId="formPlaintext">
                                            <Form.Control
                                                type="text"
                                                placeholder="Enter plaintext to generate HMAC"
                                                value={plaintext}
                                                onChange={(e) => setPlaintext(e.target.value)}
                                            />
                                            <Form.Label>Plaintext</Form.Label>
                                        </Form.Group>
                                        <Button variant="primary" className="w-100 mt-2" onClick={hmacText}>
                                            Compute HMAC
                                        </Button>

                                        {encryptedText && (
                                            <div className="bg-success-soft text-green border-0 mt-2 mb-2 rounded-2 p-3">
                                                <code>{encryptedText}</code>
                                            </div>
                                        )}

                                    </Card.Body>
                                    <Card.Footer>
                                        <CLIUsage cmd="transit hmac [key_id] [plaintext]" />
                                    </Card.Footer>
                                </Card>

                            </Tab>
                            <Tab eventKey="hmac-verify" title="Verify">
                                <Card className='mt-3'>
                                    <Card.Header>HMAC Verify Operation</Card.Header>
                                    <Card.Body>
                                        <Form.Group className='form-floating' controlId="formPlaintext">
                                            <Form.Control
                                                type="text"
                                                placeholder="Enter plaintext to generate HMAC"
                                                value={plaintext}
                                                onChange={(e) => setPlaintext(e.target.value)}
                                            />
                                            <Form.Label>Plaintext</Form.Label>
                                        </Form.Group>

                                        <Form.Group className='form-floating mt-2' controlId="formPlaintext">
                                            <Form.Control
                                                type="text"
                                                placeholder="Enter plaintext to generate HMAC"
                                                value={encryptedText}
                                                disabled
                                            />
                                            <Form.Label>HMAC</Form.Label>
                                        </Form.Group>

                                        <Button variant="primary" className="w-100 mt-2" onClick={hmacTextVerify}>
                                            Verify HMAC
                                        </Button>


                                        {signVerifyResponseText && (
                                            <div className="bg-success-soft text-green border-0 mt-2 mb-2 rounded-2 p-3">
                                                <code>{signVerifyResponseText}</code>
                                            </div>
                                        )}
                                    </Card.Body>

                                    <Card.Footer>
                                        <CLIUsage cmd="transit hmac-verify [key_id] [plaintext] [hmac]" />
                                    </Card.Footer>

                                </Card>
                            </Tab>
                        </Tabs>



                    </>
                );
            default:
                return <></>;
        }
    };

    return (
        <div className="">
            <Container className="p-0 ">
                <Row>
                    <Col >


                        {error && <div className="bg-danger-soft text-danger border-0 mb-2 rounded-2 p-3">{error}</div>}

                        {/* <Card className=''> */}
                            {/* <Card.Header>
                                </Card.Header> */}
                            {/* <Card.Body> */}
                                {/* <b>Explore Cryptography Operations on the keys created by Transit Engine</b> */}

                                <Form className='p-2'>
                                <p className='small'>The Transit engine provides encryption and decryption as a service. It is designed to encrypt data in transit. Users can encrypt and decrypt data without storing it. Users can encrypt and decrypt data without storing it. The engine supports various encryption algorithms.</p>

                                    <Form.Group className='form-floating' controlId="formKey">
                                        <Form.Control
                                            as="select"
                                            value={selectedKey}
                                            onChange={handleKeyChange}
                                        >
                                            <option value="" disabled>Select Key</option>
                                            {keys.map((key) => (
                                                <option key={key.id} value={key.id}>
                                                    {key.path}/{key.key}/{key.version} [{key.key_type}]
                                                </option>
                                            ))}
                                        </Form.Control>
                                        <Form.Label>Select Key</Form.Label>
                                    </Form.Group>
                                </Form>
                            {/* </Card.Body>
                        </Card> */}

                        <div className="mt-2">
                            {renderOperations()}
                        </div>



                    </Col>

                </Row>
            </Container>
        </div>
    );
};

export default TransitEncryption;