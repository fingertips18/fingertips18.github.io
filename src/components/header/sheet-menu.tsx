import { Link, useLocation } from "react-router-dom";
import { LucideMenu, MoveLeft } from "lucide-react";
import { useEffect, useState } from "react";
import { useLenis } from "lenis/react";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/shadcn/sheet";
import { Button } from "@/components/shadcn/button";
import { ROOTMENU } from "@/constants/collections";
import { useResize } from "@/lib/hooks/useResize";
import { AppRoutes } from "@/routes/app-routes";
import { Hint } from "@/components/common/hint";
import { cn } from "@/lib/utils";

import { ModeToggle } from "./mode-toggle";

interface SheetMenuProps {
  active: string;
}

const SheetMenu = ({ active }: SheetMenuProps) => {
  const lenis = useLenis();
  const location = useLocation();
  const [open, setOpen] = useState(false);
  const { width } = useResize();

  useEffect(() => {
    if (width > 1024) {
      setOpen(false);
    }
  }, [width]);

  const onOpenChange = (open: boolean) => {
    if (!lenis) return;

    if (open) {
      lenis.stop();
    } else {
      lenis.start();
    }
  };

  const onClick = (id: string) => {
    const section = document.getElementById(id);
    if (section) {
      section.scrollIntoView({
        behavior: "smooth",
        block: "start",
        inline: "nearest",
      });
    }
  };

  return (
    <Sheet
      open={open}
      onOpenChange={(open) => {
        onOpenChange(open);
        setOpen(open);
      }}
    >
      <Hint asChild label="Menu">
        <SheetTrigger asChild>
          <Button
            variant={"ghost"}
            size={"icon"}
            aria-label="menu-toggle"
            className="lg:hidden hover:drop-shadow-primary-glow"
          >
            <LucideMenu className="w-5 h-5" />
          </Button>
        </SheetTrigger>
      </Hint>
      <SheetContent data-lenis-prevent className="overflow-y-auto no-scrollbar">
        <SheetHeader className="mt-4 !items-start">
          <SheetTitle className="text-sm">Menu</SheetTitle>
          <SheetDescription className="text-xs text-start">
            Discover my portfolio, skills, projects, and how to connect.
          </SheetDescription>
        </SheetHeader>

        <nav className="w-full flex justify-end mt-10 flex-1">
          {location.pathname === "/" ? (
            <ul className="space-y-6 text-end">
              {ROOTMENU.map((m, i) => (
                <li
                  key={`${m.label}-${i}`}
                  className={cn(
                    "capitalize font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow lg:hover:text-accent",
                    active === m.label && "text-accent"
                  )}
                  onClick={() => onClick(m.label)}
                >
                  {m.label}
                </li>
              ))}
            </ul>
          ) : (
            <Link
              to={AppRoutes.root}
              className="flex items-center gap-x-2 hover:text-accent"
            >
              <MoveLeft className="size-4" /> Go home
            </Link>
          )}
        </nav>

        <SheetFooter className="fixed bottom-4 right-4">
          <ModeToggle />
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};

export { SheetMenu };
