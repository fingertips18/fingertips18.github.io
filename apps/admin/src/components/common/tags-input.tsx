import { ChevronsUpDown, Stars, X } from 'lucide-react';
import {
  type ChangeEvent,
  type ComponentProps,
  type KeyboardEvent,
  useMemo,
  useRef,
  useState,
} from 'react';

import { Badge } from '@/components/shadcn/badge';
import { Button } from '@/components/shadcn/button';
import { Input } from '@/components/shadcn/input';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/shadcn/popover';
import { cn } from '@/lib/utils';

import { Command, CommandInput } from '../shadcn/command';

interface TagsInputProps
  extends Omit<ComponentProps<'input'>, 'value' | 'onChange' | 'onKeyDown'> {
  value?: string[];
  onChange?: (tags: string[]) => void;
  suggestions?: string[];
  className?: string;
}

export function TagsInput({
  value = [],
  onChange,
  suggestions = [],
  ...props
}: TagsInputProps) {
  const [input, setInput] = useState<string>('');
  const [tags, setTags] = useState<string[]>(value);
  const [open, setOpen] = useState<boolean>(false);
  const inputRef = useRef<HTMLInputElement>(null);

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

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const nextValue = e.target.value;
    setInput(nextValue);

    // Only open if there are suggestions
    setOpen(
      nextValue.trim().length > 0 &&
        suggestions.some(
          (s) =>
            s.toLowerCase().includes(nextValue.toLowerCase()) &&
            !tags.includes(s),
        ),
    );
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== 'Enter' || !input.trim()) return;
    e.preventDefault();
    addTag(input);
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
              className='p-1 rounded-full size-5 flex-center cursor-pointer hover:bg-accent/25 transition-colors'
            >
              <X aria-hidden='true' />
            </button>
          </Badge>
        ))}
      </div>

      <div className='flex-center gap-x-2'>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger className='w-full'>
            <Button
              variant='outline'
              role='combobox'
              aria-expanded={open}
              className='justify-between'
            >
              Select tech stack...
              <ChevronsUpDown aria-hidden='true' className='opacity-50' />
            </Button>
          </PopoverTrigger>
          <PopoverContent
            align='start'
            side='bottom'
            className={cn(
              'space-y-2 transition-opacity duration-150 p-1',
              filteredSuggestions.length > 0
                ? 'opacity-100'
                : 'opacity-0 pointer-events-none',
            )}
            style={{
              width: 'var(--radix-popover-trigger-width)',
            }}
          >
            <Command>
              <CommandInput
                ref={inputRef}
                value={input}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                {...props}
              ></CommandInput>
            </Command>

            <div className='flex-start gap-x-1 opacity-50 px-2 py-1'>
              <Stars aria-hidden='true' className='size-4' />
              <span className='text-sm font-medium'>Suggestions</span>
            </div>

            <ul className='space-y-1 list-none'>
              {filteredSuggestions.map((suggestion) => (
                <li
                  key={suggestion}
                  onMouseDown={(e) => e.preventDefault()}
                  onClick={() => {
                    addTag(suggestion);
                    setOpen(false);
                  }}
                  className='px-2 py-1 rounded-md cursor-pointer hover:bg-accent hover:text-accent-foreground'
                >
                  {suggestion}
                </li>
              ))}
            </ul>
          </PopoverContent>
        </Popover>

        <Button
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
