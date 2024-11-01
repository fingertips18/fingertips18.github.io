import { Link } from "react-router-dom";

import { useMounted } from "@/lib/hooks/useMounted";
import { Hint } from "@/components/common/hint";
import { CONTACTS } from "@/constants/contact";
import { cn } from "@/lib/utils";

interface OtherContactsProps {
  isVisible: boolean;
}

const OtherContacts = ({ isVisible }: OtherContactsProps) => {
  const isMounted = useMounted();

  return (
    <ul
      className={cn(
        "flex-center gap-x-6 transition-opacity duration-500 ease-in-out",
        isMounted ? "opacity-100" : "opacity-0",
        isVisible ? "opacity-100" : "opacity-0"
      )}
    >
      {CONTACTS.map((c) => {
        const Icon = c.icon;

        return (
          <Hint key={c.href} asChild label={c.label} side="top">
            <Link
              to={c.href}
              target="_blank"
              className="rounded-full border border-muted-foreground p-2.5
                hover:scale-105 hover:-translate-y-2 transition-all
                ease-in-out cursor-pointer hover:bg-muted-foreground group hover:drop-shadow-foreground-glow"
            >
              <Icon className="w-4 h-4 ease-in-out group-hover:text-background pointer-events-none" />
            </Link>
          </Hint>
        );
      })}
    </ul>
  );
};

export { OtherContacts };
