import { Link } from "react-router-dom";

import { SOCIALS } from "@/constants/collections";
import { Hint } from "@/components/common/hint";
import { cn } from "@/lib/utils";

interface SocialButtonsProps {
  isMounted: boolean;
}

const SocialButtons = ({ isMounted }: SocialButtonsProps) => {
  return (
    <ul
      className={cn(
        "flex-center gap-x-6 transition-opacity duration-500 ease-in-out",
        isMounted ? "opacity-100" : "opacity-0"
      )}
    >
      {SOCIALS.map((s) => {
        const Icon = s.icon;

        return (
          <Hint key={s.href} asChild label={s.label} side="top">
            <Link
              to={s.href}
              target="_blank"
              className="rounded-full border border-muted-foreground p-2.5
              hover:scale-105 hover:-translate-y-2 transition-all
              ease-in-out cursor-pointer hover:bg-muted-foreground group hover:drop-shadow-foreground-glow"
            >
              <Icon className="w-4 h-4 transition-colors ease-in-out group-hover:text-background pointer-events-none" />
            </Link>
          </Hint>
        );
      })}
    </ul>
  );
};

export default SocialButtons;
