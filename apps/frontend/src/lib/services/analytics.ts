import { APIRoutes } from '@/routes/api-routes';

export const AnalyticsService = {
  pageView: async ({
    location,
    title,
  }: {
    location: string;
    title: string;
  }) => {
    try {
      const res = await fetch(APIRoutes.pageView, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          location,
          title,
        }),
      });

      if (!res.ok) {
        throw new Error(
          `failed to page view (status: ${res.status} - ${res.statusText})`,
        );
      }

      return;
    } catch (error) {
      console.error('AnalyticsService.pageView error: ', error);
      return;
    }
  },
};
