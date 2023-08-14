import { Routes, Route, BrowserRouter as Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import React, { lazy, Suspense } from 'react';

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
      <Suspense fallback="...">
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
      </Suspense>
    </Router>
  );
};

export default AppRouter;
