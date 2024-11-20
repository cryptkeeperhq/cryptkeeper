// src/components/DeletedSecrets.js

import React, { useState, useEffect } from 'react';
import { Card, Button, Form, Container, ListGroup, ListGroupItem, Row, Col } from 'react-bootstrap';
import { useApi } from '../api/api'

const DeletedSecrets = ({ token, paths }) => {
    const { get, post, put, del } = useApi();
    const [deletions, setDeletions] = useState([]);
    const [path, setPath] = useState('');
    const [message, setMessage] = useState('');


    const getDeletedSecrets = async () => {

        try {
            const data = await get(`/paths/${path}/deleted`);
            setDeletions(data || [])
        } catch (error) {
            setMessage(error.message);
        }

    };

    const restoreSecret = async (deletion) => {
        try {
            const data = await post(`/paths/${deletion.path_id}/deleted/${deletion.secret_id}/restore?id=${deletion.id}`, {});
            setDeletions(deletions.filter(deletion => deletion.id !== id));
            setMessage('Secret restored.');
        } catch (error) {
            setMessage(error.message);
        }

    };

    return (
        <div>
            <Card>
                <Card.Header>Deleted Secrets</Card.Header>
                <Card.Body>
                {message && <p className='ps-3 text-warning fw-bold'><strong>{message}</strong></p>}
                    <div className='input-group  mb-1'>
                        <div className="form-floating ">
                            <select
                                className="form-control"
                                id="createPath"
                                value={path}
                                onChange={(e) => setPath(e.target.value)}
                            >
                                <option value="" disabled>Select Path</option>
                                {paths.map((path) => (
                                    <option key={path.id} value={path.id}>
                                        {path.path}
                                    </option>
                                ))}
                            </select>
                            <label htmlFor="createPath">Path:</label>
                        </div>


                        <Button variant="primary" onClick={getDeletedSecrets}>Show</Button>
                    </div>

                </Card.Body>

                <ListGroup variant='flush'>
                    {deletions.map(deletion => (
                        <ListGroupItem key={deletion.id}>
                            <div className="d-flex justify-content-between align-items-center">
                                <div>
                                    <div>Path: <code>{deletion.path}</code>, Key: <code>{deletion.key}</code>, Version: {deletion.version}</div>
                                    <small>Deleted At: {new Date(deletion.deleted_at).toString()}</small>
                                </div>
                                <button className='btn btn-sm btn-success' onClick={() => restoreSecret(deletion)}>Restore</button>
                            </div>
                        </ListGroupItem>
                    ))}
                </ListGroup>
            </Card>
        </div>
    );
};

export default DeletedSecrets;
