import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useApi } from '../../api/api'
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import { ButtonGroup, Offcanvas } from 'react-bootstrap';

const CreateWorkflow = () => {
    const { get, post, put, del } = useApi();

    const [showModal, setShowModal] = useState(false);
    const [newWorkflowName, setNewWorkflowName] = useState('');

    const createFlow = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
        setNewWorkflowName(''); // Clear the input field when the modal is closed
    };

    const handleCreateWorkflow = () => {
        post(`/workflows`, { name: newWorkflowName })
            .then(response => {
                const data = response.data;
                console.log(data);
                // loadWorkflows();
                handleCloseModal(); // Close the modal after creating the workflow
            })
            .catch(error => console.error('Error creating workflow:', error));
    };

    return (
        <>
            <Button className=' rounded-2 btn-dark btn-sm w-100' variant="success" onClick={createFlow}>Create New</Button>
            <Offcanvas show={showModal} onHide={handleCloseModal} className={ showModal ? 'right fade in': 'right fade'}>
                <Offcanvas.Header closeButton>
                    <Offcanvas.Title>Create New Workflow</Offcanvas.Title>
                </Offcanvas.Header>
                <Offcanvas.Body>
                    <Form.Group>
                        <Form.Label>Workflow Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter workflow name"
                            value={newWorkflowName}
                            onChange={(e) => setNewWorkflowName(e.target.value)}
                        />
                        <ButtonGroup className='mt-3 w-100'>
                        {/* <Button variant="secondary" onClick={handleCloseModal}>
                        Cancel
                    </Button> */}
                    <Button variant="primary" onClick={handleCreateWorkflow}>
                        Create
                    </Button>
                    </ButtonGroup>
                    </Form.Group>
                </Offcanvas.Body>
                    
            </Offcanvas>
            </>
    );
};

export default CreateWorkflow;
