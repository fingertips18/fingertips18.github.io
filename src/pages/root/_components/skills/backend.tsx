import { BACKEND } from "@/constants/skills";

import { SkillIcon } from "./skill-icon";

const Backend = () => {
  return (
    <div className="max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused w-max">
        {BACKEND.concat(BACKEND).map((b) => (
          <SkillIcon key={`${b.label}-1`} Icon={b.icon} hexColor={b.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Backend };
