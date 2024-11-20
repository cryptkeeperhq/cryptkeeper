import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';

const AddOutput = ({ handleAddOutputVariable }) => {
  const [newOutputVariable, setNewOutputVariable] = useState({
    id: '',
    label: '',
    jsonPath: '',
    dataType: '',
  });

  const handleInputChangeForOutput = (event) => {
    const { name, value } = event.target;
    setNewOutputVariable({
      ...newOutputVariable,
      [name]: value,
    });
  };

  const handleAddOutputClick = () => {
    handleAddOutputVariable(newOutputVariable);
    setNewOutputVariable({
      id: '',
      label: '',
      jsonPath: '',
      dataType: '',
    });
  };

  // Define a list of popular data types
  const popularDataTypes = ['String', 'Number', 'Boolean', 'Date'];

  return (
    <div className=''>
      <form>
        <div className=" mb-1 p-2 ">
        <div><b><label className="form-label ">Label (attributeName)</label></b></div>
          <input
            type="text"
            name="label"
            className="form-control w-100"
            value={newOutputVariable.label}
            onChange={handleInputChangeForOutput}
          />
        </div>
        <div className=" mb-1 p-2 ">
        
        <div><b><label className="form-label ">JSON Path</label></b></div>
          <input
            type="text"
            name="jsonPath"
            className="form-control w-100"
            value={newOutputVariable.jsonPath}
            onChange={handleInputChangeForOutput}
          />
        </div>
        <Form.Group className="mb-1 p-2 ">
        <div><b><label className="form-label ">Data Type</label></b></div>
          <Form.Control
            as="select"
            name="dataType"
            value={newOutputVariable.dataType}
            onChange={handleInputChangeForOutput}
          >
            <option value="">Select Data Type</option>
            {popularDataTypes.map((type) => (
              <option key={type} value={type}>
                {type}
              </option>
            ))}
          </Form.Control>
        </Form.Group>
        <Button
          onClick={handleAddOutputClick}
          className="btn btn-sm btn-dark w-100 rounded-2"
          variant="dark"
        >
          Add New Output
        </Button>
      </form>
    </div>
  );
};

export default AddOutput;
