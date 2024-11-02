import { useLocation } from "react-router-dom";
import { useLenis } from "lenis/react";

import { ROOTMENU } from "@/constants/collections";
import { cn } from "@/lib/utils";

interface SpreadMenuProps {
  active: string;
  isMounted: boolean;
}

const SpreadMenu = ({ active, isMounted }: SpreadMenuProps) => {
  const lenis = useLenis();
  const location = useLocation();

  const onClick = (id: string) => {
    const section = document.getElementById(id);
    if (section) {
      lenis?.scrollTo(section);
    }
  };

  return (
    <nav
      className={cn(
        "hidden lg:flex-center px-4 flex-grow transition-opacity duration-500 ease-in-out",
        location.pathname === "/"
          ? "opacity-100"
          : "opacity-0 pointer-events-none"
      )}
    >
      <ul
        className={cn(
          "flex-center gap-x-10 transition-opacity duration-1000 ease-in-out",
          isMounted ? "opacity-100" : "opacity-0"
        )}
      >
        {ROOTMENU.map((m, i) => (
          <li
            key={`${m.label}-${i}`}
            className={cn(
              "capitalize text-sm font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow hover:text-accent",
              active === m.label && "text-accent"
            )}
            onClick={() => onClick(m.label)}
          >
            {m.label}
          </li>
        ))}
      </ul>
    </nav>
  );
};

export { SpreadMenu };
