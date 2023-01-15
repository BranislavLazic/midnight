import React from 'react';
import StatusTable from '../../components/StatusTable';
import Nav from '../../components/Nav';

const StatusPage = () => {
  return (
    <div className="flex h-screen w-full flex-col overflow-hidden">
      <Nav editBoardButtonShown />
      <StatusTable />
    </div>
  );
};

export default StatusPage;
