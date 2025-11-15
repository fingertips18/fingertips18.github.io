import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Back } from '@/components/common/back';
import { Button } from '@/components/shadcn/button';
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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/shadcn/select';
import { Textarea } from '@/components/shadcn/textarea';
import { cn } from '@/lib/utils';

import { Preview } from './preview';
import { Tags } from './tags';

const ProjectType = {
  web: 'web',
  mobile: 'mobile',
  game: 'game',
} as const;

const MAX_BYTES = 10 * 1024 * 1024; // 10MB

const formSchema = z.object({
  preview: z
    .custom<FileList>((v) => v instanceof FileList, {
      error: 'No image provided.',
    })
    .refine((files) => files.length === 1, {
      error: 'Exactly one image must be uploaded.',
    })
    .refine((files) => files[0]?.type === 'image/webp', {
      error: 'Only .webp images are allowed.',
    })
    .refine((files) => files[0]?.size <= MAX_BYTES, {
      error: 'Image must be less than 10MB',
    }),
  blurhash: z
    .string()
    .min(1, { message: 'Blurhash cannot be empty.' })
    .refine(
      (value) => value.startsWith('data:image/') && value.includes(';base64,'),
      {
        message: 'Invalid preview format. Expected a base64 data URL.',
      },
    ),
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
  tags: z
    .array(
      z.string().min(1, {
        error: 'Tag item cannot be empty',
      }),
    )
    .min(1, { error: 'At least one tag item is required' }),
  type: z.enum(ProjectType, {
    error: 'Please select a valid project type.',
  }),
  link: z.url({ error: 'Please provide a valid URL.' }),
});

type Schema = z.infer<typeof formSchema>;

export function Form() {
  const form = useForm<Schema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      preview: undefined,
      blurhash: '',
      title: '',
      subTitle: '',
      description: '',
      tags: [],
      type: undefined,
      link: '',
    },
  });

  const onSubmit = (values: Schema) => {
    console.log(values);
  };

  return (
    <ShadcnForm {...form}>
      <form
        onSubmit={(e) => void form.handleSubmit(onSubmit)(e)}
        className='flex-1 space-y-6'
      >
        <Preview
          control={form.control}
          name='preview'
          maxSize={MAX_BYTES}
          onBlurhashChange={(blurhash: string) =>
            form.setValue('blurhash', blurhash)
          }
          hasError={
            !!form.formState.errors.preview || !!form.formState.errors.blurhash
          }
        />

        <div className='flex-center max-lg:flex-col gap-x-4 gap-y-6'>
          <FormField
            control={form.control}
            name='title'
            render={({ field }) => (
              <FormItem className='w-full'>
                <FormLabel>Title</FormLabel>
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

        <Tags
          control={form.control}
          name='tags'
          hasError={
            Array.isArray(form.formState.errors.tags)
              ? form.formState.errors.tags.length > 0
              : !!form.formState.errors.tags
          }
        />

        <FormField
          control={form.control}
          name='type'
          render={({ field }) => (
            <FormItem className='w-full'>
              <FormLabel>Type</FormLabel>
              <FormDescription>
                Select the main platform or category your project belongs to.
              </FormDescription>
              <FormControl>
                <Select
                  onValueChange={field.onChange}
                  value={field.value}
                  name={field.name}
                  disabled={field.disabled}
                >
                  <SelectTrigger
                    className={cn(
                      'w-full',
                      form.formState.errors.type && 'border-destructive',
                    )}
                  >
                    <SelectValue placeholder='Select a type' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>Project Type</SelectLabel>
                      {Object.values(ProjectType).map((t) => (
                        <SelectItem key={t} value={t} className='capitalize'>
                          {t}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='link'
          render={({ field }) => (
            <FormItem className='w-full'>
              <FormLabel>Link</FormLabel>
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

        <div className='flex-end flex-col-reverse sm:flex-row gap-2'>
          <Back
            type='button'
            variant='outline'
            label='Cancel'
            withIcon={false}
            className='w-full sm:w-fit'
          >
            Cancel
          </Back>
          <Button type='submit' className='w-full sm:w-fit cursor-pointer'>
            Submit
          </Button>
        </div>
      </form>
    </ShadcnForm>
  );
}
