import React, { useState } from 'react';
import { Card, Form, Button, ListGroup, ListGroupItem, FloatingLabel } from 'react-bootstrap';
import { useApi } from '../api/api'
import { Link } from 'react-router-dom';

const SearchSecrets = () => {
  const { get, post, put, del } = useApi();
  const [path, setPath] = useState('');
  const [metadata, setMetadata] = useState('');
  const [createdAfter, setCreatedAfter] = useState('');
  const [createdBefore, setCreatedBefore] = useState('');
  const [results, setResults] = useState([]);
  const [searchParams, setSearchParams] = useState({ path: '', key: '', metadata: '', createdAfter: '', createdBefore: '' });

  const handleSearch = async () => {
    const query = new URLSearchParams();
    // if (path) query.append('path', path);
    // if (metadata) query.append('metadata', metadata);
    // if (createdAfter) query.append('created_after', createdAfter);
    // if (createdBefore) query.append('created_before', createdBefore);




    try {
      const query = new URLSearchParams();
      if (searchParams.path) query.append('path', searchParams.path);
      if (searchParams.key) query.append('key', searchParams.key);
      if (searchParams.metadata) query.append('metadata', searchParams.metadata);
      if (searchParams.createdAfter) query.append('created_after', searchParams.createdAfter);
      if (searchParams.createdBefore) query.append('created_before', searchParams.createdBefore);

      try {
        const data = await get(`/search-secrets?${query.toString()}`);
        setResults(data || [])
      } catch (error) {
        console.log(error.message);
      }


    } catch (error) {
      console.log(error.message);
    }
  };

  return (
    <div>
      {/* <Card.Header>Advanced Search</Card.Header>
      <Card.Body> */}
      <Form>
        <FloatingLabel controlId="floatingInput" label="Path" className="mb-3">
          <Form.Control
            type="text"
            value={searchParams.path}
            onChange={(e) => setSearchParams({ ...searchParams, path: e.target.value })}
          />
        </FloatingLabel>

        <FloatingLabel controlId="floatingInput" label="Key" className="mb-3">
          <Form.Control
            type="text"
            value={searchParams.key}
            onChange={(e) => setSearchParams({ ...searchParams, key: e.target.value })}
          />
        </FloatingLabel>

        <FloatingLabel controlId="floatingInput" label="Metadata (JSON)" className="mb-3">
          <Form.Control
            type="text"
            value={searchParams.metadata}
            onChange={(e) => setSearchParams({ ...searchParams, metadata: e.target.value })}
          />
        </FloatingLabel>

        <FloatingLabel controlId="floatingInput" label="Created After" className="mb-3">
          <Form.Control
            type="date"
            value={searchParams.createdAfter}
            onChange={(e) => setSearchParams({ ...searchParams, createdAfter: e.target.value })}
          />
        </FloatingLabel>

        <FloatingLabel controlId="floatingInput" label="Created Before" className="mb-3">
          <Form.Control
            type="date"
            value={searchParams.createdBefore}
            onChange={(e) => setSearchParams({ ...searchParams, createdBefore: e.target.value })}
          />
        </FloatingLabel>
        <Button variant="primary" className='w-100 mt-2' onClick={handleSearch}>Search</Button>
      </Form>
      {/* </Card.Body> */}
      {results.length > 0 && (
        <>
          <h5 className='m-3'>Search Results</h5>
          <ListGroup variant="flush">
            {results.map(secret => (
              <ListGroupItem key={secret.id}>
                <div>
                  <Link className='text-primary text-decoration-none' to={`/user/secrets/${secret.id}/${secret.version}`}>
                    <strong>{secret.path}{secret.key}</strong>
                  </Link></div>
                <div><strong>Version:</strong> {secret.version} | <strong>Created At:</strong> {new Date(secret.created_at).toLocaleString()}</div>
                <div><strong>Metadata:</strong> <code>{JSON.stringify(secret.metadata)}</code></div>
              </ListGroupItem>
            ))}
          </ListGroup>
        </>
      )}
    </div>
  );
};

export default SearchSecrets;
