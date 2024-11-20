import React, { useState, useEffect, useRef } from 'react';
import { Alert, Nav, Card, Container, Form, Button, Row, Col } from 'react-bootstrap';
// import PolicyEditor from './PolicyEditor';
import { FaTrash } from 'react-icons/fa';
import Title from '../common/Title';
import { PolicyManagementHelp } from '../help/Help';
import CodeMirror from "@uiw/react-codemirror";
import { vscodeDark } from "@uiw/codemirror-theme-vscode";
import { json } from '@codemirror/lang-json';
import { useApi } from '../../api/api'
import PolicyAuditLogs from '../PolicyAuditLogs';
import { oneDark } from '@codemirror/theme-one-dark';

const PolicyManagement = ({ setTitle, setHelp }) => {
    const { get, post, put, del } = useApi();
    const [policy, setPolicy] = useState({});
    const [groups, setGroups] = useState([]);
    const [users, setUsers] = useState([]);
    const [approles, setAppRoles] = useState([]);
    const [paths, setPaths] = useState([]);
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [pathID, setPathID] = useState('');
    const [policyID, setPolicyID] = useState(null);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [hcl, setHCL] = useState('');

    const [updatedHCL, setUpdatedHCL] = useState('')


    const [action, setAction] = useState("read");

    const [option, setOption] = useState("default")

    const hasFetchDone = useRef(false);
    useEffect(() => {
        if (!hasFetchDone.current) {
            hasFetchDone.current = true;
            fetchPaths()
        }
    }, []);

    useEffect(() => {
        setTitle({ heading: "Policy Management", subheading: "Manage policies and permission for the paths" })
        setHelp(
            <>
                <PolicyManagementHelp />
            </>
        )
    }, []);

    const fetchPaths = async () => {
        try {
            const data = await get(`/paths`);
            setPaths(data || []);
        } catch (error) {
            setError(error.message);
        }
    };

    const fetchPolicy = async (pathId) => {
        try {
            const data = await get(`/paths/${pathId}/policy`);
            setPolicy(data || {});
            if (data) {
                setPolicyID(data.id);
                setName(data.name);
                setDescription(data.description);
                setHCL(data.hcl || "")
                setUpdatedHCL(data.hcl || "")
            }
        } catch (error) {
            setError(error.message);
        }
    };


    const handleDelete = async () => {
        if (policy && policy.id) {
            try {
                const data = await del(`/policies/${policyID}`);
                setMessage('Policy deleted successfully.');
                setPathID('')
                setPolicy({})
                setPolicyID(null);
                setName("");
                setDescription("");
            } catch (error) {
                setError(error.message);
            }

        }
    };

    const savePolicy = async () => {
        const payload = { id: policyID, name, description, path_id: pathID, hcl: updatedHCL };

        try {
            const data = await post(`/policies`, payload);
            fetchPolicy(pathID);
            setMessage('Policy saved/udpated successfully.');
        } catch (error) {
            setError(error.message);
        }
    };


    const handlePathChange = (e) => {
        setPathID(e.target.value);
        fetchPolicy(e.target.value);
    };

    const onChange = React.useCallback((val, viewUpdate) => {
        setUpdatedHCL(val);
    }, []);


    return (
        <div className="">
            <Container className="p-0">
                <Row>


                    <Col >

                        <Nav variant="underline" activeKey={`${option}`} className='mb-3'>
                            <Nav.Item onClick={() => setOption("default")}>
                                <Nav.Link eventKey="default" className={`ms-2 me-auto text-center text-center `} >
                                    <span className="fw-bold">Manage</span>
                                </Nav.Link>
                            </Nav.Item>
                            <Nav.Item onClick={() => setOption("logs")}>
                                <Nav.Link eventKey="logs" className={`ms-2 me-auto text-center text-center`} >
                                    <span className="fw-bold">Audit Logs</span>
                                </Nav.Link>
                            </Nav.Item>

                        </Nav>

                        {option == "default" && <>
                            <Card>
                                <Card.Header>
                                    <div className="d-flex align-items-center">
                                        Create/Edit Policy
                                        {policyID && pathID && <Button variant="transparent" className="ms-auto text-danger ms-2" onClick={handleDelete}><FaTrash /></Button>}
                                    </div>
                                </Card.Header>
                                <Card.Body>

                                    {message && <div className='bg-success-soft text-success  p-3 rounded-2 mb-2'>{message}</div>}
                                    {error && <div className='bg-danger-soft text-danger  p-3 rounded-2 mb-2'>{error}</div>}

                                    <Form>
                                        <Form.Group className='form-floating '>
                                            <Form.Control
                                                as="select"
                                                name="path"
                                                value={pathID}
                                                onChange={handlePathChange}
                                            >
                                                <option value="">Select Path</option>
                                                {paths.map((path) => (
                                                    <option key={path.id} value={path.id}>{path.path}</option>
                                                ))}
                                            </Form.Control>
                                            <Form.Label>Path</Form.Label>
                                        </Form.Group>

                                        <Form.Group controlId="policyName" className='form-floating mt-1'>
                                            <Form.Control
                                                type="text"
                                                value={name}
                                                onChange={(e) => setName(e.target.value)}
                                                placeholder="Enter policy name"
                                            />
                                            <Form.Label>Policy Name</Form.Label>
                                        </Form.Group>

                                        <Form.Group controlId="policyDescription" className='form-floating mt-1'>
                                            <Form.Control
                                                type="text"
                                                value={description}
                                                onChange={(e) => setDescription(e.target.value)}
                                                placeholder="Enter policy description"
                                            />
                                            <Form.Label>Policy Description</Form.Label>
                                        </Form.Group>


                                        <div className='mt-2 mb-2'>


                                            {pathID &&
                                                <CodeMirror
                                                    value={hcl}
                                                    extensions={[json()]}
                                                    theme={oneDark}
                                                    onChange={onChange}
                                                    height="400px"
                                                    options={{
                                                        inlineSuggest: true,
                                                        fontSize: "12px",
                                                        formatOnType: true,
                                                        indentUnit: 2,
                                                        smartIndent: true,
                                                        autoClosingBrackets: true,
                                                        minimap: { scale: 10 },
                                                    }}

                                                />
                                            }
                                        </div>

                                        <Button variant="primary" className='w-100 mt-1' onClick={savePolicy}>
                                            {policyID ? 'Update Policy' : 'Create Policy'}
                                        </Button>



                                    </Form>



                                </Card.Body>
                            </Card>

                            <div className='mt-2 rounded-2 p-3 bg-dark text-white mb-3'>
                                <strong>Debug:</strong><br />
                                Path ID: <code>{pathID}</code> | Policy ID: <code>{policyID}</code>
                            </div>


                        </>}


                        {option == "logs" && <>
                            {/* <PolicyEditor /> */}
                            <PolicyAuditLogs />
                        </>}
                    </Col>

                </Row>

            </Container>
        </div>
    );
};

export default PolicyManagement;
