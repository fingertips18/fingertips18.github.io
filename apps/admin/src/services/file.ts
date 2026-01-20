import { APIRoute } from '@/constants/api';
import { mapFileUpload } from '@/types/file';

/**
 * FileService provides methods for uploading files to the backend via UploadThing.
 */
export const FileService = {
  /**
   * Uploads a file by first requesting an upload URL from the backend, then uploading the file to storage.
   *
   * @param file - The File object to upload
   * @param signal - Optional AbortSignal to cancel the upload request
   * @returns The URL of the uploaded file, or null if the upload fails
   */
  upload: async ({
    file,
    signal,
  }: {
    file: File;
    signal?: AbortSignal;
  }): Promise<string | null> => {
    try {
      const response = await fetch(`${APIRoute.file}/upload`, {
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
          `failed to upload file (status: ${response.status} - ${response.statusText})`,
        );
      }

      const data = await response.json();

      const fileUpload = mapFileUpload(data.file);

      const formData = new FormData();
      Object.entries(fileUpload.fields).forEach(([k, v]) => {
        formData.append(k, v);
      });
      formData.append('file', file, file.name);

      const uploadResponse = await fetch(fileUpload.URL, {
        method: 'POST',
        body: formData,
        signal,
      });

      if (!uploadResponse.ok) {
        throw new Error(
          `failed to upload file to storage (status: ${uploadResponse.status})`,
        );
      }

      return fileUpload.fileURL;
    } catch (error) {
      console.error('FileService.upload error: ', error);
      return null;
    }
  },
};
