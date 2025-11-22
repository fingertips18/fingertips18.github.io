import { APIRoute } from '@/constants/api';
import { type Project, toJSONProject } from '@/types/project';

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
};
