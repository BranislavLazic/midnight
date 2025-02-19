import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Button, Table } from 'flowbite-react';
import { BASE_URL } from '../../../constants.cjs';
import { useIntl } from 'react-intl';
import { authorized } from '../../lib/authorized.js';
import { PencilSimple, Plus, Trash } from '@phosphor-icons/react';

const ServiceTable = () => {
  const intl = useIntl();
  const [services, setServices] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await authorized.get(`${BASE_URL}/v1/services`);
        setServices(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  const handleDelete = async (id) => {
    try {
      await authorized.delete(`${BASE_URL}/v1/services/${id}`);
      const { data } = await authorized.get(`${BASE_URL}/v1/services`);
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
                <Plus className="h-4 w-4" weight="bold" />
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
                      <PencilSimple className="h-4 w-4" weight="bold" />
                    </Button>
                  </Link>
                  <Button color="gray" onClick={() => handleDelete(s.id)}>
                    <Trash className="h-4 w-4" weight="bold" />
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
