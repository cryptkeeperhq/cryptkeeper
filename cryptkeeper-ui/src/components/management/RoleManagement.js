import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Form, Button, ListGroup, Alert, Container } from 'react-bootstrap';
import AppRoles from '../AppRole';
import Title from '../common/Title';
import { RoleManagementHelp } from '../help/Help';

const RoleManagement = ({setTitle, setHelp}) => {

    useEffect(() => {
    }, []);

    useEffect(() => {
        setTitle({ heading: "App Roles", subheading: "Manage policies and permission for the paths"})
        setHelp(
            <>
            <RoleManagementHelp />
            </>
        )
    }, []);

    return (
        <div className=''>


            <Container className='p-0'>
                <Row >

                    <Col className='mb-3'>
                        <AppRoles />
                    </Col>
                    
                </Row>

            </Container>







        </div>
    );
};

export default RoleManagement;