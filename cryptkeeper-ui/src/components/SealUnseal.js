import React, { useState, useEffect } from 'react';
import { Button, Form, Alert, Card, Container, Row, Col } from 'react-bootstrap';
import { useApi } from '../api/api'
import { FaLock, FaUnlock } from 'react-icons/fa';

const SealUnseal = () => {
    const [shares, setShares] = useState(['', '', '']);
    const [message, setMessage] = useState('');
    const [isSealed, setIsSealed] = useState(true);
    const { get, post, put, del } = useApi();


    useEffect(() => {
        getSealStatus()
    }, []);

    const getSealStatus = async () => {

        try {
            const response = await get('/admin/seal/status');
            setIsSealed(response.status);

        } catch (e) {
            setMessage(`Error sealing vault`);
        }
    };

    const handleSeal = async () => {

        try {
            const response = await post('/admin/seal', {});
            setMessage('Vault sealed successfully.');
            setIsSealed(true);
        } catch (e) {
            setMessage(`Error sealing vault`);
        }
    };

    const handleUnseal = async () => {
        try {
            const response = await post('/admin/unseal', { shares });
            setMessage('Vault unsealed successfully.');
            setIsSealed(false);
        } catch (e) {
            setMessage(`Error sealing vault`);
        }

    };

    const handleShareChange = (index, value) => {
        const newShares = [...shares];
        newShares[index] = value;
        setShares(newShares);
    };

    return (
        <Container className="mb-3 p-0">
            <Row>
                <Col>
                    <Card>
                        {/* <Card.Header>CryptKeeper Status</Card.Header> */}
                        <Card.Body>
                            {message && <Alert variant="info">{message}</Alert>}
                            <div>{isSealed ?
                                <div>
                                    <div className='d-flex align-items-center'>
                                        <FaUnlock className='me-2' />
                                        <div className='text-danger bg-danger-soft p-2 rounded-2 fw-bold  me-auto'>CryptKeeper is current Sealed</div>
                                    </div>
                                    <Form className='mt-3'>
                                        {shares.map((share, index) => (
                                            <Form.Group key={index} className="mb-1 form-floating">
                                                <Form.Control
                                                    type="text"
                                                    value={share}
                                                    onChange={(e) => handleShareChange(index, e.target.value)}
                                                    disabled={!isSealed}
                                                />
                                                <Form.Label>Share {index + 1}</Form.Label>
                                            </Form.Group>
                                        ))}
                                        <Button variant="success" onClick={handleUnseal} disabled={!isSealed}>
                                            Unseal CryptKeeper
                                        </Button>
                                    </Form>
                                </div> :
                                <div className='d-flex align-items-center'>
                                    <FaUnlock className='me-2' />
                                    <div className='text-success bg-success-soft p-2 rounded-2 fw-bold  me-auto'>CryptKeeper is current unsealed</div>
                                    <Button variant="transparent" className="ms-auto text-danger fw-bold" onClick={handleSeal} disabled={isSealed}>
                                        Seal?
                                    </Button>
                                </div>}
                            </div>



                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Container>
    );
};

export default SealUnseal;