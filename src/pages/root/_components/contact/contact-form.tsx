import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import emailjs from "@emailjs/browser";
import { Loader2 } from "lucide-react";
import { useTransition } from "react";
import { toast } from "sonner";
import { z } from "zod";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/common/components/shadcn/form";
import { Textarea } from "@/common/components/shadcn/textarea";
import { Input } from "@/common/components/shadcn/input";

const formSchema = z.object({
  email: z
    .string()
    .min(1, { message: "Email address is required" })
    .email({ message: "Invalid email address" }),
  name: z.string().min(1, { message: "Name is required" }),
  subject: z.string().min(1, { message: "Subject is required" }),
  message: z
    .string()
    .max(500, { message: "Message must be 500 characters long" })
    .optional(),
});

const ContactForm = () => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      name: "",
      subject: "",
    },
  });
  const [pending, startTransition] = useTransition();

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    startTransition(() => {
      emailjs
        .send(
          import.meta.env.VITE_EMAILJS_SERVICE_ID,
          import.meta.env.VITE_EMAILJS_TEMPLATE_ID,
          values,
          {
            publicKey: import.meta.env.VITE_EMAILJS_PUBLIC_KEY,
          }
        )
        .then(() => toast.success("Message sent. Thanks for reaching out!"))
        .catch(() =>
          toast.error("Something went wrong. Please try again later.")
        );
    });
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="border border-primary/50 rounded-md w-full lg:w-3/5 bg-primary/10
        p-4 lg:p-6 transition-all duration-500 ease-in-out hover:shadow-2xl hover:shadow-primary/50 space-y-4"
      >
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="email">Email Address</FormLabel>
              <FormControl>
                <Input
                  placeholder="example@domain.com"
                  {...field}
                  id="email"
                  autoComplete="email"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="name">Name</FormLabel>
              <FormControl>
                <Input
                  placeholder="John Doe"
                  {...field}
                  id="name"
                  autoComplete="name"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="subject"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="subject">Subject</FormLabel>
              <FormControl>
                <Input
                  placeholder="Subject of Your Inquiry"
                  {...field}
                  id="subject"
                  name="subject"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="message"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="message">Message</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="What's on your mind?"
                  {...field}
                  id="message"
                  name="message"
                  className="resize-none"
                  rows={6}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <button
          type="submit"
          disabled={pending}
          className="py-2 w-full bg-gradient-to-r from-[#310055] to-[#DC97FF]
          hover:brightness-125 transition-all rounded-md active:scale-95
          hover:drop-shadow-purple-glow font-semibold text-white"
        >
          {pending ? <Loader2 className="w-4 h-4 animate-spin" /> : "Submit"}
        </button>
      </form>
    </Form>
  );
};

export { ContactForm };
