// src/components/SecretVersion.js

import React, { useState, useEffect } from 'react';
import { Card, ListGroup, ListGroupItem, Alert, Dropdown, Tabs, Tab, Modal, Button } from 'react-bootstrap';
import { FaFileAlt, FaTrashAlt, FaEye, FaEyeSlash, FaEllipsisV, FaLink, FaDownload, FaTerminal } from 'react-icons/fa';
import MetadataItem from './MetadataItem';
import RenderNestedObject from './utils/UnflattenAndDisplay';
import { useApi } from '../api/api';
import moment from "moment";
import { Link } from 'react-router-dom';
import CLIUsage from './CLIUsage';
import TransitEncryption from './TransitEncryption';

const SecretDetails = ({ path, secret, onDelete }) => {
    const { get, del, downloadGet } = useApi();
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [secretValue, setSecretValue] = useState("");
    const [showDetail, setShowDetail] = useState(false);

    useEffect(() => {
        setMessage("");
        setError('');
        setSecretValue(secret.value);
    }, [path, secret]);

    const deleteSecret = async (path, key, version) => {
        try {
            await del(`/secrets/delete?path=${encodeURIComponent(path)}&key=${encodeURIComponent(key)}&version=${version}`);
            setMessage('Secret deleted.');
            onDelete();
        } catch (error) {
            setError(error.message);
        }
    };

    const showSecretValue = async (path, key, version) => {
        if (secretValue !== '') {
            setSecretValue('');
            return;
        }

        try {
            const result = await get(`/secrets/version?path=${encodeURIComponent(path)}&key=${encodeURIComponent(key)}&version=${version}`);
            setSecretValue(result.value);
        } catch (error) {
            setError(error.message);
        }
    };

    const downloadCertificate = async () => {
        if (path.engine_type !== "pki") {
            alert('No certificate available for download.');
            return;
        }

        try {
            const data = await downloadGet(`/pki/download-certificate?secret_id=${secret.id}`);
            const blob = await data.blob();
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = secret.key + '.p12';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (error) {
            setError(error.message);
        }
    };

    const downloadCA = async () => {
        if (path.engine_type !== "pki") {
            alert('No certificate available for download.');
            return;
        }

        try {
            const data = await downloadGet(`/pki/download-ca?path_id=${secret.path_id}`);
            const blob = await data.blob();
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = 'ca.pem';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (error) {
            setError(error.message);
        }
    };

    const unflattenObject = (flatObj) => {
        const result = {};
        try {
            flatObj = JSON.parse(flatObj);
            Object.keys(flatObj).forEach((key) => {
                const keys = key.split('.');
                keys.reduce((acc, currentKey, index) => {
                    if (index === keys.length - 1) {
                        acc[currentKey] = flatObj[key];
                        return acc;
                    }
                    if (!acc[currentKey]) {
                        acc[currentKey] = {};
                    }
                    return acc[currentKey];
                }, result);
            });
        } catch (error) {
            console.log("Error parsing Object: " + error.message);
        }
        return result;
    };

    return (
        <div>



            <Card>
                <Card.Header>
                    <div className='d-flex align-items-center'>
                        <FaFileAlt size={18} className='pe-2' /> Secret Details
                        <span className='ms-auto'>
                            <Link className='text-muted text-decoration-none' to={`/user/secrets/${secret.id}/${secret.version}`}>
                                <FaLink className='d-line me-1' />
                            </Link>
                        </span>
                        <span>
                            {onDelete && (
                                <Dropdown >
                                    <Dropdown.Toggle variant="transparent" id="dropdown-basic" className='p-0 ms-2'>
                                        <FaEllipsisV />
                                    </Dropdown.Toggle>
                                    <Dropdown.Menu >
                                        
                                        {path.engine_type === "pki" && (
                                            <>
                                                <Dropdown.Item className='small' onClick={downloadCertificate}><FaDownload /> Download Certificate</Dropdown.Item>
                                                <Dropdown.Divider />
                                                <Dropdown.Item className='small' onClick={downloadCA}><FaDownload /> Download CA</Dropdown.Item>
                                                <Dropdown.Divider />
                                            </>
                                        )}
                                        {path.engine_type === "transit" && (
                                            <>
                                            <Dropdown.Item className='small' onClick={() => setShowDetail(true)}><FaTerminal /> Explore Operations</Dropdown.Item>
                                            <Dropdown.Divider />
                                            </>
                                        )}

                                        <Dropdown.Item className='small text-danger' onClick={(e) => { e.stopPropagation(); deleteSecret(path.path, secret.key, secret.version); }}>
                                            <FaTrashAlt /> Delete Version
                                        </Dropdown.Item>


                                    </Dropdown.Menu>
                                </Dropdown>
                            )}
                        </span>
                    </div>
                </Card.Header>

                {message && <p className='ps-3 text-success fw-bold'><strong>{message}</strong></p>}
                {error && <p className='ps-3 text-danger fw-bold'><strong>{error}</strong></p>}

                <ListGroup variant='flush'>
                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Path</span>
                        <code className='ms-auto'>{path ? path.path : secret.path}</code>
                    </ListGroupItem>

                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Key</span>
                        <code className='ms-auto'>{secret.key}</code>
                    </ListGroupItem>

                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Version</span>
                        <code className='ms-auto'>{secret.version}</code>
                    </ListGroupItem>

                    {path.engine_type !== "transit"  && (
                        <ListGroupItem>
                            <div className='d-flex align-items-center'>
                                <span>Value {secret.is_multi_value ? " (Multi-Value)" : ""}</span>
                                {secret.is_one_time && <span className="badge bg-warning ms-2">One-Time</span>}
                                <span className='ms-auto'>
                                    <span onClick={(e) => { e.stopPropagation(); showSecretValue(path.path, secret.key, secret.version); }}>
                                        {secretValue === '' ? <FaEye size={16} /> : <FaEyeSlash size={16} />}
                                    </span>
                                </span>
                            </div>
                            <div>
                                {secretValue !== '' ? (
                                    secret.is_multi_value ? (
                                        <code><RenderNestedObject data={unflattenObject(secretValue)} /></code>
                                    ) : (
                                        <code>{secretValue}</code>
                                    )
                                ) : (
                                    <code>**********</code>
                                )}
                            </div>
                        </ListGroupItem>
                    )}

                    {secret.tags && secret.tags.length > 0 && (
                        <ListGroupItem className='d-flex align-items-center'>
                            {secret.tags.map((tag, index) => (
                                <span key={index} className='badge bg-primary-soft text-primary me-1 mb-1'>{tag}</span>
                            ))}
                        </ListGroupItem>
                    )}



                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Created</span>
                        <span className='ms-auto badge bg-light text-dark'>{moment(secret.created_at).format("MMMM Do YYYY, h:mm:ss a")}</span>
                    </ListGroupItem>

                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Updated</span>
                        <span className='ms-auto badge bg-light text-dark'>{moment(secret.updated_at).format("MMMM Do YYYY, h:mm:ss a")}</span>
                    </ListGroupItem>

                    <ListGroupItem className='d-flex align-items-center'>
                        <span>Last Rotated</span>
                        <span className='ms-auto badge bg-light text-dark'>{moment(secret.last_rotated_at).format("MMMM Do YYYY, h:mm:ss a")}</span>
                    </ListGroupItem>

                    {secret.expires_at && (
                        <ListGroupItem className='d-flex align-items-center'>
                            <span>Expires At</span>
                            <span className='ms-auto badge bg-secondary text-dark'>{moment(secret.expires_at).format("MMMM Do YYYY, h:mm:ss a")}</span>
                            {secret.expires_at && new Date(secret.expires_at) < new Date() && <span className="badge bg-danger ms-2">Expired</span>}
                        </ListGroupItem>

                    )}

<ListGroupItem className='d-flex align-items-center'>
                        <span>Created By</span>
                        <span className='ms-auto badge bg-light text-dark'>{secret.created_by}</span>
                    </ListGroupItem>

                </ListGroup>

                {secret.rotation_interval && (
                    <Card.Body>
                        <Alert variant='success' className='small p-2'>
                            <b>Rotation is enabled every <code>{secret.rotation_interval}</code>.<br />The secret was last rotated at <code>{secret.last_rotated_at}</code></b>
                        </Alert>
                    </Card.Body>
                )}

                <Card.Footer>
                    <CLIUsage cmd="get [path] [key] --version=0" />
                </Card.Footer>
            </Card>


            <Modal size="lg" show={showDetail} onHide={() => setShowDetail(false)}>
                <Modal.Header closeButton>
                    <Modal.Title>
                        Explore Cryptography Operations on the keys created by Transit Engine

                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <TransitEncryption keyId={secret.id} />
                </Modal.Body>
                <Modal.Footer>


                    <Button variant="secondary" onClick={() => setShowDetail(false)}>Close</Button>
                </Modal.Footer>
            </Modal>


        </div>
    );
};

export default SecretDetails;
