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

interface LinkProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Link<T extends FieldValues>({ control, name }: LinkProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Link</FormLabel>
          <FormDescription>
            Enter the full URL for the deployed project (e.g.,
            https://example.com)
          </FormDescription>
          <FormControl>
            <Input placeholder='https://example.com' {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
