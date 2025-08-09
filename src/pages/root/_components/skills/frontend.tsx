import { SkillIcon } from "@/components/common/skill-icon";
import { FRONTEND } from "@/constants/skills";
import { useVisibility } from "@/lib/hooks/useVisibility";
import { cn } from "@/lib/utils";

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
        {FRONTEND.map((f, i) => (
          <SkillIcon
            key={`frontend-${f.label}-${i}`}
            Icon={f.icon}
            hexColor={f.hexColor}
          />
        ))}
        {FRONTEND.map((f, i) => (
          <SkillIcon
            key={`frontend-${f.label}-${i}`}
            Icon={f.icon}
            hexColor={f.hexColor}
            ariaHidden="true"
          />
        ))}
      </ul>
    </div>
  );
};

export { Frontend };
