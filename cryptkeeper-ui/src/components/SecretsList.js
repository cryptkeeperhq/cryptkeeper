// import React, { useState, useEffect } from 'react';
// import { Card, ListGroup, ListGroupItem, InputGroup } from 'react-bootstrap';
// import RotateSecret from './RotateSecret';
// import SecretAccessControl from './SecretAccessControl';
// import SecretMetadataUpdate from './SecretMetadataUpdate';

// const SecretsList = () => {
//   const [secrets, setSecrets] = useState([]);
//   const [path, setPath] = useState('');
//   const [secretAccesses, setSecretAccesses] = useState({});

//   useEffect(() => {
//     // Fetch secret accesses for each secret
//     const fetchAccesses = async () => {
//       const accesses = {};
//       for (const secret of secrets) {
//         const response = await fetch(`/api/paths/${secret.path_id}/secrets/${secret.id}/accesses`, {
//           headers: {
//             'Authorization': `Bearer ${token}`
//           }
//         });
//         const data = await response.json();
//         accesses[secret.id] = data;
//       }
//       setSecretAccesses(accesses);
//     };

//     fetchAccesses();
//   }, [secrets, token]);

//   const deleteSecret = (path, version) => {
//     fetch(`/api/paths/${path_id}/secrets?path=${encodeURIComponent(path)}&version=${version}`, {
//       method: 'DELETE',
//       headers: {
//         'Authorization': `Bearer ${token}`
//       }
//     })
//       .then(() => {
//         setSecrets(secrets.filter(secret => !(secret.path === path && secret.version === version)));
//         alert('Secret deleted.');
//       })
//       .catch(error => {
//         console.error('Error deleting secret:', error);
//       });
//   };

//   useEffect(() => {
//     if (path) {
//       console.log(`Fetching secrets for path: ${path}`);
//       fetch(`/api/paths/${path}/secrets?path=${encodeURIComponent(path)}`, {
//         headers: {
//           'Authorization': `Bearer ${token}`
//         }
//       })
//         .then(res => res.json())
//         .then(data => {
//           console.log('Fetched secrets:', data);
//           if (Array.isArray(data)) {
//             setSecrets(data);
//           } else {
//             setSecrets([]);
//           }
//         })
//         .catch(error => {
//           console.error('Error fetching secrets:', error);
//           setSecrets([]);
//         });
//     }
//   }, [path, token]);


//   return (
//     <div className="">
//       <Card>
//         <Card.Header>Secret List</Card.Header>
//         <div className="form-floating border-0">
//           <input type="text" className="form-control border-0" id="createPath" placeholder="Path" value={path} onChange={e => setPath(e.target.value)} />
//           <label htmlFor="createPath">Path</label>
//         </div>
//         <ListGroup variant="flush">
//           {secrets.map(secret => (
//             <ListGroupItem key={secret.id} className="list-group-item">
//               <div className="">
//                 <div className='d-flex justify-content-between align-items-center'>
//                   <h6>
//                     {secret.path} (v{secret.version}): {secret.value}
//                   </h6>
//                   <button className="btn btn-danger btn-sm " onClick={(e) => { e.stopPropagation(); deleteSecret(secret.path, secret.version); }}>Delete</button>

//                 </div>

//                 <Card className='mt-3'>
//                   <Card.Header className='bg-light text-dark'>Membership Info</Card.Header>
//                   <Card.Body>
//                     <div>
//                       <strong>Owner:</strong> {secretAccesses[secret.id]?.find(a => a.access_level === 'owner')?.username}
//                     </div>
//                     <div>
//                       <strong>Groups with access:</strong>
//                       <ul>
//                         {secretAccesses[secret.id]?.filter(a => a.group_name).map(access => (
//                           <li key={access.group_id}>{access.group_name}</li>
//                         ))}
//                       </ul>
//                     </div>
//                     <div>
//                       <strong>Users with access:</strong>
//                       <ul>
//                         {secretAccesses[secret.id]?.filter(a => a.username).map(access => (
//                           <li key={access.user_id}>{access.username}</li>
//                         ))}
//                       </ul>
//                     </div></Card.Body></Card>

//                 <span><RotateSecret path={secret.path} /></span>
//                 <span><SecretAccessControl secretId={secret.id} /></span>
//                 <span><SecretMetadataUpdate secret={secret} /></span>



//               </div>
//             </ListGroupItem>
//           ))}
//         </ListGroup>
//       </Card>
//     </div>
//   );
// };

// export default SecretsList;
