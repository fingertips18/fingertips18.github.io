export type ImageFile = {
  key: string;
  fileName: string;
  fileType: string;
  fileURL: string;
  contentDisposition: string;
  pollingJWT: string;
  pollingURL: string;
  customId?: string;
  URL: string;
  fields: Record<string, string>;
};

export function mapImageFile(dto: unknown): ImageFile {
  if (typeof dto !== 'object' || dto === null) {
    throw new Error('Invalid image DTO');
  }

  const d = dto as Record<string, unknown>;

  const fields =
    d.fields && typeof d.fields === 'object'
      ? Object.fromEntries(
          Object.entries(d.fields).map(([k, v]) => [k, ensureString(v, k)]),
        )
      : {};

  return {
    key: ensureString(d.key, 'key'),
    fileName: ensureString(d.file_name, 'file_name'),
    fileType: ensureString(d.file_type, 'file_type'),
    fileURL: ensureString(d.file_url, 'file_url'),
    contentDisposition: ensureString(
      d.content_disposition,
      'content_disposition',
    ),
    pollingJWT: ensureString(d.polling_jwt, 'polling_jwt'),
    pollingURL: ensureString(d.polling_url, 'polling_url'),
    customId:
      d.custom_id != null ? ensureString(d.custom_id, 'custom_id') : undefined,
    URL: ensureString(d.url, 'url'),
    fields,
  };
}

function ensureString(value: unknown, name: string): string {
  if (typeof value === 'string') return value;
  if (typeof value === 'number' || typeof value === 'boolean')
    return String(value);
  throw new Error(`Expected property '${name}' to be a string`);
}
