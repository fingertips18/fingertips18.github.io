import { FolderPlus } from 'lucide-react';

import { Back } from '@/components/common/back';

import { Form } from './_components/form';

export default function AddProjectPage() {
  return (
    <section className='content padding flex flex-col gap-y-6 lg:gap-y-8 overflow-y-auto'>
      <div className='flex-between gap-x-4'>
        <Back />
        <div className='flex-end gap-x-2 text-primary'>
          <FolderPlus aria-hidden='true' className='size-4 lg:size-6' />
          <h1 className='font-bold text-sm lg:text-2xl lg:tracking-wider uppercase'>
            Add New Project
          </h1>
        </div>
      </div>

      <Form />
    </section>
  );
}
