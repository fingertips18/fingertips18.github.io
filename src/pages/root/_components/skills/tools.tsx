import { TOOLS } from "@/constants/collections";

import { SkillIcon } from "./skill-icon";

const Tools = () => {
  return (
    <div className="flex gap-x-4 max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused">
        {TOOLS.map((t) => (
          <SkillIcon key={`${t.label}-1`} Icon={t.icon} hexColor={t.hexColor} />
        ))}
      </ul>
      <ul
        className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused"
        aria-hidden="true"
      >
        {TOOLS.map((t) => (
          <SkillIcon key={`${t.label}-2`} Icon={t.icon} hexColor={t.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Tools };
