import 'react-vertical-timeline-component/style.min.css';

import { BriefcaseBusiness } from 'lucide-react';
import { useRef } from 'react';
import { VerticalTimeline } from 'react-vertical-timeline-component';

import { QUERYELEMENT, ROOTSECTION } from '@/constants/enums';
import { EXPERIENCES } from '@/constants/experiences';
import { useObserver } from '@/lib/hooks/useObserver';
import { cn } from '@/lib/utils';

import { TimelineItem } from './timeline-item';

const Experience = () => {
  const sectionRef = useRef<HTMLElement | null>(null);
  const { isVisible } = useObserver({ elementRef: sectionRef });

  return (
    <section
      id={ROOTSECTION.experience}
      ref={sectionRef}
      className={cn(
        'min-h-dvh flex items-center flex-col gap-y-2 lg:gap-y-6 border-b pt-14 pb-6 px-2 lg:px-0',
        QUERYELEMENT.rootSection,
      )}
    >
      <div className='flex items-center gap-x-2 w-full pt-6 lg:relative'>
        <span className='w-[32px] lg:w-[128px] h-1 rounded-full bg-muted-foreground tracking-widest' />
        <h2 className='text-lg lg:text-4xl font-bold'>WORK EXPERIENCE</h2>
        <BriefcaseBusiness className='w-5 lg:w-8 h-5 lg:h-8 sm:absolute xs:right-6 lg:right-4 xl:right-0 opacity-50' />
      </div>
      <p className='text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full'>
        Here are details of my experience as a software developer, including my
        roles across various companies and projects.
      </p>
      <VerticalTimeline
        lineColor='hsl(var(--foreground) / 0.6)'
        className={cn(
          'mt-4 lg:mt-20 transition-opacity duration-500 ease-in-out',
          isVisible ? 'opacity-100 visible' : 'opacity-0 invisible',
        )}
      >
        {EXPERIENCES.map((e) => (
          <TimelineItem key={e.company} {...e} />
        ))}
      </VerticalTimeline>
    </section>
  );
};

export { Experience };
