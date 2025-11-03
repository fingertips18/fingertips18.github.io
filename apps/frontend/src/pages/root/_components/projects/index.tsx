import { useLenis } from 'lenis/react';
import { Terminal } from 'lucide-react';
import { useMemo, useRef } from 'react';
import { Link } from 'react-router-dom';

import {
  ProjectCard,
  ProjectCardSkeleton,
} from '@/components/common/project-card';
import { QUERYELEMENT, ROOTSECTION } from '@/constants/enums';
import { PROJECTS } from '@/constants/projects';
import { useObserver } from '@/lib/hooks/useObserver';
import { useRootSectionStore } from '@/lib/stores/useRootSectionStore';
import { cn, shuffleArray } from '@/lib/utils';
import { AppRoutes } from '@/routes/app-routes';

const Projects = () => {
  const sectionRef = useRef<HTMLElement | null>(null);
  const { isVisible } = useObserver({ elementRef: sectionRef });
  const lenis = useLenis();
  const { onActive } = useRootSectionStore();

  // Use useMemo to compute shuffled projects once
  const shuffledProjects = useMemo(() => shuffleArray(PROJECTS), []);

  const handleScroll = () => {
    if (!lenis) return;

    onActive(ROOTSECTION.projects);
    lenis.scrollTo(0);
  };

  return (
    <section
      id={ROOTSECTION.projects}
      ref={sectionRef}
      className={cn(
        'min-h-dvh flex items-center flex-col gap-y-2 lg:gap-y-6 border-b pt-14 pb-6 px-2 lg:px-0',
        QUERYELEMENT.rootSection,
      )}
    >
      <div className='flex items-center justify-end gap-x-2 w-full pt-6 lg:relative'>
        <Terminal className='w-5 lg:w-8 h-5 lg:h-8 sm:absolute xs:left-6 lg:left-4 xl:left-0 opacity-50' />
        <h2 className='text-lg lg:text-4xl font-bold'>PROJECTS</h2>
        <span className='w-8 lg:w-32 h-1 rounded-full bg-muted-foreground tracking-widest' />
      </div>

      <div className='flex-center flex-col gap-y-1'>
        <p className='text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full'>
          I’ve developed various projects, ranging from web applications to
          Android apps. Here are a few highlights.
        </p>

        <Link
          to={AppRoutes.projects}
          onClick={handleScroll}
          className='mt-2 text-sm hover:text-accent hover:drop-shadow-purple-glow underline-offset-4 hover:underline'
        >
          View All
        </Link>
      </div>

      <div
        style={{
          gridAutoRows: '1fr',
        }}
        className={cn(
          `w-full grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 mt-8 gap-4
          transition-opacity duration-500 ease-in-out`,
          isVisible ? 'opacity-100' : 'opacity-0',
        )}
      >
        {isVisible ? (
          <>
            {shuffledProjects.slice(0, 6).map((p) => (
              <ProjectCard key={p.name} {...p} />
            ))}
          </>
        ) : (
          <>
            {[...Array(6)].map((_, i) => (
              <ProjectCardSkeleton key={`project-item-skeleton-${i}`} />
            ))}
          </>
        )}
      </div>
    </section>
  );
};

export { Projects };
