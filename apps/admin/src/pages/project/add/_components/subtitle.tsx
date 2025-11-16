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

interface SubtitleProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Subtitle<T extends FieldValues>({
  control,
  name,
}: SubtitleProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Subtitle</FormLabel>
          <FormDescription>
            Add a short tagline or secondary title that describes your project.
          </FormDescription>
          <FormControl>
            <Input placeholder='e.g. The Pulse of the Web' {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
