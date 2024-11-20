import React, { useState, useEffect, useRef } from 'react';
import { Card, Container, Row, Col, ListGroup, Button, Form, ButtonGroup, Nav, Badge, ListGroupItem } from 'react-bootstrap';
import DeletedSecrets from './DeletedSecrets'
import Datetime from 'react-datetime';
import "react-datetime/css/react-datetime.css";
import { FaCheck, FaGit, FaKey, FaNetworkWired, FaTimes, FaTrash } from 'react-icons/fa';
import Title from './common/Title';
import { ApprovalRequestHelp } from './help/Help';
import CreateSecretForm from './CreateSecretForm';
import SecretDetails from './SecretDetails'
import { useApi } from '../api/api'


const ApprovalRequests = ({setTitle, setHelp}) => {
    const { get, post, put, del } = useApi();
    const [requests, setRequests] = useState([]);
    const [statusFilter, setStatusFilter] = useState('pending'); // Default to 'pending'
    const [error, setError] = useState('');
    const [message, setMessage] = useState('');

    const hasFetchApprovals = useRef(false);
    useEffect(() => {
        if (!hasFetchApprovals.current) {
            hasFetchApprovals.current = true;
            fetchApprovals(statusFilter)
        }
    }, [statusFilter]);

    useEffect(() => {
        setTitle({ heading: "Secret - Approval Workflow", subheading: "Approve creation of secrets"})

        setHelp((
            <>
            <ApprovalRequestHelp />
            </>
        ))

    }, []);


    // useEffect(() => {
    //     fetchApprovals()
    // }, [statusFilter]);

    const setStatus = async (status) => {
        setStatusFilter(status)
        fetchApprovals(status)
    }
    
    const fetchApprovals = async (statusFilter) => {
        try {
            const data = await get(`/approval-requests?status=${statusFilter}`);
            setRequests(data || [])
        } catch (error) {
            setError(error.message);
        }
    }

    const handleApprove = async (request) => {
        try {
            const data = await post(`/approval-requests/approve?path=${encodeURIComponent(request.details.path)}&key=${encodeURIComponent(request.details.key)}`, { request_id: request.id });
            setRequests(requests.filter(r => r.id !== request.id))
        } catch (error) {
            setError(error.message);
        }
    };

    const handleReject = async (request) => {
        try {
            const data = await post(`/approval-requests/reject?path=${encodeURIComponent(request.details.path)}&key=${encodeURIComponent(request.details.key)}`, { request_id: request.id });
            setRequests(requests.filter(r => r.id !== request.id))
        } catch (error) {
            setError(error.message);
        }
    };

    return (
        <div className="">
            <Container className='p-0'>
                <Row>


                    <Col  className='mb-3'>

                    {message && <p className='mb-3 bg-success-soft text-success p-3 rounded-2 fw-bold' variant='info'>{message}</p>}
                    {error && <p className='mb-3 bg-danger-soft text-danger p-3 rounded-2 fw-bold' variant='danger'>{error}</p>}

                    <Nav variant="underline" activeKey={`${statusFilter}`} className='mb-3'>
                                    <Nav.Item onClick={() => setStatus("pending")}>
                                        <Nav.Link eventKey="pending" className={`ms-2 me-auto text-center text-center `} >
                                            <span className="fw-bold">Pending</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item onClick={() => setStatus("approved")}>
                                        <Nav.Link eventKey="approved" className={`ms-2 me-auto text-center text-center`} >
                                            <span className="fw-bold">Approved</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item onClick={() => setStatus("rejected")}>
                                        <Nav.Link eventKey="rejected" className={`ms-2 me-auto text-center text-center`} >
                                            <span className="fw-bold">Rejected</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                </Nav>

                        <Card className="mb-3">
                            {/* <Card.Header>Approval Requests</Card.Header> */}
                            {/* <Card.Body> */}
                                
                                <ListGroup variant="flush">
                                    {requests.map(request => (
                                        <ListGroup.Item key={request.id} className='d-flex align-items-center'>
                                            <div className='w-100'>

                                                <div className='d-flex justify-content-between align-items-start'>
                                                    <div className=''>
                                                        <FaKey />
                                                    </div>

                                                    <div className='ms-2 me-auto'>
                                                    <Badge  bg="primary">User {request.user_id}</Badge> wants to <Badge  bg="dark">{request.action}</Badge> the secret 
                                                    </div>

                                                    <div className='' >
                                                        {statusFilter == "pending" && <ButtonGroup size='sm' className='p-0 m-0'>
                                                            <Button variant="success"  size="sm" onClick={() => handleApprove(request)}><FaCheck /></Button>
                                                            <Button variant="danger"  size="sm" onClick={() => handleReject(request)}><FaTimes /></Button>
                                                        </ButtonGroup>}
                                                    </div>
                                                </div>

                                                {/* <pre>{JSON.stringify(request.details, null, 2)}</pre> */}
                                                <div className="mt-2 border rounded-2 p-2 bg-light">
                                                {/* <SecretDetails path="" secret={request.details} /> */}

                                                        <div className='d-flex justify-content-between align-items-start'>
                                                            <div className='me-auto'><b>Key</b><br/>{request.details.key}</div>
                                                            <div className='me-auto'><b>Value</b><br/>{request.details.value}</div>
                                                        </div>
                                                </div>

                                                
                                            </div>
                                            


                                        </ListGroup.Item>
                                    ))}
                                </ListGroup>
                            {/* </Card.Body> */}
                        </Card>

                        {/* <CreateSecretForm approvalRequired={true} /> */}
                    </Col>



                </Row>
            </Container>
        </div>



    );
};

export default ApprovalRequests;
