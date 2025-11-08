import { Plus } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { Route } from '@/routes/route';

export default function SkillPage() {
  return (
    <section className='content padding flex flex-col'>
      <div className='flex-between gap-x-4'>
        <Button asChild>
          <Link to={`${Route.skill}/add`} className='flex-center gap-x-2'>
            <Plus aria-hidden='true' className='size-4' />
            Add Skill
          </Link>
        </Button>
      </div>

      <div className='flex-1 flex-center'>
        <h6>Skill</h6>
      </div>
    </section>
  );
}
