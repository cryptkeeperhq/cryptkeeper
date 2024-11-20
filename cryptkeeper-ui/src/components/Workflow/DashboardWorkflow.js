import React, { useState, useEffect, useRef } from 'react';
// import './Home.css'; // Import CSS for styling
import Workflows from './Workflows';
import CreateWorkflow from './CreateWorkflow'
import CardBody from 'react-bootstrap/esm/CardBody';
import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Container from 'react-bootstrap/Container'
import { ListGroup, ListGroupItem } from 'react-bootstrap';

function DashboardWorkflow({ setTitle  }) {

  // const [email, setEmail] = useState(null)
  // const [name, setName] = useState(null)
  // const [photo, setPhoto] = useState(null)

  useEffect(() => {
    // // const token = JSON.parse(localStorage.getItem('user'));
    // setEmail(user.email)
    // setName(user.name)
    // setPhoto(user.picture)

    
    setTitle({ heading: "Workflows", subheading: "Automate everything!" })

  }, []);

  return (
    <div className="home">

      
      <div className="container">
        <div className="row">
          <div className="col-md-8 col-xl-10">
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>

            <Workflows />
          </div>
          <Col>
                  <ListGroup className='mb-3'>
                  <ListGroupItem>Menu Item 1</ListGroupItem>
                  <ListGroupItem>Menu Item 2</ListGroupItem>
                  </ListGroup>
                  <CreateWorkflow />
          </Col>
        </div>
      </div>
    </div>
  );
}

export default DashboardWorkflow;
