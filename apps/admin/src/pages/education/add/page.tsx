import { NotebookPen } from 'lucide-react';

import { Back } from '@/components/common/back';

export default function AddEducationPage() {
  return (
    <section className='content padding flex flex-col'>
      <div className='flex-between gap-x-4'>
        <Back />
        <div className='flex-end gap-x-2 text-primary'>
          <NotebookPen aria-hidden='true' className='size-4 lg:size-6' />
          <h1 className='font-bold text-sm lg:text-2xl lg:tracking-wider uppercase'>
            Add New Education
          </h1>
        </div>
      </div>

      <div className='flex-1 flex-center'>
        <h6>Add Education</h6>
      </div>
    </section>
  );
}
