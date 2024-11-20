import React, { useState, useEffect } from 'react';
import { Card, InputGroup } from 'react-bootstrap';
import { useApi } from '../api/api'


const SecretAccessControl = ({ token, secret }) => {
  const { get, post, put, del } = useApi();
  const [userId, setUserId] = useState('');
  const [groupId, setGroupId] = useState('');
  const [accessLevel, setAccessLevel] = useState('');
  const [message, setMessage] = useState('');
  const [users, setUsers] = useState([]);
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const data = await get(`/users`);
        setUsers(data);
      } catch (error) {
        console.log(error.message);
      }
    };

    const fetchGroups = async () => {
      try {
        const data = await get(`/groups`);
        setGroups(data);
      } catch (error) {
        console.log(error.message);
      }
    };

    fetchUsers();
    fetchGroups();
  }, []);

  const assignAccess = async () => {
    const access = {
      secret_id: secret.id,
      access_level: accessLevel,
    };
    if (userId) {
      access.user_id = userId;
    }
    if (groupId) {
      access.group_id = groupId;
    }

    try {
      const data = await post(`/paths/${secret.path_id}/secrets/assign_access?path=${encodeURIComponent(secret.path)}`, access);
      setMessage('Access assigned successfully.');
    } catch (error) {
      setMessage(error.message);
    }
  };

  return (
    <div className='mt-2'>
      <Card className='p-0 '>
        {/* <Card.Header className='bg-light text-dark'>Assign Access Control</Card.Header> */}
        <Card.Body>
          {message && <p className='text-warning fw-bold'>{message}</p>}
          <InputGroup>
            <div className="form-floating">
              <select
                className="form-control"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
              >
                <option value="">Select User</option>
                {users.map(user => (
                  <option key={user.id} value={user.id}>
                    {user.username}
                  </option>
                ))}
              </select>
              <label>User ID:</label>
            </div>
            <div className="form-floating">
              <select
                className="form-control"
                value={groupId}
                onChange={(e) => setGroupId(e.target.value)}
              >
                <option value="">Select Group</option>
                {groups.map(group => (
                  <option key={group.id} value={group.id}>
                    {group.name}
                  </option>
                ))}
              </select>
              <label>Group ID:</label>
            </div>
            <div className="form-floating">
              <select
                className="form-control"
                value={accessLevel}
                onChange={(e) => setAccessLevel(e.target.value)}
              >
                <option value="">Select Access Level</option>
                <option value="read">Read</option>
                <option value="write">Write</option>
                <option value="owner">Owner</option>
              </select>
              <label>Access Level:</label>
            </div>
          </InputGroup>
          <button className="btn btn-sm btn-primary w-100" onClick={assignAccess}>Assign Access</button>
        </Card.Body></Card>
    </div>
  );
};

export default SecretAccessControl;
