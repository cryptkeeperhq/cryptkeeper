// import React, { useState, useEffect } from 'react';
// import { Card, Container, Breadcrumb, ListGroup, ListGroupItem, Row, Col, Form, Button, FloatingLabel, Alert } from 'react-bootstrap';
// import SecretAccessControl from './SecretAccessControl';
// import RotateSecret from './RotateSecret';
// import SecretMetadataUpdate from './SecretMetadataUpdate';
// import Accordion from 'react-bootstrap/Accordion';
// import Datetime from 'react-datetime';
// import SecretExistingAccess from './SecretExistingAccess';
// import "react-datetime/css/react-datetime.css";
// import { FaUser, FaKey, FaHistory, FaTrashAlt, FaEdit, FaPlus } from 'react-icons/fa';

// const UserSecretsAlt = () => {
//     const [secrets, setSecrets] = useState([]);
//     const [selectedSecretPath, setSelectedSecretPath] = useState(null);
//     const [selectedKey, setSelectedKey] = useState(null);
//     const [selectedVersion, setSelectedVersion] = useState(null);
//     const [searchParams, setSearchParams] = useState({ path: '', metadata: '', createdAfter: '', createdBefore: '' });
//     const [keys, setKeys] = useState([]);
//     const [versions, setVersions] = useState([]);
//     const [expiresAt, setExpiresAt] = useState('');
//     const [sharedLink, setSharedLink] = useState(null);

//     useEffect(() => {
//         fetch('/api/user/secrets', {
//             headers: {
//                 'Authorization': `Bearer ${token}`
//             }
//         })
//             .then(res => res.json())
//             .then(data => {
//                 if (Array.isArray(data)) {
//                     setSecrets(data);
//                 } else {
//                     setSecrets([]);
//                 }
//             })
//             .catch(error => {
//                 console.error('Error fetching user secrets:', error);
//                 setSecrets([]);
//             });
//     }, []);

//     const handleSearch = () => {
//         const query = new URLSearchParams();
//         if (searchParams.path) query.append('path', searchParams.path);
//         if (searchParams.key) query.append('key', searchParams.key);
//         if (searchParams.metadata) query.append('metadata', searchParams.metadata);
//         if (searchParams.createdAfter) query.append('created_after', searchParams.createdAfter);
//         if (searchParams.createdBefore) query.append('created_before', searchParams.createdBefore);

//         fetch(`/api/search-secrets?${query.toString()}`, {
//             headers: {
//                 'Authorization': `Bearer ${token}`
//             }
//         })
//             .then(res => res.json())
//             .then(data => {
//                 if (Array.isArray(data)) {
//                     setSecrets(data);
//                 } else {
//                     setSecrets([]);
//                 }
//             })
//             .catch(error => {
//                 console.error('Error searching secrets:', error);
//                 setSecrets([]);
//             });
//     };

//     const handleSecretClick = (path) => {
//         setSelectedSecretPath(path);
//         setSelectedKey(null);
//         setSelectedVersion(null);
//         setSharedLink(null);
//         const secretKeys = [...new Set(secrets.filter(secret => secret.path === path).map(secret => secret.key))];
//         setKeys(secretKeys);
//     };

//     const handleKeyClick = (key) => {
//         setSelectedKey(key);
//         setSelectedVersion(null);
//         setSharedLink(null);
//         const secretVersions = secrets.filter(secret => secret.path === selectedSecretPath && secret.key === key);
//         setVersions(secretVersions);
//     };

//     const handleVersionClick = (version) => {
//         setSharedLink(null);
//         setSelectedVersion(version);
//     };

//     const createSharedLink = (secretId) => {
//         if (!expiresAt) {
//             alert('Please select an expiration date and time.');
//             return;
//         }

//         fetch('/api/create-shared-link', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`
//             },
//             body: JSON.stringify({ secret_id: secretId, expires_at: expiresAt })
//         })
//             .then(res => res.json())
//             .then(data => {
//                 setSharedLink(`Shared link created: ${window.location.origin}/shared/${data.link_id}`);
//             })
//             .catch(error => console.error('Error creating shared link:', error));
//     };

//     const deleteSecret = (path, key, version) => {
//         fetch(`/api/secrets?path=${encodeURIComponent(path)}&key=${encodeURIComponent(key)}&version=${version}`, {
//             method: 'DELETE',
//             headers: {
//                 'Authorization': `Bearer ${token}`
//             }
//         })
//             .then(() => {
//                 setSecrets(secrets.filter(secret => !(secret.path === path && secret.key === key && secret.version === version)));
//                 alert('Secret deleted.');
//             })
//             .catch(error => {
//                 console.error('Error deleting secret:', error);
//             });
//     };

//     const uniqueSecretPaths = [...new Set(secrets.map(secret => secret.path))];

//     return (
//         <div className="mt-4">
//             <Container>
//                 <Row>
//                     <Col xs={12} lg={12} className='mb-3'>
//                         <Accordion>
//                             <Accordion.Item eventKey="0">
//                                 <Accordion.Header>Advanced Search</Accordion.Header>
//                                 <Accordion.Body>
//                                     <Form>
//                                         <FloatingLabel controlId="floatingInput" label="Path" className="mb-3">
//                                             <Form.Control
//                                                 type="text"
//                                                 value={searchParams.path}
//                                                 onChange={(e) => setSearchParams({ ...searchParams, path: e.target.value })}
//                                             />
//                                         </FloatingLabel>

//                                         <FloatingLabel controlId="floatingInput" label="Key" className="mb-3">
//                                             <Form.Control
//                                                 type="text"
//                                                 value={searchParams.key}
//                                                 onChange={(e) => setSearchParams({ ...searchParams, key: e.target.value })}
//                                             />
//                                         </FloatingLabel>

//                                         <FloatingLabel controlId="floatingInput" label="Metadata (JSON)" className="mb-3">
//                                             <Form.Control
//                                                 type="text"
//                                                 value={searchParams.metadata}
//                                                 onChange={(e) => setSearchParams({ ...searchParams, metadata: e.target.value })}
//                                             />
//                                         </FloatingLabel>

//                                         <FloatingLabel controlId="floatingInput" label="Created After" className="mb-3">
//                                             <Form.Control
//                                                 type="date"
//                                                 value={searchParams.createdAfter}
//                                                 onChange={(e) => setSearchParams({ ...searchParams, createdAfter: e.target.value })}
//                                             />
//                                         </FloatingLabel>

//                                         <FloatingLabel controlId="floatingInput" label="Created Before" className="mb-3">
//                                             <Form.Control
//                                                 type="date"
//                                                 value={searchParams.createdBefore}
//                                                 onChange={(e) => setSearchParams({ ...searchParams, createdBefore: e.target.value })}
//                                             />
//                                         </FloatingLabel>
//                                         <Button variant="primary" className='w-100 mt-2' onClick={handleSearch}>Search</Button>
//                                     </Form>
//                                 </Accordion.Body>
//                             </Accordion.Item>
//                         </Accordion>
//                     </Col>

// <Col xs={12} lg={12}>
// <Breadcrumb>
      
//         <Breadcrumb.Item></Breadcrumb.Item>
//         {selectedSecretPath && <Breadcrumb.Item>{selectedSecretPath}</Breadcrumb.Item>}
//         {selectedKey && <Breadcrumb.Item>{selectedKey}</Breadcrumb.Item>}
//         {selectedVersion && <Breadcrumb.Item>{selectedVersion.version}</Breadcrumb.Item>}

//     </Breadcrumb>

// </Col>
//                     <Col xs={12} lg={2}>
//                         <Card>
//                             <Card.Header>Your Secrets</Card.Header>
//                             <ListGroup variant="flush">
                                
//                                 {uniqueSecretPaths.map(path => (
//                                     <ListGroupItem
//                                         key={path}
//                                         className={`list-group-item ${selectedSecretPath === path ? 'active' : ''}`}
//                                         onClick={() => handleSecretClick(path)}
//                                         style={{ cursor: 'pointer' }}
//                                     >
//                                         <FaKey /> <b>{path}</b>
//                                     </ListGroupItem>
//                                 ))}
//                             </ListGroup>
//                         </Card>
//                     </Col>

//                     <Col xs={12} lg={2}>
//                         {selectedSecretPath && (
//                             <Card className="">
//                                 <Card.Header>Keys for {selectedSecretPath}</Card.Header>
//                                 <ListGroup variant="flush">
//                                     {keys.map(key => (
//                                         <ListGroupItem
//                                             key={key}
//                                             className={`list-group-item ${selectedKey === key ? 'active' : ''}`}
//                                             onClick={() => handleKeyClick(key)}
//                                             style={{ cursor: 'pointer' }}
//                                         >
//                                             <b>{key}</b>
//                                         </ListGroupItem>
//                                     ))}
//                                 </ListGroup>
//                             </Card>
//                         )}
//                     </Col>

//                     <Col xs={12} lg={2}>
//                         {selectedKey && (
//                             <Card className="">
//                                 <Card.Header>Versions for {selectedSecretPath}/{selectedKey}</Card.Header>
//                                 <ListGroup variant="flush">
//                                     {versions.map(secret => (
//                                         <ListGroupItem
//                                             key={secret.version}
//                                             className={`list-group-item ${selectedVersion && selectedVersion.version === secret.version ? 'active' : ''}`}
//                                             onClick={() => handleVersionClick(secret)}
//                                             style={{ cursor: 'pointer' }}
//                                         >
//                                             <b>
//                                                 Version {secret.version}
//                                                 {secret.is_one_time && <span className="badge bg-warning ms-2">One-Time</span>}
//                                                 {secret.expires_at && new Date(secret.expires_at) < new Date() && <span className="badge bg-danger ms-2">Expired</span>}
//                                             </b>
//                                         </ListGroupItem>
//                                     ))}
//                                 </ListGroup>
//                             </Card>
//                         )}
//                     </Col>

//                     <Col xs={12} lg={6}>
//                         {selectedVersion && (
//                             <Card className="">
//                                 <Card.Header>Details for {selectedSecretPath}/{selectedKey} (v{selectedVersion.version})</Card.Header>
//                                 <Card.Body>
//                                     <button className="btn btn-danger btn-sm float-end" onClick={(e) => { e.stopPropagation(); deleteSecret(selectedVersion.path, selectedVersion.key, selectedVersion.version); }}>Delete</button>

//                                     <h6 className='m-0'>Path: <code>{selectedSecretPath}</code></h6>
//                                     <h6 className='m-0'>Key: <code>{selectedKey}</code></h6>
//                                     <div><strong>Version: <code>{selectedVersion.version}</code></strong></div>
//                                     <div><strong>Value: <code>{selectedVersion.value}</code></strong></div>
//                                     <div><strong>Metadata:</strong> <pre>{JSON.stringify(selectedVersion.metadata)}</pre></div>


// { selectedVersion.rotation_interval && <>
// <Alert variant='success' className=' small p-2'><b>Rotation is enabled every <code>{selectedVersion.rotation_interval}</code>.<br/>The secret was last rotated at <code>{selectedVersion.last_rotated_at}</code></b></Alert>
// </>}

//                                     <Accordion className='mt-3'>
//                                         <Accordion.Item eventKey="0">
//                                             <Accordion.Header>Rotate Secret</Accordion.Header>
//                                             <Accordion.Body>
//                                                 <RotateSecret path={selectedSecretPath} selected_key={selectedKey} />
//                                             </Accordion.Body>
//                                         </Accordion.Item>
//                                         <Accordion.Item eventKey="1">
//                                             <Accordion.Header>Access Control</Accordion.Header>
//                                             <Accordion.Body>
//                                                 <SecretAccessControl secretId={selectedVersion.id} />
//                                                 <SecretExistingAccess secret={selectedVersion} />
//                                             </Accordion.Body>
//                                         </Accordion.Item>
//                                         <Accordion.Item eventKey="2">
//                                             <Accordion.Header>Metadata Update</Accordion.Header>
//                                             <Accordion.Body>
//                                                 <SecretMetadataUpdate secret={selectedVersion} />
//                                             </Accordion.Body>
//                                         </Accordion.Item>

//                                         <Accordion.Item eventKey="3">
//                                             <Accordion.Header>Share One Time Link</Accordion.Header>
//                                             <Accordion.Body>

//                                                 { sharedLink && 
//                                                 <Alert className='p-2' variant='info'>{sharedLink}</Alert>
//                                                 }
//                                                 <Datetime
//                                                     value={expiresAt}
//                                                     onChange={setExpiresAt}
//                                                     dateFormat="YYYY-MM-DD"
//                                                     timeFormat="HH:mm:ss"
//                                                     inputProps={{ placeholder: 'YYYY-MM-DD HH:mm:ss' }}
//                                                 />
//                                                 <Button variant="primary btn-sm w-100 mt-1" onClick={() => createSharedLink(selectedVersion.id)}>Create Shared Link</Button>
//                                             </Accordion.Body>
//                                         </Accordion.Item>
//                                     </Accordion>
//                                 </Card.Body>
//                             </Card>
//                         )}
//                     </Col>
//                 </Row>
//             </Container>
//         </div>
//     );
// };

// export default UserSecretsAlt;
