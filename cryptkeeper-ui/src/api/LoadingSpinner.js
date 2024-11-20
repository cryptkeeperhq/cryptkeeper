// LoadingSpinner.js
import React from 'react';
import { useLoading } from './LoadingContext';
import Spinner from 'react-bootstrap/Spinner';

const LoadingSpinner = () => {
    const { isLoading } = useLoading();

    if (!isLoading) return null;

    return (
        <div style={spinnerStyle}>
            <Spinner animation="border" />
        </div>
    );
};

const spinnerStyle = {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    zIndex: 9999,
};

export default LoadingSpinner;