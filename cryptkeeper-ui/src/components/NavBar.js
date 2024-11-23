import React, { useState, useEffect, useRef } from 'react';
import { Navbar, Nav, Container, Card, Accordion } from 'react-bootstrap';
import { BrowserRouter as Router, Route, Routes, NavLink, Navigate } from 'react-router-dom';
import { useTheme } from '../theme/ThemeContext';
import HelpCard from './HelpCard';
import logo from '../assets/logo.webp';
import { FaKey, FaFolder, FaFileAlt, FaHistory, FaTrashAlt, FaEdit, FaUser, FaPlus, FaHome, FaCheck, FaShieldAlt, FaList, FaSignOutAlt, FaCog, FaCogs, FaCheckCircle, FaUsers, FaUsersCog, FaKeybase, FaEnvelope, FaListAlt, FaPlusCircle, FaNetworkWired, FaDochub, FaFileWord, FaServer, FaCertificate, FaBars, FaLock, FaChevronLeft, FaDashcube, FaDatabase, FaBriefcase } from 'react-icons/fa';
import { Link } from 'react-router-dom';
import Notifications from './Notifications';

const NavBar = ({ token, logout }) => {
    const { theme } = useTheme();
    const [sidebarVisible, setSidebarVisible] = useState(true);
    const sidebarClass = 'close';
    const toggleSidebar = () => {
        const body = document.body;
        if (sidebarVisible) {
            body.classList.add(sidebarClass);
        } else {
            body.classList.remove(sidebarClass);
        }
        setSidebarVisible(!sidebarVisible);
    };
    const username = localStorage.getItem("user");
    const name = localStorage.getItem("name");

    const handleLogout = () => {
        logout()
        // localStorage.removeItem('token');
    };

    const [adminMenuOpen, setAdminMenuOpen] = useState(false);

    const toggleAdminMenu = (e) => {
        e.preventDefault();
        setAdminMenuOpen(!adminMenuOpen);
    };

    return (
        <>
            {/* <div className=''>
                <button type="button" id="sidebar_menu_icon"
                    onClick={() => {
                        toggleSidebar();
                    }}
                    className="p-2 ms-1 ps-3 border-0 shadow-none  toggle-button  btn-sm   btn btn-transparent"><FaBars /></button>
            </div> */}
            <Navbar.Brand className=''>
                <div className='mt-4 ms-2 d-flex align-items-center '>
                    <NavLink to="/" className="">
                        <div className='text-center' style={{ width: "50px" }}><img className='rounded-2 w-75' src={logo} /></div>
                    </NavLink>
                    <div className='mt-1 w-100 me-auto'>
                        <div className='d-flex w-100 align-items-start'>
                            <div className=''>
                                <div style={{ lineHeight: "0.6em" }}>Welcome</div>
                                <div className='mt-1'><b>{name}</b> <Link className='ms-1 p-0 ' to="/" onClick={handleLogout}><FaSignOutAlt /></Link></div>
                            </div>
                            <div className='ms-auto me-2'><Notifications /></div>
                        </div>
                    </div>
                </div>
            </Navbar.Brand>
            {/* <Navbar.Collapse id="basic-navbar-nav"> */}
            <ul className="nav flex-column mt-3">
                {token && <>
                    {/* <NavLink className="nav-link" to="/" end><FaHome className='me-2' />Home</NavLink> */}

                    <div className="">

                        {adminMenuOpen ? (
                            <>
                                {/* <NavLink className="nav-link" to="/" ><FaChevronLeft className='me-2' /><span>Back</span></NavLink> */}
                                <div className='p-2 ps-3 sidenav-menu-heading' onClick={toggleAdminMenu}><FaChevronLeft className='me-2' /> <strong>Back</strong></div>

                            </>
                        ) :
                            <>

                                <Accordion defaultActiveKey="0" flush>
                                <Accordion.Item eventKey="0">
                                        <Accordion.Header>Home</Accordion.Header>
                                        <Accordion.Body className='p-0'>
                                        <NavLink className="nav-link" to="/"><FaBriefcase className='me-2' /><span>Dashboard</span></NavLink>
                                        <NavLink className="nav-link" to="/audit-logs"><FaListAlt className='me-2' /><span>Audit Logs</span></NavLink>



                                        </Accordion.Body>
                                    </Accordion.Item>
                                    <Accordion.Item eventKey="1">
                                        <Accordion.Header>Secrets Management</Accordion.Header>
                                        <Accordion.Body className='p-0'>
                                            <NavLink className="nav-link" to="/user/secrets/kv"><FaList className='me-2' /><span>KV</span></NavLink>
                                            <NavLink className="nav-link" to="/user/secrets/transit"><FaKey className='me-2' /><span>Transit</span></NavLink>
                                            <NavLink className="nav-link" to="/user/secrets/pki"><FaCertificate className='me-2' /><span>PKI</span></NavLink>
                                            <NavLink className="nav-link" to="/user/secrets/database"><FaDatabase className='me-2' /><span>Database</span></NavLink>


                                        </Accordion.Body>
                                    </Accordion.Item>
                                    <Accordion.Item eventKey="2">
                                        <Accordion.Header>System</Accordion.Header>
                                        <Accordion.Body className='p-0'>
                                            <NavLink className="nav-link" to="/secrets/create"><FaPlusCircle className='me-2' /><span>Create Secret</span></NavLink>
                                            <NavLink className="nav-link" to="/approval-requests"><FaCheckCircle className='me-2' /><span>Approval Requests</span></NavLink>
                                            <NavLink className="nav-link" to="/paths"><FaNetworkWired className='me-2' /><span>Paths</span></NavLink>
                                            <NavLink className="nav-link" to="/policies"><FaFileAlt className='me-2' /><span>Policies</span></NavLink>
                                            <NavLink className="nav-link" to="/pki"><FaCertificate className='me-2' /><span>Create PKI CA</span></NavLink>

                                            {/* <NavLink className="nav-link" to="/workflows"><FaNetworkWired className='me-2' /><span>Workflows</span></NavLink> */}

                                        </Accordion.Body>
                                    </Accordion.Item>

                                    <Accordion.Item eventKey="3">
                                        <Accordion.Header>Auth</Accordion.Header>
                                        <Accordion.Body className='p-0'>
                                            <NavLink className="nav-link" to="/management"><FaUsersCog className='me-2' /><span>Users and Groups</span></NavLink>
                                            <NavLink className="nav-link" to="/roles"><FaServer className='me-2' /><span>App Roles</span></NavLink>
                                            <NavLink className="nav-link" to="/certificates"><FaCertificate className='me-2' /><span>Certificates</span></NavLink>

                                        </Accordion.Body>
                                    </Accordion.Item>
                                </Accordion>



                                {/* <div className='p-2 ps-3 sidenav-menu-heading'><strong>Secrets Management</strong></div>
                                <NavLink className="nav-link" to="/user/secrets/kv"><FaList className='me-2' /><span>KV</span></NavLink>
                                <NavLink className="nav-link" to="/user/secrets/transit"><FaKey className='me-2' /><span>Transit</span></NavLink>
                                <NavLink className="nav-link" to="/user/secrets/pki"><FaCertificate className='me-2' /><span>PKI</span></NavLink>
                                <NavLink className="nav-link" to="/user/secrets/database"><FaDatabase className='me-2' /><span>Database</span></NavLink>

                                <div className='p-2 ps-3 sidenav-menu-heading'><strong>System</strong></div>
                                <NavLink className="nav-link" to="/secrets/create"><FaPlusCircle className='me-2' /><span>Create Secret</span></NavLink>
                                <NavLink className="nav-link" to="/approval-requests"><FaCheckCircle className='me-2' /><span>Approval Requests</span></NavLink>
                                <NavLink className="nav-link" to="/paths"><FaNetworkWired className='me-2' /><span>Paths</span></NavLink>
                                <NavLink className="nav-link" to="/policies"><FaFileAlt className='me-2' /><span>Policies</span></NavLink>

                                <div className='p-2 ps-3 sidenav-menu-heading'><strong>Manage</strong></div>
                                <NavLink className="nav-link" to="/pki"><FaCertificate className='me-2' /><span>Create PKI CA</span></NavLink>

                                <NavLink className="nav-link" to="/audit-logs"><FaListAlt className='me-2' /><span>Audit Logs</span></NavLink>
                                {/* <NavLink className="nav-link" to="/workflows"><FaNetworkWired className='me-2' /><span>Workflows</span></NavLink> */}

                                {/* <div className='p-2 ps-3 sidenav-menu-heading'><strong>Auth</strong></div>

                                <NavLink className="nav-link" to="/management"><FaUsersCog className='me-2' /><span>Users and Groups</span></NavLink>
                                <NavLink className="nav-link" to="/roles"><FaServer className='me-2' /><span>App Roles</span></NavLink>
                                <NavLink className="nav-link" to="/certificates"><FaCertificate className='me-2' /><span>Certificates</span></NavLink>  */}



                                {/* <div className='p-2 ps-3 sidenav-menu-heading'><strong>Test</strong></div>
                                <NavLink className='nav-link' to="/transit/encryption"><FaEnvelope className='me-2' /><span>Transit Encryption</span></NavLink> */}


                                {/* <div className='p-2 ps-3 sidenav-menu-heading' onClick={toggleAdminMenu}><strong>Admin / Management</strong></div> */}

                                {/* <NavLink className="nav-link" to="/admin" onClick={toggleAdminMenu}><FaLock className='me-2' /><span>Admin</span></NavLink> */}
                            </>
                        }
                    </div>
                    {token ? (
                        <>
                        </>
                    ) : (
                        <NavLink className="nav-link" to="/login">Login</NavLink>
                    )}
                </>}
            </ul>
            <HelpCard />
        </>
    );
};

export default NavBar;