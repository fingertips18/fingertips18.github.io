import { useRef } from 'react';
import type { Control, FieldValues, Path } from 'react-hook-form';

import { Combobox } from '@/components/common/combobox';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormMessage,
} from '@/components/shadcn/form';
import { APIRoute } from '@/constants/api';
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
  hasError?: boolean;
}

export function Tags<T extends FieldValues>({
  control,
  name,
  hasError = false,
}: TagsProps<T>) {
  const triggerRef = useRef<HTMLButtonElement>(null);

  const { data, loading } = useFetch<[GitHubResponse, DevToTag[]]>({
    url: [APIRoute.githubTags, APIRoute.devToTags],
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
      render={({ field }) => {
        const { ref, ...fields } = field;
        void ref; // Explicitly ignore react-hook-form ref; using triggerRef instead

        return (
          <FormItem>
            {/* 
              Using a <span> instead of a native <label> because 
              Cmdk / headless Select overrides IDs and the input element, 
              so native <label htmlFor> associations don’t work correctly. 
              This span acts as the accessible label and also triggers the select 
              when clicked via the ref.
            */}
            <span
              id={field.name}
              data-slot='form-label'
              data-error={!!hasError}
              className='text-sm font-medium data-[error=true]:text-destructive w-fit cursor-default'
              onClick={() => triggerRef.current?.click()}
            >
              Tags
            </span>
            <FormDescription>
              Add tags for technologies, frameworks, languages, etc.
            </FormDescription>
            <FormControl>
              <Combobox
                triggerRef={triggerRef}
                aria-labelledby={field.name}
                placeholder='e.g. javascript, typescript, react, python, docker'
                suggestions={uniqueSuggestions}
                defaultSuggestions={defaults}
                emptyMessage='No tags found.'
                selectPlaceholder='Select tags...'
                disabled={loading}
                hasError={hasError}
                {...fields}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        );
      }}
    />
  );
}
