import { APIRoutes } from '@/routes/api-routes';

export const EmailService = {
  send: async ({
    name,
    email,
    subject,
    message,
  }: {
    name: string;
    email: string;
    subject: string;
    message?: string;
  }) => {
    const hasError: boolean = false;

    try {
      const res = await fetch(APIRoutes.sendEmail, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          email,
          subject,
          message,
        }),
      });

      if (!res.ok) {
        throw new Error(
          `failed to send email (status: ${res.status} - ${res.statusText})`,
        );
      }

      return { hasError, message: 'Message sent. Thanks for reaching out!' };
    } catch (error) {
      console.error('EmailService.send error: ', error);
      return {
        hasError: true,
        message:
          'Oops! Something went wrong while sending your message. Please try again later.',
      };
    }
  },
};
