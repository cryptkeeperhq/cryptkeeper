import React from 'react';
import { Form } from 'react-bootstrap';
import YAML from 'js-yaml';

const UploadSecrets = ({ onUpload }) => {
    
    const flattenObject = (obj, prefix = '') => {
        return Object.keys(obj).reduce((acc, k) => {
            const pre = prefix.length ? `${prefix}.` : '';
            if (typeof obj[k] === 'object' && obj[k] !== null && !Array.isArray(obj[k])) {
                Object.assign(acc, flattenObject(obj[k], pre + k));
            } else {
                acc[pre + k] = obj[k];
            }
            return acc;
        }, {});
    };

    const handleFileUpload = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    let parsedData;
                    if (file.type === "application/json") {
                        parsedData = JSON.parse(e.target.result);
                    } else if (file.type === "application/x-yaml" || file.type === "text/yaml") {
                        parsedData = YAML.load(e.target.result);
                    } else {
                        throw new Error("Unsupported file type");
                    }
                    const flattenedData = flattenObject(parsedData);
                    const multiValues = Object.keys(flattenedData).map(key => ({
                        key: key,
                        value: flattenedData[key]
                    }));
                    onUpload(multiValues);
                } catch (error) {
                    alert("Error parsing file: " + error.message);
                }
            };
            reader.readAsText(file);
        }
    };

    return (
        <Form.Group className="mt-2">
            <Form.Control
                type="file"
                accept=".json,.yaml,.yml"
                onChange={handleFileUpload}
            />
            <Form.Label>Upload JSON or YAML</Form.Label>
        </Form.Group>
    );
};

export default UploadSecrets;