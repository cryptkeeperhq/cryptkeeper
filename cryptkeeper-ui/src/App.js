import React, { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import './App.scss';


import { BrowserRouter as Router, Route, Routes, Link, Navigate } from 'react-router-dom';
import CreateSecret from './components/CreateSecret';
import Register from './components/Register';
import Login from './components/Login';
import Logout from './components/Logout';


import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import { Card, Container, ListGroup, ListGroupItem, Row, Col, Button } from 'react-bootstrap';
import UserSecrets from './components/UserSecrets';
import Home from './components/Home';
import SharedSecret from './components/SharedSecret';
import NavBar from './components/NavBar';
import ThemeToggle from './theme/ThemeToggle';
import ApprovalRequests from './components/ApprovalRequests';
import AuditLogs from './components/AuditLogs';
import { ThemeProvider, useTheme } from './theme/ThemeContext';
import TransitEncryption from './components/TransitEncryption';

import { useApi } from './api/api';

import UserManagementPage from './components/management/UserManagementPage';
import RoleManagement from './components/management/RoleManagement';
import PathManagement from './components/management/PathManagement';
import PolicyManagement from './components/management/PolicyManagement';

import { LoadingProvider } from './api/LoadingContext';
import LoadingSpinner from './api/LoadingSpinner';
import Secret from './components/Secret';
import Title from './components/common/Title';
import PKIManagement from './components/management/PKIManagement';
import { FaQuestionCircle } from 'react-icons/fa';

import WorkflowDetails from './components/Workflow/WorkflowDetails';
import Workflows from './components/Workflow/Workflows';
import DashboardWorkflow from './components/Workflow/DashboardWorkflow';
import { ReactKeycloakProvider, useKeycloak } from '@react-keycloak/web'
import keycloak from './components/utils/keycloak';

import CertificateManagement from './components/management/CertificateManagement';

const ThemedContainer = ({ children }) => {
    const { theme } = useTheme();


    return (
        <div id="main" className={`min-vh-100 ${theme === 'dark' ? 'bg-dark text-white' : 'bg-light text-dark'}`}>
            {children}
        </div>
    );
};

function App() {

    const { get, post, put, del } = useApi();

    // const navigate = useNavigate();

    const [token, setToken] = useState(null);
    const [loading, setLoading] = useState(true)
    const [authenticated, setAuthenticated] = useState(false)
    


    const [title, setTitle] = useState({
        heading: "",
        subheading: ""
    });

    const [help, setHelp] = useState("");
    const [showHelp, setShowHelp] = useState(false)



    useEffect(() => {
        console.log(keycloak.authenticated || false)
    }, [keycloak])


    useEffect(() => {
        const savedToken = localStorage.getItem('token');
        if (savedToken) {
            setToken(savedToken);
        }
        setLoading(false)
        console.log(keycloak)

        // keycloak.init({ onLoad: 'login-required' }).then((authenticated) => {
        //     if (authenticated) {
        //         console.log('User is authenticated');
        //     } else {
        //         console.log('User is not authenticated');
        //     }
        // });


    }, []);

    const handleLogin = (newToken) => {
        localStorage.setItem('token', newToken);
        setToken(newToken);
    };

    const doLogin = () => {
        keycloak.login({ redirectUri: 'http://localhost:3000/' })
    }

    const handleLogout = () => {
        localStorage.removeItem('token');
        setToken(null);

        if(keycloak.authenticated) {
            keycloak.logout()
        }
    };

    const toggleHelp = () => {
        setShowHelp(!showHelp)
    };


    // const setTitleDetails = (details) => {
    //     setTitle(details)
    // };



    const keycloakProviderInitConfig = {
        onLoad: 'check-sso',
    }

    const onKeycloakEvent = (event, error) => {
        console.log('onKeycloakEvent', event, error)

        // onInitError
        // onAuthError

        if (event === 'onReady' && !keycloak.authenticated) {
            setLoading(false)
        }

        if (event === 'onAuthSuccess') {
            if (keycloak.authenticated) {
                setAuthenticated(true)

            } else {
                console.error(error)
            }
        }

    }

    const onKeycloakTokens =  (tokens) => {
        console.log('onKeycloakTokens', tokens)
        if (keycloak.authenticated) {
            keyCloakLogin(tokens)
        }
    }


    const keyCloakLogin = async(tokens) => {
       
            try {
                console.log(tokens)
                const data = await post(`/auth/keycloak`, tokens);
                console.log(data)
                localStorage.setItem('user', data.username);
                localStorage.setItem('name', data.name);
                localStorage.setItem('permissions', JSON.stringify(data.permissions || []));
                localStorage.setItem('token', data.token);
                setToken(data.token);
                setLoading(false)
            } catch (error) {
                console.log(error)
            }
    }

    const initOptions = {
        // pkceMethod: 'S256', 
        onLoad: "check-sso",
        checkLoginIframe: true,
        enableLogging: true,
        // onLoad: "login-required"
    }


    return (

        <ReactKeycloakProvider authClient={keycloak} initOptions={initOptions} initConfig={keycloakProviderInitConfig} onEvent={onKeycloakEvent} onTokens={onKeycloakTokens}>
            <ThemeProvider>
                <Router>
                    <ThemedContainer>
                        {loading ? <div className='p-3 h4'>Loading...</div> : <>
                            <LoadingProvider>
                                <LoadingSpinner />
                                <div className="App vh-100">
                                    <Container fluid >
                                        <Row>
                                            {token && <Col id="sidebar" >
                                                <div className=''>
                                                    <NavBar token={token} logout={handleLogout} />
                                                </div></Col>}


                                            <Col id="content" className='vh-100 overflow-auto'>
                                                {token &&
                                                    <div className=''>
                                                        <Container className=''>
                                                            <Row>
                                                                <Col>
                                                                    <div className=''>
                                                                        <Title heading={title.heading} subheading={title.subheading} />
                                                                    </div>
                                                                </Col>
                                                            </Row>
                                                        </Container>
                                                    </div>
                                                }

                                                <Container>
                                                    <Row><Col>
                                                        <Routes>
                                                            <Route path="/register" element={<Register />} />
                                                            <Route path="/login" element={<Login doLogin={doLogin} setToken={handleLogin} />} />
                                                            <Route path="/logout" element={<Logout setToken={handleLogout} />} />
                                                            <Route path="/" element={token ? <Home setHelp={setHelp} setTitle={setTitle} /> : <Login doLogin={doLogin}  setToken={handleLogin} />} />
                                                            <Route path="/shared/:linkID" element={<SharedSecret setHelp={setHelp} setTitle={setTitle} />} />
                                                            <Route path="/user/secrets/:engine" element={token ? <UserSecrets setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/secrets/create" element={token ? <CreateSecret setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/management" element={token ? <UserManagementPage setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/approval-requests" element={token ? <ApprovalRequests setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/audit-logs" element={token ? <AuditLogs setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/roles" element={token ? <RoleManagement setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/certificates" element={token ? <CertificateManagement setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/paths" element={token ? <PathManagement setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/policies" element={token ? <PolicyManagement setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />

                                                            <Route path="/pki" element={token ? <PKIManagement setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />

                                                            <Route path="/transit/encryption" element={token ? <TransitEncryption setHelp={setHelp} setTitle={setTitle} /> : <Navigate to="/login" />} />
                                                            <Route path="/user/secrets/:id" element={<Secret setHelp={setHelp} setTitle={setTitle} />} />
                                                            <Route path="/user/secrets/:id/:version" element={<Secret setHelp={setHelp} setTitle={setTitle} />} />

                                                            <Route path="/workflows" element={<DashboardWorkflow setTitle={setTitle} />} />
                                                            <Route path="/workflows/:uuid" element={<WorkflowDetails setTitle={setTitle} />} />
                                                        </Routes>
                                                    </Col></Row>

                                                </Container>
                                            </Col>

                                            <Button variant='' className='rounded-circle bg-transparent' style={{ position: "fixed", bottom: "10px", right: "10px", width: "100px", height: "100px" }} onClick={toggleHelp}><FaQuestionCircle className='text-primary' size={48} /></Button>

                                            {showHelp && token && help &&
                                                <>
                                                    <Col className='d-md-none d-lg-block bg-secondary vh-100 overflow-auto col-md-3 ' style={{ maxWidth: "450px" }}>
                                                        <div id="helpbar" className='mt-3 p-2'>
                                                            {/* <Title heading={title.heading} subheading={title.subheading} /> */}
                                                            {help}
                                                        </div></Col>

                                                </>}
                                        </Row>
                                    </Container>
                                </div>
                            </LoadingProvider>

                        </> }
                    </ThemedContainer>
                </Router>
            </ThemeProvider>
        </ReactKeycloakProvider>
    );
}

export default App;
