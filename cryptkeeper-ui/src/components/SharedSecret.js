import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Card, ListGroup, ListGroupItem } from 'react-bootstrap';
import SecretDetails from './SecretDetails';
import Title from './common/Title';
import { useApi } from '../api/api'



const SharedSecret = ({setTitle}) => {
    const { get, post, put, del } = useApi();
    const { linkID } = useParams();
    const [secret, setSecret] = useState(null);
    const [error, setError] = useState('');

    useEffect(() => {
        setTitle({ heading: "Shared Secret", subheading: "You were shared"})
        fetchSecret()
    }, [linkID]);

    const fetchSecret = async () => {
        try {
            const data = await get(`/access-shared-link/${linkID}`);
            setSecret(data);
        } catch (error) {
            setError(error.message);
        }
    }

    if (error) {
        return <div>{error}</div>;
    }

    if (!secret) {
        return <div>Loading...</div>;
    }

    return (
        <div className='vh-100 d-flex align-items-center justify-content-center'>
        <div className='p-3 mx-auto w-100'>
            <SecretDetails path="" secret={secret} />
        </div>
        </div>
    );
};

export default SharedSecret;
