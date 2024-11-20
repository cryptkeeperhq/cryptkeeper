import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Form, Button, Alert, InputGroup, Tabs, Tab } from 'react-bootstrap';
// import CreatePolicy from '../CreatePolicy';
// import PolicyEditor from './PolicyEditor';
// import PolicyManagement from './PolicyManagement';
import Title from '../common/Title';
import { PathManagementHelp } from '../help/Help';
import { useApi } from '../../api/api'
import CodeMirror from "@uiw/react-codemirror";
import { vscodeDark } from "@uiw/codemirror-theme-vscode";
import { json } from '@codemirror/lang-json';
import { FaSearch } from 'react-icons/fa';
import AddCA from './AddCA';
import AddTemplate from './AddTemplate';
import AddCertificateRequest from './AddCertificateRequest';


const PKIManagement = ({ setTitle }) => {
    const { get, post, put, del } = useApi();
    const [message, setMessage] = useState(null);


    useEffect(() => {
        setTitle({ heading: "PKI Management", subheading: "Manage your PKI CAs and Templates" })
    }, []);



    return (
        <div className="">
            <Container className="p-0">
                <Row>

                    <Col >
                        {message && <div className='bg-success-soft text-success p-3 rounded-2 mb-2'>{message}</div>}

                        <AddCA />
                        {/* <AddTemplate /> */}
                        {/* <Tabs
                            defaultActiveKey="ca"
                            id="ca-mgt"
                            variant='underline'
                            className="mb-3"
                        >
                            <Tab eventKey="ca" title="New CA">
                                <AddCA />
                            </Tab>
                            <Tab eventKey="request" title="Request New Certificate">
                                <AddCertificateRequest />
                            </Tab>
                            <Tab eventKey="template" title="New Certificate Template">
                                <AddTemplate />
                            </Tab>
                            
                        </Tabs> */}






                    </Col>
                </Row>
            </Container>
        </div>
    );
};

export default PKIManagement;