import Nav from '../../components/Nav';
import { Sidebar } from 'flowbite-react';
import { useIntl } from 'react-intl';
import { Link } from 'react-router-dom';

const Dashboard = ({ children }) => {
  const intl = useIntl();
  return (
    <div className="flex h-screen w-full flex-col overflow-hidden">
      <Nav />
      <div className="flex gap-4 min-h-screen p-4 bg-white dark:bg-gray-900">
        <Sidebar>
          <Sidebar.Items>
            <Sidebar.ItemGroup>
              <li>
                <Link
                  to="/dashboard"
                  className="relative flex items-center rounded-lg p-2 text-base font-normal text-gray-900 hover:bg-gray-100 dark:text-white dark:hover:bg-gray-700"
                >
                  <span className="absolute inset-y-0 left-0 w-2 bg-cyan-600 rounded-lg"></span>
                  <span className="ml-4">
                    {intl.formatMessage({ id: 'services' })}
                  </span>
                </Link>
              </li>
            </Sidebar.ItemGroup>
          </Sidebar.Items>
        </Sidebar>
        {children}
      </div>
    </div>
  );
};

export default Dashboard;
