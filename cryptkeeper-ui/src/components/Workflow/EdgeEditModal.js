import React, { useState } from 'react';
import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';

const EdgeEditModal = ({ edge, onSave, onCancel }) => {
    const [newLabel, setNewLabel] = useState(edge.label);
    const [newType, setNewType] = useState(edge.type);
    const [newAnimated, setNewAnimated] = useState(edge.animated || false);
    const [newStroke, setNewStroke] = useState(edge.style ? edge.style.stroke : '');

    const handleSave = () => {
        onSave(newLabel, newType, newAnimated, { stroke: newStroke });
    };

    return (
        <div
            className="pt-2 mb-3"
            style={{ display: 'block', position: 'initial' }}
            onHide={onCancel}
        >
            <div className="container">
            <div className="row">
                    <div className="col-1">
                        <i style={{}} className={`fa fa-edit`}></i>
                    </div>
                    <div className="col">
                        <button
                            className="p-0 m-0 btn float-end"
                            onClick={() => onCancel()}
                        >
                            <i className="fa fa-close"></i>
                        </button>
                        <h4 className="pt-1">Edit Edge</h4>
                    </div>
                </div>
            <Form>
                <Form.Group controlId="label">
                    <Form.Label>Label</Form.Label>
                    <Form.Control
                        type="text"
                        value={newLabel}
                        onChange={(e) => setNewLabel(e.target.value)}
                    />
                </Form.Group>
                <Form.Group controlId="type">
                    <Form.Label>Type</Form.Label>
                    <Form.Control
                        as="select"
                        value={newType}
                        onChange={(e) => setNewType(e.target.value)}
                    >
                        <option value="default">Default</option>
                        <option value="straight">Straight</option>
                        <option value="step">Step</option>
                        <option value="smoothstep">Smooth Step</option>
                    </Form.Control>
                </Form.Group>
                <Form.Group className='mt-3' controlId="stroke">
                <Form.Label>Stroke (Hex Color)</Form.Label>
                <input className="nodrag ms-3" type="color" onChange={(e) => setNewStroke(e.target.value)} defaultValue={newStroke} />
                    
                    

                    {/* <Form.Control
                        type="text"
                        value={newStroke}
                        onChange={(e) => setNewStroke(e.target.value)}
                    /> */}
                </Form.Group>
                <Form.Group className='mt-3' controlId="animated">
                    <Form.Check
                        type="checkbox"
                        label="Animated"
                        checked={newAnimated}
                        onChange={(e) => setNewAnimated(e.target.checked)}
                    />
                </Form.Group>                
            </Form>
            <div className="mt-3">
                <Button variant="transparent" className="btn-sm" onClick={onCancel}>
                    Cancel
                </Button>
                <Button variant="dark" className="float-end btn-sm" onClick={handleSave}>
                    Save Changes
                </Button>
            </div>
            </div>
        </div>
    );
};


export default EdgeEditModal;
