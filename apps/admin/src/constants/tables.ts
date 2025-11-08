import {
  Brain,
  FolderKanban,
  GraduationCap,
  type LucideIcon,
} from 'lucide-react';

import { Route } from '@/routes/route';

type Table = {
  label: string;
  url: string;
  icon: LucideIcon;
};
export const TABLES: Table[] = [
  {
    label: 'Project',
    url: Route.project,
    icon: FolderKanban,
  },
  {
    label: 'Education',
    url: Route.education,
    icon: GraduationCap,
  },
  {
    label: 'Skill',
    url: Route.skill,
    icon: Brain,
  },
];

type Map = Partial<{
  [K in (typeof Route)[keyof typeof Route]]: {
    label: string;
    url: string;
    subPaths: {
      label: string;
      url: string;
    }[];
  };
}>;
export const TABLE_BREADCRUMB_MAP: Map = {
  [Route.project]: {
    label: 'Project',
    url: Route.project,
    subPaths: [
      {
        label: 'New Project',
        url: `${Route.project}/add`,
      },
      {
        label: 'Edit Project',
        url: `${Route.project}/edit`,
      },
    ],
  },
  [Route.education]: {
    label: 'Education',
    url: Route.education,
    subPaths: [
      {
        label: 'New Education',
        url: `${Route.education}/add`,
      },
      {
        label: 'Edit Education',
        url: `${Route.education}/edit`,
      },
    ],
  },
  [Route.skill]: {
    label: 'Skill',
    url: Route.skill,
    subPaths: [
      {
        label: 'New Skill',
        url: `${Route.skill}/add`,
      },
      {
        label: 'Edit Skill',
        url: `${Route.skill}/edit`,
      },
    ],
  },
};
