import React from 'react';
import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';

const OutputPickerModal = ({ showOutputPicker, handleOutputPickerClose, filteredNodes, handleOutputSelect }) => {

  // if (filteredNodes.length > 0 && filteredNodes[0].data !== undefined) {
  //   const availableAttributes = filteredNodes[0].data.availableAttributes || []; 
  //   console.log(availableAttributes)
  // }

  console.log(filteredNodes)

  return (
    <Modal show={showOutputPicker} onHide={handleOutputPickerClose}>
      <Modal.Header closeButton>
        <Modal.Title>Select Output Variable</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <ul className='list-group'>
          {filteredNodes.map((node) => (
            <li className='list-group-item' key={node.id}>
              <i style={{}} className={`me-2 fa ${node.data.icon}`}></i>
              <strong>{node.data.label}</strong>
              <ul className='list-group list-group-flush p-0 m-0'>
                {node.data.outputs != null && (
                  <>
                {node.data.outputs.map((outputVariable) => (
                  <a href="#" key={outputVariable.id} className='list-group-item' onClick={() =>
                    handleOutputSelect(node.id, outputVariable.id)
                  }>
                    {/* <pre>{JSON.stringify(outputVariable)}</pre> */}
                    {outputVariable.label}
                  </a>
                ))}
                </>
                )}

              </ul>
            </li>
          ))}
        </ul>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleOutputPickerClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default OutputPickerModal;
