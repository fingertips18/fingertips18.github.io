import type { ReactNode } from 'react';

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/shadcn/tooltip';

interface HintProps {
  children: ReactNode;
  label: string;
  asChild?: boolean;
  side?: 'top' | 'right' | 'bottom' | 'left';
  align?: 'center' | 'end' | 'start';
}

export function Hint({ children, label, asChild, side, align }: HintProps) {
  return (
    <Tooltip delayDuration={0}>
      <TooltipTrigger asChild={asChild}>{children}</TooltipTrigger>
      <TooltipContent side={side} align={align}>
        <p>{label}</p>
      </TooltipContent>
    </Tooltip>
  );
}
