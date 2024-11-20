import React, { useState, useEffect, useRef } from 'react';
import { Button, ButtonGroup, Card, Container, FormLabel, ListGroup, ListGroupItem, Row, Col, InputGroup, Form, Alert } from 'react-bootstrap';



const TagsForm = ({ updateTags }) => {

    const [tags, setTags] = useState([]);
    const handleAddTag = () => {
        setTags([...tags, ""]);
    };

    const handleRemoveTag = (index) => {
        const newTags = tags.filter((_, i) => i !== index);
        setTags(newTags);
        updateTags(newTags);
    };

    const handleTagChange = (index, value) => {
        const newTags = tags.map((tag, i) => (i === index ? value : tag));
        setTags(newTags);
        updateTags(newTags);
    };



    return (
        <Card className='mt-3'>
            <Card.Header>Manage Tags</Card.Header>
            <Card.Body>
                <div className="">
                    {/* <label className="form-label me-1">Tags</label> */}
                    {tags.map((tag, index) => (
                        <InputGroup className="mb-2" key={index}>
                            <Form.Control
                                type="text"
                                placeholder="Enter tag"
                                value={tag}
                                onChange={(e) => handleTagChange(index, e.target.value)}
                            />
                            <Button variant="danger" className='btn-sm' onClick={() => handleRemoveTag(index)}>Remove</Button>
                        </InputGroup>
                    ))}
                    <Button variant="info" className='btn-sm rounded-5 p-2' onClick={handleAddTag}><small>Add Tag</small></Button>
                </div>

            </Card.Body>
        </Card>
    );
};

export default TagsForm;