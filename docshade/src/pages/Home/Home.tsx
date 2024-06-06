import React, { useEffect } from 'react';
import UploadForm from '../../components/UploadForm/UploadForm';
import ServiceInfo from '../../components/ServiceInfo/ServiceInfo';
import { checkHealth } from '../../api/health';
import { Container } from './Home.styles';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const Home: React.FC = () => {
  useEffect(() => {
    const fetchHealthStatus = async () => {
      try {
        const response = await checkHealth();
        if (response.status !== 'OK') {
          toast.error('Service is not active');
        }
      } catch (error) {
        toast.error('Service is not active');
      }
    };

    fetchHealthStatus();
  }, []);

  return (
    <Container>
      <UploadForm />
      <ServiceInfo />
    </Container>
  );
};

export default Home;
