import React from 'react';
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
    applyEdgeChanges, applyNodeChanges
} from 'reactflow';

// import { useStore } from '../store';

export default function StartNode({ id, data }) {
    return (
        <>
            
            <div className='shadow border-0  rounded-pill'
                style={{
                    backgroundColor: "#ffffff",
                    padding: '0px',
                    background: 'white',
                }}
            >
                <div className="container" style={{ maxWidth: "250px" }}>
                    <div className="row">
                        
                        <div className="ms-3 mt-3 col-1 node-icon center" >
                            <i style={{ color: "green" }} className="mt-1 w-100 fa-2x  fa fa-flag"></i>
                        </div>
                        <div className="mt-3 ms-3 col node-label">
                            <h6 className="p-0 m-0 mb-0">{data.label}</h6>
                            <p className="small">{data.description}</p>
                        </div>
                    </div>
                </div>

            </div>
            <Handle
                type="source"
                position={Position.Bottom}
                id="b"
            />
        </>
    );
};