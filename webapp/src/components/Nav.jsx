import { Button, Checkbox, Dropdown, Label, Navbar } from 'flowbite-react';

import { Link, useLocation, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { BASE_URL } from '../../constants.cjs';
import { useIntl } from 'react-intl';
import { useEffect, useState } from 'react';
import { PencilSimple, SignOut } from '@phosphor-icons/react';

const Nav = ({
  editBoardButtonShown = false,
  onEnvironmentChange = () => {},
}) => {
  const intl = useIntl();
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const authUser = sessionStorage.getItem('authUser');
  const [environments, setEnvironments] = useState([]);
  const [selectedEnvironmentNames, setSelectedEnvironmentNames] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await axios.get(`${BASE_URL}/v1/environments`);
        setEnvironments(data);
        setSelectedEnvironmentNames(data.map(({ id }) => id));
      } catch (e) {
        console.error(e);
      }
    })();
  }, []);

  useEffect(() => {
    const envQueryParams = '?env=' + selectedEnvironmentNames.join('&env=');
    onEnvironmentChange(envQueryParams);
  }, [selectedEnvironmentNames]);

  const handleLogout = () => {
    try {
      sessionStorage.removeItem('authUser');
      sessionStorage.removeItem('accessToken');
      navigate('/');
    } catch (e) {
      console.error(e);
    }
  };

  const handleEnvChange = (evt) => {
    const envId = Number(evt.target.id);
    if (evt.target.checked) {
      const { id } = environments.find(({ id }) => id === envId);
      setSelectedEnvironmentNames((current) =>
        Array.from(new Set([...current, id]))
      );
    } else {
      setSelectedEnvironmentNames((current) =>
        current.filter((id) => id !== envId)
      );
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
        {pathname === '/' && (
          <Dropdown
            label={intl.formatMessage({ id: 'environments' })}
            dismissOnClick={false}
          >
            <div className="flex flex-col gap-2 p-2">
              {environments.map((env, idx) => (
                <div key={idx} className="flex items-center gap-2">
                  <Checkbox
                    id={env.id}
                    defaultChecked
                    onChange={handleEnvChange}
                  />
                  <Label htmlFor={`env-${idx}`}>{env.name}</Label>
                </div>
              ))}
            </div>
          </Dropdown>
        )}
        {editBoardButtonShown && (
          <Link to="/dashboard">
            <Button color="light">
              <div className="flex gap-1 items-center">
                <PencilSimple className="h-4 w-4" weight="bold" />
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
