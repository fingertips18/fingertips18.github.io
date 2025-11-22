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

interface SubtitleProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  disabled?: boolean;
}

export function Subtitle<T extends FieldValues>({
  control,
  name,
  disabled,
}: SubtitleProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      disabled={disabled}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Subtitle</FormLabel>
          <FormDescription>
            Add a short tagline or secondary title that describes your project.
          </FormDescription>
          <div className={cn(disabled && 'cursor-not-allowed')}>
            <FormControl>
              <Input placeholder='e.g. The Pulse of the Web' {...field} />
            </FormControl>
          </div>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
