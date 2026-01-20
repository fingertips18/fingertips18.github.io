import { ensureDate, ensureNumber, ensureString } from '.';

export const FileRole = {
  image: 'image',
} as const;

export type File = {
  id: string;
  parentTable: string;
  parentID: string;
  role: (typeof FileRole)[keyof typeof FileRole];
  name: string;
  url: string;
  type: string;
  size: number;
  createdAt: Date;
  updatedAt: Date;
};

/**
 * Validates that a value is a valid file role.
 *
 * @param value - The value to validate as a file role
 * @returns The validated file role value
 * @throws {Error} If the value is not a string or is not a valid file role
 */
function ensureRole(value: unknown): 'image' {
  const roleValue = ensureString({
    value,
    name: 'role',
  }) as (typeof FileRole)[keyof typeof FileRole];

  if (!Object.values(FileRole).includes(roleValue)) {
    throw new Error(`Invalid file role: ${roleValue}`);
  }

  return roleValue;
}

/**
 * Maps a DTO object to a File domain object.
 *
 * @param dto - The data transfer object to map
 * @returns The mapped File domain object with validated properties
 * @throws {Error} If the DTO is not a valid object or contains invalid properties
 */
export function mapFile(dto: unknown): File {
  if (typeof dto !== 'object' || dto === null) {
    throw new Error('Invalid file DTO');
  }

  const d = dto as Record<string, unknown>;

  return {
    id: ensureString({ value: d.id, name: 'id' }),
    parentTable: ensureString({ value: d.parent_table, name: 'parent_table' }),
    parentID: ensureString({ value: d.parent_id, name: 'parent_id' }),
    role: ensureRole(d.role),
    name: ensureString({ value: d.name, name: 'name' }),
    url: ensureString({ value: d.url, name: 'url' }),
    type: ensureString({ value: d.type, name: 'type' }),
    size: ensureNumber({ value: d.size, name: 'size' }),
    createdAt: ensureDate({ value: d.created_at, name: 'created_at' }),
    updatedAt: ensureDate({ value: d.updated_at, name: 'updated_at' }),
  };
}

/**
 * Converts a File domain object to a JSON-serializable format.
 *
 * @param file - The File domain object to convert
 * @returns A record with snake_case keys suitable for API serialization, excluding undefined values
 */
export function toJSONFile(file: Partial<File>): Record<string, unknown> {
  const result: Record<string, unknown> = {
    id: file.id,
    parent_table: file.parentTable,
    parent_id: file.parentID,
    role: file.role,
    name: file.name,
    url: file.url,
    type: file.type,
    size: file.size,
    created_at: file.createdAt ? file.createdAt.toISOString() : undefined,
    updated_at: file.updatedAt ? file.updatedAt.toISOString() : undefined,
  };

  // Filter out undefined values
  return Object.fromEntries(
    Object.entries(result).filter(([, v]) => v !== undefined),
  );
}

// -------------------- UPLOADTHING types below --------------------

/**
 * Represents a file that has been successfully uploaded via UploadThing.
 */
export type FileUpload = {
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

/**
 * Validates and converts an unknown value to a fields object with string values.
 *
 * @param value - The value to validate and convert as fields
 * @returns An object with string values, or an empty object if value is not a valid object
 * @throws {Error} If any field value is not a string
 */
function ensureFields(value: unknown): { [k: string]: string } {
  if (!value || typeof value !== 'object') {
    throw new Error("Expected property 'fields' to be an object");
  }

  return Object.fromEntries(
    Object.entries(value).map(([k, v]) => [
      k,
      ensureString({ value: v, name: k }),
    ]),
  );
}

/**
 * Maps a DTO object to a FileUpload domain object.
 *
 * @param dto - The data transfer object to map
 * @returns The mapped FileUpload domain object with validated properties
 * @throws {Error} If the DTO is not a valid object or contains invalid properties
 */
export function mapFileUpload(dto: unknown): FileUpload {
  if (typeof dto !== 'object' || dto === null) {
    throw new Error('Invalid file DTO');
  }

  const d = dto as Record<string, unknown>;

  return {
    key: ensureString({ value: d.key, name: 'key' }),
    fileName: ensureString({ value: d.file_name, name: 'file_name' }),
    fileType: ensureString({ value: d.file_type, name: 'file_type' }),
    fileURL: ensureString({ value: d.file_url, name: 'file_url' }),
    contentDisposition: ensureString({
      value: d.content_disposition,
      name: 'content_disposition',
    }),
    pollingJWT: ensureString({ value: d.polling_jwt, name: 'polling_jwt' }),
    pollingURL: ensureString({ value: d.polling_url, name: 'polling_url' }),
    customId:
      d.custom_id != null
        ? ensureString({ value: d.custom_id, name: 'custom_id' })
        : undefined,
    URL: ensureString({ value: d.url, name: 'url' }),
    fields: ensureFields(d.fields),
  };
}
