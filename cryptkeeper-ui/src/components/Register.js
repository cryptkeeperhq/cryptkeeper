// src/components/Register.js

import React, { useState } from 'react';
import { Card, InputGroup } from 'react-bootstrap';
import { useApi } from '../api/api'


const Register = () => {
    const { get, post, put, del } = useApi();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [message, setMessage] = useState('');

    const register = async () => {
        try {
            const data = await post(`/register`, { username, password, email, name });
            setMessage('User registered successfully!');
        } catch (error) {
            setMessage(error.message);
        }
    };

    return (
        <div>
            {message && <p className='mb-3 text-info fw-bold' variant='info'>{message}</p>}

            <Card className=''>
                <Card.Header>
                    <h2 className='m-0'>Register</h2>
                    <small className='fw-normal'>Get started with CryptKeeper</small>


                </Card.Header>
                <Card.Body>
                    <div className='form-floating'><input type="text" className='form-control' placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} /><label>Username</label></div>
                    <div className='form-floating mt-1'><input type="password" className='form-control' placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} /><label>Password</label></div>
                    <div className='form-floating  mt-1'><input type="text" className='form-control' placeholder="Name" value={name} onChange={e => setName(e.target.value)} /><label>Name</label></div>

                    <div className='form-floating  mt-1'><input type="text" className='form-control' placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} /><label>Email</label></div>

                    <button className='mt-2 btn-primary btn w-100' onClick={register}>Register</button>
                </Card.Body>
            </Card>


        </div>
    );
};

export default Register;
