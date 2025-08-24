import { useLenis } from 'lenis/react';
import { useLocation, useNavigate } from 'react-router-dom';

import { ROOTMENU } from '@/constants/collections';
import { ROOTSECTION } from '@/constants/enums';
import { cn } from '@/lib/utils';
import { AppRoutes } from '@/routes/app-routes';

interface SpreadMenuProps {
  active: ROOTSECTION;
  isMounted: boolean;
}

const SpreadMenu = ({ active, isMounted }: SpreadMenuProps) => {
  const lenis = useLenis();
  const location = useLocation();
  const navigate = useNavigate();

  const onClick = (id: string, hash: string) => {
    if (location.pathname === AppRoutes.root) {
      if (!lenis) return;

      const section = document.getElementById(id);

      if (!section) return;

      lenis.scrollTo(section);
    } else {
      navigate(AppRoutes.root + hash);
    }
  };

  return (
    <nav className='hidden lg:flex items-center justify-center px-4 flex-grow'>
      <ul
        className={cn(
          'flex-center gap-x-10 transition-opacity duration-1000 ease-in-out',
          isMounted ? 'opacity-100' : 'opacity-0',
        )}
      >
        {ROOTMENU.map((m) => (
          <li
            key={m.hash}
            className={cn(
              'capitalize text-sm font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow hover:text-accent',
              active === m.label && 'text-accent',
            )}
            onClick={() => onClick(m.label, m.hash)}
          >
            {m.label}
          </li>
        ))}
      </ul>
    </nav>
  );
};

export { SpreadMenu };
