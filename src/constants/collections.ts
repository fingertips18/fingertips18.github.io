import {
  SiCodewars,
  SiGithub,
  SiStackoverflow,
} from '@icons-pack/react-simple-icons';
import { SiLinkedin } from 'react-icons/si';

import { ROOTSECTION } from './enums';

export const ROOTMENU = [
  {
    label: ROOTSECTION.about,
    id: `#${ROOTSECTION.about}`,
  },
  {
    label: ROOTSECTION.skills,
    id: `#${ROOTSECTION.skills}`,
  },
  {
    label: ROOTSECTION.experience,
    id: `#${ROOTSECTION.experience}`,
  },
  {
    label: ROOTSECTION.projects,
    id: `#${ROOTSECTION.projects}`,
  },
  {
    label: ROOTSECTION.education,
    id: `#${ROOTSECTION.education}`,
  },
  {
    label: ROOTSECTION.contact,
    id: `#${ROOTSECTION.contact}`,
  },
];

export const BUILDS = ['Mobile Applications', 'Web Applications', 'Games'];

export const SOCIALS = [
  {
    icon: SiGithub,
    label: 'GitHub',
    href: 'https://github.com/fingertips18',
  },
  {
    icon: SiLinkedin,
    label: 'LinkedIn',
    href: 'https://linkedin.com/in/ghiantan',
  },
  {
    icon: SiStackoverflow,
    label: 'Stack Overflow',
    href: 'https://stackoverflow.com/users/18320841/fingertips',
  },
  {
    icon: SiCodewars,
    label: 'Codewars',
    href: 'https://codewars.com/users/Fingertips',
  },
];
