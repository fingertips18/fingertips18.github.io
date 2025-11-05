import {
  Brain,
  FolderKanban,
  GraduationCap,
  type LucideIcon,
} from 'lucide-react';

import { Route } from '@/routes/route';

type Table = {
  title: string;
  url: string;
  icon: LucideIcon;
};
export const TABLES: Table[] = [
  {
    title: 'Project',
    url: Route.project,
    icon: FolderKanban,
  },
  {
    title: 'Education',
    url: Route.education,
    icon: GraduationCap,
  },
  {
    title: 'Skill',
    url: Route.skill,
    icon: Brain,
  },
];

type TITLE = {
  [K in (typeof Route)[keyof typeof Route]]: string;
};
export const TABLE_TITLE: TITLE = {
  [Route.root]: 'Dashboard',
  [Route.project]: 'Project',
  [Route.education]: 'Education',
  [Route.skill]: 'Skill',
};
