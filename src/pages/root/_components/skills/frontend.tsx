import { useVisibility } from "@/lib/hooks/useVisibility";
import { FRONTEND } from "@/constants/skills";
import { cn } from "@/lib/utils";

import { SkillIcon } from "./skill-icon";

const Frontend = () => {
  const { isVisible } = useVisibility();

  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul
        className={cn(
          "flex gap-x-4 animate-loop-scroll group-hover:paused w-max",
          !isVisible && "paused"
        )}
      >
        {FRONTEND.concat(FRONTEND).map((f, i) => (
          <SkillIcon
            key={`frontend-${f.label}-${i}`}
            Icon={f.icon}
            hexColor={f.hexColor}
          />
        ))}
      </ul>
    </div>
  );
};

export { Frontend };
