import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Collapse from 'react-bootstrap/Collapse';

const ExistingOutputVariables = ({ outputVariables, handleDeleteOutputVariable }) => {
    const [collapsedItems, setCollapsedItems] = useState([]);

    const toggleCollapse = (index) => {
        if (collapsedItems.includes(index)) {
            setCollapsedItems(collapsedItems.filter((item) => item !== index));
        } else {
            setCollapsedItems([...collapsedItems, index]);
        }
    };

    return (
        <div className="mt-3">
            {outputVariables.length > 0 && <h6 className='text-center'>Existing Variables</h6>}
            <ul className="list-group">
                {outputVariables.map((outputVariable, index) => (
                    <li key={index} className="list-group-item p-0">
                        <a
                            className="btn btn-link"
                            type="button"
                            onClick={() => toggleCollapse(index)}
                        >
                            {outputVariable.label}
                        </a>
                        <button
                            onClick={() => handleDeleteOutputVariable(index)}
                            className="btn btn-link text-danger btn-sm float-end"
                        >
                            <i className="fa fa-trash"></i>
                        </button>
                        <Collapse in={collapsedItems.includes(index)}>
                            <ul>
                                <li><strong>JSON Path:</strong><pre>{outputVariable.jsonPath}</pre></li>
                                <li><strong>Data Type:</strong><pre>{outputVariable.dataType}</pre></li>
                            </ul>
                        </Collapse>

                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ExistingOutputVariables;
