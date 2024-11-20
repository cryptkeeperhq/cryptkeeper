import React, { useState, useEffect } from 'react';
import { ListGroup, ListGroupItem } from 'react-bootstrap';

const Permissions = () => {

    const permissions = JSON.parse(localStorage.getItem("permissions") || [])

  return (
    <div>
      {/* <h2>My Permissions</h2>
      {permissions.length > 0 ? (
        <ListGroup>
          {permissions.map((permission) => (
            <ListGroupItem key={permission.path}>{permission.path} <pre className='m-0'>{permission.permission}</pre></ListGroupItem>
          ))}
        </ListGroup>
      ) : (
        <p>No Permisssions</p>
      )} */}
    </div>
  );
};

export default Permissions;
