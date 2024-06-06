import styled, { keyframes } from 'styled-components';

export const Form = styled.form`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border-radius: 8px;
  max-width: 600px;
  margin: 20px auto;
  background: #f9f9f9;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);

  @media (max-width: 768px) {
    padding: 10px;
  }
`;

export const Dropzone = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  border: 2px dashed #ccc;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  transition: border 0.2s ease-in-out;

  &:hover {
    border-color: #007bff;
  }

  @media (max-width: 768px) {
    padding: 20px;
  }
`;

export const DropzoneText = styled.p`
  font-size: 16px;
  color: #888;
  text-align: center;

  @media (max-width: 768px) {
    font-size: 14px;
  }
`;

export const DropzoneLink = styled.span`
  color: #007bff;
  cursor: pointer;
  text-decoration: underline;

  &:hover {
    color: #0056b3;
  }
`;

export const Input = styled.input`
  display: none;
`;

export const ErrorMessage = styled.div`
  color: red;
  margin: 5px 0;
  font-size: 14px;
`;

export const Title = styled.h1`
  font-size: 24px;
  margin-bottom: 20px;

  @media (max-width: 768px) {
    font-size: 20px;
  }
`;

export const UploadButton = styled.button`
  background-color: #007bff;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 20px;
  font-size: 16px;
  width: 200px;

  &:hover {
    background-color: #0056b3;
  }

  &:disabled {
    background-color: #007bff;
    opacity: 0.6;
    cursor: not-allowed;
  }

  @media (max-width: 768px) {
    width: 100%;
    padding: 10px;
  }
`;

const spin = keyframes`
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
`;

export const Spinner = styled.div`
  display: inline-block;
  animation: ${spin} 1s linear infinite;
  font-size: 24px;
  margin-top: 20px;
`;

export const DocumentList = styled.ul`
  list-style-type: none;
  padding: 0;
  margin: 20px 0;
  width: 100%;
`;

export const DocumentItem = styled.li`
  margin: 10px 0;
  padding: 10px;
  background: #f4f4f4;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  a {
    text-decoration: none;
    color: #007bff;

    &:hover {
      text-decoration: underline;
    }
  }

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
  }
`;
