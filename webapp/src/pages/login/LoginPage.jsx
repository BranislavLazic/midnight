import Nav from '../../components/Nav.jsx';
import { Button, Label, TextInput } from 'flowbite-react';
import axios from 'axios';
import { BASE_URL } from '../../../constants.cjs';
import { useNavigate } from 'react-router-dom';
import { useIntl } from 'react-intl';

const LoginPage = () => {
  const intl = useIntl();
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    const email = e.target.email.value;
    const password = e.target.password.value;
    try {
      const { data } = await axios.post(`${BASE_URL}/v1/login`, {
        email,
        password
      });
      sessionStorage.setItem('authUser', JSON.stringify(data?.authUser));
      sessionStorage.setItem('accessToken', data?.accessToken);
      navigate('/dashboard');
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <>
      <Nav />
      <div
        className='flex flex-col justify-center items-center w-full p-4 min-h-screen overflow-hidden bg-white dark:bg-gray-900'>
        <div className='flex flex-col justify-center items-center p-4 gap-4 dark:bg-gray-800 rounded'>
          <h1 className='text-xl font-bold dark:text-white md:text-3xl'>
            {intl.formatMessage({ id: 'signIn' })}
          </h1>
          <form
            className='flex flex-col gap-4 min-w-[24rem]'
            onSubmit={handleLogin}
          >
            <div className='flex flex-col gap-2'>
              <Label
                htmlFor='email'
                value={intl.formatMessage({ id: 'email' })}
              />
              <TextInput id='email' name='email' type='email' />
            </div>
            <div className='flex flex-col gap-2'>
              <Label
                htmlFor='password'
                value={intl.formatMessage({ id: 'password' })}
              />
              <TextInput id='password' name='password' type='password' />
            </div>
            <Button type='submit'>{intl.formatMessage({ id: 'login' })}</Button>
          </form>
        </div>
      </div>
    </>
  );
};

export default LoginPage;
