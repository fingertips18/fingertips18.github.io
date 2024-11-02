import { useVisibility } from "@/lib/hooks/useVisibility";
import { BACKEND } from "@/constants/skills";
import { cn } from "@/lib/utils";

import { SkillIcon } from "../../../../components/common/skill-icon";

const Backend = () => {
  const { isVisible } = useVisibility();

  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul
        className={cn(
          "flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused w-max",
          !isVisible && "paused"
        )}
      >
        {BACKEND.map((b, i) => (
          <SkillIcon
            key={`backend-${b.label}-${i}`}
            Icon={b.icon}
            hexColor={b.hexColor}
          />
        ))}
        {BACKEND.map((b, i) => (
          <SkillIcon
            key={`backend-${b.label}-${i}`}
            Icon={b.icon}
            hexColor={b.hexColor}
            ariaHidden="true"
          />
        ))}
      </ul>
    </div>
  );
};

export { Backend };
