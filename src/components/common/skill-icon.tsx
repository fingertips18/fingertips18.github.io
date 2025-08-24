import { IconType as IconsPackType } from '@icons-pack/react-simple-icons';
import { useState } from 'react';
import { IconType as ReactIconsType } from 'react-icons/lib';

import { cn } from '@/lib/utils';

interface SkillIconProps {
  Icon: IconsPackType | ReactIconsType;
  hexColor: string;
  ariaHidden?: React.AriaAttributes['aria-hidden'];
}

const SkillIcon = ({ Icon, hexColor, ariaHidden }: SkillIconProps) => {
  const [hovered, setHovered] = useState(false);

  return (
    <li
      aria-hidden={ariaHidden}
      className={cn(
        'rounded-full p-4 border bg-foreground/5',
        hovered ? 'border-foreground/15' : 'border-border',
      )}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <Icon
        color={hovered ? hexColor : undefined}
        className={cn('size-6 md:size-12', !hovered && 'opacity-50')}
      />
    </li>
  );
};

export { SkillIcon };
