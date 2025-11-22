import { Check, ChevronsUpDown, X } from 'lucide-react';
import { type ComponentProps, type Ref, useMemo, useState } from 'react';

import { Badge } from '@/components/shadcn/badge';
import { Button } from '@/components/shadcn/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/shadcn/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/shadcn/popover';
import { cn } from '@/lib/utils';

interface ComboboxProps
  extends Omit<
    ComponentProps<typeof CommandInput>,
    'value' | 'onChange' | 'onKeyDown'
  > {
  value?: string[];
  onChange?: (tags: string[]) => void;
  suggestions?: string[];
  defaultSuggestions?: string[];
  className?: string;
  emptyMessage?: string;
  selectPlaceholder?: string;
  hasError?: boolean;
  triggerRef?: Ref<HTMLButtonElement>;
}

/**
 * Combobox component for multi-tag input with suggestions.
 * Note: All tags are normalized to lowercase for consistent comparison and storage.
 */
export function Combobox({
  value: tags = [],
  onChange,
  suggestions = [],
  defaultSuggestions = [],
  emptyMessage = 'No results found.',
  selectPlaceholder = 'Select...',
  hasError,
  triggerRef,
  disabled,
  ...props
}: ComboboxProps) {
  const [input, setInput] = useState<string>('');
  const [open, setOpen] = useState<boolean>(false);

  const addTag = (tag: string) => {
    const normalized = tag.toLowerCase();
    if (!normalized || tags.some((t) => t.toLowerCase() === normalized)) return;
    const newTags = [...tags, normalized];
    setInput('');
    onChange?.(newTags);
  };

  const removeTag = (tag: string) => {
    const newTags = tags.filter((t) => t !== tag);
    onChange?.(newTags);
  };

  const filteredSuggestions = useMemo(() => {
    if (!input) return defaultSuggestions;

    const allSuggestions = [
      ...new Set(
        [...suggestions, ...defaultSuggestions].filter(
          (s) => s.trim().length > 0,
        ),
      ),
    ];
    return allSuggestions.filter(
      (s) =>
        s.toLowerCase().includes(input.toLowerCase()) &&
        !tags.some((t) => t.toLowerCase() === s.toLowerCase()),
    );
  }, [input, suggestions, defaultSuggestions, tags]);

  return (
    <div className='space-y-2'>
      <div className='flex flex-wrap gap-2'>
        {tags.map((tag) => (
          <Badge
            key={tag}
            className={cn('flex items-center', disabled && 'opacity-50')}
          >
            {tag}
            <button
              type='button'
              onClick={() => removeTag(tag)}
              disabled={disabled}
              aria-label={`Remove ${tag}`}
              className='p-1 rounded-full size-5 flex-center cursor-pointer hover:bg-accent/25 transition-colors'
            >
              <X aria-hidden='true' />
            </button>
          </Badge>
        ))}
      </div>

      <div
        className={cn('flex-center gap-x-2', disabled && 'cursor-not-allowed')}
      >
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              data-state={open ? 'open' : 'closed'}
              ref={triggerRef}
              type='button'
              variant='outline'
              role='combobox'
              aria-expanded={open}
              disabled={disabled}
              className={cn(
                'justify-between flex-1 text-muted-foreground data-[state=open]:border-ring data-[state=open]:ring-ring/50 data-[state=open]:ring-[3px]',
                hasError && 'border-destructive',
              )}
            >
              {selectPlaceholder}
              <ChevronsUpDown aria-hidden='true' className='opacity-50' />
            </Button>
          </PopoverTrigger>
          <PopoverContent
            className='p-0'
            style={{
              width: 'var(--radix-popover-trigger-width)',
            }}
          >
            <Command>
              <CommandInput
                value={input}
                onValueChange={(value) => setInput(value)}
                {...props}
              />
              <CommandList>
                <CommandEmpty>{emptyMessage}</CommandEmpty>
                <CommandGroup>
                  {filteredSuggestions.map((suggestion) => (
                    <CommandItem
                      key={suggestion}
                      value={suggestion}
                      onSelect={(value) => addTag(value)}
                    >
                      {suggestion}
                      <Check
                        aria-hidden='true'
                        className={cn(
                          'ml-auto',
                          tags.some(
                            (t) => t.toLowerCase() === suggestion.toLowerCase(),
                          )
                            ? 'opacity-100'
                            : 'opacity-0',
                        )}
                      />
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>

        <Button
          type='button'
          onClick={() => addTag(input)}
          disabled={disabled || !input.trim()}
          className='cursor-pointer'
        >
          Add Tag
        </Button>
      </div>
    </div>
  );
}
