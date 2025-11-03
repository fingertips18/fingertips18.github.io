import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';

import { RootLayout } from '@/pages/root/layout';
import { RootPage } from '@/pages/root/page';

import { Route as AppRoute } from './route';

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={AppRoute.root} element={<RootLayout />}>
      {/* Root */}
      <Route index element={<RootPage />} />
    </Route>,
  ),
);
