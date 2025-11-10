import { Check, ChevronsUpDown, X } from 'lucide-react';
import { type ComponentProps, useMemo, useState } from 'react';

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
  defaultSuggestions?: string[];
  className?: string;
  emptyMessage?: string;
  selectPlaceholder?: string;
}

export function Combobox({
  value: tags = [],
  onChange,
  suggestions = [],
  defaultSuggestions = [],
  emptyMessage = 'No results found.',
  selectPlaceholder = 'Select...',
  ...props
}: ComboboxProps) {
  const [input, setInput] = useState<string>('');
  const [open, setOpen] = useState<boolean>(false);

  const addTag = (tag: string) => {
    if (!tag || tags.some((t) => t.toLowerCase() === tag.toLowerCase())) return;
    const newTags = [...tags, tag];
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
      ...new Set([...suggestions, ...defaultSuggestions]),
    ];
    return allSuggestions.filter(
      (s) => s.toLowerCase().includes(input.toLowerCase()) && !tags.includes(s),
    );
  }, [input, suggestions, defaultSuggestions, tags]);

  return (
    <div className='space-y-2'>
      <div className='flex flex-wrap gap-2'>
        {tags.map((tag) => (
          <Badge key={tag} className='flex items-center'>
            {tag}
            <button
              type='button'
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
          type='button'
          onClick={() => addTag(input)}
          disabled={!input.trim()}
          className='cursor-pointer'
        >
          Add Tag
        </Button>
      </div>
    </div>
  );
}
