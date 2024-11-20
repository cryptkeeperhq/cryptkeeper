// import React, { useState, useEffect } from 'react';
// import { Container, Row, Col, Card, Form, Button, Alert } from 'react-bootstrap';
// import { get, post, put, del } from '../../api/api';

// const PolicyEditor = () => {
//     const [policies, setPolicies] = useState([]);
//     const [selectedPolicy, setSelectedPolicy] = useState(null);
//     const [message, setMessage] = useState('');
//     const [paths, setPaths] = useState([]);

//     useEffect(() => {
//         fetchPolicies();
//         fetchPaths();
//     }, []);

//     const fetchPolicies = async () => {
//         try {
//             const data = await get(`/policies`);
//             setPolicies(data || [])
//         } catch (error) {
//             console.log(error.message);
//         }
//     };

//     const fetchPaths = async () => {
//         const response = await get('/paths');
//         setPaths(response.data || []);
//     };


//     const handleSave = async () => {

//         try {
//             selectedPolicy.rules = JSON.parse(selectedPolicy.rules)
//             const data = await post(`/policies`, selectedPolicy);
//             setMessage('Policy saved successfully.');
//             const updatedPolicies = policies.map(policy => policy.id === data.id ? data : policy);
//             if (!selectedPolicy.id) {
//                 updatedPolicies.push(data);
//             }
//             setPolicies(updatedPolicies);
//         } catch (error) {
//             setMessage(error.message);
//         }

//     };

//     const handleDelete = () => {
//         if (selectedPolicy && selectedPolicy.id) {
//             fetch(`/api/api/policies/${selectedPolicy.id}`, {
//                 method: 'DELETE',
//                 headers: {
//                     'Authorization': `Bearer ${token}`
//                 }
//             })
//                 .then(() => {
//                     setMessage('Policy deleted successfully.');
//                     setPolicies(policies.filter(policy => policy.id !== selectedPolicy.id));
//                     setSelectedPolicy(null); // Deselect policy after deleting
//                 })
//                 .catch(error => {
//                     console.error('Error deleting policy:', error);
//                     setMessage('Error deleting policy.');
//                 });
//         }
//     };

//     const handleSelectPolicy = (policy) => {
//         setSelectedPolicy({ ...policy });
//     };

//     const handleChange = (e) => {
//         const { name, value } = e.target;
//         setSelectedPolicy({ ...selectedPolicy, [name]: value });
//     };

//     return (
//         <Container>
//             <Row>
//                 <Col md={4}>
//                     <Card>
//                         <Card.Header>Policies</Card.Header>
//                         <Card.Body>
//                             <ul className="list-group">
//                                 {policies.map(policy => (
//                                     <li key={policy.id} className="list-group-item" onClick={() => handleSelectPolicy(policy)}>
//                                         {policy.name}
//                                     </li>
//                                 ))}
//                             </ul>
//                         </Card.Body>
//                     </Card>
//                 </Col>
//                 <Col md={8}>
//                     {selectedPolicy && (
//                         <Card>
//                             <Card.Header>Edit Policy</Card.Header>
//                             <Card.Body>
//                                 {message && <div className='bg-primary-soft text-primary p-3 rounded-2 mb-2'>{message}</div>}
//                                 <Form>
//                                     <Form.Group controlId="policyName">
//                                         <Form.Label>Name</Form.Label>
//                                         <Form.Control
//                                             type="text"
//                                             name="name"
//                                             value={selectedPolicy.name}
//                                             onChange={handleChange}
//                                         />
//                                     </Form.Group>
//                                     <Form.Group controlId="policyName">
//                                         <Form.Label>Description</Form.Label>
//                                         <Form.Control
//                                             type="text"
//                                             name="description"
//                                             value={selectedPolicy.description}
//                                             onChange={handleChange}
//                                         />
//                                     </Form.Group>
//                                     <Form.Group controlId="policyRules" className="mt-3">
//                                         <Form.Label>Rules (JSON)</Form.Label>
//                                         <Form.Text
//                                             as="textarea"
//                                             rows={10}
//                                             name="rules"
//                                             style={{ fontFamily: "Menlo" }}
//                                             value={JSON.stringify(selectedPolicy.rules, null, 4)}
//                                             onChange={handleChange}
//                                         />
//                                     </Form.Group>
//                                     <Button className="mt-3" onClick={handleSave}>Save Policy</Button>
//                                     <Button variant="danger" className="mt-3 ms-2" onClick={handleDelete}>Delete Policy</Button>

//                                 </Form>
//                             </Card.Body>
//                         </Card>
//                     )}
//                 </Col>
//             </Row>
//         </Container>
//     );
// };

// export default PolicyEditor;
