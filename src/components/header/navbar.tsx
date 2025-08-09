import { useCallback, useEffect, useMemo } from "react";

import { QUERYELEMENT } from "@/constants/enums";
import { useElementsByQuery } from "@/lib/hooks/useElementsByQuery";
import { useMounted } from "@/lib/hooks/useMounted";
import { useRootSectionStore } from "@/lib/stores/useRootSectionStore";

import { SheetMenu } from "./sheet-menu";
import { SpreadMenu } from "./spread-menu";

const Navbar = () => {
  const { active, onActive } = useRootSectionStore((state) => state);
  const rootSections = useElementsByQuery(`.${QUERYELEMENT.rootSection}`);
  const isMounted = useMounted();

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
    if (!sectionOffsets) return;

    for (let i = 0; i < sectionOffsets.length; i++) {
      if (window.scrollY >= sectionOffsets[i].offset) {
        onActive(sectionOffsets[i].id);
      }
    }
  }, [sectionOffsets, onActive]);

  useEffect(() => {
    window.addEventListener("scroll", handleActiveSection);

    return () => window.removeEventListener("scroll", handleActiveSection);
  }, [handleActiveSection]);

  return (
    <>
      <SpreadMenu active={active} isMounted={isMounted} />
      <SheetMenu active={active} />
    </>
  );
};

export { Navbar };
