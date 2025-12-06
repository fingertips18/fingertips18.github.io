import { APIRoute } from '@/constants/api';
import { type Project, toJSONProject } from '@/types/project';
import { ProjectType } from '@/types/project';

export const ProjectService = {
  create: async ({
    project,
    signal,
  }: {
    project: Partial<Project>;
    signal?: AbortSignal;
  }): Promise<string | null> => {
    try {
      const response = await fetch(`${APIRoute.project}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(toJSONProject(project)),
        signal,
      });

      if (!response.ok) {
        throw new Error(
          `failed to create project (status: ${response.status} - ${response.statusText})`,
        );
      }

      const data = await response.json();

      const projectId = data.id as string;

      return projectId;
    } catch (error) {
      console.error('ProjectService.create error: ', error);
      return null;
    }
  },
  list: async ({
    page,
    pageSize,
    sortBy,
    sortAscending,
    type,
    signal,
  }: {
    page?: number;
    pageSize?: number;
    sortBy?: 'created_at';
    sortAscending?: boolean;
    type?: keyof typeof ProjectType;
    signal?: AbortSignal;
  } = {}): Promise<Project[]> => {
    try {
      const searchParams = new URLSearchParams();
      searchParams.set('page', String(page || '1'));
      searchParams.set('page_size', String(pageSize || '10'));
      if (sortBy) {
        searchParams.set('sort_by', sortBy);
      }
      if (sortAscending !== undefined) {
        searchParams.set('sort_ascending', String(sortAscending));
      }
      if (type) {
        searchParams.set('type', type);
      }
      const response = await fetch(
        `${APIRoute.project}s?${searchParams.toString()}`,
        {
          method: 'GET',
          signal,
        },
      );

      if (!response.ok) {
        throw new Error(
          `failed to list projects (status: ${response.status} - ${response.statusText})`,
        );
      }

      const data = await response.json();

      const projects = data as Project[];

      return projects;
    } catch (error) {
      console.error('ProjectService.list error: ', error);
      return [];
    }
  },
};
