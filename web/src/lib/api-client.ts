import type {
  GitHubProfile,
  Repository,
  Subcommit,
  UserRepository,
} from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit & { params?: Record<string, string> } = {},
): Promise<T> {
  const { params, ...fetchOptions } = options;

  let url = `${API_URL}${endpoint}`;
  if (params) {
    url += `?${new URLSearchParams(params).toString()}`;
  }

  const response = await fetch(url, {
    ...fetchOptions,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...fetchOptions.headers,
    },
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: "Unknown error" }));
    throw new ApiError(response.status, body.error || "Request failed");
  }

  return response.json();
}

export const api = {
  getLoginUrl: () => `${API_URL}/auth/github/login`,

  authStatus: () => request<{ isLoggedIn: boolean }>("/auth/status"),

  logout: () => request<void>("/auth/logout", { method: "POST" }),

  getProfile: () => request<GitHubProfile>("/user/profile"),

  getRepositories: () =>
    request<{ repositories: UserRepository[] }>("/repositories"),

  searchRepos: (query: string) =>
    request<{ repositories: Repository[] }>("/user/repos/search", {
      params: { q: query },
    }),

  analyzeRepository: (repoUrl: string) =>
    request<{ message: string; repoId: number }>("/analyze", {
      method: "POST",
      body: JSON.stringify({ repoUrl }),
    }),

  getSubcommitsTimeline: (repoId: string) =>
    request<{ subcommits: Subcommit[]; repoId: string }>(
      "/subcommits-timeline",
      { params: { repo_id: repoId } },
    ),
};
