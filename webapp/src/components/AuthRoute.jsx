import { Navigate, Outlet } from 'react-router-dom';

const AuthRoute = () => {
  const isAuthenticated =
    sessionStorage.getItem('authUser') && sessionStorage.getItem('accessToken');
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" />;
};

export default AuthRoute;
