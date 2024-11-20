// import React, { useState, useEffect } from 'react';
// import { Container, Form, Button, Alert } from 'react-bootstrap';
// import { get, post, put, del } from '../api/api';

       
// const CreatePolicy = () => {
//     const [name, setName] = useState('');
//     const [description, setDescription] = useState('');
//     // const [rules, setRules] = useState('');
//     const [paths, setPaths] = useState([]);
//     const [selectedPath, setSelectedPath] = useState('');
//     const [message, setMessage] = useState('');
//     const [error, setError] = useState('');

//     const [rules, setRules] = useState(JSON.stringify({
//         path: '',
//         rules: {
//             read: { groups: [], users: [] },
//             write: { groups: [], users: [] },
//             delete: { groups: [], users: [] },
//         },
//     }, null, 4));


//     useEffect(() => {
//         fetchPaths();
//     }, []);

//     const fetchPaths = async() => {
//         try {
//             const data = await get(`/paths`);
//             setPaths(data || []);
//         } catch (error) {
//             console.log(error.message);
//         }
//     }



//     const handleCreatePolicy = async () => {
//         try {
//             const data = await post(`/policies`, { name, description, rules });
//             setMessage('Policy created successfully.');
//             assignPolicyToGroup(data.id);
//         } catch (error) {
//             console.log(error.message);
//         }

//         const response = await fetch('/api/policies', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`
//             },
//             body: JSON.stringify({ name, description, rules })
//         });

//     };

//     const assignPolicyToGroup = async (policyID) => {
//         const response = await fetch('/api/assign-policy', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`
//             },
//             body: JSON.stringify({ path_id: parseInt(selectedPath), policy_id: policyID })
//         });

//         if (response.ok) {
//             setMessage('Policy assigned to group successfully.');
//         } else {
//             const errorText = await response.text();
//             setError(`Error assigning policy to group: ${errorText}`);
//         }
//     };

//     return (
//         <Container className='mt-4'>
//             {message && <Alert variant='success'>{message}</Alert>}
//             {error && <Alert variant='danger'>{error}</Alert>}
//             <Form>
//                 <Form.Group controlId="policyName">
//                     <Form.Label>Policy Name</Form.Label>
//                     <Form.Control
//                         type="text"
//                         value={name}
//                         onChange={(e) => setName(e.target.value)}
//                         placeholder="Enter policy name"
//                     />
//                 </Form.Group>

//                 <Form.Group controlId="policyDescription" className='mt-2'>
//                     <Form.Label>Policy Description</Form.Label>
//                     <Form.Control
//                         type="text"
//                         value={description}
//                         onChange={(e) => setDescription(e.target.value)}
//                         placeholder="Enter policy description"
//                     />
//                 </Form.Group>

//                 <Form.Group controlId="policyRules" className='mt-2'>
//                     <Form.Label>Policy Rules (JSON)</Form.Label>
//                     <Form.Control
//                         as="textarea"
//                         rows={5}
//                         value={rules}
//                         onChange={(e) => setRules(e.target.value)}
//                         placeholder='Enter policy rules in JSON format, e.g., {"actions": ["read", "write"], "resources": ["/path/to/resource"]}'
//                     />
//                 </Form.Group>

//                 <Form.Group controlId="selectPath" className='mt-3'>
//                     <Form.Label>Select Path</Form.Label>
//                     <Form.Control
//                         as="select"
//                         value={selectedPath}
//                         onChange={(e) => setSelectedPath(e.target.value)}
//                     >
//                         <option value="">Select Path</option>
//                         {paths.map(path => (
//                             <option key={path.id} value={path.id}>
//                                 {path.path}
//                             </option>
//                         ))}
//                     </Form.Control>
//                 </Form.Group>

//                 <Button variant="primary" className='mt-3' onClick={handleCreatePolicy}>
//                     Create Policy and Assign to Group
//                 </Button>
//             </Form>
//         </Container>
//     );
// };

// export default CreatePolicy;
