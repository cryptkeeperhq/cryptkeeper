import React from 'react';
import { Button, Nav } from 'react-bootstrap';
import { useTheme } from './ThemeContext';
import { FaMoon } from 'react-icons/fa';

const ThemeToggle = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button variant='transparent' className='bg-transparent w-100' onClick={toggleTheme}>
      <FaMoon size={20} className='pe-2' /> {theme === 'light' ? 'Dark' : 'Light'}
    </Button>
  );
};

export default ThemeToggle;
