import { useParams } from 'react-router-dom';
import React, { useState, useEffect, useRef } from 'react';
import { Nav, Card, Container, Dropdown, Breadcrumb, ListGroup, ListGroupItem, Row, Col, Form, Button, FloatingLabel, Alert, Accordion, Tabs, Tab, Modal } from 'react-bootstrap';
import SecretAccessControl from './SecretAccessControl';
import RotateSecret from './RotateSecret';
import SecretMetadataUpdate from './SecretMetadataUpdate';

import SecretExistingAccess from './SecretExistingAccess';
import "react-datetime/css/react-datetime.css";
import { FaKey, FaList, FaDatabase, FaCertificate, FaFolder, FaHistory, FaTrashAlt, FaEdit, FaUser, FaPlus, FaShare, FaLink, FaEllipsisV } from 'react-icons/fa';
import SecretLineage from './SecretLineage';
import SecretDetails from './SecretDetails'
import SecretHistory from './SecretHistory';
import DeletedSecrets from './DeletedSecrets';
import Title from './common/Title';
import moment from "moment";
import { useApi } from '../api/api'
import SearchSecrets from './SearchSecrets'
import { Link } from 'react-router-dom';
import ShareSecret from './ShareSecret';
import { PathManagementHelp } from './help/Help';
import CreateSecretForm from './CreateSecretForm';
import MetadataItem from './MetadataItem';
import CLIUsage from './CLIUsage';


const Secret = ({ setTitle, secret_id, secret_version }) => {
    const { id: paramId, version: paramVersion } = useParams();

    const id = secret_id || paramId;
    const version = secret_version || paramVersion;


    const { get, post, put, del } = useApi();
    const [secret, setSecret] = useState(null);
    const [path, setPath] = useState(null)
    const [error, setError] = useState('');

    const [currentTab, setCurrentTab] = useState("details");

    const [selectedKey, setSelectedKey] = useState(null);



    // const hasFetchedSecret = useRef(false);
    // useEffect(() => {
    //     if (!hasFetchedSecret.current) {
    //         hasFetchedSecret.current = true;
    //         fetchSecret()
    //     }
    // }, [id]);


    useEffect(() => {
        fetchSecret()
    }, [id, version]);




    const fetchSecret = async () => {
        // console.log("version", version)
        try {
            const data = await get(`/secrets/secret/${id}?version=${version}`);
            setSecret(data);
            fetchPath(data.path_id)
            if (setTitle) {
                setTitle({ heading: "Secret Details", subheading: data.path + data.key })
            }
        } catch (error) {
            setError(error.message);
        }
    }

    const fetchPath = async (path) => {

        try {
            const data = await get(`/paths/${path}`);
            setPath(data || {});
        } catch (error) {
            console.log(error.message);
        }


    };

    const rotated = async (data) => {
        setSecret(data)
        // getSecrets(path)
    }


    const handleVersionClick = (version) => {
        setSecret(version);
        setCurrentTab("details");


    };


    function showSecretDetails() {
        return (
            <div>
                <Tabs
                    defaultActiveKey="details"
                    activeKey={currentTab}
                    id="secret-mgt"
                    variant='underline'
                    className="mb-3 ps-2"
                    onSelect={(key) => setCurrentTab(key)}
                >
                    <Tab eventKey="details" title={"Details (v" + secret.version + ")"}>
                        <>
                            <div className=' '><SecretDetails path={path} secret={secret} onDelete={() => getSecrets(path)} /></div>
                            <Accordion className='mt-3' >
                                <Accordion.Item eventKey="0">
                                    <Accordion.Header><FaEdit className='me-2' /> Create New Version (Rotate Secret)</Accordion.Header>
                                    <Accordion.Body>
                                        <RotateSecret onRotated={(data) => rotated(data)} secret={secret} path={path} selected_key={selectedKey} />
                                    </Accordion.Body>
                                </Accordion.Item>
                            </Accordion>

                            {/* <Button variant="danger" className='mt-3'>Destroy Secret</Button> */}

                        </>
                    </Tab>
                    <Tab eventKey="metadata" title="Metadata">

                        {Object.keys(secret.metadata).length !== 0 && (
                            <ListGroup variant='flush'>
                                <ListGroupItem>
                                    <MetadataItem metadata={secret.metadata} />
                                </ListGroupItem>
                            </ListGroup>
                        )}


                        <Accordion className='mt-3'>

                            <Accordion.Item eventKey="2">
                                <Accordion.Header><FaEdit className='me-2' /> Metadata Update</Accordion.Header>
                                <Accordion.Body>
                                    <SecretMetadataUpdate onUpdate={() => getSecrets(path)} path={path} secret={secret} />
                                </Accordion.Body>
                            </Accordion.Item>



                        </Accordion>

                    </Tab>

                    <Tab eventKey="versions" title="Versions">

                        <div className=''><SecretHistory onVersionSelect={handleVersionClick} path={path} secret={secret} /></div>
                    </Tab>


                    <Tab eventKey="share" title="Share Secret">

                        <div>
                            {path.engine_type == "kv" ? <>
                                <ShareSecret path={path} secret={secret} />
                            </> : <Alert variant='info'>Engine does not support sharing secrets!</Alert>}
                        </div>

                    </Tab>
                </Tabs>

            </div>
        )

    }



    if (error) {
        return <div>{error}</div>;
    }

    if (!secret) {
        return <div>Loading...</div>;
    }

    return (
        <div className=''>
            {secret && path && showSecretDetails()}

         
        </div>
    );
};

export default Secret;
