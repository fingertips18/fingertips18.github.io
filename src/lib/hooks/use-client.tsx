import { useEffect, useState } from "react";

export function useClient() {
  const [mounted, setMounted] = useState(false);

  useEffect(() => setMounted(true), []);

  return mounted;
}
