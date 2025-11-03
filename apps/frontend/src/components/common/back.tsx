import { MoveLeft } from 'lucide-react';
import { Link, To } from 'react-router-dom';

import { ROOTSECTION } from '@/constants/enums';
import { useRootSectionStore } from '@/lib/stores/useRootSectionStore';

interface BackProps {
  section?: ROOTSECTION;
  to: To;
  label?: string;
}

const Back = ({
  section = ROOTSECTION.about,
  to,
  label = 'Back',
}: BackProps) => {
  const { onActive } = useRootSectionStore();

  const handleBack = () => {
    onActive(section);
  };

  return (
    <Link
      to={to}
      onClick={handleBack}
      className='flex items-center gap-x-2 hover:text-accent text-sm w-fit'
    >
      <MoveLeft className='size-4' /> {label}
    </Link>
  );
};

export { Back };
