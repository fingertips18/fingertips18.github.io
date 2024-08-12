import { OTHERS } from "@/constants/collections";

import { SkillIcon } from "./skill-icon";

const Others = () => {
  return (
    <div className="flex gap-x-4 max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll group-hover:paused">
        {OTHERS.map((o) => (
          <SkillIcon key={`${o.label}-1`} Icon={o.icon} hexColor={o.hexColor} />
        ))}
      </ul>
      <ul
        className="flex gap-x-4 animate-loop-scroll group-hover:paused"
        aria-hidden="true"
      >
        {OTHERS.map((o) => (
          <SkillIcon key={`${o.label}-2`} Icon={o.icon} hexColor={o.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Others };
