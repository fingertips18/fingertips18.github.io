import { Check, ChevronsUpDown, X } from 'lucide-react';
import { type ComponentProps, useEffect, useMemo, useState } from 'react';

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
  extends Omit<ComponentProps<'input'>, 'value' | 'onChange' | 'onKeyDown'> {
  value?: string[];
  onChange?: (tags: string[]) => void;
  suggestions?: string[];
  defaultSuggestions: string[];
  className?: string;
  emptyMessage: string;
  selectPlaceholder?: string;
}

export function Combobox({
  value = [],
  onChange,
  suggestions = [],
  defaultSuggestions,
  emptyMessage,
  selectPlaceholder = 'Select...',
  ...props
}: ComboboxProps) {
  const [input, setInput] = useState<string>('');
  const [tags, setTags] = useState<string[]>(value);
  const [open, setOpen] = useState<boolean>(false);

  useEffect(() => {
    setTags(value);
  }, [value]);

  const addTag = (tag: string) => {
    if (!tag || tags.includes(tag)) return;
    const newTags = [...tags, tag];
    setTags(newTags);
    setInput('');
    onChange?.(newTags);
  };

  const removeTag = (tag: string) => {
    const newTags = tags.filter((t) => t !== tag);
    setTags(newTags);
    onChange?.(newTags);
  };

  const filteredSuggestions = useMemo(() => {
    if (!input) return [];
    return suggestions.filter(
      (s) => s.toLowerCase().includes(input.toLowerCase()) && !tags.includes(s),
    );
  }, [input, suggestions, tags]);

  return (
    <div className='space-y-2'>
      <div className='flex flex-wrap gap-2'>
        {tags.map((tag) => (
          <Badge key={tag} className='flex items-center'>
            {tag}
            <button
              onClick={() => removeTag(tag)}
              aria-label={`Remove ${tag}`}
              className='p-1 rounded-full size-5 flex-center cursor-pointer hover:bg-accent/25 transition-colors'
            >
              <X aria-hidden='true' />
            </button>
          </Badge>
        ))}
      </div>

      <div className='flex-center gap-x-2'>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant='outline'
              role='combobox'
              aria-expanded={open}
              className='justify-between flex-1'
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
                  {(filteredSuggestions.length > 0
                    ? filteredSuggestions
                    : defaultSuggestions
                  ).map((suggestion) => (
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
                          tags.includes(suggestion)
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
          onClick={() => addTag(input.trim().toLowerCase())}
          disabled={!input.trim()}
          className='cursor-pointer'
        >
          Add Tag
        </Button>
      </div>
    </div>
  );
}
