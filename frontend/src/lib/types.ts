export type SubcommitType =
  | "FEATURE"
  | "BUG"
  | "REFACTOR"
  | "DOCS"
  | "CHORE"
  | "MILESTONE"
  | "WARNING";

export interface Subcommit {
  id: number;
  createdAt: string;
  title: string;
  idea: string;
  description: string;
  commitSha: string;
  type: SubcommitType;
  epic: string;
  files: string[];
}

export interface Repository {
  id: number;
  createdAt: string;
  name: string;
  url: string;
  lastAnalyzedCommit: string;
}

export interface GitHubProfile {
  id: number;
  login: string;
  name: string;
  avatarUrl: string;
  email: string;
}

export interface UserRepository {
  id: string;
  name: string;
  url: string;
  addedAt: string;
}
