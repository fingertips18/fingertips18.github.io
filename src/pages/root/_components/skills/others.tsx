import { useVisibility } from "@/lib/hooks/useVisibility";
import { OTHERS } from "@/constants/skills";
import { cn } from "@/lib/utils";

import { SkillIcon } from "../../../../components/common/skill-icon";

const Others = () => {
  const { isVisible } = useVisibility();

  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul
        className={cn(
          "flex gap-x-4 animate-loop-scroll group-hover:paused w-max",
          !isVisible && "paused"
        )}
      >
        {OTHERS.map((o, i) => (
          <SkillIcon
            key={`others-${o.label}-${i}`}
            Icon={o.icon}
            hexColor={o.hexColor}
          />
        ))}
        {OTHERS.map((o, i) => (
          <SkillIcon
            key={`others-${o.label}-${i}`}
            Icon={o.icon}
            hexColor={o.hexColor}
            ariaHidden="true"
          />
        ))}
      </ul>
    </div>
  );
};

export { Others };
