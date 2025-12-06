import { useEffect, useState } from 'react';
import type { Accept } from 'react-dropzone';
import Cropper, { type Area } from 'react-easy-crop';

import { Button } from '@/components/shadcn/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/shadcn/dialog';
import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
  type DropzoneProps,
} from '@/components/shadcn/dropzone';
import { Slider } from '@/components/shadcn/slider';
import { fileToImage, imageToFile, loadCroppedImage } from '@/lib/image';
import { toast } from '@/lib/toast';
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
  disabled,
  ...props
}: ImageUploaderProps) {
  const [rawImage, setRawImage] = useState<string | null>(null);
  const [editedImage, setEditedImage] = useState<string | null>(null);
  const [crop, setCrop] = useState<{ x: number; y: number }>({ x: 0, y: 0 });
  const [rotation, setRotation] = useState<number>(0);
  const [zoom, setZoom] = useState<number>(1);
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null);
  const [droppedFiles, setDroppedFiles] = useState<File[]>([]);
  const [editorOpen, setEditorOpen] = useState<boolean>(false);

  // Sync preview with value prop
  useEffect(() => {
    if (!value || value.length === 0) {
      // Schedule the state update to avoid synchronous setState
      queueMicrotask(() => setRawImage(null));
      return;
    }

    const file = value[0];
    const reader = new FileReader();

    reader.onload = (e) => {
      if (typeof e.target?.result === 'string') {
        setRawImage(e.target.result);
      }
    };

    reader.onerror = () => {
      setRawImage(null);
    };

    reader.readAsDataURL(file);

    return () => {
      reader.abort();
    };
  }, [value]);

  // Revoking object URLs to prevent memory leaks
  useEffect(() => {
    return () => {
      if (editedImage?.startsWith('blob:')) {
        URL.revokeObjectURL(editedImage);
      }
    };
  }, [editedImage]);

  // Cleanup effect when the editor is closed
  useEffect(() => {
    if (!editorOpen) {
      setRawImage(null);
      setCrop({ x: 0, y: 0 });
      setRotation(0);
      setZoom(1);
      setCroppedAreaPixels(null);
      setDroppedFiles([]);
    }
  }, [editorOpen]);

  const handleDrop = async (files: File[]) => {
    setDroppedFiles(files);

    if (files.length === 0) return;

    const image = await fileToImage(files[0]);
    setRawImage(image.src);
    setEditorOpen(true);
  };

  const handleCropComplete = (_: Area, croppedAreaPixels: Area) => {
    setCroppedAreaPixels(croppedAreaPixels);
  };

  const handleCloseEditor = () => {
    setCrop({ x: 0, y: 0 });
    setRotation(0);
    setZoom(1);
    setCroppedAreaPixels(null);
    setEditorOpen(false);
    setDroppedFiles([]);
  };

  const handleConfirmEditor = async () => {
    if (droppedFiles.length === 0) {
      handleCloseEditor();
      return;
    }

    try {
      const image = await fileToImage(droppedFiles[0]);

      const croppedImage = await loadCroppedImage({
        image,
        pixelCrop: croppedAreaPixels || {
          x: 0,
          y: 0,
          height: image.height,
          width: image.width,
        },
        rotation,
      });

      if (!croppedImage) {
        throw new Error('Unable to crop image');
      }

      const file = await imageToFile({
        image: croppedImage,
        filename: droppedFiles[0].name,
      });

      // Convert File[] to FileList for form compatibility
      const dataTransfer = new DataTransfer();
      dataTransfer.items.add(file);
      await onChange?.(dataTransfer.files);

      setEditedImage(croppedImage.src);
      onBlur?.();
    } catch {
      handleCloseEditor();
      toast({
        level: 'error',
        title: 'Unable to edit',
        description: 'There was an issue editing your image. Try again later.',
      });
    }

    setEditorOpen(false);
  };

  return (
    <>
      <div onBlur={onBlur} className={cn(disabled && 'cursor-not-allowed')}>
        <Dropzone
          id={id}
          accept={accept || { 'image/webp': ['.webp'] }}
          onDrop={(props) => void handleDrop(props)}
          src={value ? Array.from(value) : undefined}
          disabled={disabled}
          className={cn(editedImage && 'p-0', className)}
          {...props}
        >
          <DropzoneContent className='relative aspect-video'>
            {editedImage && !editorOpen && (
              <img
                src={editedImage}
                alt='Preview'
                className='absolute inset-0 size-full object-cover object-center group-hover:scale-105 transition-transform duration-600 ease-in-out rounded-md'
              />
            )}
          </DropzoneContent>
          <DropzoneEmptyState
            className={cn(
              editedImage &&
                'z-20 relative bg-black/25 backdrop-blur-sm p-8 rounded-md m-2',
            )}
          />
        </Dropzone>
      </div>

      {rawImage && editorOpen && (
        <Dialog open={editorOpen} onOpenChange={() => setEditorOpen(false)}>
          <DialogContent className='flex-col'>
            <DialogHeader>
              <DialogTitle>Customize Your Image</DialogTitle>
              <DialogDescription>
                Adjust your image before saving or uploading.
              </DialogDescription>
            </DialogHeader>
            <div className='relative aspect-video size-full'>
              <Cropper
                image={rawImage}
                crop={crop}
                rotation={rotation}
                zoom={zoom}
                aspect={16 / 9}
                onCropChange={setCrop}
                onRotationChange={setRotation}
                onCropComplete={handleCropComplete}
                onZoomChange={setZoom}
              />
            </div>
            <div className='space-y-1 w-full'>
              <div className='flex-center gap-x-2 w-full px-4 sm:px-0'>
                <span className='text-sm font-medium min-w-12'>Rotate</span>
                <Slider
                  value={[rotation]}
                  min={0}
                  max={360}
                  step={1}
                  onValueChange={(value) => setRotation(value[0])}
                  className='[&>span:first-child]:bg-gray-200'
                />
              </div>
              <div className='flex-center gap-x-2 w-full px-4 sm:px-0'>
                <span className='text-sm font-medium min-w-12'>Zoom</span>
                <Slider
                  value={[zoom]}
                  min={1}
                  max={3}
                  step={0.1}
                  onValueChange={(value) => setZoom(value[0])}
                  className='[&>span:first-child]:bg-gray-200'
                />
              </div>
            </div>
            <Button
              type='button'
              onClick={() => void handleConfirmEditor()}
              className='cursor-pointer w-full'
            >
              Save
            </Button>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}
