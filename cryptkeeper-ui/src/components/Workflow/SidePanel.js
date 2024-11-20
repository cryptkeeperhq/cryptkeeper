import React from 'react';
import { Card, ListGroup, ListGroupItem } from 'react-bootstrap';

const SidePanel = ({ nodeTypes, addNode }) => {
  const handleNodeTypeClick = (type, event, label, icon, description, borderColor, inputs, outputs) => {
    // Call the addNode function with the selected properties
    addNode(type, event, label, icon, description, borderColor, inputs, outputs);
  };

  return (
    <div className="">
      {/* <h6 className='strong p-3 text-strong border-bottom'>Nodes</h6> */}
      {Object.entries(nodeTypes).map(([group, types]) => (
        <Card flush className='mt-3 ' key={group}>
          <Card.Header as="h6">{group}</Card.Header>
          {/* <Card.Body> */}
          <ListGroup variant='flush' className="list-group list-group-flush" >
            {types.map((nodeType) => (
              <ListGroupItem  key={nodeType.label}>
                <div
                  onClick={() =>
                    handleNodeTypeClick(
                      nodeType.type,
                      nodeType.event,
                      nodeType.label,
                      nodeType.icon,
                      nodeType.description,
                      nodeType.borderColor,
                      nodeType.inputs,
                      nodeType.outputs
                    )
                  }
                  className='text-small m-0 p-0'
                >
                            <i style={{ color: nodeType.borderColor }} className={`me-2 fa ${nodeType.icon}`}></i>
                  {nodeType.label}
                </div>
              </ListGroupItem>
            ))}
          </ListGroup>
          {/* </Card.Body> */}
          </Card>
      ))}
      
    </div>
  );
};

export default SidePanel;