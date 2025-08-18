import { MoveLeft } from 'lucide-react';
import { useEffect } from 'react';
import { Link } from 'react-router-dom';

import { SkillIcon } from '@/components/common/skill-icon';
import { BACKEND, FRONTEND, OTHERS, TOOLS } from '@/constants/skills';
import { AnalyticsService } from '@/lib/services/analytics';
import { AppRoutes } from '@/routes/app-routes';

const SkillsPage = () => {
  useEffect(() => {
    // Intentionally ignore the returned promise.
    void AnalyticsService.pageView({
      location: AppRoutes.skills,
      title: 'Skills View',
    });
  }, []);

  return (
    <section className='min-h-[calc(100dvh_-_56px)] space-y-2 lg:space-y-12 p-6 lg:py-6 lg:px-4 xl:px-0 mt-14'>
      <Link
        to={AppRoutes.root}
        className='flex items-center gap-x-2 hover:text-accent text-sm'
      >
        <MoveLeft className='size-4' /> Go home
      </Link>

      <div className='h-full w-full lg:pb-8'>
        <h1 className='text-xs lg:text-sm font-bold text-center tracking-widest pt-6 lg:pb-2'>
          SKILLS
        </h1>
        <p className='text-xl lg:text-5xl text-center'>
          Innovate, Implement, <span className='text-primary'>Repeat.</span>
        </p>
        <p className='text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full mx-auto'>
          Showcasing the skills I've developed and refined over the past 3
          years.
        </p>

        <h4 className='text-xs lg:text-sm font-bold text-center tracking-widest mt-12'>
          FRONTEND
        </h4>
        <ul className='flex-center flex-wrap gap-4 mt-4'>
          {FRONTEND.map((f, i) => (
            <SkillIcon
              key={`frontend-${f.label}-${i}`}
              Icon={f.icon}
              hexColor={f.hexColor}
            />
          ))}
        </ul>

        <h4 className='text-xs lg:text-sm font-bold text-center tracking-widest mt-12'>
          BACKEND
        </h4>
        <ul className='flex-center flex-wrap gap-4 mt-4'>
          {BACKEND.map((b, i) => (
            <SkillIcon
              key={`backend-${b.label}-${i}`}
              Icon={b.icon}
              hexColor={b.hexColor}
            />
          ))}
        </ul>

        <h4 className='text-xs lg:text-sm font-bold text-center tracking-widest mt-12'>
          OTHERS
        </h4>
        <ul className='flex-center flex-wrap gap-4 mt-4'>
          {OTHERS.map((o, i) => (
            <SkillIcon
              key={`others-${o.label}-${i}`}
              Icon={o.icon}
              hexColor={o.hexColor}
            />
          ))}
        </ul>

        <h4 className='text-xs lg:text-sm font-bold text-center tracking-widest mt-12'>
          TOOLS
        </h4>
        <ul className='flex-center flex-wrap gap-4 mt-4'>
          {TOOLS.map((t, i) => (
            <SkillIcon
              key={`tools-${t.label}-${i}`}
              Icon={t.icon}
              hexColor={t.hexColor}
            />
          ))}
        </ul>
      </div>
    </section>
  );
};

export { SkillsPage };
