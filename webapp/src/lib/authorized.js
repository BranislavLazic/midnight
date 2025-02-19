import axios from 'axios';

export const authorized = axios.create();

authorized.interceptors.request.use(
  (config) => {
    const accessToken = sessionStorage.getItem('accessToken');
    config.headers = {
      Authorization: `Bearer ${accessToken}`
    };
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

authorized.interceptors.response.use(
  async (response) => {
    return response;
  },
  async (error) => {
    if (error.response.status === 403 || error.response.status === 401) {
      sessionStorage.clear();
      window.location.pathname = '/login';
    }
    return Promise.reject(error);
  }
);
