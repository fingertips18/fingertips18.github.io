import { Moon, Sun } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/common/components/shadcn/dropdown-menu";
import { Button } from "@/common/components/shadcn/button";
import { useResize } from "@/lib/hooks/use-resize";
import { useTheme } from "@/lib/hooks/use-theme";
import { Hint } from "@/common/components/hint";

const ModeToggle = () => {
  const { setTheme } = useTheme();
  const { width } = useResize();

  const lg = width > 1024;

  return (
    <DropdownMenu modal={false}>
      <Hint asChild label="Mode">
        <DropdownMenuTrigger asChild>
          <Button
            variant={lg ? "ghost" : "default"}
            size="icon"
            className="rounded-full outline-none border-none focus-visible:border-none 
            focus-visible:ring-0 focus-visible:ring-transparent focus-visible:ring-offset-0 
            hover:drop-shadow-primary-glow"
          >
            <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span className="sr-only">Toggle theme</span>
          </Button>
        </DropdownMenuTrigger>
      </Hint>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme("light")}>
          Light
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("dark")}>
          Dark
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("system")}>
          System
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export { ModeToggle };
