// src/components/SecretHistory.js
import moment from "moment";

import React, { useState, useEffect } from 'react';
import { Card, Container, ListGroup, ListGroupItem, Row, Col, InputGroup } from 'react-bootstrap';
import MetadataItem from "./MetadataItem";
import { useApi } from '../api/api'
import { Link } from "react-router-dom";
import { FaHistory, FaLink } from "react-icons/fa";
import CLIUsage from "./CLIUsage";

const SecretHistory = ({ token, path, secret, onVersionSelect }) => {
    const { get, post, put, del } = useApi();
    // const [path, setPath] = useState('');
    const [history, setHistory] = useState([]);

    const getSecretHistory = async () => {
        try {
            const data = await get(`/secrets/history?path=${path.path}&key=${secret.key}`);
            setHistory(data)
        } catch (error) {
            console.log(error.message);
            setHistory([]);
        }


    };

    const onSelect = async (secret) => {
        onVersionSelect(secret)
        // onVersionSelect(path.path, secret.key, )
        // return

    }


    useEffect(() => {
        getSecretHistory();
    }, [token, secret]);


    return (
        <div>
            <Card>
                <Card.Header><FaHistory className='me-2' /> Version History</Card.Header>


                <ListGroup variant='flush'>
                    {history.map(hsecret => (
                        <>
                            {/* {secret.id !== hsecret.id && */}
                            <ListGroupItem key={hsecret.version} className='small' onClick={(e) => { e.stopPropagation(); onSelect(hsecret); }} >
                                <Link className='float-end text-decoration-none' to={`/user/secrets/${hsecret.id}/${hsecret.version}`}><FaLink className='d-line me-1' /></Link>

                                <div>Version: <code>{hsecret.version}</code></div>
                                <MetadataItem metadata={hsecret.metadata} />
                                {/* Metadata: <code>{JSON.stringify(secret.metadata)}</code><br/> */}
                                <div>Created by <strong>{hsecret.created_by}</strong> on {moment(hsecret.created_at).format("MMMM Do YYYY, h:mm:ss a")}</div>
                                <div>Updated {moment(hsecret.udpated_at).format("MMMM Do YYYY, h:mm:ss a")}</div>

                            </ListGroupItem>
                            {/* } */}
                        </>

                    ))}
                </ListGroup>
                <Card.Footer>
                <CLIUsage cmd="versions [path] [key]" />

                


                </Card.Footer>
            </Card>
        </div>
    );
};

export default SecretHistory;
