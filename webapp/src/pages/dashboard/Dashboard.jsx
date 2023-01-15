import Nav from '../../components/Nav';
import { Sidebar } from 'flowbite-react';

const Dashboard = ({ children }) => {
  return (
    <div className="flex h-screen w-full flex-col overflow-hidden">
      <Nav />
      <div className="flex gap-4 min-h-screen p-4 bg-white dark:bg-gray-900">
        <Sidebar>
          <Sidebar.Items>
            <Sidebar.ItemGroup>
              <Sidebar.Item href="/dashboard">Services</Sidebar.Item>
            </Sidebar.ItemGroup>
          </Sidebar.Items>
        </Sidebar>
        {children}
      </div>
    </div>
  );
};

export default Dashboard;
