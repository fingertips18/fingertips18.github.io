import type { Control, FieldValues, Path } from 'react-hook-form';

import { Combobox } from '@/components/common/combobox';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/shadcn/form';
import { useFetch } from '@/hooks/useFetch';

type GitHubResponse = { items: { topics: string[] }[] };
type DevToTag = { name: string };

interface TagsProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Tags<T extends FieldValues>({ control, name }: TagsProps<T>) {
  const { data, loading } = useFetch<GitHubResponse | DevToTag[]>({
    url: [
      '/github/search/repositories?q=stars:>500&sort=stars&order=desc',
      'https://dev.to/api/tags',
    ],
    method: 'GET',
    toastOptions: {
      errorTitle: 'Failed to load tags suggestions',
    },
  });

  const suggestions: string[] = [];

  if (data) {
    // Handle GitHub response
    if ('items' in data) {
      data.items.forEach((repo) => {
        suggestions.push(...repo.topics);
      });
    }

    // Handle Dev.to response
    if (Array.isArray(data) && data.length && 'name' in data[0]) {
      data.forEach((tag) => {
        suggestions.push(tag.name);
      });
    }
  }

  const uniqueSuggestions = Array.from(new Set(suggestions));

  return (
    <FormField
      control={control}
      name={name}
      disabled={loading}
      render={({ field }) => (
        <FormItem className='w-full'>
          <FormLabel>Tags</FormLabel>
          <FormDescription>
            Add any tags relevant to your project, including technologies,
            frameworks, languages, etc.
          </FormDescription>
          <FormControl>
            <Combobox
              placeholder='e.g. js, ts, python, ruby, go, c#, java'
              suggestions={uniqueSuggestions}
              defaultSuggestions={[
                'js',
                'ts',
                'python',
                'ruby',
                'go',
                'c#',
                'java',
                'c++',
                'php',
                'rust',
                'kotlin',
                'swift',
                'dart',
                'scala',
                'elixir',
                'react',
                'vue',
                'angular',
                'svelte',
                'node',
                'express',
                'graphql',
                'docker',
                'kubernetes',
                'aws',
              ]}
              emptyMessage='No tag found.'
              selectPlaceholder='Select tags...'
              disabled={loading}
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
