import { useEffect, useState } from "react";

interface TypingTextsProps {
  texts: string[];
  speedForward?: number;
  speedBackward?: number;
  delayForward?: number;
  delayBackward?: number;
}

const TypingTexts = ({
  texts,
  speedForward = 120,
  speedBackward = 50,
  delayForward = 2500,
  delayBackward = 800,
}: TypingTextsProps) => {
  const [displayText, setDisplayText] = useState("");
  const [currentTextIndex, setCurrentTextIndex] = useState(0);
  const [index, setIndex] = useState(0);
  const [typingForward, setTypingForward] = useState(true);

  useEffect(() => {
    const text = texts[currentTextIndex];

    if (typingForward) {
      if (index < text.length) {
        const timer = setTimeout(() => {
          setDisplayText((prev) => prev + text[index]);
          setIndex((prev) => prev + 1);
        }, speedForward);

        return () => clearTimeout(timer);
      } else {
        // Delay before starting to type backwards
        const delayTimer = setTimeout(
          () => setTypingForward(false),
          delayForward
        );
        return () => clearTimeout(delayTimer);
      }
    } else {
      if (index > 0) {
        const timer = setTimeout(() => {
          setDisplayText((prev) => prev.slice(0, -1));
          setIndex((prev) => prev - 1);
        }, speedBackward);
        return () => clearTimeout(timer);
      } else {
        // Move to the next text after typing backwards
        const nextTextDelayTimer = setTimeout(() => {
          setTypingForward(true);
          setCurrentTextIndex((prev) => (prev + 1) % texts.length); // Move to the next text
          setIndex(0);
        }, delayBackward);
        return () => clearTimeout(nextTextDelayTimer);
      }
    }
  }, [
    index,
    typingForward,
    currentTextIndex,
    texts,
    speedForward,
    speedBackward,
    delayForward,
    delayBackward,
  ]);

  return (
    <div className="flex items-start justify-center font-semibold text-lg lg:text-2xl">
      I create
      <div className="text-[#8B2FC9] ml-1.5">
        {displayText}
        <span className="inline-block w-0.5 h-[1em] animate-blink ease-in-out">
          |
        </span>
      </div>
    </div>
  );
};

export { TypingTexts };
