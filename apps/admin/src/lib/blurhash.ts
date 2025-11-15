import { decode, encode, isBlurhashValid } from 'blurhash';

/**
 * Converts a File object to an Image element.
 * @param file - The File object to convert
 * @returns A Promise that resolves to an HTMLImageElement when the file is successfully loaded, or rejects with an Error if the file cannot be read or loaded
 * @throws {Error} If the FileReader result is not a string or if the image fails to load
 */
export function fileToImage(file: File): Promise<HTMLImageElement> {
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
  canvas.width = image.width;
  canvas.height = image.height;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create canvas context');
  }

  context.drawImage(image, 0, 0);

  return context.getImageData(0, 0, image.width, image.height);
}

/**
 * Loads an image from the specified URL.
 * @param imageURL - The URL of the image to load
 * @returns A promise that resolves to the loaded HTMLImageElement, or rejects with an error if the image fails to load
 * @throws {Error} When the image fails to load, with a message indicating the failed URL
 */
export function loadImage(imageURL: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
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

  return canvas.toDataURL('image/webp'); // returns "data:image/png;base64,..."
}
