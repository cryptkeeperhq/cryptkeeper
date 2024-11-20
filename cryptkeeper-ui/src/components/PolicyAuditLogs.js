import React, { useState, useEffect, useRef } from 'react';
import { Card, Container, ListGroup, ListGroupItem, Row, Col, Alert } from 'react-bootstrap';
import Title from './common/Title';
import { useApi } from '../api/api';

const PolicyAuditLogs = ({setTitle }) => {
    const { get } = useApi();
    const [logs, setLogs] = useState([]);
    const [error, setError] = useState('');

    const hasFetched = useRef(false);
    useEffect(() => {
        setTitle={ heading: "Policy Audit Logs", subheading: "View the audit trail for policy changes" }
        if (!hasFetched.current) {
            hasFetched.current = true;
            fetchPolicyAuditLogs()
        }
    }, []);


    const fetchPolicyAuditLogs = async () => {
        try {
            const data = await get('/policies/audit-logs');
            setLogs(data || []);
        } catch (error) {
            setError('Error fetching policy audit logs');
            console.error(error.message);
        }
    };

    return (
        <div className='mt-4'>
            <Container>
                <Row>
                    <Col>
                        <Card>
                            <Card.Header>Policy Audit Logs</Card.Header>
                            {error && <Alert variant="danger">{error}</Alert>}
                            <ListGroup variant="flush">
                                {logs.map(log => (
                                    <ListGroupItem key={log.id}>
                                        <div className='d-flex align-items-start justify-content-start'>
                                            <div className='me-auto text-start'>
                                                <div><span className="badge bg-primary">{log.action}</span> <span className="badge bg-secondary">{log.username}</span></div>
                                                <div><pre>{JSON.stringify(log.details, null, 2)}</pre></div>
                                            </div>
                                            <span className='badge bg-info text-muted small'>{new Date(log.timestamp).toLocaleString()}</span>
                                        </div>
                                    </ListGroupItem>
                                ))}
                            </ListGroup>
                        </Card>
                    </Col>
                </Row>
            </Container>
        </div>
    );
};

export default PolicyAuditLogs;