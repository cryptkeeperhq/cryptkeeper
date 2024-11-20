import React, { useState, useEffect } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import { Navbar, Nav, Container } from 'react-bootstrap';
import { Image, Stack, Modal, Card, Button } from 'react-bootstrap';
import { FaBell, FaHandsHelping, FaKey, FaSignOutAlt, FaUser } from 'react-icons/fa';
import { Link, Navigate } from 'react-router-dom';
import Notifications from '../Notifications';
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import logo from '../../assets/logo.webp'

function Title({ heading, subheading }) {
  const username = localStorage.getItem("user")
  const token = localStorage.getItem('token');

  const handleLogout = () => {
    localStorage.removeItem('token');
  };

  return (

    <Navbar variant='fixed-top' expand="lg">
      <Container className='ms-0 me-4'>
        <Navbar.Brand className='d-flex align-items-center'>
        {/* <div className='text-center float-start me-2' style={{ width: "60px" }}><img className='rounded-circle w-75' src={logo} /></div> */}

        {/* <FaKey size={40} className='me-2 text-muted' /> */}

<div className='ms-auto'>
  
          {heading && <h1 className='m-0 mt-3 p-0'>
          {heading}</h1>}
          {subheading && <p className='fs-6 p-0' dangerouslySetInnerHTML={{ __html: subheading }}></p>}

          {/* <b dangerouslySetInnerHTML={{ __html: heading }} /> */}
          {/* <Breadcrumb className='small' style={{ fontSize: "14px"}}>
      <Breadcrumb.Item className='' href="/#/">Home</Breadcrumb.Item>
      
        
        {heading && <Breadcrumb.Item active>{heading}</Breadcrumb.Item>}
        {subheading && <Breadcrumb.Item active>{subheading}</Breadcrumb.Item>}
      
    </Breadcrumb> */}
    </div>

        </Navbar.Brand>
        {/* <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto ">
            <Notifications />

            {username && (
              <Nav.Item className="ms-3 d-flex align-items-top">
                <div className="text-end">
                  <div>
                  <Link className='p-0 nav-link d-inline' to="/logout" onClick={handleLogout}><FaSignOutAlt /></Link></div>

                </div>
              </Nav.Item>
            )}

          </Nav>
        </Navbar.Collapse> */}
      </Container>
    </Navbar>

  );
}

export default Title;
