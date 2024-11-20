import React, { useState, useEffect, useRef } from 'react';
import { Button, Card, Container, FormLabel, ListGroup, ListGroupItem, Row, Col, InputGroup, Form, Alert } from 'react-bootstrap';
import Title from './common/Title';
import { CreateSecretHelp } from './help/Help';
import CreateSecretForm from './CreateSecretForm';
import { FaInfoCircle } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';

const CreateSecret = ({setTitle, setHelp}) => {


    const navigate = useNavigate();

    useEffect(() => {
        setTitle({ heading: "Create Secret", subheading: "Easy Peasy"})
        setHelp((
            <>
            <CreateSecretHelp />
            </>
        ))

    }, []);
    
    const onCreate = async (data) => {
       navigate(`/user/secrets/${data.id}`)
    };

    return (
        <div className="">
            <Container className='p-0'>
 
                <Row>


                    <Col><CreateSecretForm approvalRequired={false} onCreate={onCreate} /></Col>



                </Row>

               
            </Container>
        </div>
    );
};

export default CreateSecret;
