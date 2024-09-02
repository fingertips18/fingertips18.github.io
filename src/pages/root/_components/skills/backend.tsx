import { BACKEND } from "@/constants/skills";

import { SkillIcon } from "./skill-icon";

const Backend = () => {
  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused w-max">
        {BACKEND.concat(BACKEND).map((b, i) => (
          <SkillIcon
            key={`backend-${b.label}-${i}`}
            Icon={b.icon}
            hexColor={b.hexColor}
          />
        ))}
      </ul>
    </div>
  );
};

export { Backend };
