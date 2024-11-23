import { useLoading } from './LoadingContext';

const API_URL = process.env.REACT_APP_API_URL || '/api';

const handleResponse = async (response) => {

    const resp = await response
    const status = resp.status

    if (status == 401) {
        console.log(resp)
        // throw new Error(resp.text())
        localStorage.removeItem('token');
        setToken(null);
        return
    }

    if (status == 201 || status == 203) {
        return {};
    }

    const data = resp.json();
    if (!response.ok) {
        const error = (data && data.message) || response.statusText;
        throw new Error(error);
    }
    return data;
};

const getHeaders = () => {
    const savedToken = localStorage.getItem('token');

    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${savedToken}`,
    };
};

export const useApi = () => {

    // const { setIsLoading } = useLoading();

    const get = async (endpoint) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'GET',
                headers: getHeaders(),
            });
            return handleResponse(response);
        } finally {
            // setIsLoading(false);
        }
    };

    const post = async (endpoint, body) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'POST',
                headers: getHeaders(),
                body: JSON.stringify(body),
            });
            return handleResponse(response);
        } finally {
            // setIsLoading(false);
        }
    };

    const downloadPost = async (endpoint, body) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'POST',
                headers: getHeaders(),
                body: JSON.stringify(body),
            });
            return response;
        } finally {
            // setIsLoading(false);
        }
    };   
    
    
    const downloadGet = async (endpoint) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'GET',
                headers: getHeaders(),
            });
            return response;
        } finally {
            // setIsLoading(false);
        }
    };   

    const put = async (endpoint, body) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'PUT',
                headers: getHeaders(),
                body: JSON.stringify(body),
            });
            return handleResponse(response);
        } finally {
            // setIsLoading(false);
        }
    };

    const del = async (endpoint) => {
        // setIsLoading(true);
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                method: 'DELETE',
                headers: getHeaders(),
            });
            return handleResponse(response);
        } finally {
            // setIsLoading(false);
        }
    };

    return { get, post, put, del, downloadPost, downloadGet };
};
// export const get = async (endpoint) => {
//     const response = await fetch(`${API_URL}${endpoint}`, {
//         method: 'GET',
//         headers: getHeaders(),
//     });
//     return handleResponse(response);
// };

// export const post = async (endpoint, body) => {
//     const response = await fetch(`${API_URL}${endpoint}`, {
//         method: 'POST',
//         headers: getHeaders(),
//         body: JSON.stringify(body),
//     });
//     return handleResponse(response);
// };

// export const put = async (endpoint, body) => {
//     const response = await fetch(`${API_URL}${endpoint}`, {
//         method: 'PUT',
//         headers: getHeaders(),
//         body: JSON.stringify(body),
//     });
//     return handleResponse(response);
// };

// export const del = async (endpoint) => {
//     const response = await fetch(`${API_URL}${endpoint}`, {
//         method: 'DELETE',
//         headers: getHeaders(),
//     });
//     return handleResponse(response);
// };



 // const fetchGroups = async () => {
    //     const response = await axios.get('/api/groups', {
    //         headers: { Authorization: `Bearer ${token}` },
    //     });
    //     setGroups(response.data || []);
    // };

    // const fetchUsers = async () => {
    //     const response = await axios.get('/api/users', {
    //         headers: { Authorization: `Bearer ${token}` },
    //     });
    //     setUsers(response.data || []);
    // };

    // const fetchAppRoles = async () => {
    //     const response = await axios.get('/api/approles', {
    //         headers: { Authorization: `Bearer ${token}` },
    //     });
    //     setAppRoles(response.data || []);
    // };

