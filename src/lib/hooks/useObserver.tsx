import { RefObject, useEffect, useState } from 'react';

interface useObserverProps {
  elementRef: RefObject<HTMLElement | null>;
  threshold?: number;
  root?: Element | Document | null;
  rootMargin?: string;
  triggerOnce?: boolean;
}

const useObserver = ({
  elementRef,
  threshold = 0.1,
  root = null,
  rootMargin = '0px',
  triggerOnce = true,
}: useObserverProps) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const currentElement = elementRef.current;

    if (!currentElement) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);

          if (!triggerOnce) return;

          observer.disconnect(); // Stop observing after it becomes visible
        } else if (!triggerOnce) {
          setIsVisible(false); // Keep tracking visibility if not one-time
        }
      },
      { threshold, root, rootMargin },
    );

    observer.observe(currentElement);

    return () => observer.unobserve(currentElement);
  }, [elementRef, threshold, root, rootMargin, triggerOnce]);

  return { isVisible };
};
export { useObserver };
