import React, { useState } from 'react';
import { useApi } from '../api/api'


const ScanSecrets = () => {
  const { get, post, put, del } = useApi();

  const [text, setText] = useState('');
  const [secrets, setSecrets] = useState([]);

  const scanForSecrets = async () => {
    try {
      const data = await post(`secrets/scan`, { text });
      setSecrets(data);
    } catch (error) {
      console.log(error.message);
    }
  };

  return (
    <div>
      <h2>Scan for Secrets</h2>
      <textarea
        placeholder="Enter text to scan for secrets"
        value={text}
        onChange={(e) => setText(e.target.value)}
        rows="10"
        cols="50"
      />
      <br />
      <button onClick={scanForSecrets}>Scan</button>
      <div>
        {secrets.length > 0 && (
          <div>
            <h3>Potential Secrets Found:</h3>
            <ul>
              {secrets.map((secret, index) => (
                <li key={index}>{secret}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};

export default ScanSecrets;
