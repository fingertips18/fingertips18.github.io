import { Link } from "react-router-dom";

import { Button } from "@/components/shadcn/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/shadcn/dialog";
import { FORMLINK } from "@/constants/projects";

const AppRequestButton = () => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          variant={"link"}
          className="h-auto w-auto px-2.5 py-0.5 text-sm font-bold"
        >
          Request App
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Request Access for Apps</DialogTitle>
        </DialogHeader>
        <DialogDescription>
          Request access to my mobile apps by filling out this form. Please
          provide your name, email, and select the app/s you’re interested in.
          I’ll get back to you with the download details shortly!
        </DialogDescription>
        <DialogFooter>
          <Button asChild variant={"link"}>
            <Link to={FORMLINK} target="_blank">
              Fill out form
            </Link>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export { AppRequestButton };
