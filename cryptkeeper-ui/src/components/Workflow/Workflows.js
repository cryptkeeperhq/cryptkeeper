import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Container from 'react-bootstrap/Container'
import Spinner from 'react-bootstrap/Spinner'; 
import { useApi } from '../../api/api'

const Workflows = () => {
    const { get, post, put, del } = useApi();

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [workflows, setWorkflows] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [newWorkflowName, setNewWorkflowName] = useState('');

    function formatDate(dateString) {
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        return new Date(dateString).toLocaleDateString(undefined, options);
      }

      
    useEffect(() => {
        loadWorkflows();
    }, []);

    const loadWorkflows = () => {
        get('/workflows')
            .then((response) => {
                setWorkflows(response || []);
                console.log(response)
                setLoading(false)
            })
            .catch((error) => {
                setLoading(false)
                console.error('Error fetching workflows:', error);
                setError(error)
            });
    }

    const createFlow = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
        setNewWorkflowName(''); // Clear the input field when the modal is closed
    };

    const handleCreateWorkflow = () => {
        post(`/create`, { name: newWorkflowName })
            .then(response => {
                const data = response.data;
                console.log(data);
                loadWorkflows();
                handleCloseModal(); // Close the modal after creating the workflow
            })
            .catch(error => console.error('Error creating workflow:', error));
    };

    return (
        <>
        {loading ? (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
          <p>Loading...</p>
        </div>
      ) : (
        <Container className='p-0'>
                <Row xs={1} md={2} lg={3} xl={4} className="mt-2">
                    {workflows.map((workflow) => (
                        <Col key={workflow.id}>
                            
                            
                            {/* <Card style={{ height: "200px" }}> */}
                            <div className='bg-white rounded-2 shadow-sm h-100 p-4 mb-4'>
                                {/* <Card.Img variant="top" src={item.imageUrl} /> */}
                                {/* <Card.Body> */}
                                <small>{formatDate(workflow.created_at)}</small>
                                <div>
                                <Link to={`/workflows/${workflow.id}`} style={{ textDecoration: "none" }} className=''>
                                    
                                        {workflow.name && ( <b>{workflow.name}</b> )}
                                        {workflow.name == "" && ( <b>{workflow.id}</b>)}
                                    
                                        </Link>
                                        </div>
                                {/* </Card.Body>
                                <Card.Footer> */}
                                    
                                    {/* </Card.Footer>
                            </Card> */}
                            </div>
                            
                            
                        </Col>
                    ))}


                </Row>
            </Container>
      )}
            
        </>

    );
};

export default Workflows;
