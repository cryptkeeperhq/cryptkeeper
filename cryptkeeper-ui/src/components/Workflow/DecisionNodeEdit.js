import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/Button';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Collapse from 'react-bootstrap/Collapse';

const DecisionNodeEdit = ({ selectedNode, unselectNodeId, saveEditedNode, nodes, edges }) => {
    const [editedNodeName, setEditedNodeName] = useState(selectedNode.data.label);
    const [editedNodeDescription, setEditedNodeDescription] = useState(selectedNode.data.description);

    const availableAttributes = nodes[0].data.availableAttributes || []; // Access availableAttributes from the first node
    const [attributeConfigs, setAttributeConfigs] = useState([]);
    const [newAttributeConfig, setNewAttributeConfig] = useState({
        attribute: '',
        operator: '==',
        value: '',
    });

    const [collapsedItems, setCollapsedItems] = useState([]);

    const toggleCollapse = (index) => {
        if (collapsedItems.includes(index)) {
            setCollapsedItems(collapsedItems.filter((item) => item !== index));
        } else {
            setCollapsedItems([...collapsedItems, index]);
        }
    };



    useEffect(() => {
        setAttributeConfigs(selectedNode.data.attributeConfigs || []);
    }, [selectedNode]);

    // State to manage the collapsed state of each attribute
    const [collapsedAttributes, setCollapsedAttributes] = useState({});

    // Function to toggle the collapsed state of an attribute
    const toggleAttributeCollapse = (attribute) => {
        setCollapsedAttributes((prevCollapsedAttributes) => ({
            ...prevCollapsedAttributes,
            [attribute]: !prevCollapsedAttributes[attribute],
        }));
    };


    const handleNodeTypeClick = (event) => {
        unselectNodeId(event);
    };


    const saveEditedNodeName = () => {
        saveEditedNode({
            ...selectedNode,
            data: {
                ...selectedNode.data,
                label: editedNodeName,
                description: editedNodeDescription,
                attributeConfigs: attributeConfigs
            },
        }, true);
    };

    const handleAttributeChange = (index, field, value) => {
        const updatedAttributeConfigs = [...attributeConfigs];
        updatedAttributeConfigs[index][field] = value;
        setAttributeConfigs(updatedAttributeConfigs);
    };

    const handleAddAttribute = () => {
        const selectedAttribute = newAttributeConfig.attribute;
        console.log(selectedAttribute);
        console.log(newAttributeConfig)
        setAttributeConfigs([...attributeConfigs, { ...newAttributeConfig }]);
        setNewAttributeConfig({
            attribute: '',
            operator: '==',
            value: '',
        });
    };

    const handleDeleteAttribute = (index) => {
        const updatedAttributeConfigs = [...attributeConfigs];
        updatedAttributeConfigs.splice(index, 1);
        setAttributeConfigs(updatedAttributeConfigs);
    };


    return (
        <div className="pt-2 mb-3">
            <div className="container">


                <div className="row">
                    <div className="col p-3">
                        <form>
                            <div className="form-floating mb-1">
                                <input
                                    id="floatingInput"
                                    type="text"
                                    className="form-control"
                                    placeholder="name@example.com"
                                    value={editedNodeName}
                                    onChange={(event) => {
                                        setEditedNodeName(event.target.value);
                                    }}
                                />
                                <label htmlFor="floatingInput" className="form-label">
                                    Node Name:
                                </label>
                            </div>

                            <div className="form-floating mb-1">
                                <input
                                    id="floatingInput"
                                    type="text"
                                    className="form-control"
                                    placeholder="name@example.com"
                                    value={editedNodeDescription}
                                    onChange={(event) => {
                                        setEditedNodeDescription(event.target.value);
                                    }}
                                />
                                <label htmlFor="floatingInput" className="form-label">
                                    Node Description:
                                </label>
                            </div>
                        </form>
                    </div>
                </div>

                <Tabs variant="underline" defaultActiveKey="attributes" id="node-details-tabs" className='nav-fill'>

                    <Tab eventKey="attributes" title="Attributes">
                        <div className='mt-3 mb-3'>

                            <div className="card mt-3 mb-3">

                                <ul className="list-group list-group-flush">
                                    {attributeConfigs.map((config, index) => (
                                        <li key={index} className="list-group-item p-0">
                                            <a
                                                className="btn btn-link"
                                                type="button"
                                                onClick={() => toggleCollapse(index)}
                                            >
                                                {config.attribute}
                                            </a>
                                            <button
                                                onClick={() => handleDeleteAttribute(index)}
                                                className="btn btn-link text-danger btn-sm float-end"
                                            >
                                                <i className="fa fa-trash"></i>
                                            </button>

                                            <Collapse in={collapsedItems.includes(index)}>
                                                <ul>
                                                    <li><strong>Operator:</strong><pre>{config.operator}</pre></li>
                                                    <li><strong>Value:</strong><pre>{config.value}</pre></li>
                                                </ul>
                                            </Collapse>
                                        </li>
                                    ))}
                                </ul>
                                <div className='p-3 bg-light'>
                                    <b>Add New Condition</b>
                                    <div className="mb-2 mt-2">
                                        <select
                                            className='form-select'
                                            value={newAttributeConfig.attribute}
                                            onChange={(e) => setNewAttributeConfig({ ...newAttributeConfig, attribute: e.target.value })}
                                        >
                                            {availableAttributes.map((attribute, attrIndex) => (
                                                <option key={attrIndex} value={attribute}>
                                                    {attribute}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                    <div className="mb-2">
                                        <select
                                            className='form-select'
                                            value={newAttributeConfig.operator}
                                            onChange={(e) => setNewAttributeConfig({ ...newAttributeConfig, operator: e.target.value })}
                                        >
                                            <option value="==">==</option>
                                            <option value="!=">!=</option>
                                            <option value=">">{'>'}</option>
                                            <option value="<">{'<'}</option>
                                        </select>
                                    </div>
                                    <div className="mb-2">
                                        <input
                                            className='form-control'
                                            type="text"
                                            value={newAttributeConfig.value}
                                            onChange={(e) => setNewAttributeConfig({ ...newAttributeConfig, value: e.target.value })}
                                        />
                                    </div>

                                    <button
                                        className="btn btn-success w-100 btn-sm"
                                        onClick={handleAddAttribute}
                                    >
                                        Add
                                    </button>

                                </div>
                            </div>
                        </div>
                        <div className="row">
                            <div className="col">
                                <Button
                                    onClick={saveEditedNodeName}
                                    className="btn btn-sm btn-dark w-100 rounded-2"
                                    variant="primary"
                                >
                                    Save changes
                                </Button>
                            </div>
                        </div>
                    </Tab>

                    {/* ... Rest of your Tabs ... */}

                </Tabs>
            </div>
        </div>
    );
};

export default DecisionNodeEdit;
