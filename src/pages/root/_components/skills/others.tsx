import { OTHERS } from "@/constants/skills";

import { SkillIcon } from "./skill-icon";

const Others = () => {
  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll group-hover:paused w-max">
        {OTHERS.concat(OTHERS).map((o, i) => (
          <SkillIcon
            key={`others-${o.label}-${i}`}
            Icon={o.icon}
            hexColor={o.hexColor}
          />
        ))}
      </ul>
    </div>
  );
};

export { Others };
