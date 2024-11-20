import React, { useState } from 'react';
import { useApi } from '../../api/api';
import { Card } from 'react-bootstrap';

function AddTemplate() {
  const { post } = useApi();

  const [name, setName] = useState('');
  const [commonName, setCommonName] = useState('');
  const [organization, setOrganization] = useState('');
  const [validityPeriod, setValidityPeriod] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await post('/pki/template', {
        name,
        common_name: commonName,
        organization,
        validity_period: parseInt(validityPeriod, 10),
      });
      alert('Certificate Template added successfully!');
    } catch (error) {
      alert('Error adding Certificate Template');
    }
  };

  return (
   <Card className='mb-3'>
    <Card.Header>Add Template</Card.Header>
    <Card.Body>
    <form onSubmit={handleSubmit}>
      <div className='form-floating mb-2'>
        <input
          type="text"
          className='form-control'
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <label>Name:</label>
      </div>
{/* 
      <div className='form-floating mb-2'>
        <input
          type="text"
          className='form-control'
          value={commonName}
          onChange={(e) => setCommonName(e.target.value)}
          required
        />
        <label>Common Name:</label>
      </div> */}

      <div className='form-floating mb-2'>
        <input
          type="text"
          className='form-control'
          value={organization}
          onChange={(e) => setOrganization(e.target.value)}
        />
        <label>Organization:</label>
      </div>

      <div className='form-floating mb-2'>
        <input
          type="number"
          className='form-control'
          value={validityPeriod}
          onChange={(e) => setValidityPeriod(e.target.value)}
          required
        />
        <label>Validity Period (days):</label>
      </div>

      <button className='btn btn-primary w-100' type="submit">Add Template</button>
    </form>
    </Card.Body>
   </Card>
  );
}

export default AddTemplate;
