import React, { useState } from 'react';
import { ListGroupItem, Button } from 'react-bootstrap';

const MetadataItem = ({ metadata }) => {
    const [showFull, setShowFull] = useState(false);

    const toggleShowFull = () => {
        setShowFull(!showFull);
    };

    const metadataString = JSON.stringify(metadata, null, 2);
    const truncatedMetadata = metadataString.length > 200 ? metadataString.substring(0, 200) + '...' : metadataString;
    const isTruncated = metadataString.length > 200;

    return (
        <>
        { Object.keys(metadata).length === 0 ? '' : 
    
        <div>
            <span>Metadata</span>
            <span className='ms-auto'>
                <pre className='m-0'>{showFull || !isTruncated ? metadataString : truncatedMetadata}</pre>
                {isTruncated && (
                    <Button variant="link" onClick={toggleShowFull} className="p-0">
                        {showFull ? 'Show less' : 'Show more'}
                    </Button>
                )}
            </span>
        </div>
        }
        </>
    );
};

export default MetadataItem;
