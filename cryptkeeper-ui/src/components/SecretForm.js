// import React, { useState } from 'react';

// const SecretForm = () => {
//     const [path, setPath] = useState('');
//     const [value, setValue] = useState('');
//     const [isOneTime, setIsOneTime] = useState(false);
//     const [expiresAt, setExpiresAt] = useState('');

//     const createSecret = () => {
//         const secretData = {
//             path,
//             value,
//             is_one_time: isOneTime,
//             expires_at: expiresAt ? new Date(expiresAt) : null
//         };

//         fetch('/api/secrets', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`
//             },
//             body: JSON.stringify(secretData)
//         })
//             .then(res => {
//                 if (res.ok) {
//                     alert('Secret created.');
//                     setPath('');
//                     setValue('');
//                     setIsOneTime(false);
//                     setExpiresAt('');
//                 } else {
//                     alert('Failed to create secret.');
//                 }
//             })
//             .catch(error => console.error('Error creating secret:', error));
//     };

//     return (
//         <div className="container">
//             <h2>Create Secret</h2>
//             <form>
//                 <div className="form-group">
//                     <label htmlFor="path">Path</label>
//                     <input type="text" className="form-control" id="path" value={path} onChange={e => setPath(e.target.value)} />
//                 </div>
//                 <div className="form-group">
//                     <label htmlFor="value">Value</label>
//                     <input type="text" className="form-control" id="value" value={value} onChange={e => setValue(e.target.value)} />
//                 </div>
//                 <div className="form-group">
//                     <label htmlFor="isOneTime">One-Time Secret</label>
//                     <input type="checkbox" className="form-control" id="isOneTime" checked={isOneTime} onChange={e => setIsOneTime(e.target.checked)} />
//                 </div>
//                 <div className="form-group">
//                     <label htmlFor="expiresAt">Expiration Date</label>
//                     <input type="datetime-local" className="form-control" id="expiresAt" value={expiresAt} onChange={e => setExpiresAt(e.target.value)} />
//                 </div>
//                 <button type="button" className="btn btn-primary" onClick={createSecret}>Create Secret</button>
//             </form>
//         </div>
//     );
// };

// export default SecretForm;
