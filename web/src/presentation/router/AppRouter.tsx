import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { ProtectedRoute } from './ProtectedRoute';
import { LoginPage } from '../pages/auth/LoginPage';
import { RegisterPage } from '../pages/auth/RegisterPage';
import { HomePage } from '../pages/home/HomePage';
import { DashboardPage } from '../pages/dashboard/DashboardPage';
import { TransactionsPage } from '../pages/transactions/TransactionsPage';
import { TransactionFormPage } from '../pages/transactions/TransactionFormPage';
import { AccountsPage } from '../pages/accounts/AccountsPage';
import { StatsPage } from '../pages/stats/StatsPage';
import { ImportPage } from '../pages/import/ImportPage';
import { ImportUploadPage } from '../pages/import/ImportUploadPage';
import { ImportPreviewPage } from '../pages/import/ImportPreviewPage';
import { ImportResultPage } from '../pages/import/ImportResultPage';
import { SettingsPage } from '../pages/settings/SettingsPage';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <HomePage />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: 'transactions',
        element: <TransactionsPage />,
      },
      {
        path: 'transactions/new',
        element: <TransactionFormPage />,
      },
      {
        path: 'accounts',
        element: <AccountsPage />,
      },
      {
        path: 'stats',
        element: <StatsPage />,
      },
      {
        path: 'import',
        element: <ImportPage />,
      },
      {
        path: 'import/upload',
        element: <ImportUploadPage />,
      },
      {
        path: 'import/preview',
        element: <ImportPreviewPage />,
      },
      {
        path: 'import/result',
        element: <ImportResultPage />,
      },
      {
        path: 'settings',
        element: <SettingsPage />,
      },
    ],
  },
]);

export function AppRouter() {
  return <RouterProvider router={router} />;
}
