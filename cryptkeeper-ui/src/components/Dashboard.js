import React, { useEffect, useState } from 'react';
import { Container, Row, Col, Card, ListGroup, ListGroupItem, Accordion } from 'react-bootstrap';
import Title from './common/Title';
import { FaUser, FaKey, FaHistory, FaTrashAlt, FaEdit, FaShare, FaShareAlt, FaPlus, FaEye, FaInfo } from 'react-icons/fa';
import { useApi } from '../api/api'
import AuditLogs from './AuditLogs';
import CLIUsage from './CLIUsage';
import SearchSecrets from './SearchSecrets';

const Dashboard = () => {
    const { get, post, put, del } = useApi();

    const [summary, setSummary] = useState({});
    const [recentActivity, setRecentActivity] = useState([]);

    useEffect(() => {
        fetchSummary();
        fetchRecentActivity();
    }, []);

    const fetchSummary = async () => {
        try {
            const data = await get('/dashboard/summary');
            setSummary(data || []);
        }
        catch (e) {
            console.error('Failed to fetch dashboard summary');
        }
    };

    const fetchRecentActivity = async () => {
        try {
            const data = await get('/dashboard/recent-activity');
            setRecentActivity(data || []);
        }
        catch (e) {
            console.error('Failed to fetch dashboard summary');
        }

    };

    const renderIcon = (action) => {
        switch (action) {
            case 'view':
                return <FaEye className="text-success" />;
            case 'create':
                return <FaPlus className="text-success" />;
            case 'updated':
                return <FaEdit className="text-warning" />;
            case 'update':
                return <FaEdit className="text-warning" />;
            case 'delete':
                return <FaTrashAlt className="text-danger" />;
            case 'rotate':
                return <FaHistory className="text-info" />;
            case 'shared':
                return <FaShareAlt className="text-warning" />;
            case 'shared_link_access':
                return <FaShare className="text-warning" />;
            default:
                return <FaKey />;
        }
    };


    return (
        <div className="mb-4">
            <Container className='p-0'>
                <Row className="">
                    <Col>
                        <Card>
                            <Card.Body>
                                <Card.Title as="h6">Secrets</Card.Title>
                                <Card.Text>
                                    <div className='fs-1'>
                                        {summary.secrets_count}</div>
                                </Card.Text>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col>
                        <Card>
                            <Card.Body>
                                <Card.Title as="h6">Paths</Card.Title>
                                <Card.Text>
                                    <div className='fs-1'>{summary.paths_count}</div>
                                </Card.Text>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col>
                        <Card>
                            <Card.Body>
                                <Card.Title as="h6">Policies</Card.Title>
                                <Card.Text>
                                    <div className='fs-1'>{summary.policies_count}</div>
                                </Card.Text>
                            </Card.Body>
                        </Card>
                    </Col>
                    <Col>
                        <Card>
                            <Card.Body>
                                <Card.Title as="h6">Shared Secrets</Card.Title>
                                <Card.Text>
                                    <div className='fs-1'>0</div>
                                </Card.Text>
                            </Card.Body>
                        </Card>
                    </Col>
                </Row>

                <Row>
                    <Col xl={12} className='mt-3'>
                        <Accordion >
                            <Accordion.Item eventKey="0">
                                <Accordion.Header>Search Secrets</Accordion.Header>
                                <Accordion.Body>
                                    <SearchSecrets />
                                </Accordion.Body>
                            </Accordion.Item>
                        </Accordion>


                    </Col>
                    <Col>
                        <Card className='mt-3'>
                            <Card.Header>Recent Activity</Card.Header>
                            {recentActivity.length == 0 && <Card.Body>No activity to show</Card.Body>}
                            <ListGroup variant="flush">
                                {recentActivity.map(log => (
                                    <ListGroupItem className='p-2 d-flex align-items-start justify-content-start' key={log.id}>
                                        <span style={{ width: "40px", height: "40px", lineHeight: "35px", fontSize: "18px" }} className='p-0 ms-1 me-2 d-block text-center'>{renderIcon(log.action)}</span>
                                        <div className='me-auto text-start'>
                                            <span className=' small'>{new Date(log.timestamp).toLocaleString()}</span>
                                            <div><span className="badge me-1 bg-green-soft text-green rounded-pill">{log.action}</span> <span className="rounded-pill badge bg-primary-soft text-primary"><FaUser size={10} className="me-1" /> {log.username}</span></div>
                                            <div ><pre className='p-0 m-0'>{JSON.stringify(log.details)}</pre></div>

                                        </div>

                                    </ListGroupItem>
                                ))}
                            </ListGroup>
                        </Card>
                    </Col>
                </Row>
                <Row>
                    <Col>

                        <div className='mt-3'>
                            <CLIUsage cmd="rotate login --role_id=44a34100-5c61-41f5-8683-0682384c49f7	 --secret_id=b4d323c7-e00b-4f84-9019-dda7d1b0a3ff" />
                        </div>
                        <div className='mt-1'>
                            <CLIUsage cmd="export CRYPTKEEPER_TOKEN='your_jwt_token'" />
                        </div>

                    </Col>
                </Row>
            </Container>
        </div>
    );
};

export default Dashboard;