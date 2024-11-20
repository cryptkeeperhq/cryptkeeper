import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Form, Button, Alert, InputGroup } from 'react-bootstrap';
// import CreatePolicy from '../CreatePolicy';
// import PolicyEditor from './PolicyEditor';
// import PolicyManagement from './PolicyManagement';
import Title from '../common/Title';
import { PathManagementHelp } from '../help/Help';
import { useApi } from '../../api/api'
import CodeMirror from "@uiw/react-codemirror";
import { vscodeDark } from "@uiw/codemirror-theme-vscode";
import { json } from '@codemirror/lang-json';
import { FaSearch } from 'react-icons/fa';
import AddCA from './AddCA';
import AddTemplate from './AddTemplate';
import AddCertificateRequest from './AddCertificateRequest';


const PathManagement = ({ setTitle, setHelp }) => {
    const { get, post, put, del } = useApi();
    // const [paths, setPaths] = useState([]);
    const [path, setPath] = useState('');
    const [engineType, setEngineType] = useState(null);
    const [metadata, setMetadata] = useState({});
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [paths, setPaths] = useState([]);

    const [cas, setCas] = useState([]);
    


    const [id, setID] = useState(0);

    async function fetchCas() {
        const data = await get('/pki/ca'); // Replace with your actual API endpoint
        setCas(data || []);
    }



    const default_connection_string = `postgresql://{{username}}:{{password}}@localhost:5432/cryptkeeper?sslmode=disable`
    const default_role_template = `CREATE ROLE {{name}} WITH LOGIN PASSWORD '{{password}}';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO {{name}};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO {{name}};`


    useEffect(() => {
        setTitle({ heading: "Path Management", subheading: "Manage paths to store your secrets" })
        setHelp(
            <>
            <PathManagementHelp />
            </>
        )
        fetchPaths()
        fetchCas()
    }, []);



    const createPath = async () => {
        // Ensure the path starts with a '/'
        // let formattedPath = path.startsWith('/') ? path : `/${path}`;
        // // Ensure the path does not end with a '/'
        // formattedPath = formattedPath.endsWith('/') ? formattedPath.slice(0, -1) : formattedPath;

        let formattedPath = path;
        
        try {
            if (id == 0) {
                const data = await post(`/paths`, { path: formattedPath, engine_type: engineType, metadata });
            } else {
                const data = await put(`/paths`, { id: id, path: formattedPath, metadata });
            }
            setMessage('Path Updated.');
        } catch (error) {
            setError(error.message);
            console.log(error)
        }

    };

    const fetchPaths = async () => {
        try {
            const data = await get(`/paths`);
            setPaths(data || [])
        } catch (error) {
            setMessage("No paths found with that name");
        }
    };

    const fetchPath = async () => {
        try {
            // Ensure the path starts with a '/'
            let formattedPath = path.startsWith('/') ? path : `/${path}`;
            // Ensure the path does not end with a '/'
            formattedPath = formattedPath.endsWith('/') ? formattedPath.slice(0, -1) : formattedPath;

            const data = await get(`/paths?path=${encodeURIComponent(formattedPath)}`);
            if (data.length > 0) {
                setPath(data[0].path);
                setEngineType(data[0].engine_type);
                setMetadata(data[0].metadata);
                setID(data[0].id)
                renderMetadataFields(data[0].engine_type)
            }
        } catch (error) {
            setMessage("No paths found with that name");
        }
    };


    const handlePathChange = (e) => {
        console.log(e.target.value)

        if (e.target.value == "") {
            setPath('');
            setEngineType('');
            setMetadata({});
            setID(0)
            renderMetadataFields('')
            return
        }

        for (var v in paths) {
            if (paths[v].id == e.target.value) {

                console.log(paths[v])
                setPath(paths[v].path);
                setEngineType(paths[v].engine_type);
                setMetadata(paths[v].metadata);
                setID(paths[v].id)
                renderMetadataFields(paths[v].engine_type)
            }
        }


    };



    const renderMetadataFields = (engineType) => {
        switch (engineType) {
            case 'database':
                return (
                    <>
                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                type="text"
                                placeholder="Connection String"
                                value={metadata.connection_string || default_connection_string}
                                onChange={(e) => setMetadata({ ...metadata, connection_string: e.target.value })}
                            />
                            <Form.Label>Connection String</Form.Label>
                        </Form.Group>


                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                type="text"
                                as="textarea"
                                style={{ minHeight: "200px" }}
                                placeholder="Role Template"
                                value={metadata.role_template || default_role_template}
                                onChange={(e) => setMetadata({ ...metadata, role_template: e.target.value })}
                            />
                            <Form.Label>Role Template</Form.Label>
                        </Form.Group>
                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                type="text"
                                placeholder="Username"
                                value={metadata.username || ''}
                                onChange={(e) => setMetadata({ ...metadata, username: e.target.value })}
                            />
                            <Form.Label>Username</Form.Label>
                        </Form.Group>
                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                type="password"
                                placeholder="Password"
                                value={metadata.password || ''}
                                onChange={(e) => setMetadata({ ...metadata, password: e.target.value })}
                            />
                            <Form.Label>Password</Form.Label>
                        </Form.Group>
                    </>
                );
            case 'kv':
                return (
                    <>
                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                type="number"
                                placeholder="Maximum Number of Versions"
                                value={metadata.max_versions || ''}
                                onChange={(e) => setMetadata({ ...metadata, max_versions: e.target.value })}
                            />
                            <Form.Label>Maximum Number of Versions</Form.Label>
                        </Form.Group>
                        <Form.Group className="form-floating mt-2">
                            <Form.Check
                                type="checkbox"
                                label="Automated Secret Deletion"
                                checked={metadata.auto_delete || false}
                                onChange={(e) => setMetadata({ ...metadata, auto_delete: e.target.checked })}
                            />
                        </Form.Group>
                    </>
                );
            case 'transit':
                return (
                    <Form.Group className="form-floating mt-2">
                        <Form.Control
                            as="select"
                            value={metadata.key_source || ''}
                            onChange={(e) => setMetadata({ ...metadata, key_source: e.target.value })}
                        >
                            <option value="software">Software</option>
                            <option value="hsm">HSM</option>
                        </Form.Control>
                        <Form.Label>Key Source</Form.Label>
                    </Form.Group>
                );
            case 'pki':
                return (
                    <>
                    
                        <Form.Group className="form-floating mt-2">
                            <Form.Control
                                as="select"
                                value={metadata.root_ca || ''}
                                onChange={(e) => setMetadata({ ...metadata, root_ca: e.target.value })}
                            >
                                <option value="">Select Default CA</option>
                                {cas.map((ca) => (
                                    <option key={ca.id} value={ca.id}>
                                        {ca.name}
                                    </option>
                                ))}
                            </Form.Control>
                            <Form.Label>Select Certificate Authority (CA)</Form.Label>

                            <ul className='mt-2 small'>
                                <li>For each "path" (e.g., /path/something), we create a subordinate (sub) CA.</li>
                                <li>This sub-CA is signed by the root CA or an intermediate CA, creating a certificate hierarchy.</li>
                                <li>The sub-CA's private key and certificate are stored in the database</li>
                                <li>The idea is to isolate certificates under different paths, so each path essentially has its own CA, enhancing separation and scalability. </li>
                                <li><b>Separation of Certificates by Path:</b> By creating a sub-CA for each path, we are logically isolating certificates within different domains or organizational units. This can be useful for large organizations or multi-tenant environments where each path could represent a specific department, customer, or application.
                                </li>
                                <li><b>Scalability:</b> Using separate sub-CAs allows you to scale your PKI system horizontally. Each path-based sub-CA can issue certificates independently, reducing the load on the root CA.</li>
                                <li><b>Flexibility:</b> Different paths could have their own policies, expirations, and naming conventions, allowing for flexible certificate management per path.</li>
                                <li><b>Security:</b> By using sub-CAs instead of directly issuing all certificates from a root CA, we reduce the risk of root CA key compromise. Compromising one sub-CA only affects a single path rather than the entire hierarchy.</li>

                            </ul>
                            
                        </Form.Group>
                        <Form.Group className="form-floating mt-2">
                        <Form.Control
                            type="number"
                            placeholder="Max Lease Time (Days)"
                            value={metadata.max_lease_time || ''}
                            onChange={(e) => setMetadata({ ...metadata, max_lease_time: e.target.value })}
                        />
                        <Form.Label>Max Lease Time (Days)</Form.Label>
                    </Form.Group>

                    
                    </>
                );
            default:
                return (
                    <Form.Group className="form-floating mt-2">
                        <Form.Control
                            as="textarea"
                            rows={3}
                            placeholder="Metadata (JSON)"
                            value={JSON.stringify(metadata, null, 2)}
                            onChange={(e) => {
                                try {
                                    setMetadata(JSON.parse(e.target.value));
                                    setMessage('');
                                } catch (err) {
                                    setMessage('Invalid JSON format');
                                }
                            }}
                        />
                        <Form.Label>Metadata (JSON)</Form.Label>
                    </Form.Group>
                );
        }
    };

    return (
        <div className="">
            <Container className="p-0">
                <Row>

                    <Col >
                        {message && <div className='bg-success-soft text-success p-3 rounded-2 mb-2'>{message}</div>}
                        {error && <div className='bg-danger-soft text-danger p-3 rounded-2 mb-2'>{error}</div>}

                        

                        <Card className='mb-3'>
                            {/* <Card.Header>Edit Path?</Card.Header> */}
                            <Card.Body>
                                <Form.Group className='form-floating '>
                                    <Form.Control
                                        as="select"
                                        name="path"
                                        value={id}
                                        onChange={handlePathChange}
                                    >
                                        <option value="">Create New Path</option>
                                        {paths.map((path) => (
                                            <option key={path.id} value={path.id}>{path.path}</option>
                                        ))}
                                    </Form.Control>
                                    <Form.Label>Path</Form.Label>
                                </Form.Group>
                            </Card.Body>
                        </Card>
                        <Card className=''>
                            <Card.Header>{id == 0 ? "Create Path" : <div>Editing Path '{path}'</div>}</Card.Header>
                            <Card.Body>
                                <Form>

                                <Form.Group className="form-floating mb-3">
                                        <Form.Control
                                            as="select"
                                            disabled={id == 0 ? false : true}
                                            value={engineType}
                                            onChange={(e) => setEngineType(e.target.value)}
                                        >
                                            <option value="">---</option>
                                            <option value="kv">Key-Value</option>
                                            <option value="pki">PKI</option>
                                            <option value="transit">Transit</option>
                                            <option value="database">Database (Postgres Only)</option>
                                            {/* Add other engine types as needed */}
                                        </Form.Control>
                                        <Form.Label>Engine Type</Form.Label>
                                    </Form.Group>

                                    { engineType ? 
                                    <>
                                    <InputGroup className='mb-3'>
                                        <Form.Group className="form-floating">
                                            <Form.Control
                                                type="text"
                                                value={path}
                                                onChange={(e) => setPath(e.target.value)}
                                            />
                                            <Form.Label>Path</Form.Label>
                                        </Form.Group>
                                        {/* <Form.Group><div className="p-3"><FaSearch  onClick={(e) => fetchPath()} /></div></Form.Group> */}
                                    </InputGroup>
                                    


                                    {renderMetadataFields(engineType)}
                                    <Button className='w-100 mt-3' variant="primary" onClick={createPath}>
                                        {id == 0 ? "Create Path" : "Update Path"}
                                    </Button>
                                    </>: <></>}
                                </Form>
                            </Card.Body>
                        </Card>

                        <div className='mt-2 rounded-2 p-3 bg-dark text-white'>
                            <strong>Debug:</strong>
                            <pre>{JSON.stringify(metadata, null, 2)}</pre>
                        </div>
                    </Col>
                    
                </Row>
            </Container>
        </div>
    );
};

export default PathManagement;