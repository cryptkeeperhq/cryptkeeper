import React, { useState, useEffect, useRef } from 'react';
import { Form, Button, Card, Container, ListGroup, ListGroupItem, Row, Col, InputGroup, ButtonGroup, Table, Tab, Tabs } from 'react-bootstrap';
import { FaPlus, FaPlusSquare, FaTrash, FaUser, FaUsers } from 'react-icons/fa';
import Title from '../common/Title';
import { UserManagementHelp } from '../help/Help';
import { useApi } from '../../api/api'
import Register from '../Register';

const UserManagementPage = ({ setTitle, setHelp }) => {
    const { get, post, put, del } = useApi();

    const [users, setUsers] = useState([]);
    const [groups, setGroups] = useState([]);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [groupName, setGroupName] = useState('');
    const [selectedUserId, setSelectedUserId] = useState('');
    const [selectedGroupId, setSelectedGroupId] = useState('');
    const [userGroupsMap, setUserGroupsMap] = useState({});
    const [groupUsersMap, setGroupUsersMap] = useState({});
    const [message, setMessage] = useState('');

    const hasFetchDone = useRef(false);
    useEffect(() => {
        if (!hasFetchDone.current) {
            hasFetchDone.current = true;
            fetchUsers()
            fetchGroups()
        }
    }, []);

    useEffect(() => {
        setTitle({ heading: "User and Group Management", subheading: "Manage users and groups for CryptKeeper" })
        setHelp(
            <>
                <UserManagementHelp />
            </>
        )
    }, []);






    const fetchGroups = async () => {
        try {
            const data = await get(`/groups`);
            setGroups(data || []);
        } catch (error) {
            console.log(error.message);
        }
    };

    const fetchUsers = async () => {
        try {
            const data = await get(`/users`);
            setUsers(data || []);
        } catch (error) {
            console.log(error.message);
        }
    };


    const fetchUserGroups = async (userId) => {
        try {
            const data = await get(`/users/${userId}/groups`);
            setUserGroupsMap(prev => ({ ...prev, [userId]: data }));
        } catch (error) {
            console.log(error.message);
        }
    };

    const fetchGroupUsers = async (groupId) => {
        try {
            const data = await get(`/groups/${groupId}/users`);
            setGroupUsersMap(prev => ({ ...prev, [groupId]: data }));
        } catch (error) {
            console.log(error.message);
        }
    };

    const createUser = async () => {
        const userData = { username, password };

        try {
            const data = await post(`/users`, userData);
            setUsername('');
            setPassword('');
            setMessage('User created.');
            fetchUsers()
        } catch (error) {
            console.log(error.message);
        }

    };

    const createGroup = async () => {
        const groupData = { name: groupName };

        try {
            const data = await post(`/groups`, groupData);
            setGroupName('');
            setMessage('Group created.');
            fetchGroups()
        } catch (error) {
            console.log(error.message);
        }


    };

    const addUserToGroup = async () => {
        const userGroupData = {
            user_id: selectedUserId,
            group_id: selectedGroupId,
        };

        try {
            const data = await post(`/groups/add_user`, userGroupData);
            setSelectedUserId('');
            setSelectedGroupId('');
            setMessage('User added to group.');
        } catch (error) {
            console.log(error.message);
        }

    };

    const removeUserFromGroup = async () => {
        const userGroupData = {
            user_id: selectedUserId,
            group_id: selectedGroupId,
        };

        try {
            const data = await post(`/groups/remove_user`, userGroupData);
            setSelectedUserId('');
            setSelectedGroupId('');
            setMessage('User removed to group.');
        } catch (error) {
            console.log(error.message);
        }
    };


    return (
        <div className="">


            <Container className="p-0">
                <Row>

                    <Col className='mb-3'>
                        {message && <div className='bg-info text-white p-2 rounded-2 mb-2'>{message}</div>}



                        <Tabs
                            defaultActiveKey="users"
                            id="user-mgt"
                            variant='underline'
                            className="mb-3 ps-3"
                        >
                            <Tab eventKey="users" title="Users">


                       
                            <Card className='mt-3'>
                            <Card.Header>Existing Users</Card.Header>
                            <Card.Body>
                                <Table striped bordered hover variant='flush'>
                                    <thead>
                                        <tr>
                                            <th>#</th>
                                            <th>Name</th>
                                            <th>Username</th>
                                            <th>ID</th>
                                            <th></th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {users.map(user => (
                                            <tr key={user.id} >
                                                <td><FaUser className='me-2' /> </td>
                                                <td>{user.name}</td>
                                                <td>{user.username}</td>
                                                <td>{user.id}</td>
                                                <td>
                                                    <div className='link text-decoration-underline text-primary' onClick={() => fetchUserGroups(user.id)}>Fetch Groups</div>
                                                    {userGroupsMap[user.id] && <div>
                                                        {
                                                            userGroupsMap[user.id] != [] &&
                                                            <ul>
                                                                {(userGroupsMap[user.id] || []).map(group => (
                                                                    <li key={group.id}>{group.name}</li>
                                                                ))}
                                                            </ul>
                                                        }
                                                    </div>}
                                                </td>



                                            </tr>
                                        ))}
                                    </tbody>
                                </Table>

                            </Card.Body>

                        </Card>

                        <Card className='mt-3'>
                            <Card.Header>New User</Card.Header>
                            <Card.Body>
                                <Register />
                            </Card.Body>
                        </Card>

                       
                       

                            </Tab>
                            <Tab eventKey="groups" title="Groups">
                                
                            <Card className='mt-3'>
                            <Card.Header>Existing Groups</Card.Header>
                            <Card.Body>
                            
                                <ListGroup>
                                    {groups.map(group => (
                                        <ListGroupItem key={group.id} className="list-group-item" onClick={() => fetchGroupUsers(group.id)}>
                                            <FaUsers className='me-2' /> {group.name} (ID: {group.id})
                                            {
                                                groupUsersMap[group.id] && <div>
                                                    <ul>
                                                        {(groupUsersMap[group.id] || []).map(user => (
                                                            <li key={user.id}>{user.username}</li>
                                                        ))}
                                                    </ul>
                                                </div>
                                            }

                                        </ListGroupItem>
                                    ))}
                                </ListGroup>
                                </Card.Body>
                        </Card>

                            <Card className='mt-3'>
                            <Card.Header>New Group</Card.Header>
                            <Card.Body>
                                <form>
                                    <InputGroup>
                                        <div className="form-floating">
                                            <input type="text" className="form-control" id="createGroupName" placeholder="Group Name" value={groupName} onChange={e => setGroupName(e.target.value)} />
                                            <label htmlFor="createGroupName">Group Name</label>
                                        </div>
                                        <button type="button" className="btn btn-primary" onClick={createGroup}>Create Group</button>
                                    </InputGroup>
                                </form>
                            </Card.Body>
                            </Card>

                          



                        <Card className='mt-3'>
                            <Card.Header>Add User to Group</Card.Header>
                            <Card.Body>
                                <p>This allows user to access all <code>paths</code> which are assigned to the group</p>
                                <form>
                                    <InputGroup>
                                        <div className="form-floating">
                                            <select className="form-control" id="selectUser" value={selectedUserId} onChange={e => setSelectedUserId(e.target.value)}>
                                                <option value="">Select a user</option>
                                                {users.map(user => (
                                                    <option key={user.id} value={user.id}>{user.username}</option>
                                                ))}
                                            </select>
                                            <label htmlFor="selectUser">Select User</label>

                                        </div>
                                        <div className="form-floating">
                                            <select className="form-control" id="selectGroup" value={selectedGroupId} onChange={e => setSelectedGroupId(e.target.value)}>
                                                <option value="">Select a group</option>
                                                {groups.map(group => (
                                                    <option key={group.id} value={group.id}>{group.name}</option>
                                                ))}
                                            </select>
                                            <label htmlFor="selectGroup">Select Group</label>

                                        </div>
                                        {/* <Form.Group className='form-floating'>
                                                <Form.Control as="select" value={selectedRoleId} onChange={(e) => setSelectedRoleId(e.target.value)}>
                                                    <option value="">Select Role</option>
                                                    {roles.map(role => (
                                                        <option key={role.id} value={role.id}>{role.name}</option>
                                                    ))}
                                                </Form.Control>
                                                <Form.Label>Role</Form.Label>

                                            </Form.Group> */}


                                    </InputGroup>
                                    <div className='mt-2 d-flex justify-content-between align-items-center'>
                                        <button type="button" className="w-100 btn btn-primary me-3" onClick={addUserToGroup}><FaPlus /> Add</button>
                                        <button type="button" className="ms-auto btn btn-danger" onClick={removeUserFromGroup}><FaTrash /></button>

                                    </div>
                                </form>
                            </Card.Body>

                        </Card>
                            </Tab>
                        </Tabs>




                       

                    </Col>

                </Row>



            </Container>





        </div >
    );
};

export default UserManagementPage;
