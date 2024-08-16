import {
  SiCodewars,
  SiGithub,
  SiLinkedin,
  SiStackoverflow,
} from "@icons-pack/react-simple-icons";

import { ROOTSECTIONS } from "./enums";

export const ROOTMENU = [
  {
    label: ROOTSECTIONS.about,
    id: `#${ROOTSECTIONS.about}`,
  },
  {
    label: ROOTSECTIONS.skills,
    id: `#${ROOTSECTIONS.skills}`,
  },
  {
    label: ROOTSECTIONS.experience,
    id: `#${ROOTSECTIONS.experience}`,
  },
  {
    label: ROOTSECTIONS.projects,
    id: `#${ROOTSECTIONS.projects}`,
  },
  {
    label: ROOTSECTIONS.education,
    id: `#${ROOTSECTIONS.education}`,
  },
  {
    label: ROOTSECTIONS.contact,
    id: `#${ROOTSECTIONS.contact}`,
  },
];

export const BUILDS = ["Mobile Applications", "Web Applications", "Games"];

export const SOCIALS = [
  {
    icon: SiGithub,
    label: "GitHub",
    href: "https://github.com/Fingertips18",
  },
  {
    icon: SiLinkedin,
    label: "LinkedIn",
    href: "https://linkedin.com/in/ghiantan",
  },
  {
    icon: SiStackoverflow,
    label: "Stack Overflow",
    href: "https://stackoverflow.com/users/18320841/fingertips",
  },
  {
    icon: SiCodewars,
    label: "Codewars",
    href: "https://codewars.com/users/Fingertips",
  },
];
