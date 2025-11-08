import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom';

import AddEducationPage from '@/pages/education/add/page';
import EducationPage from '@/pages/education/page';
import ErrorPage from '@/pages/error/page';
import NotFoundPage from '@/pages/not-found/page';
import AddProjectPage from '@/pages/project/add/page';
import ProjectPage from '@/pages/project/page';
import RootLayout from '@/pages/root/layout';
import RootPage from '@/pages/root/page';
import AddSkillPage from '@/pages/skill/add/page';
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
      <Route path={`${AppRoute.project}/add`} element={<AddProjectPage />} />

      {/* Education */}
      <Route path={AppRoute.education} element={<EducationPage />} />
      <Route
        path={`${AppRoute.education}/add`}
        element={<AddEducationPage />}
      />

      {/* Skill */}
      <Route path={AppRoute.skill} element={<SkillPage />} />
      <Route path={`${AppRoute.skill}/add`} element={<AddSkillPage />} />

      {/* Not Found */}
      <Route path={AppRoute.notFound} element={<NotFoundPage />} />
    </Route>,
  ),
);
