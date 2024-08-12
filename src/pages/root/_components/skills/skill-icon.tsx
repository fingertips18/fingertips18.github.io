import { IconType } from "@icons-pack/react-simple-icons";
import { useState } from "react";

import { cn } from "@/lib/utils";

interface SkillIconProps {
  Icon: IconType;
  hexColor: string;
}

const SkillIcon = ({ Icon, hexColor }: SkillIconProps) => {
  const [hovered, setHovered] = useState(false);

  return (
    <li
      className={cn(
        "rounded-full p-4 border",
        hovered ? "border-foreground/15" : "border-border"
      )}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <Icon
        color={hovered ? hexColor : undefined}
        className={cn("w-6 lg:w-12 h-6 lg:h-12", !hovered && "opacity-50")}
      />
    </li>
  );
};

export { SkillIcon };
