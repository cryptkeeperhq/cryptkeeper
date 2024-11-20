import React from 'react';

const RenderNestedObject = ({ data }) => {
    if (typeof data === 'object' && data !== null) {
        return (
            <ul className='yaml'>
                {Object.keys(data).map(key => (
                    <li key={key}>
                        {key}: <RenderNestedObject data={data[key]} />
                    </li>
                ))}
            </ul>
        );
    }
    return <span>{data}</span>;
};

export default RenderNestedObject;
