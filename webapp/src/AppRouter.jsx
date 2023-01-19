import { Routes, Route, BrowserRouter as Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import ServiceTable from './components/service/ServiceTable.jsx';
import ServiceForm from './components/service/ServiceForm.jsx';
import AuthRoute from './components/AuthRoute.jsx';
import { lazy } from 'react';

const browserHistory = createBrowserHistory();

const StatusPage = lazy(() => import('./pages/status/StatusPage'));
const LoginPage = lazy(() => import('./pages/login/LoginPage'));
const Dashboard = lazy(() => import('./pages/dashboard/Dashboard'));

const AppRouter = () => {
  return (
    <Router history={browserHistory}>
      <Routes>
        <Route path="/" element={<StatusPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route exact path="/dashboard" element={<AuthRoute />}>
          <Route
            exact
            path="/dashboard"
            element={
              <Dashboard>
                <ServiceTable />
              </Dashboard>
            }
          />
        </Route>
        <Route
          path="/dashboard/services/:id"
          element={
            <Dashboard>
              <ServiceForm />
            </Dashboard>
          }
        />
        <Route
          path="/dashboard/services"
          element={
            <Dashboard>
              <ServiceForm />
            </Dashboard>
          }
        />
      </Routes>
    </Router>
  );
};

export default AppRouter;
