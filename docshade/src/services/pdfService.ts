import { uploadPDF } from '../api/pdf';

export const handlePDFUpload = async (file: File) => {
  try {
    const response = await uploadPDF(file);
    console.log('Upload successful:', response);
    return response;
  } catch (error) {
    console.error('Error uploading file:', error);
    throw error;
  }
};
