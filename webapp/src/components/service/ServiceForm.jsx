import { Label, TextInput, Button } from 'flowbite-react';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useParams } from 'react-router-dom';
import { useFormik } from 'formik';
import { validationSchema } from './validation/validation';

const initialFormValues = {
  name: '',
  url: '',
  checkIntervalSeconds: 5,
};

const ServiceForm = () => {
  const { id } = useParams();
  const [service, setService] = useState();

  const formik = useFormik({
    enableReinitialize: true,
    validationSchema: validationSchema,
    initialValues: initialFormValues,
    onSubmit: async (values) => {
      if (!service) {
        await axios.post(
          `${process.env.REACT_APP_BASE_API_URL}/v1/services`,
          values,
        );
      }
    },
  });

  useEffect(() => {
    if (id) {
      (async () => {
        try {
          const { data } = await axios.get(
            `${process.env.REACT_APP_BASE_API_URL}/v1/services/${id}`,
          );
          setService(data);
        } catch (e) {
          console.error(e);
        }
      })();
    }
  }, [id]);

  return (
    <div className="flex grow h-fit w-fit">
      <div className="flex flex-col gap-4 min-w-[32rem] p-4 dark:bg-gray-800 rounded">
        <h1 className="text-lg dark:text-white">Edit service</h1>
        <form
          className="flex flex-col gap-4 max-w-xl"
          onSubmit={(e) => e.preventDefault()}
        >
          <div className="flex flex-col gap-2">
            <Label
              htmlFor="name"
              value="Name (optional)"
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
              value="URL"
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
              htmlFor="checkIntervalSeconds"
              value="Check interval in seconds"
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
                Cancel
              </Button>
            </Link>
            <Button
              className="grow"
              type="submit"
              onClick={formik.handleSubmit}
            >
              Save
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ServiceForm;
