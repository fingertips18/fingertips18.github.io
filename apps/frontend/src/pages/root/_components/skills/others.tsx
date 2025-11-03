import { SkillIcon } from '@/components/common/skill-icon';
import { OTHERS } from '@/constants/skills';
import { useVisibility } from '@/lib/hooks/useVisibility';
import { cn } from '@/lib/utils';

const Others = () => {
  const { isVisible } = useVisibility();

  return (
    <div className='max-w-screen-lg overflow-hidden group'>
      <ul
        className={cn(
          'flex gap-x-4 animate-loop-scroll group-hover:paused w-max',
          !isVisible && 'paused',
        )}
      >
        {OTHERS.map((o, i) => (
          <SkillIcon
            key={`others-${o.label}-${i}`}
            Icon={o.icon}
            hexColor={o.hexColor}
          />
        ))}
        {OTHERS.map((o, i) => (
          <SkillIcon
            key={`others-${o.label}-${i}`}
            Icon={o.icon}
            hexColor={o.hexColor}
            ariaHidden='true'
          />
        ))}
      </ul>
    </div>
  );
};

export { Others };
