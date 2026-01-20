import { ensureDate, ensureString } from '.';
import { type File, mapFile, toJSONFile } from './file';

export const ProjectType = {
  web: 'web',
  mobile: 'mobile',
  game: 'game',
} as const;

export type Project = {
  id: string;
  previews: File[];
  blurhash: string;
  title: string;
  subTitle: string;
  description: string;
  tags: string[];
  type: (typeof ProjectType)[keyof typeof ProjectType];
  link: string;
  educationId?: string;
  createdAt: Date;
  updatedAt: Date;
};

/**
 * Validates that a value is a valid project type.
 *
 * @param value - The value to validate as a project type
 * @returns The validated project type value
 * @throws {Error} If the value is not a string or is not a valid project type
 */
function ensureType(value: unknown): 'web' | 'mobile' | 'game' {
  const typeValue = ensureString({
    value: value,
    name: 'type',
  }) as (typeof ProjectType)[keyof typeof ProjectType];

  if (!Object.values(ProjectType).includes(typeValue)) {
    throw new Error(`Invalid project type: ${typeValue}`);
  }

  return typeValue;
}

/**
 * Maps a DTO object to a Project domain object.
 *
 * @param dto - The data transfer object to map
 * @returns The mapped Project domain object with validated properties
 * @throws {Error} If the DTO is not a valid object or contains invalid properties
 */
export function mapProject(dto: unknown): Project {
  if (typeof dto !== 'object' || dto === null) {
    throw new Error('Invalid project DTO');
  }

  const d = dto as Record<string, unknown>;

  return {
    id: ensureString({ value: d.id, name: 'id' }),
    previews: Array.isArray(d.previews)
      ? d.previews.map((file) => mapFile(file))
      : [],
    blurhash: ensureString({ value: d.blurhash, name: 'blurhash' }),
    title: ensureString({ value: d.title, name: 'title' }),
    subTitle: ensureString({ value: d.sub_title, name: 'sub_title' }),
    description: ensureString({ value: d.description, name: 'description' }),
    tags: Array.isArray(d.tags)
      ? d.tags.map((tag, index) =>
          ensureString({ value: tag, name: `tags[${index}]` }),
        )
      : [],
    type: ensureType(d.type),
    link: ensureString({ value: d.link, name: 'link' }),
    educationId:
      d.education_id != null
        ? ensureString({ value: d.education_id, name: 'education_id' })
        : undefined,
    createdAt: ensureDate({ value: d.created_at, name: 'created_at' }),
    updatedAt: ensureDate({ value: d.updated_at, name: 'updated_at' }),
  };
}

/**
 * Converts a Project domain object to a JSON-serializable format.
 *
 * @param project - The Project domain object to convert
 * @returns A record with snake_case keys suitable for API serialization, excluding undefined values
 */
export function toJSONProject(
  project: Partial<Project>,
): Record<string, unknown> {
  const result: Record<string, unknown> = {
    id: project.id,
    previews: project.previews?.map((file) => toJSONFile(file)),
    blurhash: project.blurhash,
    title: project.title,
    sub_title: project.subTitle,
    description: project.description,
    tags: project.tags,
    type: project.type,
    link: project.link,
    education_id: project.educationId,
    created_at: project.createdAt ? project.createdAt.toISOString() : undefined,
    updated_at: project.updatedAt ? project.updatedAt.toISOString() : undefined,
  };

  // Filter out undefined values
  return Object.fromEntries(
    Object.entries(result).filter(([, v]) => v !== undefined),
  );
}
