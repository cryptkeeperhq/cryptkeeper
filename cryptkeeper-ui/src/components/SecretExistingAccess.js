import React, { useState, useEffect } from 'react';
import { Card, InputGroup } from 'react-bootstrap';
import { useApi } from '../api/api'


const SecretExistingAccess = ({ token, secret }) => {
  const { get, post, put, del } = useApi();
  const [secretAccesses, setSecretAccesses] = useState({});

  useEffect(() => {
    // Fetch secret accesses for each secret


    fetchAccesses();
  }, [secret, token]);

  const fetchAccesses = async () => {
    const accesses = {};
    try {
      const data = await get(`/paths/${secret.path_id}/secrets/${secret.id}/accesses`);
      accesses[secret.id] = data;
      setSecretAccesses(accesses);
    } catch (error) {
      console.log(error.message);
    }

  };


  return (
    <div className='mt-2'>

      <Card className='mt-3'>
        <Card.Header className='bg-light text-dark'>Membership Info</Card.Header>
        <Card.Body>
          <div>
            <strong>Owner:</strong> {secretAccesses[secret.id]?.find(a => a.access_level === 'owner')?.username}
          </div>
          <div>
            <strong>Groups with access:</strong>
            <ul>
              {secretAccesses[secret.id]?.filter(a => a.group_name).map(access => (
                <li key={access.group_id}>{access.group_name}</li>
              ))}
            </ul>
          </div>
          <div>
            <strong>Users with access:</strong>
            <ul>
              {secretAccesses[secret.id]?.filter(a => a.username).map(access => (
                <li key={access.user_id}>{access.username}</li>
              ))}
            </ul>
          </div></Card.Body></Card>
    </div>
  );
};

export default SecretExistingAccess;
