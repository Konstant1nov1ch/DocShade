import React, { useState, useEffect } from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { useDropzone } from 'react-dropzone';
import { uploadPDF } from '../../api/pdf';
import { Form, Dropzone, DropzoneText, DropzoneLink, Input, ErrorMessage, Title, UploadButton, Spinner, DocumentList, DocumentItem } from './UploadForm.styles';
import { toast } from 'react-toastify';
import axios from 'axios';
import { saveAs } from 'file-saver';
import config from '../../config';
import { FaSpinner } from 'react-icons/fa';

interface Document {
  name: string;
  url: string;
  expiry: number; 
}

const UploadForm: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [downloadLink, setDownloadLink] = useState<string | null>(null);
  const [fileUploaded, setFileUploaded] = useState(false);
  const [originalFilename, setOriginalFilename] = useState<string | null>(null); 
  const [documents, setDocuments] = useState<Document[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  // Load documents from local storage and clean up expired ones
  useEffect(() => {
    const storedDocuments = localStorage.getItem('documents');
    if (storedDocuments) {
      const parsedDocuments: Document[] = JSON.parse(storedDocuments);
      const now = Date.now();
      const validDocuments = parsedDocuments.filter(doc => doc.expiry > now);
      setDocuments(validDocuments);
      localStorage.setItem('documents', JSON.stringify(validDocuments));
    }
  }, []);

  useEffect(() => {
    localStorage.setItem('documents', JSON.stringify(documents));
  }, [documents]);

  const onDrop = (acceptedFiles: File[]) => {
    formik.setFieldValue('file', acceptedFiles[0]);
    setSelectedFile(acceptedFiles[0]);
  };

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    multiple: false,
    accept: {
      'application/pdf': ['.pdf']
    }
  });

  const formik = useFormik({
    initialValues: {
      file: null,
    },
    validationSchema: Yup.object({
      file: Yup.mixed().required('A pdf-file is required'),
    }),
    onSubmit: async (values, { resetForm }) => {
      if (values.file) {
        setLoading(true);
        setFileUploaded(false);
        setOriginalFilename(null);
        let toastId: number | string | null = null;
        try {
          const response = await uploadPDF(values.file);
          const { session_id } = response;

          const socket = new WebSocket(`wss://${config.backendHost}/ws/${session_id}`);

          socket.onopen = () => {
            console.log('WebSocket connection established');
            toastId = toast.info('Waiting for document processing...');
          };

          socket.onmessage = async (event) => {
            const message = JSON.parse(event.data);
            if (message.status === 'ok') {
              try {
                let downloadLink = message.download_link;
                let originalFilename = message.original_filename;
                if (downloadLink.startsWith('http://minio:9000/')) {
                  downloadLink = downloadLink.replace('http://minio:9000/', '/minio/');
                }

                originalFilename = originalFilename.replace(/\.pdf$/, '_anonimized.pdf');

                const fileResponse = await axios.get(downloadLink, {
                  responseType: 'blob',
                });

                saveAs(fileResponse.data, originalFilename);
                setDownloadLink(downloadLink);
                setOriginalFilename(originalFilename); 
                setFileUploaded(true);

                const expiry = Date.now() + 15 * 60 * 1000; // 15 minutes from now
                const newDocument = { name: originalFilename, url: downloadLink, expiry };
                setDocuments((prevDocuments) => {
                  const updatedDocuments = [newDocument, ...prevDocuments];
                  return updatedDocuments.slice(0, 10);
                });

                toast.success('Document processing completed.');
                setTimeout(() => {
                  toast.info('You have 15 minutes to download the processed document.');
                }, 900000); // 15 minutes
              } catch (downloadError) {
                toast.error('Error downloading the file.');
                console.error('Error downloading the file:', downloadError);
              } finally {
                setLoading(false);
                if (toastId) toast.dismiss(toastId);
                setFileUploaded(false); // Reset file uploaded status
              }
            } else if (message.status === 'error') {
              toast.error('Document processing failed.');
              setLoading(false);
              if (toastId) toast.dismiss(toastId);
            } else {
              toast.warn(`Unknown status: ${message.status}`);
              setLoading(false);
              if (toastId) toast.dismiss(toastId);
            }
          };

          socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            toast.error('WebSocket connection error.');
            setLoading(false);
            if (toastId) toast.dismiss(toastId);
          };

          socket.onclose = () => {
            setLoading(false);
            console.log('WebSocket connection closed');
            if (toastId) toast.dismiss(toastId);
          };

          resetForm();
        } catch (error) {
          console.error('Error uploading file:', error);
          toast.error('Error uploading file.');
          setLoading(false);
          if (toastId) toast.dismiss(toastId);
        }
      }
    },
  });

  return (
    <Form onSubmit={formik.handleSubmit}>
      <Title>Upload PDF</Title>
      <Dropzone {...getRootProps()}>
        <input {...getInputProps()} />
        <img src="/place_holder.png" alt="upload illustration" />
        {selectedFile ? (
          <DropzoneText>{selectedFile.name}</DropzoneText>
        ) : (
          <>
            <DropzoneText>Drag and drop file here</DropzoneText>
            <DropzoneText>or <DropzoneLink>select a pdf-file</DropzoneLink> from your computer</DropzoneText>
          </>
        )}
      </Dropzone>
      {formik.errors.file ? <ErrorMessage>{formik.errors.file}</ErrorMessage> : null}
      {loading ? (
        <Spinner>
          <FaSpinner />
        </Spinner>
      ) : (
        !fileUploaded && (
          <UploadButton type="submit" disabled={loading}>
            {loading ? 'Uploading...' : 'Upload'}
          </UploadButton>
        )
      )}
      <DocumentList>
        {documents.map((doc, index) => (
          <DocumentItem key={index}>
            <a href={doc.url} download={doc.name}>
              {doc.name}
            </a>
          </DocumentItem>
        ))}
      </DocumentList>
    </Form>
  );
};

export default UploadForm;
