import { ensureString } from '.';

export const ProjectType = {
  web: 'web',
  mobile: 'mobile',
  game: 'game',
} as const;

export type Project = {
  id: string;
  preview: string;
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

export function mapProject(dto: unknown): Project {
  if (typeof dto !== 'object' || dto === null) {
    throw new Error('Invalid project DTO');
  }

  const d = dto as Record<string, unknown>;

  return {
    id: ensureString(d.id, 'id'),
    preview: ensureString(d.preview, 'preview'),
    blurhash: ensureString(d.blurhash, 'blurhash'),
    title: ensureString(d.title, 'title'),
    subTitle: ensureString(d.sub_title, 'sub_title'),
    description: ensureString(d.description, 'description'),
    tags: Array.isArray(d.tags)
      ? d.tags.map((tag, index) => ensureString(tag, `tags[${index}]`))
      : [],
    type: ensureString(
      d.type,
      'type',
    ) as (typeof ProjectType)[keyof typeof ProjectType],
    link: ensureString(d.link, 'link'),
    educationId:
      d.education_id != null
        ? ensureString(d.education_id, 'education_id')
        : undefined,
    createdAt: new Date(ensureString(d.created_at, 'created_at')),
    updatedAt: new Date(ensureString(d.updated_at, 'updated_at')),
  };
}

export function toJSONProject(
  project: Partial<Project>,
): Record<string, unknown> {
  return {
    id: project.id,
    preview: project.preview,
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
}
