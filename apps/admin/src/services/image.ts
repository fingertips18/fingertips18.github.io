import { APIRoute } from '@/constants/api';
import { mapImageFile } from '@/types/image';

export const ImageService = {
  upload: async ({
    file,
    signal,
  }: {
    file: File;
    signal?: AbortSignal;
  }): Promise<string | null> => {
    try {
      const response = await fetch(`${APIRoute.image}/upload`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          files: [
            {
              name: file.name,
              size: file.size,
              type: file.type,
            },
          ],
        }),
        signal,
      });

      if (!response.ok) {
        throw new Error(
          `failed to upload project image (status: ${response.status} - ${response.statusText})`,
        );
      }

      const data = await response.json();

      const imageFile = mapImageFile(data.file);

      const formData = new FormData();
      Object.entries(imageFile.fields).forEach(([k, v]) => {
        formData.append(k, v);
      });
      formData.append('file', file, file.name);

      const uploadResponse = await fetch(imageFile.URL, {
        method: 'POST',
        body: formData,
        signal,
      });

      if (!uploadResponse.ok) {
        throw new Error(
          `failed to upload file to storage (status: ${uploadResponse.status})`,
        );
      }

      return imageFile.fileURL;
    } catch (error) {
      console.error('ImageService.upload error: ', error);
      return null;
    }
  },
};
