import React, { useState, useEffect } from 'react';
import { useApi } from '../../api/api';
import { Card } from 'react-bootstrap';

function AddCertificateRequest() {
    const { get, post, downloadPost } = useApi();

    const [commonName, setCommonName] = useState('');
    const [organization, setOrganization] = useState('');
    const [validityPeriod, setValidityPeriod] = useState('');
    const [caCertID, setCACertID] = useState('');
    const [templateID, setTemplateID] = useState('');
    const [generatedCert, setGeneratedCert] = useState(null);

    const getHeaders = () => {
        const savedToken = localStorage.getItem('token');

        return {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${savedToken}`,
        };
    };

    const [cas, setCas] = useState([]);
    const [templates, setTemplates] = useState([]);

    async function fetchCas() {
        const data = await get('/pki/ca'); // Replace with your actual API endpoint
        setCas(data || []);
    }

    async function fetchTemplates() {
        const data = await get('/pki/template'); // Replace with your actual API endpoint
        setTemplates(data || []);
    }


    useEffect(() => {
        fetchCas()
        fetchTemplates()
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const payload = {
                common_name: commonName,
                // organization,
                // validity_period: parseInt(validityPeriod, 10),
                ca_cert_id: parseInt(caCertID, 10),
                ca_template_id: parseInt(templateID, 10),
            };

            const response = await downloadPost(`/pki/request_certificate`, payload);

            // const data = await post('/pki/request_certificate', );

            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = 'certificate.p12';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            // console.log(data)

            // setGeneratedCert(data);

            // alert('Certificate Request added successfully!');

            // const blob = await data.blob();
            //         const url = URL.createObjectURL(blob);
            //         const link = document.createElement('a');
            //         link.href = url;
            //         link.download = 'certificate.pem';
            //         document.body.appendChild(link);
            //         link.click();
            //         document.body.removeChild(link);

        } catch (error) {
            alert('Error adding Certificate Request');
        }
    };

    return (
        <Card className='mb-3'>
            <Card.Header>Request New Certificate</Card.Header>
            <Card.Body>
                <form onSubmit={handleSubmit}>
                    <div className='form-floating mb-2'>
                        <input
                            type="text"
                            className='form-control'
                            value={commonName}
                            onChange={(e) => setCommonName(e.target.value)}
                            required
                        />
                        <label>Common Name:</label>
                    </div>

                    {/* <div className='form-floating mb-2'>
                        <input
                            type="text"
                            className='form-control'
                            value={organization}
                            onChange={(e) => setOrganization(e.target.value)}
                        />
                        <label>Organization:</label>
                    </div> */}

                    {/* <div className='form-floating mb-2'>
                        <input
                            type="number"
                            className='form-control'
                            value={validityPeriod}
                            onChange={(e) => setValidityPeriod(e.target.value)}
                            required
                        />
                        <label>Validity Period (days):</label>
                    </div> */}

                    <div className='form-floating mb-2'>

                        <select className="form-control" id="selectUser" value={caCertID} onChange={e => setCACertID(e.target.value)}>
                            <option value="">Select a CA</option>
                            {cas.map(ca => (
                                <option key={ca.id} value={ca.id}>{ca.name}</option>
                            ))}
                        </select>

                        {/* <input
                            type="number"
                            className='form-control'
                            value={caCertID}
                            onChange={(e) => setCACertID(e.target.value)}
                            required
                        /> */}
                        <label>CA Certificate ID:</label>
                    </div>

                    <div className='form-floating mb-2'>
                        <select className="form-control" id="selectTemplate" value={templateID} onChange={e => setTemplateID(e.target.value)}>
                            <option value="">Select a Template</option>
                            {templates.map(template => (
                                <option key={template.id} value={template.id}>{template.name} ({template.organization} / {template.validity_period})</option>
                            ))}
                        </select>

                        {/* <input
                            type="number"
                            className='form-control'
                            value={templateID}
                            onChange={(e) => setTemplateID(e.target.value)}
                            required
                        /> */}
                        <label>Template ID:</label>
                    </div>

                    <button className='btn btn-primary w-100' type="submit">Request Certificate</button>
                </form>


                <div className='mt-3 alert alert-info'>
                    Once the file is downloaded, you can verify the p12 file using below command.<br />
                    <pre className='bg-dark text-white rounded-3 p-1 ps-2'>
                        openssl pkcs12 -in certificate.p12 -clcerts -nodes -passin pass:"password"
                    </pre>
                </div>
            </Card.Body>
            <Card.Footer>
                {generatedCert && (
                    <div>
                        <h3>Download Certificate and Private Key</h3>
                        <a href={`/api/pki/download_certificate?cert=${encodeURIComponent(generatedCert.certificate)}`} className='btn btn-success' download>
                            Download Certificate
                        </a>
                        <a href={`/api/pki/download_private_key?key=${encodeURIComponent(generatedCert.private_key)}`} className='btn btn-success' download>
                            Download Private Key
                        </a>
                    </div>
                )}

            </Card.Footer>
        </Card>
    );
}

export default AddCertificateRequest;
