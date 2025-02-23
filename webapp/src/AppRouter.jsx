import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import React, { lazy, Suspense } from 'react';
import Loader from './components/Loader.jsx';

const browserHistory = createBrowserHistory();

const StatusPage = lazy(() => import('./pages/status/StatusPage.jsx'));
const Dashboard = lazy(() => import('./pages/dashboard/Dashboard.jsx'));
const ServiceTable = lazy(() =>
  import('./components/service/ServiceTable.jsx')
);
const ServiceForm = lazy(() => import('./components/service/ServiceForm.jsx'));
const LoginPage = lazy(() => import('./pages/login/LoginPage.jsx'));
const AuthRoute = lazy(() => import('./components/AuthRoute.jsx'));

const AppRouter = () => {
  return (
    <Router history={browserHistory}>
      <Suspense fallback={<Loader />}>
        <Routes>
          <Route path='/' element={<StatusPage />} />
          <Route path='/login' element={<LoginPage />} />
          <Route exact path='/dashboard' element={<AuthRoute />}>
            <Route
              exact
              path='/dashboard'
              element={
                <Dashboard>
                  <ServiceTable />
                </Dashboard>
              }
            />
          </Route>
          <Route exact path='/dashboard/services/:id' element={<AuthRoute />}>
            <Route
              path='/dashboard/services/:id'
              element={
                <Dashboard>
                  <ServiceForm />
                </Dashboard>
              }
            />
          </Route>
          <Route exact path='/dashboard/services' element={<AuthRoute />}>
            <Route
              path='/dashboard/services'
              element={
                <Dashboard>
                  <ServiceForm />
                </Dashboard>
              }
            />
          </Route>
        </Routes>
      </Suspense>
    </Router>
  );
};

export default AppRouter;
