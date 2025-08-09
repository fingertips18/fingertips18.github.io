import { useVisibility } from "@/lib/hooks/useVisibility";
import { TOOLS } from "@/constants/skills";
import { cn } from "@/lib/utils";

import { SkillIcon } from "@/components/common/skill-icon";

const Tools = () => {
  const { isVisible } = useVisibility();

  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul
        className={cn(
          "flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused w-max",
          !isVisible && "paused"
        )}
      >
        {TOOLS.map((t, i) => (
          <SkillIcon
            key={`tools-${t.label}-${i}`}
            Icon={t.icon}
            hexColor={t.hexColor}
          />
        ))}
        {TOOLS.map((t, i) => (
          <SkillIcon
            key={`tools-${t.label}-${i}`}
            Icon={t.icon}
            hexColor={t.hexColor}
            ariaHidden="true"
          />
        ))}
      </ul>
    </div>
  );
};

export { Tools };
