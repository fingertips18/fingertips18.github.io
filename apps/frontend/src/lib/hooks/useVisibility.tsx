import { useEffect, useState } from 'react';

const useVisibility = () => {
  const [isVisible, setIsVisible] = useState(true);

  const handleVisibilityChange = () => {
    if (document.visibilityState === 'hidden') {
      setIsVisible(false);
    } else {
      setIsVisible(true);
    }
  };

  useEffect(() => {
    document.addEventListener('visibilitychange', handleVisibilityChange);

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, []);

  return { isVisible };
};

export { useVisibility };
