import { useCallback, useEffect, useMemo } from "react";

import { useRootSectionStore } from "@/lib/stores/use-root-section-store";
import { useElementsByQuery } from "@/lib/hooks/use-elements-by-query";
import { useClient } from "@/lib/hooks/use-client";
import { QUERYELEMENT } from "@/constants/enums";

import { SpreadMenu } from "./spread-menu";
import { SheetMenu } from "./sheet-menu";

const Navbar = () => {
  const { active, onActive } = useRootSectionStore((state) => state);
  const rootSections = useElementsByQuery(`.${QUERYELEMENT.rootSection}`);
  const isMounted = useClient();

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
