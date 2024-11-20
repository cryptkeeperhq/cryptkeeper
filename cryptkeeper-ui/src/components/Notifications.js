import React, { useState, useEffect, useRef } from 'react';
import { OverlayTrigger, Popover, Badge, ListGroup, ListGroupItem } from 'react-bootstrap';
import { FaBell } from 'react-icons/fa';
import moment from 'moment';
import { useApi } from '../api/api'

        
const Notifications = () => {
  const { get, post, put, del } = useApi();
  const [notifications, setNotifications] = useState([]);

  const hasFetchedNotifications = useRef(false);
  useEffect(() => {
      if (!hasFetchedNotifications.current) {
        hasFetchedNotifications.current = true;
        fetchNotifications();
      }
  }, []);

  const fetchNotifications = async () => {
      try {
        const data = await get(`/notifications`);
        setNotifications(data || []);
    } catch (error) {
        console.log(error.message);
    }
  }
  const notificationCount = notifications.length;

  const popover = (
    <Popover id="notification-popover">
      {/* <Popover.Header as="h2" className='bg-transparent border-0'>Notifications</Popover.Header> */}
      <Popover.Body className='p-1'>
        {notifications.length > 0 ? (
          <ListGroup variant='flush'>
            {notifications.map((notification) => (
              <ListGroupItem className='p-1 ps-2 small d-flex align-items-center'  key={notification.id}>
                {/* <FaBell size={14} className="me-1" /> */}
                <div className='ms-1'>
                  <strong>{notification.message}</strong><br/>
                {moment(notification.created_at).format("MMMM Do YYYY, h:mm:ss a")}
                </div>
              </ListGroupItem>
            ))}
          </ListGroup>
        ) : (
          <p>No notifications</p>
        )}
      </Popover.Body>
    </Popover>
  );

  return (
    <OverlayTrigger trigger="click" placement="bottom" overlay={popover}>
      <div style={{ cursor: 'pointer' }}>
        {notificationCount > 0 && <Badge bg="danger fs-7"><FaBell size={10} className="text-white me-1" /> {notificationCount}</Badge>}
      </div>
    </OverlayTrigger>
  );
};

export default Notifications;
