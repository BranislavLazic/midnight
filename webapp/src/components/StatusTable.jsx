import { Badge, Table } from 'flowbite-react';
import axios from 'axios';
import { useEffect } from 'react';
import { useState } from 'react';

const StatusTable = () => {
  const [liveness, setLiveness] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await axios.get('http://localhost:8000/v1/status');
        setLiveness(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  return (
    <div className="flex justify-center p-4 min-h-screen overflow-hidden bg-white dark:bg-gray-900">
      <div className="min-w-[50%]">
        <Table>
          <Table.Head>
            <Table.HeadCell>Service</Table.HeadCell>
            <Table.HeadCell>Status</Table.HeadCell>
            <Table.HeadCell>Version</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            {liveness.map((l, idx) => (
              <Table.Row
                key={idx}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  <a href={l.url} target="_blank">
                    {l.url}
                  </a>
                </Table.Cell>
                <Table.Cell className="whitespace-nowrap font-medium text-gray-900 dark:text-white">
                  <Badge color="success">{l.status}</Badge>
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