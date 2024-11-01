import { useState } from "react";

import { LocalImageLoader } from "@/components/common/local-image-loader";
import { FINGERTIPS_HASH, ME_HASH } from "@/constants/hashes";
import { FINGERTIPS, ME } from "@/constants/assets";

import { Background } from "./background";

const ProfilePicture = () => {
  const [flipped, setFlipped] = useState(false);

  return (
    <div
      className="relative rounded-full w-[256px] lg:min-w-[364px] h-[256px] lg:min-h-[364px] cursor-pointer flex-center"
      onMouseEnter={() => setFlipped(true)}
      onMouseLeave={() => setFlipped(false)}
      style={{ perspective: "1000px" }}
    >
      <Background />
      <div
        className="absolute w-full h-full transition-transform duration-500 ease-in-out border lg:border-4 rounded-full border-secondary"
        style={{
          transformStyle: "preserve-3d",
          transform: `rotateY(${flipped ? 180 : 0}deg)`,
        }}
      >
        <div
          className="absolute w-full h-full flex-center rounded-full overflow-hidden"
          style={{
            backfaceVisibility: "hidden",
          }}
        >
          <LocalImageLoader
            src={ME}
            alt="Me"
            hash={ME_HASH}
            className="w-full h-full object-cover rounded-full"
          />
        </div>

        <div
          className="absolute w-full h-full flex-center rounded-full"
          style={{
            backfaceVisibility: "hidden",
            transform: "rotateY(180deg)",
          }}
        >
          <LocalImageLoader
            src={FINGERTIPS}
            alt="Fingertips"
            hash={FINGERTIPS_HASH}
            className="w-full h-full object-cover rounded-full"
          />
        </div>
      </div>
    </div>
  );
};

export { ProfilePicture };
