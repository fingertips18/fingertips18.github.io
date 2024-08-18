import { FRONTEND } from "@/constants/skills";

import { SkillIcon } from "./skill-icon";

const Frontend = () => {
  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll group-hover:paused w-max">
        {FRONTEND.concat(FRONTEND).map((f) => (
          <SkillIcon key={`${f.label}-1`} Icon={f.icon} hexColor={f.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Frontend };
