import { Link } from 'react-router-dom';

import { Hint } from '@/components/common/hint';
import { SOCIALS } from '@/constants/collections';

const Socials = () => {
  return (
    <ul className='flex items-start gap-x-2'>
      {SOCIALS.filter((s) => s.label !== 'LinkedIn').map((s) => (
        <Hint key={`footer-${s.label}`} asChild label={s.label} side='top'>
          <li
            className='rounded-full border border-primary/50 hover:border-primary bg-primary/20 
            hover:bg-primary/50 hover:drop-shadow-primary-glow transition-all size-8 lg:size-10'
          >
            <Link
              to={s.href}
              target='_blank'
              className='w-full h-full flex-center'
            >
              <s.icon className='w-2.5 h-2.5 lg:w-4 lg:h-4 pointer-events-none' />
            </Link>
          </li>
        </Hint>
      ))}
    </ul>
  );
};

export { Socials };
