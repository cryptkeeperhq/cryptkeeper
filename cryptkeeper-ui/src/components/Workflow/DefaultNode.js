import React from 'react';
import { FaCopy, FaExpand, FaTrash } from 'react-icons/fa';
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

export default function DefaultNode({ id, data, onDeleteClick, onCopyClick, onExpandClick }) {
    return (
        <>
      <NodeToolbar isVisible={data.toolbarVisible} position={Position.Right}>
        <div className='btn-group'>
          {/* Call the event handlers when buttons are clicked */}
          <button className="btn btn-secondary" onClick={onDeleteClick}>
            <FaTrash />
          </button>
          <button className="btn btn-secondary" onClick={onCopyClick}>
            <FaCopy />
          </button>
          <button className="btn btn-secondary" onClick={onExpandClick}>
            <FaExpand />
          </button>
        </div>
      </NodeToolbar>
            
            <Handle
                type="target"
                position={Position.Top}
                id="a"
            />
            <div className='shadow border-0 rounded-3'
                style={{
                    padding: '8px',
                    background: 'white',
                }}
            >
                <div className="container" style={{ maxWidth: "250px" }}>
                    <div className="row align-items-center justify-content-center">
                        {/* <div className="float-start">
                            <div style={{
                                width: '8px',
                                position: "absolute",
                                top: "10px",
                                left: "10px",
                                bottom: "10px",
                                borderRadius: '8px 8px',
                                background: data.borderColor,
                            }}></div>
                        </div> */}
                        <div style={{ maxWidth: "35px" }}  className="col me-2">
                            <div style={{  backgroundColor: data.borderColor, height: "35px", width: "35px",  }}  className="rounded-circle circle  d-flex text-center align-items-center justify-content-center">
                            <i style={{  fontSize: "18px", color: "#ffffff" }} className={`fa ${data.icon}`}></i>
                            </div>
                        </div>
                        <div className=" col node-label">
                            <div className='ms-flex align-items-center justify-content-center'>
                            <h6 className="title p-0 m-0">{data.label}</h6>
                            <p className="small p-0 m-0">{data.description}</p>
                            </div>
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