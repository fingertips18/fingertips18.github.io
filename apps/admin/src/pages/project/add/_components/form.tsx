import { zodResolver } from '@hookform/resolvers/zod';
import { Loader } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { z } from 'zod';

import { Back } from '@/components/common/back';
import { Button } from '@/components/shadcn/button';
import { Form as BaseForm } from '@/components/shadcn/form';
import { MAX_BYTES } from '@/constants/sizes';
import { useUnsavedChanges } from '@/hooks/useUnsavedChanges';
import { toast } from '@/lib/toast';
import { Route } from '@/routes/route';
import { ImageService } from '@/services/image';
import { ProjectService } from '@/services/project';
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
  const [imageLoading, setImageLoading] = useState<boolean>(false);
  const [projectLoading, setProjectLoading] = useState<boolean>(false);
  const [submitted, setSubmitted] = useState<boolean>(false);
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
  const abortRef = useRef<AbortController | null>(null);
  const navigate = useNavigate();
  useUnsavedChanges({
    hasUnsavedChanges: !submitted && form.formState.isDirty,
  });

  useEffect(() => {
    abortRef.current = new AbortController();

    return () => abortRef.current?.abort();
  }, []);

  const onSubmit = async (values: Schema) => {
    let imageURL: string | null;

    try {
      setImageLoading(true);

      const preview = values.preview[0];

      const url = await ImageService.upload({
        file: preview,
        signal: abortRef.current?.signal,
      });
      if (!url) {
        throw new Error('Image URL undefined');
      }

      imageURL = url;
      toast({
        level: 'success',
        title: 'Image upload complete 🎉',
        description: `${preview.name} uploaded successfully!`,
      });
    } catch {
      imageURL = null;
      toast({
        level: 'error',
        title: 'Upload failed',
        description: 'We couldn’t upload your image. Please try again.',
      });
    } finally {
      setImageLoading(false);
    }

    if (!imageURL) {
      // Image upload already surfaced an error toast; skip project creation.
      return;
    }

    // Proceed with the rest of the form submission using imageURL

    try {
      setProjectLoading(true);

      const projectId = await ProjectService.create({
        project: {
          preview: imageURL,
          blurhash: values.blurhash,
          title: values.title,
          subTitle: values.subTitle,
          description: values.description,
          tags: values.tags,
          type: values.type,
          link: values.link,
        },
        signal: abortRef.current?.signal,
      });

      if (!projectId) {
        throw new Error('Project ID undefined');
      }

      toast({
        level: 'success',
        title: 'Project upload complete 🎉',
        description: `${values.title} uploaded successfully!`,
      });

      setSubmitted(true);
      form.reset();
      void navigate(Route.project);
    } catch {
      toast({
        level: 'error',
        title: 'Upload failed',
        description: 'We couldn’t upload your project. Please try again.',
      });
    } finally {
      setProjectLoading(false);
    }
  };

  const loading = imageLoading || projectLoading;
  const previewIsEmpty =
    !form.getValues('preview') || form.getValues('preview').length === 0;
  const ariaLabel = imageLoading
    ? 'Uploading image, please wait'
    : projectLoading
    ? 'Creating project, please wait'
    : 'Submit';

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
          errors={form.formState.errors}
          isEmpty={previewIsEmpty}
          disabled={loading}
        />

        <div className='flex-center flex-col xl:flex-row gap-x-4 gap-y-6'>
          <Title control={form.control} name='title' disabled={loading} />
          <Subtitle control={form.control} name='subTitle' disabled={loading} />
        </div>

        <Description
          control={form.control}
          name='description'
          disabled={loading}
        />

        <Tags
          control={form.control}
          name='tags'
          errors={form.formState.errors}
          disabled={loading}
        />

        <Type
          control={form.control}
          name='type'
          errors={form.formState.errors}
          disabled={loading}
        />

        <Link control={form.control} name='link' disabled={loading} />

        <div className='flex-end flex-col-reverse sm:flex-row gap-2'>
          <Back
            type='button'
            variant='outline'
            label='Cancel'
            withIcon={false}
            disabled={loading}
            className='w-full sm:w-fit'
          />
          <Button
            type='submit'
            disabled={loading}
            aria-label={ariaLabel}
            className='w-full sm:w-fit cursor-pointer min-w-[78.85px]'
          >
            {loading ? (
              <Loader aria-hidden='true' className='size-4 animate-spin' />
            ) : (
              'Submit'
            )}
          </Button>
        </div>
      </form>
    </BaseForm>
  );
}
