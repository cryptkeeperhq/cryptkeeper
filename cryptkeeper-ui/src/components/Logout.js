import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const Logout = ({ setToken }) => {
  const navigate = useNavigate();

  useEffect(() => {
    // Remove the token
    localStorage.removeItem('token');
    setToken(null);

    // Redirect to the login page
    navigate('/login');
  }, [navigate, setToken]);

  return null;
};

export default Logout;
