import { LucideMenu } from "lucide-react";

import { Button } from "@/common/components/shadcn/button";
import { Hint } from "@/common/components/hint";
import { MENUS } from "@/constants/collections";
import { cn } from "@/lib/utils";

const Menu = () => {
  return (
    <nav className="flex-center lg:px-4">
      <ul className="hidden lg:flex-center gap-x-10">
        {MENUS.map((m, i) => (
          <a
            href={m.href}
            key={`${m.label}-${i}`}
            className={cn(
              "capitalize text-sm font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow hover:text-accent",
              m.href.length === 0 && "pointer-events-none text-muted-foreground"
            )}
          >
            {m.label}
          </a>
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
