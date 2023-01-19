import React from 'react';
import ReactDOM from 'react-dom/client';
import AppRouter from './AppRouter.jsx';
import './index.css';
import { IntlProvider } from 'react-intl';
import enUs from './lang/en-US.json';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <IntlProvider locale="en" defaultLocale="en" messages={enUs}>
      <AppRouter />
    </IntlProvider>
  </React.StrictMode>
);
