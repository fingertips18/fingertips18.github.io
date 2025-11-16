import type { Control, FieldValues, Path } from 'react-hook-form';

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { Textarea } from '@/components/shadcn/textarea';

interface DescriptionProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Description<T extends FieldValues>({
  control,
  name,
}: DescriptionProps<T>) {
  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel className='w-fit'>Description</FormLabel>
          <FormDescription>
            Write a brief summary of your project — what it does, its purpose,
            or main features.
          </FormDescription>
          <FormControl>
            <Textarea
              placeholder='e.g. A web app that helps teams manage tasks efficiently.'
              className='resize-none h-32'
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
