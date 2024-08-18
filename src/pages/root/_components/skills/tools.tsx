import { TOOLS } from "@/constants/skills";

import { SkillIcon } from "./skill-icon";

const Tools = () => {
  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused w-max">
        {TOOLS.concat(TOOLS).map((t) => (
          <SkillIcon key={`${t.label}-1`} Icon={t.icon} hexColor={t.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Tools };
