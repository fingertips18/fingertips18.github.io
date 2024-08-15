import { LocalImageLoader } from "@/common/components/local-image-loader";
import { Badge } from "@/common/components/shadcn/badge";
import { PROJECTTYPES } from "@/constants/enums";

interface ProjectPreviewProps {
  source: string;
  blurHash?: string;
  name: string;
  subtitle?: string;
  desc: string;
  stack: string[];
  type: string;
  live?: string;
}

const ProjectPreview = ({
  source,
  blurHash,
  name,
  subtitle,
  desc,
  stack,
  type,
}: ProjectPreviewProps) => {
  return (
    <>
      {type === PROJECTTYPES.web ? (
        <div className="aspect-video relative">
          <iframe
            className="w-full h-full"
            src={source}
            title={`${name} Preview`}
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            referrerPolicy="strict-origin-when-cross-origin"
            allowFullScreen
          />
        </div>
      ) : (
        <LocalImageLoader
          hash={blurHash!}
          src={source}
          alt={name}
          className="aspect-video object-cover object-center"
        />
      )}

      <div className="space-y-2 p-4 mt-2 flex-grow text-start">
        <h3 className="text-lg font-bold leading-none flex items-center flex-wrap gap-x-2 gap-y-1">
          {name}
          {subtitle && (
            <span className="font-semibold text-sm text-accent">
              {subtitle}
            </span>
          )}
        </h3>

        <p className="text-xs text-primary-foreground/50 line-clamp-4">
          {desc}
        </p>

        <h6 className="font-semibold text-xs text-primary-foreground/80">
          Tech Stack
        </h6>

        <div className="flex item-start flex-wrap gap-1.5 no-scrollbar">
          {stack.map((s) => (
            <Badge
              key={`${name}-${s}`}
              className="bg-primary/30 whitespace-nowrap"
            >
              {s}
            </Badge>
          ))}
        </div>
      </div>
    </>
  );
};

export { ProjectPreview };
