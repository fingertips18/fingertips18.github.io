import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Back } from '@/components/common/back';
import { Button } from '@/components/shadcn/button';
import { Form as BaseForm } from '@/components/shadcn/form';
import { MAX_BYTES } from '@/constants/sizes';
import { ProjectType } from '@/types/project';

import { Description } from './description';
import { Link } from './link';
import { Preview } from './preview';
import { Subtitle } from './subtitle';
import { Tags } from './tags';
import { Title } from './title';
import { Type } from './type';

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
    <BaseForm {...form}>
      <form
        onSubmit={(e) => void form.handleSubmit(onSubmit)(e)}
        className='flex-1 space-y-6'
      >
        <Preview
          control={form.control}
          name='preview'
          onBlurhashChange={(blurhash: string) =>
            form.setValue('blurhash', blurhash)
          }
          previewHasError={!!form.formState.errors.preview}
          blurError={form.formState.errors.blurhash?.message}
        />

        <div className='flex-center flex-col xl:flex-row gap-x-4 gap-y-6'>
          <Title control={form.control} name='title' />
          <Subtitle control={form.control} name='subTitle' />
        </div>

        <Description control={form.control} name='description' />

        <Tags
          control={form.control}
          name='tags'
          hasError={!!form.formState.errors.tags}
        />

        <Type
          control={form.control}
          name='type'
          hasError={!!form.formState.errors.type}
        />

        <Link control={form.control} name='link' />

        <div className='flex-end flex-col-reverse sm:flex-row gap-2'>
          <Back
            type='button'
            variant='outline'
            label='Cancel'
            withIcon={false}
            className='w-full sm:w-fit'
          />
          <Button type='submit' className='w-full sm:w-fit cursor-pointer'>
            Submit
          </Button>
        </div>
      </form>
    </BaseForm>
  );
}
