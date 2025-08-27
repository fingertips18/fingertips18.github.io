import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { Input } from '@/components/shadcn/input';
import { Textarea } from '@/components/shadcn/textarea';
import { EmailService } from '@/lib/services/email';

const formSchema = z.object({
  email: z
    .string()
    .min(1, { message: 'Email address is required' })
    .email({ message: 'Invalid email address' }),
  name: z.string().min(1, { message: 'Name is required' }),
  subject: z.string().min(1, { message: 'Subject is required' }),
  message: z
    .string()
    .max(500, { message: 'Message must be 500 characters long' })
    .optional(),
});

const ContactForm = () => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: '',
      name: '',
      subject: '',
      message: '',
    },
  });
  const [pending, setPending] = useState(false);

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setPending(true);

    const { hasError, message } = await EmailService.send(values);

    if (hasError) {
      toast.error(message);
    } else {
      form.reset();
      toast.success(message);
    }

    setPending(false);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className='border border-primary/50 rounded-md w-full lg:w-3/5 bg-primary/10 p-4 lg:p-6 transition-all duration-500 ease-in-out hover:shadow-2xl hover:shadow-primary/50 space-y-4'
      >
        <FormField
          control={form.control}
          name='email'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='email'>Email Address</FormLabel>
              <FormControl>
                <Input
                  placeholder='example@domain.com'
                  {...field}
                  id='email'
                  autoComplete='email'
                />
              </FormControl>
              <FormMessage className='leading-none' />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='name'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='name'>Name</FormLabel>
              <FormControl>
                <Input
                  placeholder='John Doe'
                  {...field}
                  id='name'
                  autoComplete='name'
                />
              </FormControl>
              <FormMessage className='leading-none' />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='subject'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='subject'>Subject</FormLabel>
              <FormControl>
                <Input
                  placeholder='Subject of Your Inquiry'
                  {...field}
                  id='subject'
                  name='subject'
                />
              </FormControl>
              <FormMessage className='leading-none' />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name='message'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='message'>Message</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="What's on your mind?"
                  {...field}
                  id='message'
                  name='message'
                  className='resize-none'
                  rows={6}
                />
              </FormControl>
              <FormMessage className='leading-none' />
            </FormItem>
          )}
        />
        <button
          type='submit'
          disabled={pending}
          className='py-2 w-full bg-gradient-to-r from-[#310055] to-[#DC97FF]
          hover:brightness-125 transition-all rounded-md active:scale-95 flex-center
          hover:drop-shadow-purple-glow font-semibold text-white disabled:brightness-90'
        >
          {pending ? <Loader2 className='w-5 h-5 animate-spin' /> : 'Submit'}
        </button>
      </form>
    </Form>
  );
};

export { ContactForm };
