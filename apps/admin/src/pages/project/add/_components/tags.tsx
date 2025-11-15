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

const defaults = [
  'typescript',
  'javascript',
  'react',
  'vue',
  'angular',
  'svelte',
  'next.js',
  'nuxt',
  'node.js',
  'deno',
  'bun',
  'express',
  'fastify',
  'nest.js',
  'python',
  'django',
  'flask',
  'fastapi',
  'java',
  'spring',
  'kotlin',
  'c#',
  '.net',
  'go',
  'rust',
  'php',
  'laravel',
  'ruby',
  'rails',
  'swift',
  'dart',
  'flutter',
  'react-native',
  'tailwindcss',
  'css',
  'html',
  'sql',
  'postgresql',
  'mysql',
  'mongodb',
  'redis',
  'graphql',
  'rest-api',
  'docker',
  'kubernetes',
  'aws',
  'azure',
  'gcp',
];

interface TagsProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
}

export function Tags<T extends FieldValues>({ control, name }: TagsProps<T>) {
  const { data, loading } = useFetch<[GitHubResponse, DevToTag[]]>({
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
    const [githubData, devToData] = data;
    const githubTopics = githubData?.items.flatMap((repo) => repo.topics) ?? [];
    const devToTags = devToData?.map((tag) => tag.name) ?? [];
    suggestions.push(...githubTopics.map((t) => t.toLowerCase()));
    suggestions.push(...devToTags.map((t) => t.toLowerCase()));
  }

  const uniqueSuggestions = Array.from(new Set(suggestions));

  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>Tags</FormLabel>
          <FormControl>
            <Combobox
              placeholder='e.g. javascript, typescript, react, python, docker'
              suggestions={uniqueSuggestions}
              defaultSuggestions={defaults}
              emptyMessage='No tags found.'
              selectPlaceholder='Select tags...'
              disabled={loading}
              {...field}
            />
          </FormControl>
          <FormDescription>
            Add tags for technologies, frameworks, languages, etc.
          </FormDescription>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
