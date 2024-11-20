import React, { useState, useEffect } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Container from 'react-bootstrap/Container'
import { Image, Stack, Modal, Card, Button } from 'react-bootstrap';
import { FaHandsHelping } from 'react-icons/fa';
import ThemeToggle from '../theme/ThemeToggle';
import { Link, Navigate } from 'react-router-dom';

function HelpCard() {


  return (
    <div id="help">
      <Card className=' bg-dark text-white border border-light text-center'>
        <Card.Body className='text-center'>
          <FaHandsHelping size={64} />
          <div className=''>
            <h5 className='m-0 p-0 mb-2'>Need Help?</h5>
            <Link to="/" className='btn btn-primary mt-2 rounded-2 w-100 '>Read More</Link>
          </div>
        </Card.Body>
      </Card>

      <div className='mt-2'>
        <ThemeToggle />
      </div>


    </div>
  );
}

export default HelpCard;
