import axios from 'axios';
import config from '../config';

export const checkHealth = async () => {
  try {
    const response = await axios.get(`https://${config.backendHost}/v1/health`);
    return { status: 'OK', data: response.data };
  } catch (error) {
    console.error('Error checking health:', error);
    return { status: 'NOT_OK', data: null };
  }
};
