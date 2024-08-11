import { LucideMenu } from "lucide-react";

import { Button } from "@/common/components/shadcn/button";
import { Hint } from "@/common/components/hint";
import { MENU } from "@/constants/collections";

const Menu = () => {
  return (
    <nav className="flex-center lg:px-4">
      <ul className="hidden lg:flex-center gap-x-10">
        {MENU.map((m, i) => (
          <li
            key={`${m.label}-${i}`}
            className="capitalize text-sm font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-glow hover:text-accent"
          >
            {m.label}
          </li>
        ))}
      </ul>

      <Hint asChild label="Menu">
        <Button
          variant={"ghost"}
          size={"icon"}
          className="lg:hidden hover:drop-shadow-glow"
        >
          <LucideMenu className="w-6 h-6" />
        </Button>
      </Hint>
    </nav>
  );
};

export { Menu };
