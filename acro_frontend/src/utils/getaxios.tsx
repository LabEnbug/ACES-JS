import axios from 'axios';
import {getToken} from "@/utils/authentication";

const baxios = axios.create();

baxios.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

export default baxios;
