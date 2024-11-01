import { SyntheticEvent, useState } from "react";
import { Blurhash } from "react-blurhash";

import { cn } from "@/lib/utils";

interface LocalImageLoaderProps {
  hash: string;
  className?: string;
  src: string;
  alt: string;
  loadLazy?: boolean;
}

const LocalImageLoader = ({
  hash,
  className,
  src,
  alt,
  loadLazy = true,
}: LocalImageLoaderProps) => {
  const [loaded, setLoaded] = useState(false);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

  return (
    <div className="relative w-full h-full">
      <div
        className={cn(
          "transition-opacity duration-500 ease-in-out overflow-hidden absolute inset-0 w-full h-full",
          loaded ? "opacity-0" : "opacity-100"
        )}
      >
        <Blurhash
          hash={hash}
          width="100%"
          height="100%"
          className="object-cover"
        />
      </div>
      <img
        src={src}
        alt={alt}
        loading={loadLazy ? "lazy" : undefined}
        decoding="async"
        width={dimensions.width}
        height={dimensions.height}
        onLoad={(e: SyntheticEvent<HTMLImageElement>) => {
          setLoaded(true);
          const { naturalWidth, naturalHeight } = e.currentTarget;
          setDimensions({ width: naturalWidth, height: naturalHeight });
        }}
        className={cn(
          "transition-opacity duration-500 ease-in-out w-full h-full",
          loaded ? "opacity-100" : "opacity-0",
          className
        )}
      />
    </div>
  );
};

export { LocalImageLoader };
