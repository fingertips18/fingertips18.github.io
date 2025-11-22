import { useCallback, useEffect } from 'react';
import { useBlocker } from 'react-router-dom';

interface UnsavedChangesProps {
  hasUnsavedChanges: boolean;
  promptMessage?: string;
}

/**
 * Hook that prevents accidental navigation or page unload when there are unsaved changes.
 *
 * @param options - Configuration for the hook.
 * @param options.hasUnsavedChanges - When true, in-app navigation is intercepted and a confirmation dialog is shown; a beforeunload listener is also registered to prompt on page refresh/close.
 * @param options.promptMessage - Optional message used for the in-page confirmation and assigned to the beforeunload event's returnValue.
 *   Defaults to "You have unsaved changes. Are you sure you want to leave?".
 *
 * @remarks
 * - Uses an in-app navigation blocker (via useBlocker) to call window.confirm when navigation is attempted and `hasUnsavedChanges` is true.
 * - Registers a `beforeunload` event listener that calls `preventDefault()` and sets `event.returnValue` to `promptMessage` when `hasUnsavedChanges` is true.
 * - Cleans up the `beforeunload` listener automatically on unmount or when dependencies change.
 * - Note: Many modern browsers ignore custom strings set on `beforeunload` and display a generic message instead.
 * - This hook must be used within a React component (client-side) and its behavior may depend on the routing/navigation library providing `useBlocker`.
 *
 * @returns void
 */
export function useUnsavedChanges({
  hasUnsavedChanges,
  promptMessage = 'You have unsaved changes. Are you sure you want to leave?',
}: UnsavedChangesProps): void {
  useBlocker(() => {
    if (hasUnsavedChanges) {
      const confirmLeave = window.confirm(promptMessage);
      return !confirmLeave;
    }

    return false;
  });

  const handleBeforeUnload = useCallback(
    (e: BeforeUnloadEvent) => {
      if (hasUnsavedChanges) {
        e.preventDefault();
        e.returnValue = promptMessage;
      }
    },
    [hasUnsavedChanges, promptMessage],
  );

  useEffect(() => {
    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [handleBeforeUnload]);
}
