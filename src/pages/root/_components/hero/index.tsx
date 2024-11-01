import { useRef } from "react";

import { QUERYELEMENT, ROOTSECTION } from "@/constants/enums";
import { useObserver } from "@/lib/hooks/useObserver";
import { useMounted } from "@/lib/hooks/useMounted";
import { useResize } from "@/lib/hooks/useResize";
import { BUILDS } from "@/constants/collections";
import { WAVE } from "@/constants/assets";
import { cn } from "@/lib/utils";

import { ProfilePicture } from "./profile-picture";
import { ResumeButton } from "./resume-button";
import { Introduction } from "./introduction";
import { TypingTexts } from "./typing-texts";
import SocialButtons from "./social-buttons";

const Hero = () => {
  const sectionRef = useRef<HTMLElement>(null);
  const { isVisible } = useObserver({ elementRef: sectionRef });
  const isMounted = useMounted();
  const { width } = useResize();

  const lg = width > 1024;

  return (
    <section
      id={ROOTSECTION.about}
      ref={sectionRef}
      className={cn(
        "min-h-dvh flex-center flex-col gap-y-12 lg:gap-y-24 p-6 lg:py-6 relative border-b lg:px-4 xl:px-0",
        QUERYELEMENT.rootSection
      )}
    >
      <div
        className={cn(
          `mt-14 flex-center lg:flex-between flex-col-reverse lg:flex-row gap-y-4 
          lg:gap-y-8 gap-x-24 w-full transition-opacity duration-1000 ease-in-out`,
          isVisible ? "opacity-100" : "opacity-0"
        )}
      >
        <div
          className={cn(
            "transition-opacity duration-500 ease-in-out",
            isMounted ? "opacity-100" : "opacity-0",
            isVisible
              ? "flex items-center lg:items-start flex-col lg:gap-2"
              : "hidden"
          )}
        >
          <div className="flex items-start justify-center gap-x-2 relative">
            <p className="lg:text-xl font-semibold">Hi there!</p>
            <img
              src={WAVE}
              alt="Wave"
              width={lg ? 30.16 : 20.11} // 181 original width
              height={lg ? 32.16 : 21.44} // 193 original height
              className="w-[20.11px] lg:w-[30.16px] h-[21.44px] lg:h-[32.16px] relative -top-0.5 lg:-top-2"
            />
          </div>
          <h1 className="text-2xl lg:text-4xl font-bold flex items-center flex-col lg:flex-row">
            I'm Ghian Carlos Tan{" "}
            <span className="text-sm lg:text-lg font-semibold text-muted-foreground lg:ml-2">
              (Fingertips)
            </span>
          </h1>
          <TypingTexts texts={BUILDS} />
          <Introduction />
          <ResumeButton />
        </div>

        <ProfilePicture isVisible={isVisible} />
      </div>

      {isVisible && <SocialButtons isMounted={isMounted} />}
    </section>
  );
};

export { Hero };
