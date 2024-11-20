import React, { useState } from 'react';
import { useApi } from '../../api/api'
import { Card } from 'react-bootstrap';

function AddCA() {
  const { get, post, put, del } = useApi();

  const [name, setName] = useState('');
  const [caCert, setCACert] = useState('');
  const [privateKey, setPrivateKey] = useState('');
  const [description, setDescription] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const caCertBase64 = btoa(caCert);
      const privateKeyBase64 = btoa(privateKey);
      await post('/pki/ca', { name, ca_cert: caCertBase64, private_key: privateKeyBase64, description });
      alert('Certificate Authority added successfully!');
    } catch (error) {
      alert('Error adding Certificate Authority');
    }
  };

  return (
    <Card className='mb-3'>

      <Card.Header>Add CA</Card.Header>
      <Card.Body>
      


<ul className='small'>
      <li>You can upload a root CA certificate and private key.</li>
<li>This root CA could be self-signed or provided by an external trusted CA.</li>
<li>This root CA certificate and key will serve as the anchor for signing and verifying all subsequent sub-CAs in your hierarchy.</li>
<li><a target="_blank" href="https://support.apple.com/guide/keychain-access/create-your-own-certificate-authority-kyca2686/mac#:~:text=To%20open%20Keychain%20Access%2C%20search,issued%20by%20the%20certificate%20authority.">
Follow this guide to create your own CA on the MAC</a></li>
<li>openssl pkcs12 -in Cryptkeeper\ CA.p12 -nodes -nocerts -out  cryptkeeper.priv.pem</li>
</ul>


      <form onSubmit={handleSubmit}>
      <div className='form-floating mb-2'>
      <input type="text" className='form-control'  value={name} onChange={(e) => setName(e.target.value)} required />
      <label>Name:</label>
      </div>

      <div className='form-floating mb-2'>
      <textarea className='form-control' style={{ minHeight: "200px"}} value={caCert} onChange={(e) => setCACert(e.target.value)} required />
      <label>CA Certificate:</label>
      </div>

      <div className='form-floating mb-2'>
      <textarea className='form-control' style={{ minHeight: "200px"}} value={privateKey} onChange={(e) => setPrivateKey(e.target.value)} required />
      <label>Private Key:</label>
      </div>

      <div className='form-floating mb-2'>
      <textarea className='form-control' value={description} onChange={(e) => setDescription(e.target.value)} />
      <label>Description:</label>
      </div>
      
      <button className='btn btn-primary w-100' type="submit">Add CA</button>
    </form>
      </Card.Body>

    </Card>
  );
}

export default AddCA;
