import { Button, Navbar } from 'flowbite-react';

import { Pencil, SignOut } from 'phosphor-react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { BASE_URL } from '../../constants.cjs';
import { useIntl } from 'react-intl';

const Nav = ({ editBoardButtonShown = false }) => {
  const intl = useIntl();
  const navigate = useNavigate();
  const authUser = localStorage.getItem('authUser');

  const handleLogout = async () => {
    try {
      await axios.post(`${BASE_URL}/v1/logout`);
      localStorage.removeItem('authUser');
      navigate('/');
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <Navbar fluid={true}>
      <Link to="/">
        <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
          Midnight
        </span>
      </Link>
      <div className="flex gap-4">
        {editBoardButtonShown && (
          <Link to="/dashboard">
            <Button>
              <div className="flex gap-1 items-center">
                <Pencil className="h-4 w-4" weight="bold" />
                <span>{intl.formatMessage({ id: 'editBoard' })}</span>
              </div>
            </Button>
          </Link>
        )}
        {authUser && (
          <Button color="gray" onClick={handleLogout}>
            <div className="flex gap-1 items-center">
              <SignOut className="h-4 w-4" weight="bold" />
              <span>{intl.formatMessage({ id: 'logout' })}</span>
            </div>
          </Button>
        )}
      </div>
    </Navbar>
  );
};
export default Nav;
