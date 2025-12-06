import { decode, encode, isBlurhashValid } from 'blurhash';
import type { Area } from 'react-easy-crop';

/**
 * Converts an HTML image element to a File object by drawing it to a canvas and converting to WebP format.
 *
 * @param image - The HTMLImageElement to convert
 * @param filename - The desired filename for the resulting File object
 * @returns A Promise that resolves to a File object containing the WebP-encoded image data
 * @throws {Error} If unable to get the 2D canvas context or create a Blob from the canvas
 *
 * @example
 * ```typescript
 * const img = document.querySelector('img') as HTMLImageElement;
 * const file = await imageToFile({ image: img, filename: 'photo.webp' });
 * ```
 */
export async function imageToFile({
  image,
  filename,
}: {
  image: HTMLImageElement;
  filename: string;
}): Promise<File> {
  const canvas = document.createElement('canvas');
  canvas.width = image.naturalWidth;
  canvas.height = image.naturalHeight;

  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('Failed to get 2D context');

  ctx.drawImage(image, 0, 0);

  const blob: Blob = await new Promise((resolve, reject) => {
    canvas.toBlob(
      (b) => (b ? resolve(b) : reject(new Error('Failed to create Blob'))),
      'image/webp',
      1,
    );
  });

  return new File([blob], filename, { type: 'image/webp' });
}

/**
 * Converts a File object to an Image element.
 * @param file - The File object to convert
 * @returns A Promise that resolves to an HTMLImageElement when the file is successfully loaded, or rejects with an Error if the file cannot be read or loaded
 * @throws {Error} If the FileReader result is not a string or if the image fails to load
 */
export async function fileToImage(file: File): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result;

      // guarantee string
      if (typeof result !== 'string') {
        return reject(new Error('FileReader result is not a string'));
      }

      const image = new Image();
      image.onload = () => resolve(image);
      image.onerror = reject;
      image.src = result;
    };

    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * Extracts image data from an HTMLImageElement by drawing it onto a canvas.
 * @param image - The HTML image element to extract data from
 * @returns The ImageData object containing pixel data from the image
 * @throws {Error} When unable to create a 2D canvas context
 */
export function getImageData(image: HTMLImageElement): ImageData {
  const canvas = document.createElement('canvas');
  canvas.width = image.naturalWidth;
  canvas.height = image.naturalHeight;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create canvas context');
  }

  context.drawImage(image, 0, 0);

  return context.getImageData(0, 0, image.naturalWidth, image.naturalHeight);
}

/**
 * Loads an image from the specified URL.
 * @param imageURL - The URL of the image to load
 * @returns A promise that resolves to the loaded HTMLImageElement, or rejects with an error if the image fails to load
 * @throws {Error} When the image fails to load, with a message indicating the failed URL
 */
export async function loadImage(imageURL: string): Promise<HTMLImageElement> {
  return await new Promise((resolve, reject) => {
    const image = new Image();

    image.onload = () => resolve(image);

    image.onerror = () =>
      reject(new Error(`Failed to load image: ${imageURL}`));

    image.src = imageURL;
  });
}

/**
 * Encodes an image to a blurhash string representation.
 * @param data - The image data as a URL string or HTMLImageElement
 * @returns A promise that resolves to the blurhash encoded string
 * @throws May throw if image loading fails or encoding encounters an error
 * @example
 * // From URL string
 * const hash = await encodeImageToBlurhash('https://example.com/image.jpg');
 *
 * // From HTMLImageElement
 * const imgElement = document.querySelector('img');
 * const hash = await encodeImageToBlurhash(imgElement);
 */
export async function encodeImageToBlurhash(
  data: string | HTMLImageElement,
): Promise<string> {
  let imageData: ImageData;

  if (typeof data === 'string') {
    const image = await loadImage(data);
    imageData = getImageData(image);
  } else {
    imageData = getImageData(data);
  }

  return encode(imageData.data, imageData.width, imageData.height, 4, 4);
}

interface DecodeProps {
  hash: string;
  width?: number;
  height?: number;
}

/**
 * Decodes a blurhash string into a base64-encoded data URL.
 *
 * @param props - The decode properties
 * @param props.hash - The blurhash string to decode
 * @param props.width - The width of the output image in pixels
 * @param props.height - The height of the output image in pixels
 * @returns A data URL string representing the decoded image in WebP format, or `undefined` if the hash is invalid or canvas context cannot be obtained
 *
 * @example
 * const dataUrl = decodeBlurhashToBase64URL({
 *   hash: 'UeKUpq9ajt7d_aqamqff_nIqf6f6_MbI_0jE',
 *   width: 320,
 *   height: 240
 * });
 * // Returns: "data:image/webp;base64,..."
 */
export function decodeBlurhashToBase64URL({
  hash,
  width = 32,
  height = 32,
}: DecodeProps): string | undefined {
  if (!isBlurhashValid(hash).result) return undefined;

  const pixels = decode(hash, width, height);

  const canvas = document.createElement('canvas');
  canvas.width = width;
  canvas.height = height;

  const context = canvas.getContext('2d');
  if (!context) return undefined;

  const imageData = context.createImageData(width, height);
  imageData.data.set(pixels);
  context.putImageData(imageData, 0, 0);

  return canvas.toDataURL('image/webp'); // returns "data:image/webp;base64,..."
}

/**
 * Converts an angle from degrees to radians.
 * @param degree - The angle in degrees to convert.
 * @returns The angle converted to radians.
 * @example
 * radianAngle(180) // returns Math.PI
 * radianAngle(90)  // returns Math.PI / 2
 */
function radianAngle(degree: number): number {
  return (degree * Math.PI) / 180;
}

/**
 * Calculates the bounding box dimensions of a rectangle after rotation.
 *
 * @param width - The original width of the rectangle
 * @param height - The original height of the rectangle
 * @param rotation - The rotation angle in degrees (0-360)
 * @returns An object containing the new width and height of the bounding box after rotation
 *
 * @example
 * const result = rotateSize({ width: 100, height: 50, rotation: 45 });
 * // result: { width: 106.07, height: 106.07 }
 */
function rotateSize({
  width,
  height,
  rotation,
}: {
  width: number;
  height: number;
  rotation: number;
}): { width: number; height: number } {
  const rotationRadianAngle = radianAngle(rotation);

  const w =
    Math.abs(Math.cos(rotationRadianAngle) * width) +
    Math.abs(Math.sin(rotationRadianAngle) * height);
  const h =
    Math.abs(Math.sin(rotationRadianAngle) * width) +
    Math.abs(Math.cos(rotationRadianAngle) * height);

  return {
    width: w,
    height: h,
  };
}

/**
 * Loads an image, applies rotation and flipping, and returns a cropped portion as an HTMLImageElement with a blob URL.
 *
 * @remarks
 * The returned image's `src` is an object URL created via `URL.createObjectURL`.
 * Callers should call `URL.revokeObjectURL(image.src)` when the image is no longer needed.
 *
 * @param options - Configuration object for image processing
 * @param options.url - The URL of the image to load
 * @param options.pixelCrop - The rectangular area to crop from the rotated image
 * @param options.rotation - The rotation angle in degrees (default: 0)
 *
 * @returns A promise that resolves to an HTMLImageElement with a blob URL (created via URL.createObjectURL) of the cropped image,
 *          or null if the canvas context could not be obtained
 *
 * @example
 * ```typescript
 * const croppedImage = await loadCroppedImage({
 *   image: 'https://example.com/image.webp',
 *   pixelCrop: { x: 10, y: 10, width: 100, height: 100 },
 *   rotation: 45,
 * });
 * ```
 */
export async function loadCroppedImage({
  image,
  pixelCrop,
  rotation = 0,
  flip = {
    horizontal: false,
    vertical: false,
  },
}: {
  image: string | HTMLImageElement;
  pixelCrop: Area;
  rotation?: number;
  flip?: { horizontal: boolean; vertical: boolean };
}): Promise<HTMLImageElement | null> {
  if (typeof image === 'string') {
    image = await loadImage(image);
  }
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');

  if (!ctx) return null;

  const rotationAngle = radianAngle(rotation);

  // Calculate bounding box of the rotated image
  const { width: bBoxWidth, height: bBoxHeight } = rotateSize({
    width: image.width,
    height: image.height,
    rotation,
  });

  // Set canvas size to match the bounding box
  canvas.width = bBoxWidth;
  canvas.height = bBoxHeight;

  // Translate canvas context to a central location to allow rotating and flipping around the center
  ctx.translate(bBoxWidth / 2, bBoxHeight / 2);
  ctx.rotate(rotationAngle);
  ctx.scale(flip.horizontal ? -1 : 1, flip.vertical ? -1 : 1);
  ctx.translate(-image.width / 2, -image.height / 2);

  // Draw rotated image
  ctx.drawImage(image, 0, 0);

  const croppedCanvas = document.createElement('canvas');
  const croppedCtx = croppedCanvas.getContext('2d');

  if (!croppedCtx) return null;

  // Set the size of the cropped canvas
  croppedCanvas.width = pixelCrop.width;
  croppedCanvas.height = pixelCrop.height;

  // Draw the cropped image onto the new canvas
  croppedCtx.drawImage(
    canvas,
    pixelCrop.x,
    pixelCrop.y,
    pixelCrop.width,
    pixelCrop.height,
    0,
    0,
    pixelCrop.width,
    pixelCrop.height,
  );

  // Convert to Blob → img element
  const blob: Blob = await new Promise((resolve, reject) =>
    croppedCanvas.toBlob(
      (file) => (file ? resolve(file) : reject(new Error('Blob failed'))),
      'image/webp',
      1,
    ),
  );

  const url = URL.createObjectURL(blob);

  // Return HTMLImageElement
  const img = new Image();
  img.src = url;

  await new Promise((resolve, reject) => {
    img.onload = resolve;
    img.onerror = reject;
  });

  return img;
}
