import React, { useCallback, useEffect, useState, useMemo, useRef } from 'react';
import { useParams } from 'react-router-dom';

import SidePanel from './SidePanel';
import DefaultNode from './DefaultNode';
import StartNode from './StartNode';
import EditPanel from './EditPanel';
import EdgeEditModal from './EdgeEditModal'
import ReactFlow, {
    MiniMap,
    Controls,
    Background,
    useNodesState,
    useEdgesState,
    addEdge,
    Handle, Position,
    NodeToolbar,
    Panel,
    MarkerType,
    isEdge,
    applyEdgeChanges, applyNodeChanges
} from 'reactflow';
import 'reactflow/dist/style.css';
import { useApi } from '../../api/api'
import 'font-awesome/css/font-awesome.min.css';


const customNodeTypes = {
    decision: [
        {
            label: 'Condition',
            type: 'decision',
            event: 'condition',
            icon: 'fa-flag',
            borderColor: '#0ABBBD',
            inputs: [],
            outputs: []
        },
        {
            label: 'Event',
            type: 'decision',
            event: 'event',
            icon: 'fa-calendar',
            borderColor: '#4C96FC',
            inputs: [],
            outputs: []
        },
    ],
    action: [
        {
            label: 'Email',
            type: 'action',
            event: 'email',
            icon: 'fa-envelope',
            borderColor: '#B8B8B8',
            inputs: [
                {
                    id: 'to',
                    label: 'To',
                    inputType: "select",
                    values: [
                        "vdparikh@gmail.com",
                        "hello@vishalparikh.me"
                    ]
                },
                {
                    id: 'textarea',
                    label: 'Textarea',
                    inputType: "textarea"
                },
                {
                    id: 'checkbox',
                    label: 'Checkbox Button',
                    inputType: "checkbox",
                    values: [
                        "vdparikh@gmail.com",
                        "hello@vishalparikh.me"
                    ]
                },
                {
                    id: 'radio',
                    label: 'Radio Button',
                    inputType: "radio",
                    values: [
                        "vdparikh@gmail.com",
                        "hello@vishalparikh.me"
                    ]
                },
                {
                    id: 'subject',
                    label: 'Subject',
                },
            ],
            outputs: [
                {
                    id: 'status',
                    label: 'status',
                },
            ],
        },
        {
            label: 'SMS',
            type: 'action',
            event: 'sms',
            icon: 'fa-comment',
            borderColor: '#B8B8B8',
            inputs: [
                {
                    id: 'phone',
                    label: 'Phone Number',
                },
                {
                    id: 'message',
                    label: 'Message',
                },
            ],
            outputs: [
                {
                    id: 'sent',
                    label: 'Sent',
                },
            ],
        },
        {
            label: 'Push Notification',
            type: 'action',
            event: 'push-notification',
            icon: 'fa-bell',
            borderColor: '#B8B8B8',
            inputs: [],
            outputs: [
                {
                    id: 'channel',
                    label: 'channel',
                    jsonPath: 'data.channel',
                    dataType: 'string',
                }
            ]
        },
        {
            label: 'Invoke API',
            type: 'action',
            event: 'api',
            icon: 'fa-code',
            borderColor: '#336699',
            inputs: [
                {
                    "id": "url",
                    "label": "url",
                    "inputType": "textbox",
                    "value": ""
                },
                {
                    id: 'method',
                    label: 'method',
                    inputType: "select",
                    values: [
                        "GET",
                        "POST"
                    ]
                },
                {
                    id: 'data',
                    label: 'data',
                    inputType: "textarea",
                }
            ],
            outputs: [

            ]
        },
        {
            label: 'Slack',
            type: 'action',
            event: 'slack',
            icon: 'fa-slack',
            borderColor: 'purple',
            inputs: [
                {
                    "id": "webhookUrl",
                    "label": "webhookUrl",
                    "inputType": "textbox",
                    "value": ""
                },
                {
                    id: 'messageType',
                    label: 'messageType',
                    inputType: "select",
                    values: [
                        "success",
                        "warning",
                        "danger"
                    ]
                },
                {
                    id: 'message',
                    label: 'message',
                    inputType: "textarea",
                },
                {
                    id: 'channel',
                    label: 'channel',
                    inputType: "textbox",
                }
            ],
            outputs: [

            ]
        },
    ],
};

const startNodes = [
    {
        deletable: false,
        id: '1',
        type: 'startNode',
        position: { x: 0, y: 0 },
        width: 100,
        height: 100,
        data: {
            label: "Start Node",
            icon: 'fa-flag',
            type: 'start',
            event: 'start',
            description: "Fun starts here",
            outputs: [],
        },
    },
];

const WorkflowDetails = ({ setTitle }) => {
    const { get, post, put, del } = useApi();

    const { uuid } = useParams(); // Extract the UUID from the URL parameter
    const [workflow, setWorkflow] = useState(null);

    const [reactFlowInstance, setReactFlowInstance] = useState(null)
    const nodeTypes = useMemo(() => ({ defaultNode: DefaultNode, startNode: StartNode }), []);
    const [showNodes, setShowNodes] = useState(false);
    const [nodes, setNodes] = useState(startNodes);
    const [edges, setEdges] = useState([]);
    const [selectedNodeId, setSelectedNodeId] = useState(null);
    const [selectedNode, setSelectedNode] = useState(null);



    const [message, setMessage] = useState(null);

    const [trigger, setTrigger] = useState(null);
    const [event, setEvent] = useState(null);

    const [events, setEvents] = useState(null)





    useEffect(() => {
        setTitle({ heading: "Workflow Details", subheading: "" })
        // Fetch the workflow details for the given UUID from the API
        get(`/workflows/workflow/${uuid}`)
            .then((response) => {
                console.log(response);
                setWorkflow(response);
                setTitle({ heading: response.name, subheading: "" })


                setNodes(response.details.nodes || startNodes);
                setEdges(response.details.edges || []);
                setTrigger(response.details.trigger);
                setEvent(response.details.event || { name: "" });
                console.log(response.details.event)


            })
            .catch((error) => {
                console.error('Error fetching workflow details:', error);
            });
    }, [uuid]); // Include 'uuid' in the dependency array to re-fetch when it changes


    useEffect(() => {
        if (events !== null) {
            return
        }

        get('/workflows/events')
            .then(response => {
                const data = response;
                console.log(data.events)
                setEvents(data.events);
                // startNodes.data.events = data;
            })
            .catch(error => console.error('Error fetching departments:', error));


    }, []);



    const onNodesChange = useCallback(
        (changes) => setNodes((nds) => applyNodeChanges(changes, nds)),
        [setNodes]
    );

    const onEdgesChange = useCallback(
        (changes) => setEdges((eds) => applyEdgeChanges(changes, eds)),
        [setEdges]
    );



    const onConnect = useCallback(
        (connection) => setEdges((eds) => addEdge(connection, eds)),
        [setEdges]
    );

    const handleNodeDoubleClick = (event, node) => {
        setSelectedNode(node);
        setSelectedNodeId(node.id);
        // setEditedNodeName(node.data.label);
        // setEditedNodeDescription(node.data.description);

        // console.log(node.data.editedNodeData);
        // if (node.data.editedNodeData !== undefined) {
        //     setEditedNodeData(node.data.editedNodeData);
        // }
        setShowNodes(false)
    };

    const unselectNodeId = (event) => {
        setSelectedNodeId(null);
        setSelectedNode(null);
    };

    const onTriggerChange = (newTrigger) => {
        console.log("trigger change", newTrigger)
        setTrigger(newTrigger)
    }

    const onEventChange = (newEvent) => {
        console.log("event change", newEvent)
        setEvent(newEvent)
    }


    // Function to save the edited node's name
    const saveEditedNode = (updatedNodeData, close = false) => {
        console.log(updatedNodeData);
        // Find the selected node in the nodes state array
        const updatedNodes = nodes.map((node) => {
            if (node.id === updatedNodeData.id) {
                console.log({
                    ...updatedNodeData.data,
                    label: updatedNodeData.data.label,
                    description: updatedNodeData.data.description,
                    editedNodeData: updatedNodeData.data.editedNodeData
                });
                // Update the name of the selected node
                return {
                    ...node,
                    data: {
                        ...updatedNodeData.data,
                        label: updatedNodeData.data.label,
                        description: updatedNodeData.data.description,
                        editedNodeData: updatedNodeData.data.editedNodeData
                    },
                };
            }
            return node;
        });

        // Update the nodes state with the edited node's name
        setNodes(updatedNodes);
        setSelectedNode(updatedNodeData);

        // Clear the editedNodeName state variable
        // setEditedNodeName('');
        // setEditedNodeDescription('')
        if (close) {
            unselectNodeId()
        }

    }

    // Define your click event handlers here
    const handleDeleteClick = (nodeId) => {
        // Implement the delete logic here
        console.log(`Delete button clicked for node with ID: ${nodeId}`);
    };

    const handleCopyClick = (nodeId) => {
        // Implement the copy logic here
        console.log(`Copy button clicked for node with ID: ${nodeId}`);
    };

    const handleExpandClick = (nodeId) => {
        // Implement the expand logic here
        console.log(`Expand button clicked for node with ID: ${nodeId}`);
    };

    const addNode = (type, event, label, icon, description, borderColor, inputs, outputs) => {
        var nodeId = `node-${nodes.length + 1}`

        const newNode = {
            id: nodeId,
            type: 'defaultNode',
            deletable: true,
            position: { x: 0, y: 0 },
            onDeleteClick: (nodeId) => handleDeleteClick(nodeId),
            onCopyClick: () => handleCopyClick(nodeId),
            onExpandClick: () => handleExpandClick(nodeId),
            data: {
                label,
                type,
                event,
                icon,
                description,
                borderColor,
                inputs: inputs || [],
                outputs: outputs || [],
            },
        };
        setNodes((prevNodes) => [...prevNodes, newNode]);
    };



    const saveFlow = () => {



        workflow.details = workflow.details || {}
        workflow.details.nodes = nodes
        workflow.details.event = event
        workflow.details.trigger = trigger
        workflow.details.edges = edges
        //  setWorkflow((prevWorkflow) => ({
        //      ...prevWorkflow,
        //      nodes: applyNodeChanges(changes, prevWorkflow.nodes),
        //  }));

        // setWorkflow((prevWorkflow) => ({
        //     ...prevWorkflow,
        //     edges: applyEdgeChanges(changes, prevWorkflow.edges),
        // }));


        console.log(workflow)
        
        setWorkflow(workflow)

        // console.log("saving to localStorage:", { event: event, trigger: trigger, nodes, edges }); // Debugging statement
        // const stateToSave = JSON.stringify({ id: uuid, event: event, trigger: trigger, nodes, edges });

        post(`/workflows`, workflow)
            .then(response => {
                const data = response;
                console.log(data)
            })
            .catch(error => console.error('Error fetching departments:', error));

        // localStorage.setItem('flowState', stateToSave);
    };

    const runFlow = () => {
        // // Create the JSON payload
        // const payload = JSON.stringify({
        //     id: uuid,
        //     name: workflow.name,
        //     details: {
        //         edges: edges,
        //         event: event,
        //         trigger: trigger,
        //         nodes, edges, // Replace with your nodes and edges data
        //     }
        // });

        workflow.details = workflow.details || {}
        workflow.details.nodes = nodes
        workflow.details.event = event
        workflow.details.trigger = trigger
        workflow.details.edges = edges
        //  setWorkflow((prevWorkflow) => ({
        //      ...prevWorkflow,
        //      nodes: applyNodeChanges(changes, prevWorkflow.nodes),
        //  }));

        // setWorkflow((prevWorkflow) => ({
        //     ...prevWorkflow,
        //     edges: applyEdgeChanges(changes, prevWorkflow.edges),
        // }));


        console.log(workflow)
        setWorkflow(workflow)

        post(`/workflows/execute`, { workflow: workflow, event: { "name": "test" }})
            .then(response => {
                setMessage("Workflow Executed Successfully!")
            })
            .catch(error => console.error('Error fetching departments:', error));

        // Define the API URL
        const apiUrl = 'execute';

        // // Send a POST request to the API
        // fetch(apiUrl, {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json',
        //     },
        //     body: payload,
        // })
        //     .then((response) => {
        //         if (!response.ok) {
        //             throw new Error('Failed to execute workflow');
        //         }

        //         setMessage("Workflow Executed Successfully!")
        //         return response.text();
        //     })
        //     .then((data) => {
        //         console.log('Workflow executed successfully:', data);
        //     })
        //     .catch((error) => {
        //         console.error('Error executing workflow:', error);
        //     });
    };

    // State to store the currently edited edge
    const [editingEdge, setEditingEdge] = useState(null);

    // Handler for double-click on edges
    const handleEdgeDoubleClick = (event, edge) => {
        if (isEdge(edge)) {
            // Set the editingEdge state to the clicked edge
            setEditingEdge(edge);
        }
    };

    // Handler for saving edge changes
    const handleEdgeSave = (newLabel, newType, newAnimated, newStroke) => {
        // Find the index of the editing edge in the edges array
        console.log(newStroke);
        const edgeIndex = edges.findIndex((edge) => edge.id === editingEdge.id);

        if (edgeIndex !== -1) {
            // Create a copy of the edges array
            const updatedEdges = [...edges];

            // Update the label and type of the editing edge
            updatedEdges[edgeIndex] = {
                ...updatedEdges[edgeIndex],
                label: newLabel,
                type: newType,
                animated: newAnimated,
                style: newStroke,
                markerEnd: {
                    type: MarkerType.ArrowClosed,
                },
            };

            // Set the updated edges array and clear the editingEdge state
            setEdges(updatedEdges);
            setEditingEdge(null);
        }
    };

    const onNodeDelete = (elementsToRemove) => {
        console.log(elementsToRemove);
        // Filter out the nodes you want to prevent from deletion
        const filteredElementsToRemove = elementsToRemove.filter((element) => {
            console.log(element.type, element.data.type);
            if (element.type == "startNode" && element.data.type == "start") {
                // Prevent deletion of startNode
                return false;
            }
            // Allow deletion of other elements
            return true;
        });

        // Perform the actual removal of elements
        // You can use the filteredElementsToRemove array
        // to update your state or perform any other actions.
        // Example: setElements(newElements);

        // Log the elements that are removed (optional)
        console.log('Elements removed:', filteredElementsToRemove);
    };

    return (
        <div className=''>
            <div className="navbar bg-gray">
                <div className="ms-2 btn-group">
                    <button className="btn btn-sm btn-primary" onClick={() => { unselectNodeId(); setShowNodes(!showNodes) }}><i className='fa fa-plus'></i></button>
                </div>
                {message && (<div className='text-primary strong ms-auto p-1 ps-3 pe-3 rounded-2'>{message}</div>)}
                <div className="btn-group me-2">
                    {/* <button className="btn btn-sm btn-outline-light" onClick={loadFlow}><i className='fa fa-folder-open'></i></button> */}
                    <button className="btn btn-sm btn-dark" onClick={saveFlow}><i className='fa fa-save'></i></button>
                    <button className="btn btn-sm btn-success" onClick={runFlow}><i className='fa fa-play'></i></button>
                </div>
            </div>


            <div className='container-fluid'>
                <div className="row" style={{ position: 'relative' }}>
                    <div className="col" style={{ height: 'calc(100vh - 200px)' }}>

                        {/* setTimeout(reactFlowInstance.fitView, 4) */}

                        <div style={{ height: 'calc(100vh - 200px)' }}>
                            <ReactFlow
                                onInit={(instance) => {
                                    setReactFlowInstance(instance);
                                    setTimeout(instance.fitView, "1000")
                                }}
                                nodes={nodes}
                                edges={edges}
                                defaultNodes={startNodes}
                                onNodesChange={onNodesChange}
                                onEdgesChange={onEdgesChange}
                                onConnect={onConnect}
                                onNodesDelete={onNodeDelete}
                                nodeTypes={nodeTypes}
                                onEdgeDoubleClick={handleEdgeDoubleClick}
                                onEdgeContextMenu={(event, edge) => event.preventDefault()} // Disable context menu
                                onlyRenderVisibleElements={true}
                                elevateEdgesOnSelect={true}
                                edgeUpdaterRadius={30}
                                paneMoveable
                                zoomOnScroll
                                arrowHeadColor="#206584"
                                defaultPosition={[0, 0]}
                                defaultZoom={1}

                                defaultEdgeOptions={{
                                    "type": "smoothstep",
                                    markerEnd: {
                                        type: MarkerType.ArrowClosed,
                                    }
                                }}
                                fitView
                                onNodeDoubleClick={handleNodeDoubleClick}
                            >

                                {editingEdge && (
                                    <Panel className='panel bg-white shadow-sm rounded-3' position="top-right">
                                        <EdgeEditModal
                                            edge={editingEdge}
                                            onSave={handleEdgeSave}
                                            onCancel={() => setEditingEdge(null)}
                                        /></Panel>
                                )}

                                {showNodes && (<Panel className='panel' position="top-left">
                                    <div className="panel side-panel  ">
                                        <SidePanel nodeTypes={customNodeTypes} addNode={addNode} />
                                    </div></Panel>)}


                                {selectedNodeId && (<Panel className='panel' position="top-right">
                                    <div className="panel bg-white shadow-sm rounded-3" style={{ width: "350px", height: "calc(100vh - 450px)", overflow: "auto" }}>
                                        <EditPanel event={event} trigger={trigger} selectedNode={selectedNode} unselectNodeId={unselectNodeId} saveEditedNode={saveEditedNode} onTriggerChange={onTriggerChange} nodes={nodes} edges={edges} events={events} onEventChange={onEventChange} />
                                    </div>
                                </Panel>)}

                                <Controls />
                                <MiniMap />
                                <Background variant="plain" gap={1} size={1} />
                            </ReactFlow>


                        </div>
                    </div>


                </div>
            </div>
        </div>
    );
};

export default WorkflowDetails;
