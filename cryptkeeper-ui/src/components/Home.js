// src/components/Login.js

import React, { useState, useEffect } from 'react';
import { Col, Container, InputGroup, Row } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import Notifications from './Notifications';
import AuditLogs from './AuditLogs';
import { FaCheckCircle, FaFileAlt, FaKey, FaMoon, FaPlay, FaSearch, FaShareAlt, FaTerminal, FaUserShield } from 'react-icons/fa';
import Permissions from './Permissions';
import Title from './common/Title';
import Dashboard from './Dashboard';
import SealUnseal from './SealUnseal';

const Home = ({ setTitle, setHelp }) => {

    var permissions = localStorage.getItem("permissions")

    useEffect(() => {
        setTitle({ heading: "CryptKeeper", subheading: "Secure and user-friendly secrets management system" })
        setHelp((
            <>
                <p className='lead'>Welcome to CryptKeeper - an advanced secrets management platform designed to securely store and manage sensitive information.</p>

                <div className="card mb-4">
                    <div className="card-header">
                        <FaPlay size={20} className='pe-2' />  Quick Start Guide
                    </div>
                    <div className="card-body">
                        <h5>Creating a Secret</h5>
                        <ol>
                            <li>Navigate to the "Create Secret" page.</li>
                            <li>Fill in the required details including path, key, value, metadata, expiration date, and one-time use option.</li>
                            <li>Click the "Create Secret" button to save the secret.</li>
                        </ol>
                        <h5>Managing Access Control</h5>
                        <ol>
                            <li>Navigate to the "Path Permissions" page.</li>
                            <li>Select the path you want to manage.</li>
                            <li>Add users or groups and assign the appropriate permission level (owner, editor, viewer).</li>
                        </ol>
                        <h5>Rotating a Secret</h5>
                        <ol>
                            <li>Navigate to the secret's detail page.</li>
                            <li>Click on the "Rotate Secret" option.</li>
                            <li>Enter the new secret value and click "Rotate" to update the secret.</li>
                        </ol>
                        <h5>Generating a Shared Link</h5>
                        <ol>
                            <li>Navigate to the secret's detail page.</li>
                            <li>Click on the "Share One Time Link" option.</li>
                            <li>Select the expiration date and generate the link.</li>
                        </ol>
                        <h5>Viewing Audit Logs</h5>
                        <ol>
                            <li>Navigate to the "Audit Logs" page.</li>
                            <li>View detailed logs of all actions performed within the system.</li>
                        </ol>
                    </div>
                </div>

                <div className="card mb-4">
                    <div className="card-header">
                        <FaTerminal size={20} className='pe-2' />  CLI Usage
                    </div>
                    <div className="card-body">
                        <h5>Detecting Secrets</h5>
                        <p>Use the following command to detect secrets in a code repository:</p>
                        <pre><code>./cryptkeeper-cli detect &lt;project path&gt;</code></pre>
                        <p>This command scans the specified project path for secrets and reports any findings.</p>
                    </div>
                </div>

            </>
        ))
    }, []);


    return (
        <div className='mt-4'>



            <Container>
                <Row>
                    <Col>

                        <SealUnseal />

                        <Dashboard />









                    </Col>
                </Row>
            </Container>

        </div>
    );
};

export default Home;
