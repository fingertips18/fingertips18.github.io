import { useEffect, useState } from 'react';

const useElementsByQuery = (query: string) => {
  const [elements, setElements] = useState<NodeListOf<HTMLElement>>();

  useEffect(() => {
    setElements(document.querySelectorAll(query));
  }, [query]);

  return elements;
};

export { useElementsByQuery };
