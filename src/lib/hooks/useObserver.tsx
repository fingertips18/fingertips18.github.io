import { RefObject, useEffect, useState } from "react";

interface useObserverProps {
  elementRef: RefObject<HTMLElement>;
  threshold?: number;
  root?: Element | Document | null;
  rootMargin?: string;
}

const useObserver = ({
  elementRef,
  threshold = 0.1,
  root = null,
  rootMargin = "0px",
}: useObserverProps) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.disconnect(); // Stop observing after it becomes visible
        }
      },
      { threshold: threshold, root: root, rootMargin: rootMargin }
    );

    const currentRef = elementRef.current;

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [elementRef, threshold, root, rootMargin]);

  return { isVisible };
};
export { useObserver };
