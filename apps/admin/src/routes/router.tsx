import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';

import EducationPage from '@/pages/education/page';
import ErrorPage from '@/pages/error/page';
import NotFoundPage from '@/pages/not-found/page';
import ProjectPage from '@/pages/project/page';
import RootLayout from '@/pages/root/layout';
import RootPage from '@/pages/root/page';
import SkillPage from '@/pages/skill/page';

import { Route as AppRoute } from './route';

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      path={AppRoute.root}
      element={<RootLayout />}
      errorElement={<ErrorPage />}
    >
      {/* Root */}
      <Route index element={<RootPage />} />

      {/* Project */}
      <Route path={AppRoute.project} element={<ProjectPage />} />

      {/* Education */}
      <Route path={AppRoute.education} element={<EducationPage />} />

      {/* Skill */}
      <Route path={AppRoute.skill} element={<SkillPage />} />

      {/* Not Found */}
      <Route path={AppRoute.notFound} element={<NotFoundPage />} />
    </Route>,
  ),
);
