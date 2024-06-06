import * as Yup from 'yup';

export const fileValidationSchema = Yup.object({
  file: Yup.mixed()
    .required('A file is required')
    .test('fileType', 'Only PDF files are allowed', (value) => {
      return value && (value as File).type === 'application/pdf';
    }),
});
