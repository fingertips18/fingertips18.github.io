import { useLenis } from 'lenis/react';
import { useCallback, useEffect, useMemo } from 'react';
import { useLocation } from 'react-router-dom';

import { QUERYELEMENT, ROOTSECTION } from '@/constants/enums';
import { useElementsByQuery } from '@/lib/hooks/useElementsByQuery';
import { useMounted } from '@/lib/hooks/useMounted';
import { useRootSectionStore } from '@/lib/stores/useRootSectionStore';
import { AppRoutes } from '@/routes/app-routes';

import { SheetMenu } from './sheet-menu';
import { SpreadMenu } from './spread-menu';

const Navbar = () => {
  const { active, onActive } = useRootSectionStore((state) => state);
  const rootSections = useElementsByQuery(`.${QUERYELEMENT.rootSection}`);
  const isMounted = useMounted();
  const location = useLocation();
  const lenis = useLenis();

  const sectionOffsets = useMemo(() => {
    const sections = [];

    if (!rootSections) return;

    for (let i = 0; i < rootSections.length; i++) {
      sections.push({
        offset: rootSections[i].offsetTop - 2,
        id: rootSections[i].id,
      });
    }

    return sections;
  }, [rootSections]);

  const handleActiveSection = useCallback(() => {
    if (
      !sectionOffsets ||
      location.pathname !== AppRoutes.root ||
      !lenis ||
      lenis.isScrolling
    )
      return;

    for (let i = 0; i < sectionOffsets.length; i++) {
      if (window.scrollY >= sectionOffsets[i].offset) {
        onActive(sectionOffsets[i].id as ROOTSECTION);
      }
    }
  }, [sectionOffsets, onActive, location.pathname, lenis]);

  useEffect(() => {
    if (active) return;

    handleActiveSection();
  }, [active, handleActiveSection]);

  useEffect(() => {
    window.addEventListener('scroll', handleActiveSection);

    return () => window.removeEventListener('scroll', handleActiveSection);
  }, [handleActiveSection]);

  return (
    <>
      <SpreadMenu active={active} isMounted={isMounted} />
      <SheetMenu active={active} />
    </>
  );
};

export { Navbar };
