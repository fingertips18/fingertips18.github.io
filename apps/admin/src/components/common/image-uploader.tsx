import { useEffect, useState } from 'react';
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
  id?: string;
  accept?: Accept;
  className?: string;
  value?: FileList;
  onChange?: (files: FileList) => Promise<void>;
  onBlur?: () => void;
}

export function ImageUploader({
  id,
  accept,
  className,
  value,
  onChange,
  onBlur,
  ...props
}: ImageUploaderProps) {
  const [preview, setPreview] = useState<string | null>(null);

  // Sync preview with value prop
  useEffect(() => {
    if (!value || value.length === 0) {
      // Schedule the state update to avoid synchronous setState
      queueMicrotask(() => setPreview(null));
      return;
    }

    const file = value[0];
    const reader = new FileReader();

    reader.onload = (e) => {
      if (typeof e.target?.result === 'string') {
        setPreview(e.target.result);
      }
    };

    reader.onerror = () => {
      setPreview(null);
    };

    reader.readAsDataURL(file);

    return () => {
      reader.abort();
    };
  }, [value]);

  const handleDrop = async (files: File[]) => {
    // Convert File[] to FileList for form compatibility
    const dataTransfer = new DataTransfer();
    files.forEach((file) => dataTransfer.items.add(file));
    await onChange?.(dataTransfer.files);
    onBlur?.();
  };

  return (
    <div onBlur={onBlur}>
      <Dropzone
        id={id}
        accept={accept || { 'image/webp': ['.webp'] }}
        onDrop={(props) => void handleDrop(props)}
        src={value ? Array.from(value) : undefined}
        className={cn(preview && 'p-0', className)}
        {...props}
      >
        <DropzoneContent className='relative aspect-square lg:aspect-video'>
          {preview && (
            <img
              src={preview}
              alt='Preview'
              className='absolute inset-0 size-full object-cover object-center group-hover:scale-105 transition-transform duration-600 ease-in-out rounded-md'
            />
          )}
        </DropzoneContent>
        <DropzoneEmptyState
          className={cn(
            preview &&
              'z-20 relative bg-black/25 backdrop-blur-sm p-8 rounded-md m-2',
          )}
        />
      </Dropzone>
    </div>
  );
}
