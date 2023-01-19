import { Label, TextInput, Button, Textarea } from 'flowbite-react';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useFormik } from 'formik';
import { validationSchema } from './validation/validation';
import { BASE_URL } from '../../../constants.cjs';
import { useIntl } from 'react-intl';

const initialFormValues = {
  name: '',
  url: '',
  responseBody: '',
  checkIntervalSeconds: 5,
};

const ServiceForm = () => {
  const { id } = useParams();
  const intl = useIntl();
  const navigate = useNavigate();
  const [service, setService] = useState();

  const formik = useFormik({
    enableReinitialize: true,
    validationSchema: validationSchema(intl),
    initialValues: initialFormValues,
    onSubmit: async (values) => {
      try {
        if (!service) {
          await axios.post(`${BASE_URL}/v1/services`, values);
          formik.resetForm();
        } else {
          await axios.put(`${BASE_URL}/v1/services/${id}`, values);
          navigate('/dashboard');
        }
      } catch (e) {
        console.error(e);
      }
    },
  });

  useEffect(() => {
    if (id) {
      (async () => {
        try {
          const { data } = await axios.get(`${BASE_URL}/v1/services/${id}`);
          setService(data);
          await formik.setValues(data);
        } catch (e) {
          console.error(e);
        }
      })();
    }
  }, [id]);

  return (
    <div className="flex grow h-fit w-fit">
      <div className="flex flex-col gap-4 min-w-[32rem] p-4 dark:bg-gray-800 rounded">
        <h1 className="text-lg dark:text-white">
          {id
            ? intl.formatMessage({ id: 'editService' })
            : intl.formatMessage({ id: 'addService' })}
        </h1>
        <form
          className="flex flex-col gap-4 max-w-xl"
          onSubmit={(e) => e.preventDefault()}
        >
          <div className="flex flex-col gap-2">
            <Label
              htmlFor="name"
              value={intl.formatMessage({ id: 'nameOptional' })}
              color={
                formik.touched.name &&
                formik.errors.name !== undefined &&
                formik.errors.name
                  ? 'failure'
                  : undefined
              }
            />
            <TextInput
              id="name"
              name="name"
              type="text"
              color={
                formik.touched.name &&
                formik.errors.name !== undefined &&
                formik.errors.name
                  ? 'failure'
                  : 'gray'
              }
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              value={formik.values.name}
            />
            {formik.touched.name &&
              formik.errors.name !== undefined &&
              formik.errors.name && (
                <span className="text-xs font-medium text-red-700 dark:text-red-500">
                  {formik.errors.name}
                </span>
              )}
          </div>
          <div className="flex flex-col gap-2">
            <Label
              htmlFor="url"
              value={intl.formatMessage({ id: 'url' })}
              color={
                formik.touched.url &&
                formik.errors.url !== undefined &&
                formik.errors.url
                  ? 'failure'
                  : undefined
              }
            />
            <TextInput
              id="url"
              name="url"
              type="text"
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              value={formik.values.url}
              color={
                formik.touched.url &&
                formik.errors.url !== undefined &&
                formik.errors.url
                  ? 'failure'
                  : 'gray'
              }
            />
            {formik.touched.url &&
              formik.errors.url !== undefined &&
              formik.errors.url && (
                <span className="text-xs font-medium text-red-700 dark:text-red-500">
                  {formik.errors.url}
                </span>
              )}
          </div>
          <div className="flex flex-col gap-2">
            <Label
              htmlFor="responseBody"
              value={intl.formatMessage({ id: 'responseBodyOptional' })}
              color={
                formik.touched.responseBody &&
                formik.errors.responseBody !== undefined &&
                formik.errors.responseBody
                  ? 'failure'
                  : undefined
              }
            />
            <Textarea
              id="responseBody"
              name="responseBody"
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              value={formik.values.responseBody}
              color={
                formik.touched.responseBody &&
                formik.errors.responseBody !== undefined &&
                formik.errors.responseBody
                  ? 'failure'
                  : 'gray'
              }
            />
            {formik.touched.responseBody &&
              formik.errors.responseBody !== undefined &&
              formik.errors.responseBody && (
                <span className="text-xs font-medium text-red-700 dark:text-red-500">
                  {formik.errors.responseBody}
                </span>
              )}
          </div>
          <div className="flex flex-col gap-2">
            <Label
              htmlFor="checkIntervalSeconds"
              value={intl.formatMessage({ id: 'checkIntervalSeconds' })}
              color={
                formik.touched.checkIntervalSeconds &&
                formik.errors.checkIntervalSeconds !== undefined &&
                formik.errors.checkIntervalSeconds
                  ? 'failure'
                  : undefined
              }
            />
            <TextInput
              id="checkIntervalSeconds"
              name="checkIntervalSeconds"
              type="number"
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              value={formik.values.checkIntervalSeconds}
              color={
                formik.touched.checkIntervalSeconds &&
                formik.errors.checkIntervalSeconds !== undefined &&
                formik.errors.checkIntervalSeconds
                  ? 'failure'
                  : 'gray'
              }
            />
            {formik.touched.checkIntervalSeconds &&
              formik.errors.checkIntervalSeconds !== undefined &&
              formik.errors.checkIntervalSeconds && (
                <span className="text-xs font-medium text-red-700 dark:text-red-500">
                  {formik.errors.checkIntervalSeconds}
                </span>
              )}
          </div>
          <div className="flex flex-row gap-4">
            <Link className="flex grow" to="/dashboard">
              <Button color="gray" className="grow">
                {intl.formatMessage({ id: 'cancel' })}
              </Button>
            </Link>
            <Button
              className="grow"
              type="submit"
              onClick={formik.handleSubmit}
            >
              {intl.formatMessage({ id: 'save' })}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ServiceForm;
