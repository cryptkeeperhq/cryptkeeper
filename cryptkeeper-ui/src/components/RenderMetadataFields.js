import React, { useState, useEffect } from 'react';
import { Button, Card, Container, FormLabel, ListGroup, ListGroupItem, Row, Col, InputGroup, Form, Alert } from 'react-bootstrap';
import TransitMetadataForm from './TransitMetadataForm';
import CodeEditor from './common/CodeEditor';
import { FaCode } from 'react-icons/fa';
import TemplatesSelector from './TemplatesSelector';

const RenderMetadataFields = ({ engineType, metadata, setMetadata }) => {
    const [editor, setEditor] = useState(false);
    const [metadataString, setMetadataString] = useState(JSON.stringify({}));

    console.log(metadata)
    const handleTemplateSelect = (template) => {
        setMetadata(template.fields || {});
    };


    useEffect(() => {
        setMetadataString(JSON.stringify(metadata, null, 2))
    }, [metadata]);



    const handleMetadataChange = (value) => {
        setMetadataString(value);
        try {
            const parsedMetadata = JSON.parse(value);
            setMetadata({ ...metadata, ...parsedMetadata });
        } catch (error) {
            console.log('Invalid JSON format');
        }
    };

    const renderMetadataFields = (engineType) => {
        return (
            <>
                <div className=''>
                    {engineType === 'pki' && (
                        <>
                            <ul className='small mt-2'>
                                <li>When you needs a certificate for a specific domain or service (e.g., home.localhost.com), Cryptkeeper generate a certificate signed by the relevant sub-CA which is part of the path.</li>
                                <li>The certificate includes identifying fields like CommonName and SubjectAltName as needed.</li>
                                <li>This approach is similar to creating "leaf certificates" in a typical PKI, which are used by applications and services.</li>


                            </ul>

                            {/* dnsNames := []string{"www.example.com", "example.org"}
	ipAddresses := []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("10.0.0.1")}
	emailAddresses := []string{"admin@example.com"}
 */}

                            <Form.Group className="form-floating mt-1">
                                <Form.Control
                                    type="text"
                                    // placeholder="Organization"
                                    // value={metadata.organization || ''}
                                    onChange={(e) => setMetadata({ ...metadata, organization: e.target.value })}
                                />
                                <Form.Label>Organization</Form.Label>
                            </Form.Group>
                            <br />
                            <b>Additional SANs</b>: A SAN certificate allows you to secure multiple domains or subject alternative names (SAN) under a single certificate. You can customize the SAN certificate to add up to 250 subject names during its validity period
                            <Form.Group className="form-floating mt-1">
                                <Form.Control
                                    type="text"
                                    // value={metadata.dns_names || ''}
                                    onChange={(e) => setMetadata({ ...metadata, dns_names: e.target.value })}
                                />
                                <Form.Label>DNS Names (Comma separated list of DNS names of your server)</Form.Label>
                            </Form.Group>

                            <Form.Group className="form-floating mt-1">
                                <Form.Control
                                    type="text"
                                    // placeholder="Additional SANs"
                                    // value={metadata.ip_addresses || ''}
                                    onChange={(e) => setMetadata({ ...metadata, ip_addresses: e.target.value })}
                                />
                                <Form.Label>IP Addresses (Comma separated list of IP addresses of your server)</Form.Label>
                            </Form.Group>

                            <Form.Group className="form-floating mt-1">
                                <Form.Control
                                    type="text"
                                    // value={metadata.email_addresses || ''}
                                    onChange={(e) => setMetadata({ ...metadata, email_addresses: e.target.value })}
                                />
                                <Form.Label>Email (Comma separated list of emails)</Form.Label>
                            </Form.Group>



                            <br />
                            <b>Certificate Validity Period</b>
                            <Form.Group className="form-floating mt-1">
                                <Form.Control
                                    type="text"
                                    // placeholder="Validity Period"
                                    // value={metadata.validity_period || metadata.max_lease_time || ''}
                                    onChange={(e) => setMetadata({ ...metadata, validity_period: e.target.value })}
                                />
                                <Form.Label>Validity Period (days)</Form.Label>
                            </Form.Group>
                        </>
                    )}
                    {engineType === 'transit' && (
                        <div className='mb-2'>
                            <TransitMetadataForm metadata={metadata} setMetadata={setMetadata} />
                        </div>
                    )}
                    {engineType == "database" && (
                        <div className="">
                            <Form.Group className='form-floating mt-1'>
                                <Form.Control type="number" value={metadata.ttl}
                                    onChange={(e) => {
                                        setMetadata({ "ttl": e.target.value })
                                    }} required />
                                <Form.Label>TTL (seconds)</Form.Label>

                            </Form.Group>
                            <div className='p-2 small bg-success-soft rounded-2 mt-2'>Database engine uses Dynamic secrets. They are generated on-demand and have a short lifespan. TTL field dicates the lifespan of the secret</div>
                        </div>
                    )}

                    {/* { engineType === "kv" && <TemplatesSelector onTemplateSelect={handleTemplateSelect} />} */}

                    {/* <Form.Group className="form-floating mt-2">
                    <Form.Control
                        as="textarea"
                        placeholder="Metadata (JSON)"
                        value={metadataString}
                        onChange={(e) => handleMetadataChange(e.target.value)}
                        style={{ minHeight: "150px" }}
                    />
                    <Form.Label>Metadata (JSON)</Form.Label>
                </Form.Group> */}


                </div>

                {editor && <CodeEditor height="250px" className="mt-2" code={metadataString} onChange={(value) => handleMetadataChange(value)} />}

                <p className='mt-1' style={{ cursor: "pointer" }} onClick={() => setEditor(!editor)}>
                    <FaCode className='float-start me-1 ms-2 text-muted' size={20} /> {!editor ? "Show" : "Hide"} Metadata Editor
                </p>



            </>
        );
    };



    return (
        <div>
            {renderMetadataFields(engineType)}
        </div>
    );
};

export default RenderMetadataFields;
