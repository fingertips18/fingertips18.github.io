import { Blurhash } from "react-blurhash";
import { useState } from "react";

import { cn } from "@/lib/utils";

interface LocalImageLoaderProps {
  hash: string;
  className?: string;
  src: string;
  alt: string;
}

const LocalImageLoader = ({
  hash,
  className,
  src,
  alt,
}: LocalImageLoaderProps) => {
  const [loaded, setLoaded] = useState(false);

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
        loading="lazy"
        width="100%"
        height="100%"
        onLoad={() => setLoaded(true)}
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
