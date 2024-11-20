import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/Button';
import Collapse from 'react-bootstrap/Collapse';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Modal from 'react-bootstrap/Modal';
import { Position } from 'reactflow';
import InputField from './InputField';
import AddOutput from './AddOutput';
import ExistingOutputVariables from './ExistingOutputVariables';
import OutputPickerModal from './OutputPickerModal';

const ActionNodeEdit = ({ selectedNode, unselectNodeId, saveEditedNode, nodes, edges }) => {
    const [editedNodeName, setEditedNodeName] = useState(selectedNode.data.label);
    const [editedNodeDescription, setEditedNodeDescription] = useState(selectedNode.data.description);

    const [newOutputVariable, setNewOutputVariable] = useState({
        id: '',
        label: '',
        jsonPath: '',
        dataType: '',
    });
    const [collapsedItems, setCollapsedItems] = useState([]);
    const [filteredNodes, setFilteredNodes] = useState([]);

    const [inputFieldTexts, setInputFieldTexts] = useState({}); // State to track input field texts individually

    const [showOutputPicker, setShowOutputPicker] = useState(false);
    const [selectedInput, setSelectedInput] = useState(null);

    useEffect(() => {
        var obj = {}
        for (var input of selectedNode.data.inputs) {
            obj[input.id] = input.value
        }
        setInputFieldTexts((prevTexts) => ({
            ...prevTexts,
            ...obj
        }));

        const priorNodeIds = findPriorNodes(selectedNode.id, edges);
        const filteredNodes = nodes.filter((node) => priorNodeIds.includes(node.id));
        setFilteredNodes(filteredNodes);
    }, [selectedNode, nodes, edges]);

    const handleInputButtonClick = (inputId) => {
        setSelectedInput(inputId);
        setShowOutputPicker(true);
    };



    const findPriorNodes = (currentNodeId, edges) => {
        const priorNodeIds = [];
        const visitedEdges = new Set();

        const findParents = (nodeId) => {
            for (const edge of edges) {
                if (edge.target === nodeId && !visitedEdges.has(edge.id)) {
                    visitedEdges.add(edge.id);
                    priorNodeIds.push(edge.source);
                    findParents(edge.source);
                }
            }
        };

        findParents(currentNodeId);

        return priorNodeIds;
    };



    const handleOutputPickerClose = () => {
        setShowOutputPicker(false);
    };

    const handleOutputSelect = (nodeId, outputVariableId) => {
        const currentText = inputFieldTexts[selectedInput] || '';
        const newText = currentText + `{{${nodeId}.${outputVariableId}}}`;
        setInputFieldTexts((prevTexts) => ({
            ...prevTexts,
            [selectedInput]: newText,
        }));

        setSelectedInput(null);
        setShowOutputPicker(false);
    };

    const toggleCollapse = (index) => {
        if (collapsedItems.includes(index)) {
            setCollapsedItems(collapsedItems.filter((item) => item !== index));
        } else {
            setCollapsedItems([...collapsedItems, index]);
        }
    };

    const handleNodeTypeClick = (event) => {
        unselectNodeId(event);
    };


    const handleInputChange = (event, inputId) => {
        const { value, type, checked } = event.target;
        setInputFieldTexts((prevTexts) => {
            console.log(prevTexts)
            if (type === 'checkbox') {
                return {
                    ...prevTexts,
                    [inputId]: checked
                        ? [...(prevTexts[inputId] || []), value] // Add value to the array
                        : (prevTexts[inputId] || []).filter((val) => val !== value), // Remove value from the array
                };
            } else {
                return {
                    ...prevTexts,
                    [inputId]: value,
                };
            }
        });
    };

    const saveEditedNodeName = () => {
        for (var input in selectedNode.data.inputs) {
            const value = inputFieldTexts[selectedNode.data.inputs[input].id];
            selectedNode.data.inputs[input].value = value;
        }

        saveEditedNode({
            ...selectedNode,
            data: {
                ...selectedNode.data,
                label: editedNodeName,
                description: editedNodeDescription,
            },
        }, true);
    };

    const handleAddOutputVariable = (newOutputVariable) => {
        if (
            newOutputVariable.label &&
            newOutputVariable.jsonPath &&
            newOutputVariable.dataType
        ) {
            const outputVariable = {
                id: newOutputVariable.label,
                label: newOutputVariable.label,
                jsonPath: newOutputVariable.jsonPath,
                dataType: newOutputVariable.dataType,
            };

            const updatedNode = {
                ...selectedNode,
                data: {
                    ...selectedNode.data,
                    outputs: [...selectedNode.data.outputs, outputVariable],
                },
            };

            saveEditedNode(updatedNode);
            setNewOutputVariable({
                id: '',
                label: '',
                jsonPath: '',
                dataType: '',
            });
        }
    };

    const handleDeleteOutputVariable = (index) => {
        const updatedOutputVariables = [...selectedNode.data.outputs];
        updatedOutputVariables.splice(index, 1);
        saveEditedNode({
            ...selectedNode,
            data: {
                ...selectedNode.data,
                outputs: updatedOutputVariables,
            },
        });
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

                <Tabs variant="underline" defaultActiveKey="input" id="node-details-tabs" className='nav-fill'>

                    <Tab eventKey="input" title="Inputs">
                        <div className='mt-3  mb-3'>

                            <OutputPickerModal
                                showOutputPicker={showOutputPicker}
                                handleOutputPickerClose={handleOutputPickerClose}
                                filteredNodes={filteredNodes}
                                handleOutputSelect={handleOutputSelect}
                            />



                            {selectedNode.data.inputs.length === 0 && (
                                <div className='alert alert-secondary p-2'>No inputs required for this node.</div>
                            )}
                            <div className=''>
                                {selectedNode.data.inputs.map((input) => (
                                    <InputField
                                        key={input.id}
                                        input={input}
                                        inputFieldTexts={inputFieldTexts}
                                        handleInputChange={handleInputChange}
                                        handleInputButtonClick={handleInputButtonClick}
                                    />
                                ))}
                            </div>
                        </div>
                       
                    </Tab>


                    <Tab eventKey="output" title="Outputs" >
                        <div className='mt-3  mb-3'>
                            <AddOutput handleAddOutputVariable={handleAddOutputVariable} />
                            
                            <ExistingOutputVariables
                                outputVariables={selectedNode.data.outputs || []}
                                handleDeleteOutputVariable={handleDeleteOutputVariable}
                            />

                        </div>
                    </Tab>
                </Tabs>

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


            </div>
        </div>
    );
};

export default ActionNodeEdit;
