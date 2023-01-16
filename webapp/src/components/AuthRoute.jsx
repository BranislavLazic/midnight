import { Navigate, Route } from 'react-router-dom';

const AuthRoute = ({ path, element }) => {
  const isAuthenticated = localStorage.getItem('authUser');
  return isAuthenticated ? (
    <Route path={path} element={element} />
  ) : (
    <Navigate to="/login" />
  );
};

export default AuthRoute;
