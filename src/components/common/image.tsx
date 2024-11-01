import React, { ImgHTMLAttributes, SyntheticEvent, useState } from "react";

interface ImageProps extends ImgHTMLAttributes<HTMLImageElement> {
  className?: string;
}

const Image = React.forwardRef<HTMLImageElement, ImageProps>(
  ({ className, ...props }) => {
    const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

    return (
      <img
        {...props}
        width={dimensions.width}
        height={dimensions.height}
        className={className}
        onLoad={(e: SyntheticEvent<HTMLImageElement>) => {
          const { naturalWidth, naturalHeight } = e.currentTarget;
          setDimensions({ width: naturalWidth, height: naturalHeight });
        }}
      />
    );
  }
);

export { Image };
