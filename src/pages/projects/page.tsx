import { useEffect } from 'react';

import { Back } from '@/components/common/back';
import { ProjectCard } from '@/components/common/project-card';
import { PROJECTS } from '@/constants/projects';
import { AnalyticsService } from '@/lib/services/analytics';
import { AppRoutes } from '@/routes/app-routes';

const ProjectsPage = () => {
  useEffect(() => {
    if (import.meta.env.DEV) return;

    // Intentionally ignore the returned promise.
    void AnalyticsService.pageView({
      location: AppRoutes.skills,
      title: 'Projects View',
    });
  }, []);

  return (
    <section className='min-h-[calc(100dvh_-_56px)] space-y-2 lg:space-y-12 p-6 lg:py-6 lg:px-4 xl:px-0 mt-14'>
      <Back to={AppRoutes.root} />

      <div className='size-full lg:pb-8'>
        <h1 className='text-xs lg:text-sm font-bold text-center tracking-widest pt-6 lg:pb-2'>
          PROJECTS
        </h1>
        <p className='text-xl lg:text-5xl text-center'>
          Design, Develop, <span className='text-primary'>Deliver.</span>
        </p>
        <p className='text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full mx-auto'>
          I’ve developed various projects, ranging from web applications to
          Android apps. Here are a few highlights.
        </p>

        <div
          style={{
            gridAutoRows: '1fr',
          }}
          className='w-full grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 mt-8 gap-4 transition-opacity duration-500 ease-in-out'
        >
          {PROJECTS.map((p) => (
            <ProjectCard key={p.name} {...p} />
          ))}
        </div>
      </div>
    </section>
  );
};

export { ProjectsPage };
