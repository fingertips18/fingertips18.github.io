import { useState } from 'react';
import type { Accept } from 'react-dropzone';

import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
  type DropzoneProps,
} from '@/components/shadcn/dropzone';
import { cn } from '@/lib/utils';

interface ImageUploaderProps
  extends Omit<DropzoneProps, 'accept' | 'onDrop' | 'src' | 'className'> {
  accept?: Accept;
  className?: string;
}

export function ImageUploader({
  accept,
  className,
  ...props
}: ImageUploaderProps) {
  const [files, setFiles] = useState<File[] | undefined>(undefined);
  const [preview, setPreview] = useState<string | null>(null);

  const handleDrop = (files: File[]) => {
    setFiles(files);

    if (files.length > 0) {
      const reader = new FileReader();
      reader.onload = (e) => {
        if (typeof e.target?.result === 'string') {
          setPreview(e.target.result);
        }
      };
      reader.readAsDataURL(files[0]);
    }
  };

  return (
    <Dropzone
      accept={accept || { 'image/webp': ['.webp'] }}
      onDrop={handleDrop}
      src={files}
      className={cn(preview && 'p-0', className)}
      {...props}
    >
      <DropzoneContent className='relative aspect-video'>
        {preview && (
          <img
            src={preview}
            alt='Preview'
            className='absolute inset-0 size-full object-cover object-center group-hover:scale-105 transition-transform duration-600 ease-in-out rounded-md'
          />
        )}
      </DropzoneContent>
      <DropzoneEmptyState className='z-20 relative bg-black/25 backdrop-blur-sm p-8 rounded-md m-2' />
    </Dropzone>
  );
}
