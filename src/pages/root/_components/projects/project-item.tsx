import { useLenis } from "lenis/react";
import { useState } from "react";
import { Link } from "react-router-dom";

import { LocalImageLoader } from "@/components/common/local-image-loader";
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
} from "@/components/shadcn/alert-dialog";
import { Badge } from "@/components/shadcn/badge";
import { Button } from "@/components/shadcn/button";
import { Skeleton } from "@/components/shadcn/skeleton";
import { PROJECTTYPE } from "@/constants/enums";
import { FORMLINK } from "@/constants/projects";
import { cn } from "@/lib/utils";

import { AppRequestButton } from "./app-request-button";
import { ProjectPreview } from "./project-preview";

interface ProjectItemProps {
  preview: string;
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
  const [loaded, setLoaded] = useState(false);

  const onDialogOpen = () => lenis?.stop();

  const onDialogClose = () => lenis?.start();

  return (
    <div
      className={cn(
        `w-full rounded-lg overflow-hidden bg-primary/5 drop-shadow-2xl 
        flex justify-between flex-col hover:drop-shadow-purple-glow cursor-pointer 
        transition-all duration-500 ease-in-out hover:-translate-y-2`,
        loaded && "border"
      )}
      onLoad={() => setLoaded(true)}
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
              {props.type === PROJECTTYPE.web ? (
                <iframe
                  className="w-full h-full rounded-md"
                  src={props.preview}
                  title={`${props.name} Preview`}
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                  referrerPolicy="strict-origin-when-cross-origin"
                  allowFullScreen
                />
              ) : (
                <LocalImageLoader
                  hash={props.blurHash!}
                  src={props.preview}
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
              {props.type === PROJECTTYPE.web ? (
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
        {props.type === PROJECTTYPE.web ? (
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

const getRandomWidth = () => {
  // Change the min and max values as needed
  const minWidth = 32; // Minimum width in pixels
  const maxWidth = 128; // Maximum width in pixels
  return Math.floor(Math.random() * (maxWidth - minWidth + 1)) + minWidth;
};

const ProjectItemSkeleton = () => {
  const [loaded, setLoaded] = useState(false);

  return (
    <div
      className={cn(
        `w-full rounded-lg overflow-hidden bg-primary/5 drop-shadow-2xl flex justify-between flex-col`,
        loaded && "border"
      )}
      onLoad={() => setLoaded(true)}
    >
      <Skeleton className="aspect-video" />
      <div className="h-4/5 lg:h-fit space-y-2.5 p-2">
        <Skeleton className="w-4/5 h-6" />
        <div className="space-y-1">
          <Skeleton className="w-full h-2" />
          <Skeleton className="w-4/5 h-2" />
          <Skeleton className="w-11/12 h-2" />
          <Skeleton className="w-3/4 h-2" />
          <Skeleton className="w-full h-2" />
        </div>
        <Skeleton className="w-[112px] h-4" />
        <div className="flex item-start flex-wrap gap-1.5">
          {[...Array(12)].map((_, i) => (
            <Skeleton
              key={`badge-skeleton-${i}`}
              style={{
                width: getRandomWidth(),
              }}
              className="rounded-full h-4"
            />
          ))}
        </div>
      </div>
      <div className="bg-primary/10 px-2 py-2.5 flex-center mt-4">
        <Skeleton className="w-24 h-5" />
      </div>
    </div>
  );
};

export { ProjectItem, ProjectItemSkeleton };
