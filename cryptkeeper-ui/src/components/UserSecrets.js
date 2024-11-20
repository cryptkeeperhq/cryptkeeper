import React, { useState, useEffect, useRef } from 'react';
import { Nav, Card, Container, Dropdown, Breadcrumb, ListGroup, ListGroupItem, Row, Col, Form, Button, FloatingLabel, Alert, Accordion, Tabs, Tab, Modal } from 'react-bootstrap';
import SecretAccessControl from './SecretAccessControl';
import RotateSecret from './RotateSecret';
import SecretMetadataUpdate from './SecretMetadataUpdate';
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

import SecretExistingAccess from './SecretExistingAccess';
import "react-datetime/css/react-datetime.css";
import { FaKey, FaList, FaDatabase, FaCertificate, FaFolder, FaHistory, FaTrashAlt, FaEdit, FaUser, FaPlus, FaShare, FaLink, FaEllipsisV, FaNode, FaStickyNote } from 'react-icons/fa';
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
import Secret from './Secret';


const UserSecrets = ({ setTitle, setHelp }) => {
    const { engine } = useParams();
    const navigate = useNavigate();


    const { get, post, put, del } = useApi();
    const [paths, setPaths] = useState([]);

    const [secrets, setSecrets] = useState([]);
    const [selectedSecretPath, setSelectedSecretPath] = useState(null);
    const [selectedKey, setSelectedKey] = useState(null);
    const [selectedVersion, setSelectedVersion] = useState(null);
    // const [searchParams, setSearchParams] = useState({ path: '', key: '', metadata: '', createdAfter: '', createdBefore: '' });
    const [currentTab, setCurrentTab] = useState("details");

    const [createSecret, setCreateSecret] = useState(false);


    const [secretFilterType, setSecretFilterType] = useState('active')
    const hasFetchedPaths = useRef(false);

    useEffect(() => {
        setHelp("")

    }, []);


    useEffect(() => {
        // if (!hasFetchedPaths.current) {
        // hasFetchedPaths.current = true;
        setTitle({ heading: " Secrets", subheading: "Engine: " + engine })

        setSelectedSecretPath(null)
        setSelectedKey(null)
        setSecrets([])
        fetchPaths();

        // }
    }, [engine]);


    const getPathIcon = (engineType) => {
        switch (engineType) {
            case 'pki':
                return (
                    <FaCertificate className='text-primary me-2' />
                );
            case 'transit':
                return (
                    <FaKey className='text-primary me-2' />
                );
            case 'kv':
                return (
                    <FaList className='text-primary me-2' />
                );
            case 'database':
                return (
                    <FaDatabase className='text-primary me-2' />
                );
            default:
                return (
                    <FaFolder className='text-primary me-2' />
                );
        }
    };
    const fetchPaths = async () => {

        try {
            const data = await get(`/user/paths?engine=${engine}`);
            setPaths(data || []);
        } catch (error) {
            console.log(error.message);
        }


    };




    const getSecrets = async (path) => {


        try {
            const data = await get(`/secrets?path=${path.path}`);
            setSecrets(data || []);
            // setSelectedKey(null);
            // setSelectedVersion(null);
        } catch (error) {
            console.log(error.message);
        }


    }

    const handlePathClick = (path) => {
        setSelectedSecretPath(path);
        getSecrets(path)
        setSelectedKey(null);
        setSelectedVersion(null);
        // setTitle({ heading: "My Secrets", subheading: path.path })
    };

    const handleKeyClick = (key) => {
        setSelectedKey(null);
        setSelectedVersion(null);
        setSelectedKey(key);

        // setTitle({ heading: "My Secrets", subheading: selectedSecretPath.path + " <span className='badge bg-primary'>" + key + "</span>" })

        if (groupedSecrets[selectedSecretPath.id][key].length > 0) {
            var v = groupedSecrets[selectedSecretPath.id][key][groupedSecrets[selectedSecretPath.id][key].length - 1]
            setSelectedVersion(v)

        }

    };


    const onSecretCreate = async (data) => {
        setCreateSecret(false)
        getSecrets(selectedSecretPath)
     };
 
 



    const groupedSecrets = secrets.reduce((acc, secret) => {
        if (!acc[secret.path_id]) {
            acc[secret.path_id] = {};
        }
        if (!acc[secret.path_id][secret.key]) {
            acc[secret.path_id][secret.key] = [];
        }
        acc[secret.path_id][secret.key].push(secret);
        return acc;
    }, {});

    return (
        <div className="mb-4">


            <Container className='p-0'>
                <Row>
                    <Col xl={12} lg={12}>
                        <Container className='p-0'>
                            <Row>
                                <Col xl={12}>
                                    <Nav variant="underline" activeKey={`${secretFilterType}`} className='m-0 mb-3'>
                                        <Nav.Item onClick={() => setSecretFilterType("active")}>
                                            <Nav.Link eventKey="active" className={`ms-2 me-auto text-center text-center`} >
                                                <span className="fw-bold">Active</span>
                                            </Nav.Link>
                                        </Nav.Item>
                                        <Nav.Item onClick={() => setSecretFilterType("deleted")}>
                                            <Nav.Link eventKey="deleted" className={`ms-2 me-auto text-center text-center`} >
                                                <span className="fw-bold">Deleted</span>
                                            </Nav.Link>
                                        </Nav.Item>
                                    </Nav>


                                </Col>
                            </Row>

                            {
                                secretFilterType == "active" &&
                                <Row>


                                    <Col xs={12} lg={3} >

                                        <Card className='mb-3'>
                                            <Card.Header>
                                            <div className='d-flex align-items-center'>
                                                                Paths
                                                                {/* <span className='ms-auto'><Button variant='success' onClick={() => navigate("/paths")} className='p-0 p-1 ps-2 pe-2'><FaPlus className='' /></Button></span> */}
                                                            </div>


                                            </Card.Header>
                                            {paths.length == 0 &&
                                                <div className='ps-3 pe-3'>
                                                    <p className='lead'>Empty Collection!</p>
                                                    <p className='text-info fw-bold'>No paths to access</p>

                                                    <p>To get started, simply click on the 'Add Path' and we will guide you through the process.</p>
                                                </div>
                                            }

                                            <ListGroup variant="flush">
                                                {paths.map(path => (
                                                    <ListGroupItem
                                                        key={path.id}
                                                        className={`p-0 list-group-item ${selectedSecretPath && selectedSecretPath === path ? 'active fw-bold' : ''}`}
                                                        style={{ cursor: 'pointer' }}
                                                    >
                                                        <div className='p-2 ps-3' onClick={() => handlePathClick(path)}>
                                                            {/* {getPathIcon(path.engine_type)} */}
                                                            <FaFolder className='me-2' />
                                                            {path.path}
                                                            {/* ({path.engine_type}) */}
                                                        </div>
                                                    </ListGroupItem>
                                                ))}
                                            </ListGroup>

                                            <Card.Footer>
                                                <CLIUsage cmd="paths" />
                                            </Card.Footer>

                                        </Card>
                                    </Col>

                                    <Col xs={12} lg={3} >

                                    {selectedSecretPath && !groupedSecrets[selectedSecretPath.id] &&
                                                 <Card className='mb-3'>
                                                 <Card.Header>
                                                     <div className='d-flex align-items-center'>
                                                         Keys
                                                         <span className='ms-auto'>
                                                             <Button variant='success' onClick={() => setCreateSecret(true)} className='p-0 p-1 ps-2 pe-2'>
                                                                 <FaPlus className='me-1' /> New</Button></span>
                                                     </div>


                                                 </Card.Header>
<Card.Body>
                                                    <p className='lead'>Empty Path Collection!</p>
                                                    <p className='text-info fw-bold'>No secrets to access</p>
                                                    <p>To get started, simply click on the 'Add Secret' and we will guide you through the process.</p>
                                                    </Card.Body>
                                                    </Card>

                                            }

                                        {selectedSecretPath && <>
                                            {groupedSecrets[selectedSecretPath.id] && (
                                                <>




                                                    <Card className='mb-3'>
                                                        <Card.Header>
                                                            <div className='d-flex align-items-center'>
                                                                Keys
                                                                <span className='ms-auto'>
                                                                    <Button variant='success' onClick={() => setCreateSecret(true)} className='p-0 p-1 ps-2 pe-2'>
                                                                        <FaPlus className='me-1' /> New</Button></span>
                                                            </div>


                                                        </Card.Header>

                                                        <ListGroup variant="flush">
                                                            {Object.keys(groupedSecrets[selectedSecretPath.id]).map(key => (
                                                                <ListGroupItem
                                                                    key={key}
                                                                    className={`list-group-item ${selectedKey === key ? 'active fw-bold' : ''}`}
                                                                >

                                                                    <div className='d-flex align-items-center'>

                                                                        <div onClick={() => handleKeyClick(key)} style={{ cursor: 'pointer' }}>
                                                                            {/* <FaKey size={10} className='text-danger me-2' /> */}

                                                                            {getPathIcon(selectedSecretPath.engine_type)}

                                                                            {key}
                                                                        </div>




                                                                        {selectedKey && selectedKey == key && (
                                                                            <span className='ms-auto'><FaHistory className='me-1' /> <small>{selectedVersion.version}</small></span>
                                                                        )}




                                                                    </div>
                                                                </ListGroupItem>
                                                            ))}
                                                        </ListGroup>
                                                        <Card.Footer>
                                                            <CLIUsage cmd="secrets [path]" />
                                                        </Card.Footer>
                                                    </Card>
                                                </>
                                            )}


                                            
                                        </>}


                                    </Col>



                                    <Col xs={12} lg={6}>
                                        {selectedKey && selectedVersion && (
                                            <Secret secret_id={selectedVersion.id} secret_version={selectedVersion.version} />
                                        )}
                                    </Col>
                                </Row>
                            }

                            {secretFilterType == "deleted" &&
                                <Row>
                                    <Col>
                                        <DeletedSecrets paths={paths} />
                                    </Col>
                                </Row>
                            }

                        </Container>

                    </Col>

                </Row>
            </Container>

            <Modal size="lg" show={createSecret} onHide={() => setCreateSecret(false)}>
                <Modal.Header closeButton>
                    <Modal.Title>Create Secret</Modal.Title>
                </Modal.Header>
                <Modal.Body>

                    <CreateSecretForm path={selectedSecretPath} approvalRequired={false} onCreate={() => onSecretCreate()} />


                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={() => setCreateSecret(false)}>Close</Button>
                </Modal.Footer>
            </Modal>



        </div>
    );
};

export default UserSecrets;
