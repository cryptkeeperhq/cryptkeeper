import React, { useState, useEffect, useRef } from 'react';
import { ListGroupItem, Badge, Card, Col, Container, ListGroup, Row, Form, Button, Pagination, Modal, InputGroup } from 'react-bootstrap';
import { FaUser, FaKey, FaHistory, FaTrashAlt, FaEdit, FaPlus, FaEye, FaShare, FaShareAlt } from 'react-icons/fa';
import Title from './common/Title';
import { useApi } from '../api/api';
import PolicyAuditLogs from './PolicyAuditLogs';

const AuditLogs = ({setTitle, setHelp}) => {
  const { get } = useApi();
  const [logs, setLogs] = useState([]);
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [totalLogs, setTotalLogs] = useState(0);
  const [username, setUsername] = useState('');
  const [action, setAction] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [showDetail, setShowDetail] = useState(false);
  const [selectedLog, setSelectedLog] = useState(null);

  // useEffect(() => {
  //   // fetchAuditLogs();
  // }, [page, username, action, startDate, endDate]);

  const hasFetchedLogs = useRef(false);
  useEffect(() => {
      if (!hasFetchedLogs.current) {
        hasFetchedLogs.current = true;
        fetchAuditLogs();
      }
  }, []);



  const fetchAuditLogs = async () => {
    try {
      const data = await get(`/audit-logs?page=${page}&limit=${limit}&username=${username}&action=${action}&start_date=${startDate}&end_date=${endDate}`);
      
      setLogs(data.logs || []);
      setTotalLogs(data.total || 0);
    } catch (error) {
      console.log(error.message);
    }
  };

  const handlePageChange = (newPage) => {
    setPage(newPage);
    fetchAuditLogs()
  };

  const handleShowDetail = (log) => {
    setSelectedLog(log);
    setShowDetail(true);
  };

  const renderIcon = (action) => {
    switch (action) {
      case 'view':
        return <FaEye className="text-danger" />;
      case 'create':
        return <FaPlus className="text-success" />;
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


  useEffect(() => {
    setTitle({ heading: "Audit Logs", subheading: "View detailed logs of all actions performed within the system"})
    setHelp("")
}, []);

  return (
    <div className=''>
      <Container>
        <Row>
          <Col>
            <Form>
              <Row>
                <Col>
                <InputGroup>
                  <Form.Group className='form-floating' controlId="filterUsername">
                    <Form.Control type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Filter by Username" />
                    <Form.Label>Username</Form.Label>
                  </Form.Group>
                
                  <Form.Group className='form-floating' controlId="filterAction">
                    <Form.Control type="text" value={action} onChange={(e) => setAction(e.target.value)} placeholder="Filter by Action" />
                    <Form.Label>Action</Form.Label>
                  </Form.Group>
                
                  <Form.Group className='form-floating' controlId="filterStartDate">
                    <Form.Control type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
                    <Form.Label>Start Date</Form.Label>
                  </Form.Group>
                
                  <Form.Group className='form-floating' controlId="filterEndDate">
                    <Form.Control type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
                    <Form.Label>End Date</Form.Label>
                  </Form.Group>
                
                  <Button variant="primary" onClick={fetchAuditLogs}>Apply Filters</Button>
                
                </InputGroup>                
                </Col>
              </Row>
            </Form>
            <Card className="mt-4">
              <Card.Header>Audit Logs</Card.Header>
              <ListGroup variant="flush">
                {logs.map(log => (
                  <ListGroupItem key={log.id} className='p-0 pt-2 d-flex align-items-start justify-content-start' onClick={() => handleShowDetail(log)}>
                    <span style={{ width: "40px", height: "40px", lineHeight: "35px", fontSize: "18px" }} className='p-0 ms-1 me-2 d-block text-center'>{renderIcon(log.action)}</span>
                    <div className='me-auto text-start'>
                      <div><span className="badge me-1 bg-green-soft text-green rounded-pill">{log.action}</span> <span className="rounded-pill badge bg-primary-soft text-primary"><FaUser size={10} className="me-1" /> {log.username}</span></div>
                      <div><pre className='mb-1'>{JSON.stringify(log.details)}</pre></div>
                    </div>
                    <span className='me-3 mt-1 ms-auto badge bg-info-soft text-dark text-muted small'>{new Date(log.timestamp).toLocaleString()}</span>
                  </ListGroupItem>
                ))}
              </ListGroup>
            </Card>
            <Pagination className="mt-3">
              {Array.from({ length: Math.ceil(totalLogs / limit) }).map((_, index) => (
                <Pagination.Item key={index + 1} active={index + 1 === page} onClick={() => handlePageChange(index + 1)}>
                  {index + 1}
                </Pagination.Item>
              ))}
            </Pagination>
          </Col>
        </Row>
      </Container>

      <Modal show={showDetail} onHide={() => setShowDetail(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Audit Log Detail</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedLog && (
            <div>
              <p><strong>Action:</strong> {selectedLog.action}</p>
              <p><strong>User:</strong> {selectedLog.username}</p>
              <p><strong>Timestamp:</strong> {new Date(selectedLog.timestamp).toLocaleString()}</p>
              <p><strong>Details:</strong> <pre>{JSON.stringify(selectedLog.details, null, 2)}</pre></p>
            </div>
          )}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowDetail(false)}>Close</Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default AuditLogs;