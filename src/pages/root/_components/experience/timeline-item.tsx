import 'react-vertical-timeline-component/style.min.css';

import { Link } from 'react-router-dom';
import { VerticalTimelineElement } from 'react-vertical-timeline-component';

import { Image } from '@/components/common/image';
import { Badge } from '@/components/shadcn/badge';

interface TimelineItemProps {
  image: string;
  position: string;
  company: string;
  link?: string;
  setup: string;
  date: string;
  highlights: string[];
  skills: string[];
  subCompanies?: {
    company: string;
    image: string;
    link: string;
  }[];
}

const TimelineItem = ({
  image,
  position,
  company,
  link,
  setup,
  date,
  highlights,
  skills,
  subCompanies,
}: TimelineItemProps) => {
  return (
    <VerticalTimelineElement
      contentStyle={{
        background: 'hsl(var(--secondary) / 0.2)',
        border: '1px solid hsl(var(--secondary) / 0.5)',
        color: 'hsl(var(--secondary-foreground))',
        display: 'flex',
        flexDirection: 'column',
        boxShadow: 'hsl(var(--primary) / 0.2) 0px 4px 24px',
        borderRadius: '8px',
      }}
      contentArrowStyle={{
        borderRight: '8px solid  hsl(var(--secondary) / 0.8)',
      }}
      date={date}
      iconStyle={{
        boxShadow: 'hsl(var(--primary)) 0px 4px 24px',
        outline: '2px solid hsl(var(--primary))',
      }}
      icon={
        link ? (
          <Link to={link} target='_blank'>
            <Image
              src={image}
              alt={company}
              className='rounded-full w-full h-full border object-cover cursor-pointer'
              loading='lazy'
            />
          </Link>
        ) : (
          <Image
            src={image}
            alt={company}
            className='rounded-full w-full h-full border object-cover'
            loading='lazy'
          />
        )
      }
    >
      <div className='flex items-start  gap-x-4'>
        <Image
          src={image}
          alt='company'
          className='rounded-sm drop-shadow-primary-glow h-16 w-16 object-scale-down'
          loading='lazy'
        />
        <div className='leading-none'>
          <h3 className='font-bold'>{position}</h3>
          <h4 className='text-sm text-secondary-foreground/80'>{company}</h4>
          <h5 className='text-xs text-secondary-foreground/40'>{date}</h5>
          <h6 className='text-xs text-secondary-foreground/40'>{setup}</h6>
        </div>
      </div>

      <div className='flex flex-col items-start space-y-2'>
        {highlights.map((h) => (
          <p key={h} className='!text-sm text-muted-foreground'>
            — {h}
          </p>
        ))}
      </div>

      {subCompanies && (
        <div className='space-y-1.5'>
          <p className='!font-semibold !text-sm'>Sub-Companies</p>
          <div className='flex item-start flex-wrap gap-2.5 lg:gap-6 mt-4'>
            {subCompanies.map((s) => (
              <Link
                to={s.link}
                key={`${s.company}-${s}`}
                className='flex-center gap-x-2 text-xs text-secondary-foreground/80'
              >
                <div
                  style={{
                    boxShadow: 'hsl(var(--primary)) 0px 4px 24px',
                  }}
                  className='rounded-full h-6 w-6 bg-white overflow-hidden flex-center p-0.5'
                >
                  <Image
                    src={s.image}
                    alt={s.company}
                    className='object-scale-down'
                    loading='lazy'
                  />
                </div>
                {s.company}
              </Link>
            ))}
          </div>
        </div>
      )}

      <div className='space-y-1.5'>
        <p className='!font-semibold !text-sm'>Skills Gained</p>
        <div className='flex item-start flex-wrap gap-1.5 mt-4'>
          {skills.map((s) => (
            <Badge key={`${company}-${s}`}>{s}</Badge>
          ))}
        </div>
      </div>
    </VerticalTimelineElement>
  );
};

export { TimelineItem };
