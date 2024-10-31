import { RefObject, useEffect, useState } from "react";

interface useObserverProps {
  elementRef: RefObject<HTMLElement>;
  threshold?: number;
}

const useObserver = ({ elementRef, threshold = 0.1 }: useObserverProps) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsVisible(entry.isIntersecting);
      },
      { threshold: threshold }
    );

    const currentRef = elementRef.current;

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
      observer.disconnect();
    };
  }, [elementRef, threshold]);

  return { isVisible };
};

export { useObserver };
