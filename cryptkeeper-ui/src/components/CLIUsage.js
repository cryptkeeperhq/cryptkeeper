import React, { useState, useEffect } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Container from 'react-bootstrap/Container'
import { Image, Stack, Modal, Card, Button } from 'react-bootstrap';
import { FaHandsHelping, FaTerminal } from 'react-icons/fa';
import ThemeToggle from '../theme/ThemeToggle';
import { Link, Navigate } from 'react-router-dom';
import CodeEditor from './common/CodeEditor';
import {  json } from '@codemirror/lang-json';

function CLIUsage({ cmd }) {

    const updatedCmd = "./cryptkeeper-cli " + cmd

    return (
        <div id="cli_usage" className="d-none">
            <div >
                {/* <b>CLI Usage: </b><br /> */}
                <CodeEditor code={updatedCmd} height={50}  />
                {/* <pre className='p-0 m-0 ps-2 '><FaTerminal className='me-2' /><span className='text-success'>./cryptkeeper-cli</span> {cmd}</pre> */}
            </div>
        </div>
    );
}

export default CLIUsage;
