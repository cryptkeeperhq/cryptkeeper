import React from 'react';
import StartNodeEdit from './StartNodeEdit';
import ActionNodeEdit from './ActionNodeEdit';
import DecisionNodeEdit from './DecisionNodeEdit';
// import OtherNodeEdit from './OtherNodeEdit';

const EditPanel = ({ event, trigger, selectedNode, unselectNodeId, saveEditedNode, onTriggerChange, nodes, edges, events, onEventChange }) => {
  // Determine the selected node's type
  const nodeType = selectedNode.data.type;

  // Define a function to render the appropriate edit section based on node type
  const renderEditSection = () => {
    switch (nodeType) {
      case 'start':
        return <StartNodeEdit event={event} trigger={trigger} selectedNode={selectedNode} unselectNodeId={unselectNodeId} saveEditedNode={saveEditedNode} nodes={nodes} edges={edges} events={events}
        onTriggerChange={(newTrigger) => onTriggerChange(newTrigger)} onEventChange={onEventChange}
        />;
      case 'action':
      case 'others':
        return <ActionNodeEdit selectedNode={selectedNode} unselectNodeId={unselectNodeId} saveEditedNode={saveEditedNode} nodes={nodes} edges={edges} />;
      case 'decision':
        return <DecisionNodeEdit selectedNode={selectedNode} unselectNodeId={unselectNodeId} saveEditedNode={saveEditedNode} nodes={nodes} edges={edges} />;
      default:
        return <div></div>
        // return <OtherNodeEdit selectedNode={selectedNode} />;
    }
  };

  const handleNodeTypeClick = (event) => {
    unselectNodeId(event);
};

  return (
    <div className="side-panel mb-3">
      {/* Common UI elements */}
      <div className='bg-dark text-white p-2 pt-3 ps-3'>
      
                        <button
                            className="btn-sm text-white p-0 m-0 btn float-end"
                            onClick={() => handleNodeTypeClick()}
                        >
                            <i className="fa fa-close"></i>
                        </button>
                    

        <h5><i style={{}} className={`fa ${selectedNode.data.icon}`}></i> Edit Node  ({selectedNode.data.type} | {selectedNode.data.event})</h5></div>
      {/* Render the appropriate edit section */}
      {renderEditSection()}
      {/* Other common UI elements */}
    </div>
  );
};

export default EditPanel;
