import type { Control, FieldValues, Path } from 'react-hook-form';

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { Input } from '@/components/shadcn/input';
import { cn } from '@/lib/utils';

interface LinkProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  disabled?: boolean;
}

export function Link<T extends FieldValues>({
  control,
  name,
  disabled,
}: LinkProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      disabled={disabled}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Link</FormLabel>
          <FormDescription>
            Enter the full URL for the deployed project (e.g.,
            https://example.com)
          </FormDescription>
          <div className={cn(disabled && 'cursor-not-allowed')}>
            <FormControl>
              <Input placeholder='https://example.com' {...field} />
            </FormControl>
          </div>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
