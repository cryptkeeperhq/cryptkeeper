// src/components/Login.js

import React, { useState, useEffect } from 'react';
import { Card, Col, Container, InputGroup, Row } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import Register from './Register'
import { useApi } from '../api/api'
import logo from '../assets/logo.webp'

const Login = ({ setTitle, setToken, doLogin }) => {
    const { get, post, put, del } = useApi();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');

    const navigate = useNavigate();

    useEffect(() => {
        if (setTitle) {
            setTitle({})
        }
    }, []);

    const login = async (e) => {
        try {
            e.preventDefault()
            const data = await post(`/auth/user`, { username, password });
            localStorage.setItem('user', data.username);
            localStorage.setItem('permissions', JSON.stringify(data.permissions || []));
            localStorage.setItem('token', data.token);
            localStorage.setItem('name', data.name);
            console.log(data)
            setToken(data.token);
            // console.log("redirect")
            navigate('/#/');
        } catch (error) {
            setMessage(error.message);
        }

    };

    return (
        <div className='mt-4'>



            <Container className='vh-100 container d-flex align-items-center justify-content-center'>
                <Row className=''>
                    <Col xl={{ span: 9, offset: 2 }}>




                        <Row>
                            <Col lg={1} md={1} xl={1}>
                                <div className='text-center mt-1'><img className='rounded-2 w-100' src={logo} />
                                </div>
                            </Col>
                            <Col xl={4}>

<div className='small p-3 text-start'>
    <h1>Welcome to CryptKeeper</h1>
    <p className=''>CryptKeeper is a robust and versatile secret management system designed to securely store, manage, and rotate secrets. It offers a range of features, including <strong>key-value storage</strong>, <strong>dynamic database credentials</strong>, <strong>transit encryption</strong>, and <strong>Public Key Infrastructure (PKI)</strong> support. </p>

    <p>With CryptKeeper, users can create and manage paths, each associated with a specific engine type, to handle different types of secrets. It supports role-based access control, detailed audit logging, and integration with automated application identities via AppRoles. CryptKeeper ensures that sensitive data is protected through strong encryption practices and allows for easy retrieval and management of secrets in a secure environment.</p>
</div>
</Col>

                            <Col xl={7}>
                                {message && <p className='mt-3 text-danger fw-bold' variant='info'>{message}</p>}
                                <Card >
                                    <Card.Header>
                                        <h2 className='m-0'>Login</h2>
                                        <small className='fw-normal'>Provide your login credentials to get started</small>
                                    </Card.Header>

                                    <Card.Body>
                                    <button className='btn btn-success w-100' onClick={doLogin}>Keycloak Login</button>
                                        <InputGroup className='mt-3'>
                                            <div className='form-floating'>
                                                <input type="text" className='form-control' placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
                                                <label>Username</label>
                                            </div>
                                            <div className='form-floating'>
                                                <input type="password" className='form-control' placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
                                                <label>Password</label>
                                            </div>
                                            <button className='btn btn-success' onClick={login}>Login</button>
                                        </InputGroup>
                                    </Card.Body>
                                </Card>
                                <div className='mt-2'>
                                    <Register />
                                </div>

                            </Col>
                        </Row>



                    </Col>
                </Row>
            </Container>

        </div>
    );
};

export default Login;
