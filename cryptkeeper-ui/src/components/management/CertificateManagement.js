import React, { useState, useEffect, useRef } from 'react';
import { Row, Col, Card, Form, Button, ListGroup, Alert, Container, InputGroup } from 'react-bootstrap';
import AppRoles from '../AppRole';
import Title from '../common/Title';
import { RoleManagementHelp } from '../help/Help';
import { useApi } from '../../api/api';

const CertificateManagement = ({ setTitle, setHelp }) => {

    useEffect(() => {
    }, []);

    useEffect(() => {
        setTitle({ heading: "Client Certificates", subheading: "Manage certificates" })
        setHelp(
            <>

            </>
        )
    }, []);

    const [certificates, setCertificates] = useState([]);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const { get, post, put, del, downloadPost, downloadGet } = useApi();

    const [message, setMessage] = useState('');

    const createCertificate = async (e) => {
        e.preventDefault();

        try {
            const data = await downloadPost(`/certificates`, { name: name, description: description });
            // setMessage(`Certificate created with Secret ID: ${data.secret_id}`);

            // const data = await downloadGet(`/pki/download-certificate?secret_id=${secret.id}`);
            const blob = await data.blob();
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = name + '.p12';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);


            // fetchAppRoles()
        } catch (error) {
            setMessage(error.message);
        }


    };

    const downloadCA = async (e) => {
        e.preventDefault();

        try {
            const data = await downloadGet(`/certificates/ca`);
            const blob = await data.blob();
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = 'ca.pem';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            // fetchAppRoles()
        } catch (error) {
            setMessage(error.message);
        }


    };


    return (
        <div className=''>


            <Container className='p-0'>
                <Row >

                    <Col className='mb-3'>
                        <div>
                            <Card>
                                <Card.Header>Create a new mTLS client certificate</Card.Header>
                                <Card.Body>
                                    <p>Mutual TLS (mTLS) authentication uses client certificates to ensure traffic between client and server is bidirectionally secure and trusted. Once you generate a certificate, you can use the common name in the Path policy to specify granular permissions</p>
                                    {message && <div className='bg-primary-soft text-primary p-3 rounded-2 mb-2' variant="info">{message}</div>}
                                    <form onSubmit={createCertificate}>

                                        <InputGroup>
                                            <div className='form-floating'>
                                                <input type="text" className='form-control' value={name} onChange={(e) => setName(e.target.value)} required />
                                                <label>Common Name</label>
                                            </div>

                                            <div className='form-floating'>
                                                <input type="text" className='form-control ' value={description} onChange={(e) => setDescription(e.target.value)} required />
                                                <label>Description</label>
                                            </div>

                                            <button className=' btn  btn-primary' type="submit">Generate Client Cert</button>
                                        </InputGroup>


                                    </form>
                                </Card.Body>
                                <Card.Footer>
                                    <button className='mt-2 btn w-100 btn-dark' onClick={downloadCA} type="submit">Download Client CA</button>
                                    <pre className='bg-dark text-white p-3 mt-2 mb-2 rounded-2'>
                                        # Veriy Certificate<br />
                                        openssl pkcs12 -in certificate.p12 -clcerts -nodes -passin pass:"password"<br /><br />
                                        # Verify CA<br />
                                        openssl x509 -in ca_1.pem -text -noout<br /><br />
                                        # Convert p12 to PEM file<br />
                                        openssl pkcs12 -in certificate.p12 -out certificate.pem -nodes -passin pass:"password"<br /><br />
                                        # Verify Certificate against CA<br />
                                        openssl verify -CAfile ~/go/src/github.com/cryptkeeperhq/cryptkeeper/scripts/certs/ca.pem test123.pem
                                    </pre>

                                </Card.Footer>
                            </Card>

                        </div>
                    </Col>

                </Row>

            </Container>







        </div>
    );
};

export default CertificateManagement;