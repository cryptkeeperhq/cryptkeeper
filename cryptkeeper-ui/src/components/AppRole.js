import React, { useState, useEffect, useRef } from 'react';
import { Card, Button, Table, Form, Alert } from 'react-bootstrap';

import { useApi } from '../api/api'

const AppRoles = () => {
    const [approles, setAppRoles] = useState([]);
    const [name, setName] = useState('');
    const [roleDescription, setRoleDescription] = useState('');
    const { get, post, put, del } = useApi();


    const [message, setMessage] = useState('');

    const hasFetchDone = useRef(false);
    useEffect(() => {
        if (!hasFetchDone.current) {
            hasFetchDone.current = true;
            fetchAppRoles()
        }
    }, []);




    const fetchAppRoles = async () => {
        try {
            const result = await get(`/approles`);
            setAppRoles(result || [])
        } catch (error) {
            setMessage(error.message);
        }

    };
    const createAppRole = async (e) => {
        e.preventDefault();

        try {
            const data = await post(`/approles`, { name: name, description: roleDescription });
            setMessage(`AppRole created with Secret ID: ${data.secret_id}`);
            fetchAppRoles()
        } catch (error) {
            setMessage(error.message);
        }

        
    };

    return (
        <div>
            <Card>
                <Card.Header>Create a new application identity</Card.Header>
                <Card.Body>
                    {message && <div className='bg-primary-soft text-primary p-3 rounded-2 mb-2' variant="info">{message}</div>}
                    <form onSubmit={createAppRole}>

                        <div className='form-floating'>
                            <input type="text" className='form-control' value={name} onChange={(e) => setName(e.target.value)} required />
                            <label>Role Name</label>
                        </div>

                        <div className='form-floating'>
                            <input type="text" className='form-control mt-1' value={roleDescription} onChange={(e) => setRoleDescription(e.target.value)} required />
                            <label>Role Description</label>

                        </div>


                        <button className='mt-2 btn w-100 btn-primary' type="submit">Create AppRole</button>

                    </form>
                </Card.Body>
            </Card>

            <Card className='mt-3'>
                <Card.Header>Existing Application Roles</Card.Header>
                <Table striped bordered hover className="mt-3">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Name</th>
                            <th>Description</th>
                            <th>Role ID</th>
                            <th>Secret ID</th>
                            <th>Last Updated</th>
                        </tr>
                    </thead>
                    <tbody>
                        {approles.map(role => (
                            <tr key={role.id}>
                                <td>{role.id}</td>
                                <td>{role.name}</td>
                                <td>{role.description}</td>
                                <td>{role.role_id}</td>
                                <td>{role.secret_id}</td>
                                <td>{role.updated_at}</td>
                            </tr>
                        ))}
                    </tbody>
                </Table>
            </Card>

        </div>
    );
};

export default AppRoles;
