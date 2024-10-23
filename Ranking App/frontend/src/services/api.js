import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080', // Your Golang backend URL
});

export const getAccounts = (page, limit, search, sort, order) => {
  return api.get('/accounts', {
    params: { page, limit, search, sort, order }
  });
};
