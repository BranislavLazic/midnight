import { Badge, Table } from 'flowbite-react';
import axios from 'axios';
import { useEffect } from 'react';
import { useState } from 'react';
import { BASE_URL } from '../../constants.cjs';

const StatusTable = () => {
  const [serviceStatuses, setServiceStatuses] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await axios.get(`${BASE_URL}/v1/status`);
        setServiceStatuses(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  return (
    <div className="flex justify-center w-full p-4 min-h-screen overflow-hidden bg-white dark:bg-gray-900">
      <div className="w-full">
        <Table>
          <Table.Head>
            <Table.HeadCell>Service</Table.HeadCell>
            <Table.HeadCell>Status</Table.HeadCell>
            <Table.HeadCell>Version</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            {serviceStatuses.map((l, idx) => (
              <Table.Row
                key={idx}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  <a href={l.url} target="_blank">
                    {l.name || l.url}
                  </a>
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  <Badge color={l.statusCode < 400 ? 'success' : 'failure'}>
                    {l.status}
                  </Badge>
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  {l?.version}
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>
    </div>
  );
};

export default StatusTable;
