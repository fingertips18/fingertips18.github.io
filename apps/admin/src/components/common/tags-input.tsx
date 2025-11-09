import { type ComponentProps, useState } from 'react';

export function TagsInput(props: ComponentProps<'input'>) {
  const [values, setValues] = useState<string[]>([]);
  const [currentInput, setCurrentInput] = useState<string>('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [filteredSuggestions, setFilteredSuggestions] = useState<string[]>([]);
}
