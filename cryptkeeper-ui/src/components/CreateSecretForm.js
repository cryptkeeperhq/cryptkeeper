import React, { useState, useEffect, useRef } from 'react';
import { Button, ButtonGroup, Card, Container, FormLabel, ListGroup, ListGroupItem, Row, Col, InputGroup, Form, Alert } from 'react-bootstrap';
import Datetime from 'react-datetime';
import moment from 'moment';
import "react-datetime/css/react-datetime.css";
import Title from './common/Title';
import { Link } from 'react-router-dom';
import { CreateSecretHelp, TransitEncryptionHelp } from './help/Help';
import TransitMetadataForm from './TransitMetadataForm';
import RenderMetadataFields from './RenderMetadataFields';
import UploadSecrets from './UploadSecrets';
import { FaCheck, FaKey, FaNode, FaRegSadTear, FaSadTear, FaSpinner, FaTimes, FaTrash } from 'react-icons/fa';
import { useApi } from '../api/api'
import TemplatesSelector from './TemplatesSelector';
import TagsForm from './TagsForm';
import CLIUsage from './CLIUsage';

const CreateSecretForm = ({ path, approvalRequired, onCreate }) => {
    const { get, post, put, del } = useApi();
    const [paths, setPaths] = useState([]);
    const [selectedPath, setSelectedPath] = useState('');
    const [key, setKey] = useState('');
    const [value, setValue] = useState('');
    const [expiresAt, setExpiresAt] = useState('');
    const [message, setMessage] = useState('');
    const [metadata, setMetadata] = useState(null);
    const [isOneTime, setIsOneTime] = useState(false);
    const [rotationInterval, setRotationInterval] = useState('');
    const [error, setError] = useState('');
    const hasFetchedPaths = useRef(false);
    const [isMultiValue, setIsMultiValue] = useState(false);
    const [multiValue, setMultiValue] = useState([{ key: "", value: "" }]);
    const [tags, setTags] = useState([]);

    const [loading, setLoading] = useState(false)
    const [keyLabel, setKeyLabel] = useState("Key")


    const [createPermission, setCreatePermission] = useState(false)

    useEffect(() => {
        if (!hasFetchedPaths.current) {
            hasFetchedPaths.current = true;
            fetchPaths();
        }
    }, []);

    useEffect(() => {
        if (path !== undefined) {
            setSelectedPath(path)
            console.log(path)
            getPathInfo(path.id)

            setMetadata({})
            // setMetadataString("{}")

            if (path.engine_type == "pki") {
                setKeyLabel("Common Name (ex: home.localhost.com)")
            }
            if (path.engine_type == "transit") {
                setKeyLabel("Key Name (ex: SSN, PAN, DOB)")
            }
            if (path.engine_type == "kv") {
                setKeyLabel("Key")
            }
            if (path.engine_type == "database") {
                setKeyLabel("Key")
            }

            setMetadata(path.metadata)


        }
    }, [path]);

    // useEffect(async () => {
    // }, [selectedPath]);

    const getPathInfo = async (pathId) => {
        console.log(selectedPath)
        await fetchPathPermission(pathId);
        setLoading(false)
    }
    const fetchPaths = async () => {
        try {
            const data = await get(`/user/paths`);
            setPaths(data || []);
        } catch (error) {
            console.log(error.message);
        }
    };

    const fetchPathPermission = async (path_id) => {
        try {
            setCreatePermission(false)
            const data = await get(`/paths/${path_id}/permissions`);
            console.log(data)
            if (data.permissions.create == true) {
                setCreatePermission(true)
            }
        } catch (error) {
            console.log(error.message);
        }
    };

    const handlePathChange = async (e) => {
        const pathId = e.target.value;
        if (pathId == "") {
            setSelectedPath("")
            return
        }

        setLoading(true)
        var _selectedPath

        for (var index = 0; index < paths.length; index++) {
            if (paths[index].id == pathId) {
                _selectedPath = paths[index]
            }
        }

        setSelectedPath(_selectedPath);
        
        getPathInfo(_selectedPath.id)

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


    function getPayload() {
        // var pathName
        // for (var index = 0; index < paths.length; index++) {
        //     if (paths[index].id == selectedPath) {
        //         pathName = paths[index].path
        //     }
        // }

        // Ensure the path starts with a '/'
        // let formattedKey = key.startsWith('/') ? key : `/${key}`;
        // // Ensure the path does not end with a '/'
        // formattedKey = formattedKey.endsWith('/') ? formattedKey.slice(0, -1) : formattedKey;

        const payload = {
            path_id: selectedPath.id,
            key: key,
            expires_at: expiresAt ? expiresAt.toISOString() : null,
            metadata: metadata ? metadata : null,
            is_one_time: isOneTime,
            rotation_interval: rotationInterval,
            is_multi_value: isMultiValue,
            tags: tags,
            path: selectedPath.path,
        };

        if (isMultiValue) {
            const multiValueObj = multiValue.reduce((acc, curr) => {
                acc[curr.key] = curr.value;
                return acc;
            }, {});
            payload.multi_value = multiValueObj;
        } else {
            payload.value = value;
        }

        return payload
    }

    const createSecret = async () => {
        // var pathName
        // for (var index = 0; index < paths.length; index++) {
        //     if (paths[index].id == selectedPath) {
        //         pathName = paths[index].path
        //     }
        // }

        const payload = getPayload()


        try {
            const data = await post(`/secrets?path=${selectedPath.path}`, payload);
            setMessage('Secret created successfully.');
            if (onCreate) {
                onCreate(data)
            }
        } catch (error) {
            setError(`Error creating secret: ${error.message}`);
        }


    };


    const createApprovalRequest = async (action, details) => {
        try {
            const data = await post(`/approval-requests`, { action, details });
            setMessage('Approval request created.');
            if (onCreate) {
                onCreate(data)
            }
        } catch (error) {
            setError(`Error creating secret: ${error.message}`);
        }
    };

    const handleCreateSecret = () => {
        createApprovalRequest('create', getPayload());
    };

    const handleUpdateSecret = () => {

        createApprovalRequest('update', getPayload());
    };

    const handleDeleteSecret = () => {
        createApprovalRequest('delete', getPayload());
    };

    return (
        <div className="">
            {message && <p className='mb-3 bg-success-soft text-success p-3 rounded-2 fw-bold' variant='info'>{message}</p>}
            {error && <p className='mb-3 bg-danger-soft text-danger p-3 rounded-2 fw-bold' variant='danger'>{error}</p>}

            {paths.length == 0 ?
                <div className='card p-3'>
                    <p className='lead'>Empty Collection!</p>
                    <p className='text-info fw-bold'>No paths to access</p>

                    <p>To get started, simply click on the 'Add Path' and we will guide you through the process.</p>
                </div>

                :
                <div className=''>

                   { !path && (
                     <Card className='mb-3'>
                     <Card.Header>Create New Secret </Card.Header>
                     <Card.Body>
                         <div className="form-floating">
                             <select
                                 className="form-control "
                                 id="createPath"
                                 value={selectedPath.id}
                                 onChange={handlePathChange}
                             >
                                 <option value="" >Select Path</option>
                                 {paths.map((path) => (
                                     <option key={path.id} value={path.id}>
                                         {path.path} ({path.engine_type})
                                     </option>
                                 ))}
                             </select>
                             <label htmlFor="createPath">Path:</label>
                         </div>
                     </Card.Body>
                     <Card.Footer>
                         <CLIUsage cmd="create [path] [key] [value]" />
                     </Card.Footer>
                 </Card>
                   )}

                    {loading && <div><FaSpinner className="spinner" /></div>}


                    { selectedPath && <div>
                        {!createPermission && <Alert variant="warning" className=''>
                        You don't have enough permissions on the path to create a secret. Go ahead and complete the form to trigger approval workflow.
                    </Alert>}

                    {!loading && <div className=''>

                        <Card className=''>
                            <Card.Header>Secret Details</Card.Header>
                            <Card.Body>
                                <div className="form-floating">
                                    <input
                                        type="text"
                                        value={key}
                                        className="form-control"
                                        id="createKey"
                                        onChange={(e) => setKey(e.target.value)}
                                    />
                                    <label htmlFor="createKey">{keyLabel}</label>
                                </div>

                                {selectedPath.engine_type === 'kv' &&
                                    <div>
                                        <div>
                                            {!isMultiValue && (
                                                <Form.Group className="form-floating mt-1">
                                                    <Form.Control
                                                        type="text"
                                                        value={value}
                                                        onChange={(e) => setValue(e.target.value)}
                                                    />
                                                    <Form.Label>Value</Form.Label>
                                                </Form.Group>
                                            )}

                                            {isMultiValue && (
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
                                            )}
                                            <Form.Group className="">
                                                <div className="p-1">
                                                    <Form.Check
                                                        type="checkbox"
                                                        label="Multi-Value Secret?"
                                                        checked={isMultiValue}
                                                        onChange={(e) => setIsMultiValue(e.target.checked)}
                                                    />
                                                </div>
                                            </Form.Group>
                                        </div>


                                    </div>}

                                <RenderMetadataFields metadata={selectedPath.metadata} engineType={selectedPath.engine_type} setMetadata={setMetadata} />
                            </Card.Body>

                        </Card>

                        {selectedPath.engine_type === 'kv' && (
                            <Card className='mt-3'>
                                <Card.Header>Additional Details</Card.Header>
                                <Card.Body>
                                    <div className="">
                                        <label className='small'>Expires At (RFC3339):</label>
                                        <Datetime
                                            value={expiresAt}
                                            onChange={setExpiresAt}
                                            dateFormat="YYYY-MM-DD"
                                            timeFormat="HH:mm:ss"
                                            inputProps={{ placeholder: 'YYYY-MM-DD HH:mm:ss' }}
                                        />
                                    </div>


                                    <div className="form-floating mt-1">
                                        <input
                                            type="text"
                                            value={rotationInterval}
                                            className="form-control"
                                            onChange={(e) => setRotationInterval(e.target.value)}
                                        />
                                        <label>Rotation Interval (e.g., 24h, 7d):</label>
                                    </div>

                                    {/* {engineType === 'kv' && */}
                                    <div>
                                        <div className="p-1">
                                            <Form.Check
                                                className='form-checkbox mt-1'
                                                type="checkbox"
                                                checked={isOneTime}
                                                onChange={(e) => setIsOneTime(e.target.checked)}
                                                label="One-Time Use Secret?"
                                            />
                                        </div>
                                    </div>
                                    {/* } */}

                                </Card.Body>
                            </Card>
                        )}


                        <TagsForm updateTags={(value) => setTags(value)} />

                        <div className='mt-3'>
                            {createPermission ?
                                <button className='btn w-100 btn-success' onClick={createSecret}>Create Secret</button> :
                                <div>
                                    <Button variant="btn w-100 btn-success" className="w-100" onClick={handleCreateSecret}>Submit Secret For Approval</Button>
                                    {/* <ButtonGroup className='mt-3'> */}
                                    {/* <Button variant="warning" onClick={handleUpdateSecret}>Update Secret</Button>
                                    <Button variant="danger" className='float-end' onClick={handleDeleteSecret}><FaTrash size={20} className='pe-2' /> Delete Secret</Button> */}
                                    {/* </ButtonGroup> */}
                                </div>
                            }
                        </div>

                    </div>}

                    </div>}

                </div>

            }

        </div>
    );
};

export default CreateSecretForm;
