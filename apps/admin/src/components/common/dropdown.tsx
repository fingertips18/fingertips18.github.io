import type { ComponentProps, ReactElement } from 'react';

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/shadcn/select';
import { cn } from '@/lib/utils';

interface DropdownProps extends ComponentProps<typeof Select> {
  id?: string;
  label: string;
  placeholder?: string;
  hasError?: boolean;
  className?: string;
  children: ReactElement<typeof SelectItem> | ReactElement<typeof SelectItem>[];
}

export function Dropdown({
  id,
  label,
  placeholder,
  hasError,
  className,
  children,
  ...props
}: DropdownProps) {
  return (
    <Select {...props}>
      <SelectTrigger
        id={id}
        className={cn(
          'w-full data-[state=open]:border-ring data-[state=open]:ring-ring/50 data-[state=open]:ring-[3px]',
          hasError && 'border-destructive',
          className,
        )}
      >
        <SelectValue placeholder={placeholder || 'Select...'} />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>{label}</SelectLabel>
          {children}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
