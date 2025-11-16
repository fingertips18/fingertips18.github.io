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

interface TitleProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Title<T extends FieldValues>({ control, name }: TitleProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Title</FormLabel>
          <FormDescription>
            A short, descriptive name for your project.
          </FormDescription>
          <FormControl>
            <Input placeholder='e.g. WebNexus' {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
