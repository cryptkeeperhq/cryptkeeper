import React, { useState, useEffect } from 'react';
import ReactFlow, {
    ReactFlowProvider, MiniMap, Controls
} from 'react-flow-renderer';
import { useApi } from '../api/api'

        
const SecretLineage = ({ token, path, selected_key, secretId }) => {
    const { get, post, put, del } = useApi();
    const [nodes, setNodes] = useState([]);
    const [edges, setEdges] = useState([]);

    useEffect(() => {
        const fetchLineage = async () => {
            try {
                const data = await get(`/secrets/lineage?secret_id=${secretId}&path=${path.path}`);
                const lineageData = data;

                    const initialPosition = { x: 0, y: 0 };
                    const offsetX = 300;
                    const offsetY = 100;
                    let currentY = initialPosition.y;

                    const nodes = lineageData.nodes.map((node, index) => {
                        const position = {
                            x: initialPosition.x + index * offsetX,
                            y: currentY
                        };

                        if (node.type === "secret") {
                            currentY += offsetY;
                        }

                        return {
                            id: node.id.toString(),
                            data: { label: node.label },
                            position: position,
                            type: node.type,
                            style: {
                                backgroundColor: node.type === "secret" ? 'rgba(255, 0, 0, 0.2)' : 'rgba(255, 0, 0, 0.1)',
                                width: 200,
                                height: 60
                            }
                        };
                    });

                    const edges = lineageData.edges.map(edge => ({
                        id: edge.id,
                        source: edge.source.toString(),
                        target: edge.target.toString(),
                        type: 'smoothstep',
                    }));

                    setNodes(nodes);
                    setEdges(edges);
                    
            } catch (error) {
                console.log(error.message);
            }


            
        };

        fetchLineage();
    }, [token, path, secretId]);

    return (
        <div className="p-0 mt-4">
            <div style={{ width: "100%", height: "500px" }}>
                <ReactFlowProvider>
                    <ReactFlow nodes={nodes} edges={edges} fitView>
                        <MiniMap />
                        <Controls />
                    </ReactFlow>
                </ReactFlowProvider>
            </div>
        </div>
    );
};

export default SecretLineage;
