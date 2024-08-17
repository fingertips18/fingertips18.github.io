import { Link } from "react-router-dom";
import { useLenis } from "lenis/react";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/common/components/shadcn/alert-dialog";
import { LocalImageLoader } from "@/common/components/local-image-loader";
import { Button } from "@/common/components/shadcn/button";
import { Badge } from "@/common/components/shadcn/badge";
import { PROJECTTYPES } from "@/constants/enums";
import { FORMLINK } from "@/constants/projects";

import { AppRequestButton } from "./app-request-button";
import { ProjectPreview } from "./project-preview";

interface ProjectItemProps {
  source: string;
  blurHash?: string;
  name: string;
  subtitle?: string;
  desc: string;
  stack: string[];
  type: string;
  live?: string;
}

const ProjectItem = (props: ProjectItemProps) => {
  const lenis = useLenis();

  const onDialogOpen = () => lenis?.stop();

  const onDialogClose = () => lenis?.start();

  return (
    <div
      className="w-full rounded-lg overflow-hidden bg-primary/5 
      border drop-shadow-2xl flex justify-between flex-col hover:drop-shadow-purple-glow
      transition-all duration-500 ease-in-out hover:-translate-y-2 cursor-pointer"
    >
      <AlertDialog>
        <AlertDialogTrigger
          onClick={onDialogOpen}
          className="h-full w-full flex-between flex-col"
        >
          <ProjectPreview {...props} />
        </AlertDialogTrigger>
        <AlertDialogContent
          data-lenis-prevent
          className="overflow-y-auto no-scrollbar h-4/5 lg:h-fit"
        >
          <AlertDialogHeader>
            <div className="aspect-video relative">
              {props.type === PROJECTTYPES.web ? (
                <iframe
                  className="w-full h-full rounded-md"
                  src={props.source}
                  title={`${props.name} Preview`}
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                  referrerPolicy="strict-origin-when-cross-origin"
                  allowFullScreen
                />
              ) : (
                <LocalImageLoader
                  hash={props.blurHash!}
                  src={props.source}
                  alt={props.name}
                  className="aspect-video object-cover object-center rounded-md"
                />
              )}
            </div>

            <AlertDialogTitle className="flex items-center flex-wrap gap-x-2 gap-y-1">
              {props.name}{" "}
              <span className="text-sm text-muted-foreground leading-none">
                {props.subtitle}
              </span>
            </AlertDialogTitle>
            <AlertDialogDescription className="text-start">
              {props.desc}
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="space-y-2.5">
            <h6 className="font-semibold text-sm text-primary-foreground/80">
              Tech Stack
            </h6>

            <div className="flex item-start flex-wrap gap-1.5 no-scrollbar">
              {props.stack.map((s) => (
                <Badge
                  key={`${props.name}-alert-${s}`}
                  className="bg-primary/30 whitespace-nowrap"
                >
                  {s}
                </Badge>
              ))}
            </div>
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel onClick={onDialogClose}>Close</AlertDialogCancel>
            <AlertDialogAction asChild>
              {props.type === PROJECTTYPES.web ? (
                <Link to={props.live!} target="_blank" onClick={onDialogClose}>
                  View Live
                </Link>
              ) : (
                <Link to={FORMLINK} target="_blank" onClick={onDialogClose}>
                  Fill out form
                </Link>
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <div className="bg-primary/20 px-2 py-2.5 flex-center">
        {props.type === PROJECTTYPES.web ? (
          <Button
            asChild
            variant={"link"}
            className="h-auto w-auto px-2.5 py-0.5 text-sm font-bold"
          >
            <Link to={props.live!} target="_blank">
              View Live
            </Link>
          </Button>
        ) : (
          <AppRequestButton />
        )}
      </div>
    </div>
  );
};

export { ProjectItem };
