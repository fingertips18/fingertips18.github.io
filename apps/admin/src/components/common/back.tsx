import { MoveLeft } from 'lucide-react';
import type { ComponentProps } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { cn } from '@/lib/utils';
import { Route } from '@/routes/route';

interface BackProps
  extends Omit<ComponentProps<typeof Button>, 'className' | 'onClick'> {
  label?: string;
  className?: string;
  onBack?: () => void;
  withIcon?: boolean;
}

export function Back({
  label,
  className,
  onBack,
  variant,
  withIcon = true,
  ...props
}: BackProps) {
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

    onBack?.();
  };

  return (
    <Button
      variant={variant || 'ghost'}
      onClick={handleBack}
      className={cn('gap-x-2 cursor-pointer', className)}
      {...props}
    >
      {withIcon && <MoveLeft aria-hidden='true' />}
      {label ?? 'Back'}
    </Button>
  );
}
