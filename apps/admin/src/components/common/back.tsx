import { MoveLeft } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { Route } from '@/routes/route';

interface BackProps {
  label?: string;
}

export function Back({ label }: BackProps) {
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const handleBack = () => {
    if (window.history.length <= 2) {
      // Split and filter to remove empty parts
      const parts = pathname.split('/').filter(Boolean);

      // Get the previous path by removing the last segment
      const previousPath =
        parts.length > 1
          ? pathname.slice(
              0,
              pathname.lastIndexOf(`/${parts[parts.length - 1]}`),
            )
          : undefined;

      void navigate(previousPath ?? Route.root);
    } else {
      void navigate(-1);
    }
  };

  return (
    <Button
      variant='ghost'
      onClick={handleBack}
      className='gap-x-2 cursor-pointer'
    >
      <MoveLeft aria-hidden='true' />
      {label ?? 'Back'}
    </Button>
  );
}
