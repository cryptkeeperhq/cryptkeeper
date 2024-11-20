import React, { useState, useEffect } from 'react';
import { Alert, Card, InputGroup } from 'react-bootstrap';
import RenderMetadataFields from './RenderMetadataFields';
import { useApi } from '../api/api'


const SecretMetadataUpdate = ({ secret, path, token, onUpdate }) => {
  const { get, post, put, del } = useApi();
  const [metadataUpdate, setMetadataUpdate] = useState({});
  const [message, setMessage] = useState("")
  const [error, setError] = useState("")

  const [metadata, setMetadata] = useState(null);


  useEffect(() => {
    setMessage('')
    setError('')
    setMetadata(secret.metadata);
  }, [secret]);

    const updateMetadata = async (secretId, path, key, version) => {

      try {
        const data = await put(`/secrets/metadata?path=${encodeURIComponent(path)}&key=${encodeURIComponent(key)}&version=${version}`, metadata);
        setMessage('Metadata updated.');
        onUpdate()
      } catch (error) {
        setError(error.message);
      }
    }



  return (
    <div className='mt-2'>
      {secret && <div>

        <Card className='p-0 '>
          {/* <Card.Header className='bg-light text-dark'>Update Metadata</Card.Header> */}
          <Card.Body className='p-0'>
          {message && <p className='text-success fw-bold'>{message}</p>}
                    {error && <p className='text-danger fw-bold'>{error}</p>}

            <RenderMetadataFields
              metadata={secret.metadata}
              engineType={path.engine_type}
              setMetadata={setMetadata} />


            {/* <div className='p-2 bg-dark text-white mt-2 rounded-2'>
              <small>Metadata Preview</small>
              <pre>{JSON.stringify(metadata || {})}</pre>
            </div>
 */}
            <form>
              <div>
                {/* <div className="form-floating">
                            <input type="text" className="form-control" id={`updatePath-${secret.id}`} placeholder="Path" value={secret.path} readOnly />
                            <label htmlFor={`updatePath-${secret.id}`}>Path</label>
                          </div>
                          <div className="form-floating">
                            <input type="text" className="form-control" id={`updateVersion-${secret.id}`} placeholder="Version" value={secret.version} readOnly />
                            <label htmlFor={`updateVersion-${secret.id}`}>Version</label>
                          </div> */}
                {/* <div className="form-floating">
                            <textarea style={{ height: "100px"}} className="form-control" id={`updateMetadata-${secret.id}`} rows="3" placeholder="Metadata (JSON)" value={metadataUpdate[secret.id] || JSON.stringify(secret.metadata, null, 2)} onChange={(e) => handleMetadataChange(secret.id, e.target.value)} />
                            <label htmlFor={`updateMetadata-${secret.id}`}>Metadata (JSON)</label>
                          </div> */}
                <button type="button" className="mt-2 btn-sm w-100 btn btn-primary" onClick={() => updateMetadata(secret.id, path.path, secret.key, secret.version)}>Update Metadata</button>
              </div>
            </form>
          </Card.Body>
        </Card>

      </div>}

    </div>
  );
};

export default SecretMetadataUpdate;
