import { Button, Navbar } from 'flowbite-react';
import { PencilIcon } from '@heroicons/react/24/solid';
import { Link } from 'react-router-dom';

const Nav = ({ editBoardButtonShown = false }) => {
  return (
    <Navbar fluid={true}>
      <Navbar.Brand href="/">
        <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
          Midnight
        </span>
      </Navbar.Brand>
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
    </Navbar>
  );
};
export default Nav;
