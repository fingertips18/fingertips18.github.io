import { Layers2, Loader } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Blurhash } from 'react-blurhash';
import type { Control, FieldErrors, FieldValues, Path } from 'react-hook-form';

import { ImageUploader } from '@/components/common/image-uploader';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { MAX_BYTES } from '@/constants/sizes';
import { encodeImageToBlurhash, fileToImage } from '@/lib/image';
import { toast } from '@/lib/toast';
import { cn } from '@/lib/utils';

interface PreviewProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  onBlurhashChange: (blurhash: string) => void;
  errors: FieldErrors<T>;
  isEmpty?: boolean;
  disabled?: boolean;
}

export function Preview<T extends FieldValues>({
  control,
  name,
  onBlurhashChange,
  errors,
  isEmpty,
  disabled,
}: PreviewProps<T>) {
  const [blurHash, setBlurHash] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const previewHasError = !!errors[name];
  const blurhashError = errors['blurhash'] as { message?: string } | undefined;
  const blurHasError = !!blurhashError;

  useEffect(() => {
    if (isEmpty) {
      setBlurHash(null);
    }
  }, [isEmpty]);

  return (
    <FormField
      control={control}
      name={name}
      disabled={disabled}
      render={({ field }) => {
        const { onChange, ...fields } = field;

        return (
          <div className='flex-center flex-col md:flex-row gap-x-4 gap-y-6 w-full'>
            <FormItem className='w-full'>
              <FormLabel className='w-fit'>Preview</FormLabel>
              <FormDescription>
                Provide a preview image for your project.
              </FormDescription>
              <FormControl>
                <ImageUploader
                  onChange={async (files) => {
                    setBlurHash(null); // reset previous image
                    setLoading(files.length > 0); // start loading if there’s a file

                    onChange(files);

                    if (files.length === 0) {
                      onBlurhashChange('');
                      setLoading(false);
                      return;
                    }

                    try {
                      const file = files[0];
                      const image = await fileToImage(file);
                      const blurhash = await encodeImageToBlurhash(image);
                      setBlurHash(blurhash || null);
                      onBlurhashChange(blurhash || '');
                    } catch {
                      setBlurHash(null);
                      onBlurhashChange('');
                      onChange(new DataTransfer().files);

                      toast({
                        level: 'error',
                        title: 'Blurhash generation failed',
                        description: 'Please try uploading the image again.',
                      });
                    } finally {
                      setLoading(false);
                    }
                  }}
                  {...fields}
                  maxFiles={1}
                  maxSize={MAX_BYTES}
                  hasError={previewHasError}
                  disabled={disabled}
                  className='aspect-video disabled:cursor-not-allowed'
                />
              </FormControl>
              <FormMessage />
            </FormItem>

            <div className='flex flex-col gap-y-2 w-full'>
              <h6
                data-error={!!blurHasError}
                className='text-sm leading-none font-medium data-[error=true]:text-destructive'
              >
                Blurhash
              </h6>
              <p className='text-muted-foreground text-sm'>
                A small, blurred preview generated from your image.
              </p>
              <div
                className={cn(
                  'aspect-video relative rounded-md overflow-hidden',
                  (loading || !blurHash) &&
                    'border border-dashed border-border flex-center',
                  blurHasError && 'border-destructive',
                  disabled && 'opacity-50',
                )}
              >
                {loading && (
                  <div className='animate-spin'>
                    <Loader aria-hidden='true' className='size-6' />
                    <span className='sr-only'>Loading blurhash...</span>
                  </div>
                )}
                {blurHash && (
                  <Blurhash
                    hash={blurHash}
                    width='100%'
                    height='100%'
                    className='object-cover'
                  />
                )}
                {!loading && !blurHash && (
                  <div className='flex-center flex-col gap-y-2'>
                    <Layers2 aria-hidden='true' className='size-6' />
                    <p className='text-muted-foreground text-sm text-center'>
                      The blurhash preview will appear once an image is
                      selected.
                    </p>
                  </div>
                )}
              </div>
              {blurHasError && (
                <p
                  data-slot='form-message'
                  className={cn('text-destructive text-sm')}
                >
                  {blurhashError?.message}
                </p>
              )}
            </div>
          </div>
        );
      }}
    />
  );
}
