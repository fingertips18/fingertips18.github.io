import { createJSONStorage, persist } from "zustand/middleware";
import { create } from "zustand";

import { ROOTSECTIONS } from "@/constants/enums";

interface RootSectionStore {
  active: string;
  onActive: (active: string) => void;
  onClear: () => void;
}

const useRootSectionStore = create(
  persist<RootSectionStore>(
    (set) => ({
      active: ROOTSECTIONS.about,
      onActive: (active: string) => set({ active }),
      onClear: () => set({ active: ROOTSECTIONS.about }),
    }),
    { name: "root-section", storage: createJSONStorage(() => sessionStorage) }
  )
);

export { useRootSectionStore };
