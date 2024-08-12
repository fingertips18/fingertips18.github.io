import { FRONTEND } from "@/constants/collections";

import { SkillIcon } from "./skill-icon";

const Frontend = () => {
  return (
    <div className="flex gap-x-4 max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll group-hover:paused">
        {FRONTEND.map((f) => (
          <SkillIcon key={`${f.label}-1`} Icon={f.icon} hexColor={f.hexColor} />
        ))}
      </ul>
      <ul
        className="flex gap-x-4 animate-loop-scroll group-hover:paused"
        aria-hidden="true"
      >
        {FRONTEND.map((f) => (
          <SkillIcon key={`${f.label}-2`} Icon={f.icon} hexColor={f.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Frontend };
