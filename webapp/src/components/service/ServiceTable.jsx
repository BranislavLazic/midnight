import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import { Button, Table } from 'flowbite-react';
import { PencilIcon, TrashIcon, PlusIcon } from '@heroicons/react/24/outline';
import { BASE_URL } from '../../../constants.cjs';
import { useIntl } from 'react-intl';

const ServiceTable = () => {
  const intl = useIntl();
  const [services, setServices] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await axios.get(`${BASE_URL}/v1/services`);
        setServices(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${BASE_URL}/v1/services/${id}`, {
        withCredentials: true,
      });
      const { data } = await axios.get(`${BASE_URL}/v1/services`);
      setServices(data);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="flex grow">
      <div className="flex flex-col gap-4 min-w-full p-4 dark:bg-gray-800 rounded">
        <div className="flex flex-row justify-between items-center">
          <h1 className="text-lg font-bold dark:text-white">
            {intl.formatMessage({ id: 'services' })}
          </h1>
          <Link to="/dashboard/services">
            <Button>
              <div className="flex gap-1 items-center">
                <PlusIcon className="h-4 w-4" fontWeight="bold" />
                <span>{intl.formatMessage({ id: 'add' })}</span>
              </div>
            </Button>
          </Link>
        </div>
        <Table>
          <Table.Head>
            <Table.HeadCell>
              {intl.formatMessage({ id: 'name' })}
            </Table.HeadCell>
            <Table.HeadCell>{intl.formatMessage({ id: 'url' })}</Table.HeadCell>
            <Table.HeadCell>
              {intl.formatMessage({ id: 'environment' })}
            </Table.HeadCell>
            <Table.HeadCell>
              {intl.formatMessage({ id: 'checkIntervalSeconds' })}
            </Table.HeadCell>
            <Table.HeadCell>
              {intl.formatMessage({ id: 'actions' })}
            </Table.HeadCell>
          </Table.Head>
          <Table.Body>
            {services.map((s, idx) => (
              <Table.Row
                key={idx}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.name}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.url}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s?.environment?.name || '-'}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.checkIntervalSeconds}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white flex gap-2">
                  <Link to={`/dashboard/services/${s.id}`}>
                    <Button color="gray">
                      <PencilIcon className="h-4 w-4" fontWeight="bold" />
                    </Button>
                  </Link>
                  <Button color="gray" onClick={() => handleDelete(s.id)}>
                    <TrashIcon className="h-4 w-4" fontWeight="bold" />
                  </Button>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>
    </div>
  );
};

export default ServiceTable;
