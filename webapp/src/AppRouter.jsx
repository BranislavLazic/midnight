import { Routes, Route, BrowserRouter as Router } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import StatusPage from './pages/status/StatusPage.jsx';

const browserHistory = createBrowserHistory();

const AppRouter = () => {
  return (
    <Router history={browserHistory}>
      <Routes>
        <Route path="/" element={<StatusPage />} />
      </Routes>
    </Router>
  );
};

export default AppRouter;
