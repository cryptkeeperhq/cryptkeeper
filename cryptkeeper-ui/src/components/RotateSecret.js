import React, { useState, useEffect } from 'react';
import { Card, Form, InputGroup } from 'react-bootstrap';
import Datetime from 'react-datetime';
import { FaInfoCircle } from 'react-icons/fa';
import { useApi } from '../api/api';
import moment from 'moment';
import { Link } from 'react-router-dom';
import UploadSecrets from './UploadSecrets';
import CLIUsage from './CLIUsage';

const RotateSecret = ({ secret, path, selected_key, token, onRotated }) => {
    const { post } = useApi();
    const [newSecretValue, setNewSecretValue] = useState('');
    const [expiresAt, setExpiresAt] = useState('');
    const [message, setMessage] = useState('');
    const [isOneTime, setIsOneTime] = useState(false);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [multiValue, setMultiValue] = useState([{ key: "", value: "" }]);

    useEffect(() => {
        setMessage('');
        setError('');
        setNewSecretValue('');
        setExpiresAt('');
        setIsOneTime(false);
    }, [path, secret, token]);

    const rotateSecret = async () => {
        setLoading(true);
        setError('');
        setMessage('');
        try {
            const formattedExpiresAt = expiresAt ? moment(expiresAt).toISOString() : null;

            const payload = {
                path_id: path?.id,
                key: secret.key,
                value: newSecretValue,
                expires_at: formattedExpiresAt,
                is_one_time: isOneTime,
                is_multi_value: secret?.is_multi_value
            };

            if (secret?.is_multi_value) {
                const multiValueObj = multiValue.reduce((acc, curr) => {
                    acc[curr.key] = curr.value;
                    return acc;
                }, {});
                payload.multi_value = multiValueObj;
            }

            const data = await post(`/secrets/rotate?path=${encodeURIComponent(path?.path)}&key=${encodeURIComponent(secret.key)}`, payload);
            setMessage(`Secret rotated successfully. New version: ${data.version}`);
            onRotated(data);
        } catch (err) {
            setError(`Failed to rotate secret: ${err.message}`);
        } finally {
            setLoading(false);
        }
    };

    const handleFileUpload = (uploadedValues) => {
        setMultiValue(uploadedValues);
    };

    const handleMultiValueChange = (index, event) => {
        const values = [...multiValue];
        values[index][event.target.name] = event.target.value;
        setMultiValue(values);
    };

    const addMultiValueField = () => {
        setMultiValue([...multiValue, { key: "", value: "" }]);
    };

    return (
        <div className='mt-2'>
            {secret?.is_multi_value && (
                <div className='text-info'>
                    <FaInfoCircle /><strong className='ps-1'>This is a multi-key secret</strong>
                </div>
            )}
            <Card className='p-0'>
                <Card.Body>
                    {message && <p className='text-success fw-bold'>{message}</p>}
                    {error && <p className='text-danger fw-bold'>{error}</p>}
                    {path?.engine_type === "kv" && (
                        <div>


                            {secret?.is_multi_value ? (
                                <div>
                                    {multiValue.map((field, index) => (
                                        <InputGroup className='mt-1'>
                                            <div className='form-floating'>
                                                <Form.Control
                                                    type="text"
                                                    placeholder="Field"
                                                    name="key"
                                                    value={field.key}
                                                    onChange={(e) => handleMultiValueChange(index, e)}
                                                />
                                                <label>Field</label>
                                            </div>
                                            <div className='form-floating'>
                                                <Form.Control
                                                    type="text"
                                                    placeholder="Value"
                                                    name="value"
                                                    value={field.value}
                                                    onChange={(e) => handleMultiValueChange(index, e)}
                                                />
                                                <label>Value</label>
                                            </div>
                                        </InputGroup>
                                    ))}
                                    <Link variant="transparent" className='mt-1 ps-2' onClick={addMultiValueField}>
                                        Add Field
                                    </Link>

                                    <UploadSecrets onUpload={handleFileUpload} />

                                </div>
                            ) :


                                <div>
                                    <Form.Check
                                        className='form-checkbox mt-1'
                                        type="checkbox"
                                        checked={isOneTime}
                                        onChange={(e) => setIsOneTime(e.target.checked)}
                                        label="One-Time Use Secret?"
                                    />

                                    <div className="form-floating mt-2">
                                        <input
                                            type="text"
                                            placeholder="Enter new secret value"
                                            className='form-control'
                                            value={newSecretValue}
                                            onChange={(e) => setNewSecretValue(e.target.value)}
                                        />
                                        <label>New Secret Value:</label>
                                    </div>

                                    <div className='form-floating mt-3'>
                                        <Datetime
                                            value={expiresAt}
                                            onChange={(date) => setExpiresAt(date)}
                                            dateFormat="YYYY-MM-DD"
                                            timeFormat="HH:mm:ss"
                                            className='form-floating'
                                        />
                                        <label>Expires At (RFC3339):</label>
                                    </div>
                                </div>
                            }

                        </div>
                    )}

                    <button
                        className='w-100 mt-2 btn btn-sm btn-primary'
                        onClick={rotateSecret}
                        disabled={loading}
                    >
                        {loading ? 'Rotating...' : 'Rotate'}
                    </button>
                </Card.Body>
            </Card>

            <CLIUsage cmd="rotate [path] [key] [value]" />

        </div>
    );
};

export default RotateSecret;
