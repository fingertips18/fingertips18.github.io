import { useCallback, useEffect, useMemo } from "react";
import { LucideMenu } from "lucide-react";
import { useLenis } from "lenis/react";

import { useRootSectionStore } from "@/lib/stores/use-root-section-store";
import { useElementsByQuery } from "@/lib/hooks/use-elements-by-query";
import { Button } from "@/common/components/shadcn/button";
import { ROOTMENU } from "@/constants/collections";
import { QUERYELEMENTS } from "@/constants/enums";
import { Hint } from "@/common/components/hint";
import { cn } from "@/lib/utils";

const Menu = () => {
  const { active, onActive } = useRootSectionStore((state) => state);
  const rootSections = useElementsByQuery(`.${QUERYELEMENTS.rootSection}`);
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

  const onClick = (id: string) => {
    const section = document.getElementById(id);
    if (section) {
      lenis?.scrollTo(section);
    }
  };

  return (
    <nav className="flex-center lg:px-4">
      <ul className="hidden lg:flex-center gap-x-10">
        {ROOTMENU.map((m, i) => (
          <li
            key={`${m.label}-${i}`}
            className={cn(
              "capitalize text-sm font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow hover:text-accent",
              active === m.label && "text-accent",
              m.id.length === 0 && "pointer-events-none text-muted-foreground"
            )}
            onClick={() => onClick(m.label)}
          >
            {m.label}
          </li>
        ))}
      </ul>

      <Hint asChild label="Menu">
        <Button
          variant={"ghost"}
          size={"icon"}
          className="lg:hidden hover:drop-shadow-primary-glow"
        >
          <LucideMenu className="w-6 h-6" />
        </Button>
      </Hint>
    </nav>
  );
};

export { Menu };
