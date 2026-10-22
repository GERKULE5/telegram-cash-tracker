import axios from 'axios';

export const api = axios.create({
  baseURL: '/api',
  timeout: 15_000,
});

api.interceptors.response.use(
  response => response,
  error => {
    if (
      axios.isAxiosError<{ error: string }>(error) &&
      error.response?.data?.error
    ) {
      error.message = error.response.data.error;
    }
    return Promise.reject(error);
  },
);
