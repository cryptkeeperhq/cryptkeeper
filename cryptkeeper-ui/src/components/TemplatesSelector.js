import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useApi } from '../api/api'

const TemplatesSelector = ({ onTemplateSelect }) => {
    const { get, post, put, del } = useApi();

    const [templates, setTemplates] = useState([]);
    const [selectedTemplate, setSelectedTemplate] = useState('');

    useEffect(() => {
        const fetchTemplates = async () => {
            try {
                const data = await get('/templates');
                setTemplates(data);
            }
            catch (e) {
                console.error('Failed to fetch dashboard summary');
            }

        };

        fetchTemplates();
    }, []);

    const handleTemplateSelect = (e) => {
        const template = templates.find(t => t.name === e.target.value);
        setSelectedTemplate(e.target.value);
        onTemplateSelect(template);
    };

    return (
        <Form.Group className="form-floating mt-2">
            <Form.Control
                as="select"
                value={selectedTemplate}
                onChange={handleTemplateSelect}
            >
                <option value="" disabled>Select Template</option>
                {templates.map(template => (
                    <option key={template.name} value={template.name}>
                        {template.name}
                    </option>
                ))}
            </Form.Control>
            <Form.Label>Select Template</Form.Label>
        </Form.Group>
    );
};

export default TemplatesSelector;