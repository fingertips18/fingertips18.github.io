import { useLenis } from 'lenis/react';
import { Link } from 'react-router-dom';

import { Skeleton } from '@/components/shadcn/skeleton';
import { DARKLOGO, LIGHTLOGO } from '@/constants/assets';
import { useMounted } from '@/lib/hooks/useMounted';
import { useResize } from '@/lib/hooks/useResize';
import { Theme, useTheme } from '@/lib/hooks/useTheme';
import { AppRoutes } from '@/routes/app-routes';

const Logo = () => {
  const { theme } = useTheme();
  const isMounted = useMounted();
  const lenis = useLenis();
  const { width } = useResize();

  if (!isMounted) {
    return <Skeleton className='w-20 lg:w-28 h-6 lg:h-8' />;
  }

  const logo =
    theme === Theme.dark || theme === Theme.system ? DARKLOGO : LIGHTLOGO;

  const onClick = () => lenis?.scrollTo(0);

  const lg = width > 1024;

  return (
    <Link
      to={AppRoutes.root}
      className='hover:scale-105 transition-all duration-500 ease-in-out hover:drop-shadow-primary-glow'
      onClick={onClick}
    >
      <img
        src={logo}
        alt='Logo'
        width={lg ? 89.9 : 74.91}
        height={lg ? 21.2 : 17.66}
        className='object-contain'
        loading='eager'
      />
    </Link>
  );
};

export { Logo };
