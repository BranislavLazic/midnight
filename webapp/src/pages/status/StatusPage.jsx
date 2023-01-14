import React from 'react';
import { Navbar } from 'flowbite-react';
import StatusTable from '../../components/StatusTable';

const StatusPage = () => {
  return (
    <div className="flex h-screen w-full flex-col overflow-hidden">
      <Navbar fluid={true}>
        <Navbar.Brand href="/">
          <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
            Midnight
          </span>
        </Navbar.Brand>
      </Navbar>
      <StatusTable />
    </div>
  );
};

export default StatusPage;
