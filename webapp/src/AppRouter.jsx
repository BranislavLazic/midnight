import { Routes, Route, BrowserRouter as Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import StatusPage from './pages/status/StatusPage.jsx';
import Dashboard from './pages/dashboard/Dashboard.jsx';
import ServiceTable from './components/service/ServiceTable.jsx';
import ServiceForm from './components/service/ServiceForm.jsx';

const browserHistory = createBrowserHistory();

const AppRouter = () => {
  return (
    <Router history={browserHistory}>
      <Routes>
        <Route path="/" element={<StatusPage />} />
        <Route
          path="/dashboard"
          element={
            <Dashboard>
              <ServiceTable />
            </Dashboard>
          }
        />
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
