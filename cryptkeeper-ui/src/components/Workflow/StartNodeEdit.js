import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/Button';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';

const StartNodeEdit = ({ event, trigger, selectedNode, unselectNodeId, saveEditedNode, nodes, edges, events, onTriggerChange, onEventChange }) => {
    const [editedNodeName, setEditedNodeName] = useState(selectedNode.data.label);
    const [editedNodeDescription, setEditedNodeDescription] = useState(selectedNode.data.description);
    const [outputs, setOutputs] = useState([]);

    const handleNodeTypeClick = (event) => {
        unselectNodeId(event);
    };



    const [selectedEvent, setSelectedEvent] = useState(event.name || null);
    const [selectedAttributes, setSelectedAttributes] = useState([]);

    const [selectedTrigger, setSelectedTrigger] = useState(trigger || "");

    const handleTriggerChange = (e) => {
        setSelectedTrigger(e.target.value);
        onTriggerChange(e.target.value);
    };


    const handleEventSelection = (event) => {
        var outputs = []
        for (var attribute of event.fields) {
            const outputVariable = {
                id: attribute,
                label: attribute,
                jsonPath: attribute,
                dataType: "String",
            };

            outputs.push(outputVariable)
        }

        setOutputs(outputs)
        setSelectedEvent(event.eventType);
        setSelectedAttributes(event.fields);

        onEventChange(event);
    };

    const handleSaveChanges = () => {
        // Save the changes made in the form
        saveEditedNode({
            ...selectedNode,
            data: {
                ...selectedNode.data,
                label: event.eventType || editedNodeName,
                description: editedNodeDescription,
                selectedEvent: selectedEvent, // Update with the selected event
                // availableAttributes: selectedAttributes, // Update with the selected attributes
                outputs: outputs,
            },
        }, true);
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

                <Tabs variant="underline" defaultActiveKey="events" id="node-details-tabs" className='nav-fill'>
                    <Tab eventKey="events" title="Events">
                        <div className="mt-3 mb-3">
                            <ul className="list-group">
                                {/* {selectedNode.data.events.map((event, index) => (
                                    <li
                                        key={index}
                                        className={`list-group-item ${selectedEvent === event.name ? 'active' : ''}`}
                                        onClick={() => handleEventSelection(event)}
                                    >
                                        {event.name}
                                    </li>
                                ))} */}

                                {events.map((event, index) => (
                                    <li
                                        key={index}
                                        className={`list-group-item ${selectedEvent === event.eventType ? 'active' : ''}`}
                                        onClick={() => handleEventSelection(event)}
                                    >
                                        {event.eventType}
                                    </li>
                                ))}
                            </ul>
                        </div>
                    </Tab>
                    <Tab eventKey="triggers" title="Triggers">
                        <div className="mt-3 mb-3">
                            <div className="form-group">
                                {/* <label>Select Trigger:</label> */}
                                <select
                                    className="form-control"
                                    value={selectedTrigger}
                                    onChange={handleTriggerChange}
                                >
                                    <option value="">Not Selected</option>
                                    <option value="manual">Manual</option>
                                    <option value="on_event">Run on Event</option>
                                </select>
                            </div>
                        </div>
                    </Tab>

                </Tabs>

                {/* Save Changes button */}
                <div className="row">
                    <div className="col">
                        <Button
                            onClick={handleSaveChanges}
                            className="btn btn-sm btn-dark w-100 rounded-2"
                            variant="primary"
                        >
                            Save Changes
                        </Button>
                    </div>
                </div>



            </div>
        </div>
    );
};

export default StartNodeEdit;
