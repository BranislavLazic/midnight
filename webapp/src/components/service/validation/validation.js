import * as Yup from 'yup';

export const validationSchema = (intl) =>
  Yup.object({
    name: Yup.string().max(
      200,
      intl.formatMessage({ id: 'nameMaxCharsError' })
    ),
    url: Yup.string()
      .required(intl.formatMessage({ id: 'urlRequiredError' }))
      .max(4096, intl.formatMessage({ id: 'urlMaxCharsError' })),
    checkIntervalSeconds: Yup.number()
      .required(intl.formatMessage({ id: 'checkIntervalRequiredError' }))
      .min(3, intl.formatMessage({ id: 'checkIntervalMinValueError' }))
      .max(86400, intl.formatMessage({ id: 'checkIntervalMaxValueError' })),
  });
