import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import {
  Form as ShadcnForm,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { Input } from '@/components/shadcn/input';
import { Textarea } from '@/components/shadcn/textarea';

const formSchema = z.object({
  title: z
    .string()
    .min(6, {
      error: 'Titles must be at least 6 characters long.',
    })
    .max(50, {
      error: 'Titles must not exceed 50 characters.',
    }),
  subTitle: z
    .string()
    .min(6, {
      error: 'Subtitles must be at least 6 characters long.',
    })
    .max(50, {
      error: 'Subtitles must not exceed 50 characters.',
    }),
  description: z
    .string()
    .min(10, {
      message: 'Description must be at least 10 characters long.',
    })
    .max(300, {
      message: 'Description must not exceed 300 characters.',
    }),
});

type Schema = z.infer<typeof formSchema>;

export function Form() {
  const form = useForm<Schema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: '',
      subTitle: '',
      description: '',
    },
  });

  const onSubmit = (values: Schema) => {
    console.log(values);
  };

  return (
    <ShadcnForm {...form}>
      <form
        onSubmit={void form.handleSubmit(onSubmit)}
        className='flex-1 space-y-6 lg:space-y-8'
      >
        <div className='flex-center max-lg:flex-col gap-x-4 gap-y-6'>
          <FormField
            control={form.control}
            name='title'
            render={({ field }) => (
              <FormItem className='w-full'>
                <FormLabel>Title</FormLabel>
                <FormDescription>
                  Enter a short, descriptive name for your project.
                </FormDescription>
                <FormControl>
                  <Input placeholder='e.g. WebNexus' {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name='subTitle'
            render={({ field }) => (
              <FormItem className='w-full'>
                <FormLabel>Subtitle</FormLabel>
                <FormDescription>
                  Add a short tagline or secondary title that describes your
                  project.
                </FormDescription>
                <FormControl>
                  <Input placeholder='e.g. The Pulse of the Web' {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <FormField
          control={form.control}
          name='description'
          render={({ field }) => (
            <FormItem className='w-full'>
              <FormLabel>Description</FormLabel>
              <FormDescription>
                Write a brief summary of your project — what it does, its
                purpose, or main features.
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
      </form>
    </ShadcnForm>
  );
}
