import * as Yup from 'yup';

export const validationSchema = () =>
  Yup.object({
    name: Yup.string().max(200, 'Name cannot have more than 255 characters.'),
    url: Yup.string()
      .required('URL is required')
      .max(4096, 'Name URL have more than 4096 characters.'),
    checkIntervalSeconds: Yup.number()
      .required('Check interval is required')
      .min(3, 'Check interval cannot be less than 3 seconds.')
      .max(86400, 'Check interval cannot be greater than 86400 seconds.'),
  });
