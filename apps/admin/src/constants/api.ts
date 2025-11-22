const api = '/api';

export const APIRoute = {
  githubTags: '/github/search/repositories?q=stars:>500&sort=stars&order=desc',
  devToTags: 'https://dev.to/api/tags',
  project: `${api}/project`,
  image: `${api}/image`,
} as const;
