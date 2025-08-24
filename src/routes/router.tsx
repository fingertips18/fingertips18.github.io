import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';

import ErrorPage from '@/pages/error/page';
import NotFoundPage from '@/pages/not-found/page';
import { ProjectsPage } from '@/pages/projects/page';
import RootLayout from '@/pages/root/layout';
import RootPage from '@/pages/root/page';
import { SkillsPage } from '@/pages/skills/page';

import { AppRoutes } from './app-routes';

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      path={AppRoutes.root}
      element={<RootLayout />}
      errorElement={<ErrorPage />}
    >
      {/* Root */}
      <Route index element={<RootPage />} />

      {/* Skills */}
      <Route path={AppRoutes.skills} element={<SkillsPage />} />

      {/* Projects */}
      <Route path={AppRoutes.projects} element={<ProjectsPage />} />

      {/* Github 404 */}
      <Route path={AppRoutes.github404} element={<NotFoundPage />} />

      {/* Not Found */}
      <Route path={AppRoutes.notFound} element={<NotFoundPage />} />
    </Route>,
  ),
);
