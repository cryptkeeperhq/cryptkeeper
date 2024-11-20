import React, { useState, useEffect } from 'react';
import { Button, Alert, Card, InputGroup, ListGroup, ListGroupItem } from 'react-bootstrap';
import RenderMetadataFields from './RenderMetadataFields';
import { useApi } from '../api/api'

import Datetime from 'react-datetime';
import CLIUsage from './CLIUsage';
import { FaPlus } from 'react-icons/fa';

const ShareSecret = ({ path, secret }) => {
  const { get, post, put, del } = useApi();
  const [message, setMessage] = useState("")
  const [error, setError] = useState("")
  const [expiresAt, setExpiresAt] = useState('');
  const [sharedLink, setSharedLink] = useState(null);

  useEffect(() => {
    setExpiresAt('')
    setMessage('')
    setError('')
    setSharedLink(null)
  }, [path, secret]);


  const createSharedLink = async (path, secret) => {
    if (!expiresAt) {
      alert('Please select an expiration date and time.');
      return;
    }

    try {
      const data = await post(`/secrets/share?path=${path.path}&key=${secret.key}`, { secret_id: secret.id, expires_at: expiresAt, version: secret.version });
      setSharedLink(`Shared link created: ${window.location.origin}/shared/${data.link_id}`);
    } catch (error) {
      setError(error.message);
    }


  };

  return (
    <div>
      {secret && <div>

        <Card>
          <Card.Header><FaPlus className='me-2' /> Share One Time Link</Card.Header>
          <Card.Body>
            {sharedLink && <div className='p-2 text-success bg-success-soft fw-bold mb-2'>{sharedLink}</div>}

            {message && <p className='text-success fw-bold'>{message}</p>}
            {error && <p className='text-danger fw-bold'>{error}</p>}

            <Datetime
              value={expiresAt}
              onChange={setExpiresAt}
              dateFormat="YYYY-MM-DD"
              timeFormat="HH:mm:ss"
              inputProps={{ placeholder: 'YYYY-MM-DD HH:mm:ss' }}
            />
            <Button variant="primary btn-sm w-100 mt-1" onClick={() => createSharedLink(path, secret)}>Create Shared Link</Button>
          </Card.Body>
          <Card.Footer>
          <CLIUsage cmd="share [path] [key] --expires-at [date]" />
          </Card.Footer>
        </Card>
        

        <ListGroup className='mt-3'>
                                <ListGroupItem>Shared Links here (TODO)</ListGroupItem>
                            </ListGroup>


      </div>}


    </div>
  );
};

export default ShareSecret;
