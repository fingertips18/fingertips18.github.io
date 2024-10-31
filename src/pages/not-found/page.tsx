import { Link } from "react-router-dom";

import { Button } from "@/components/shadcn/button";
import { AppRoutes } from "@/routes/app-routes";

const NotFoundPage = () => {
  return (
    <section className="h-[calc(100dvh_-_56px)] flex-center flex-col gap-y-2 lg:gap-y-1.5 leading-tight">
      <h6 className="lg:text-lg font-bold">404 Page Not Found</h6>
      <p className="text-xs lg:text-sm text-muted-foreground">
        Woops! Looks like this page doesn't exist.
      </p>
      <Button asChild className="rounded-full">
        <Link to={AppRoutes.root}>Go back</Link>
      </Button>
    </section>
  );
};

export default NotFoundPage;
