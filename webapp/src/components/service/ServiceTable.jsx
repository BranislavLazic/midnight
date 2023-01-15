import { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import { Button, Table } from 'flowbite-react';
import { PencilIcon, TrashIcon } from '@heroicons/react/24/solid';

const ServiceTable = () => {
  const [services, setServices] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await axios.get('http://localhost:8000/v1/services');
        setServices(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  return (
    <div className="flex grow">
      <div className="flex flex-col gap-4 min-w-full p-4 dark:bg-gray-800 rounded">
        <h1 className="text-lg dark:text-white">Services</h1>
        <Table>
          <Table.Head>
            <Table.HeadCell>Name</Table.HeadCell>
            <Table.HeadCell>Url</Table.HeadCell>
            <Table.HeadCell>Check interval (seconds)</Table.HeadCell>
            <Table.HeadCell>Actions</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            {services.map((s, idx) => (
              <Table.Row
                key={idx}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.Name}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.url}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {s.checkIntervalSeconds}
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white flex gap-2">
                  <Link to={`/dashboard/services/${s.id}`}>
                    <Button color="gray">
                      <PencilIcon className="h-4 w-4" />
                    </Button>
                  </Link>
                  <Button color="gray">
                    <TrashIcon className="h-4 w-4" />
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
