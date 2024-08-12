import { BACKEND } from "@/constants/collections";

import { SkillIcon } from "./skill-icon";

const Backend = () => {
  return (
    <div className="flex gap-x-4 max-w-screen-lg overflow-hidden group">
      <ul className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused">
        {BACKEND.map((b) => (
          <SkillIcon key={`${b.label}-1`} Icon={b.icon} hexColor={b.hexColor} />
        ))}
      </ul>
      <ul
        className="flex gap-x-4 animate-loop-scroll direction-reverse group-hover:paused"
        aria-hidden="true"
      >
        {BACKEND.map((b) => (
          <SkillIcon key={`${b.label}-2`} Icon={b.icon} hexColor={b.hexColor} />
        ))}
      </ul>
    </div>
  );
};

export { Backend };
