import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Combobox } from '@/components/common/combobox';
import { ImageUploader } from '@/components/common/image-uploader';
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

const ProjectType = {
  web: 'web',
  mobile: 'mobile',
  game: 'game',
} as const;

const MAX_BYTES = 10 * 1024 * 1024; // 10MB

const formSchema = z.object({
  preview: z
    .custom<FileList>((v) => v instanceof FileList, {
      error: 'No images provided.',
    })
    .refine((files) => files.length === 1, {
      error: 'Exactly one image must be uploaded.',
    })
    .refine((files) => files[0]?.type === 'image/webp', {
      error: 'Only .webp images allowed',
    })
    .refine((files) => files[0]?.size <= MAX_BYTES, {
      error: 'Image must be less than 10MB',
    }),
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
  stack: z
    .array(
      z.string().min(1, {
        error: 'Stack item cannot be empty',
      }),
    )
    .min(1, { error: 'At least one stack item is required' }),
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
      title: '',
      subTitle: '',
      description: '',
      stack: [],
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
        <FormField
          control={form.control}
          name='preview'
          render={({ field }) => {
            const { value, onChange, onBlur, ...fields } = field;

            return (
              <FormItem className='w-full'>
                <FormLabel>Preview</FormLabel>
                <FormDescription>
                  Provide a preview image for your project.
                </FormDescription>
                <FormControl>
                  <ImageUploader
                    value={value}
                    onChange={(files) => onChange(files)}
                    onBlur={onBlur}
                    {...fields}
                    maxFiles={1}
                    maxSize={MAX_BYTES}
                    className='h-[312px]'
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            );
          }}
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

        <FormField
          control={form.control}
          name='stack'
          render={({ field }) => (
            <FormItem className='w-full'>
              <FormLabel>Stack</FormLabel>
              <FormDescription>
                Add the technologies, frameworks, or languages used in your
                project.
              </FormDescription>
              <FormControl>
                <Combobox
                  placeholder='e.g. ts, js, go, ruby, c#'
                  defaultSuggestions={['js', 'ts', 'go', 'c#', 'c++']}
                  emptyMessage='No stack found.'
                  selectPlaceholder='Select tech stack...'
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
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
                  <SelectTrigger className='w-full'>
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
      </form>
    </ShadcnForm>
  );
}
