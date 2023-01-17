import { Button, Navbar } from 'flowbite-react';
import {
  PencilIcon,
  ArrowLeftOnRectangleIcon,
} from '@heroicons/react/24/solid';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { BASE_URL } from '../../constants.cjs';

const Nav = ({ editBoardButtonShown = false }) => {
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
      <Navbar.Brand href="/">
        <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
          Midnight
        </span>
      </Navbar.Brand>
      <div className="flex gap-4">
        {editBoardButtonShown && (
          <Link to="/dashboard">
            <Button>
              <div className="flex gap-1 items-center">
                <PencilIcon className="h-4 w-4" />
                <span>Edit board</span>
              </div>
            </Button>
          </Link>
        )}
        {authUser && (
          <Button color="gray" onClick={handleLogout}>
            <div className="flex gap-1 items-center">
              <ArrowLeftOnRectangleIcon className="h-4 w-4" />
              <span>Logout</span>
            </div>
          </Button>
        )}
      </div>
    </Navbar>
  );
};
export default Nav;
