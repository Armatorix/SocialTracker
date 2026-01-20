export interface User {
  id: number;
  user_id: string;
  email: string;
  username: string;
  role: 'admin' | 'creator';
  created_at: string;
  updated_at: string;
}

export interface SocialAccount {
  id: number;
  user_id: number;
  platform: 'twitter' | 'facebook' | 'instagram' | 'youtube' | 'tiktok';
  account_name: string;
  account_id?: string;
  last_pull_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Content {
  id: number;
  user_id: number;
  social_account_id?: number;
  platform: 'twitter' | 'facebook' | 'instagram' | 'youtube' | 'tiktok';
  link: string;
  original_text?: string;
  description?: string;
  tags?: string[];
  external_post_id?: string;
  posted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ContentWithUser extends Content {
  username: string;
  email: string;
}

export interface CreateSocialAccountRequest {
  platform: string;
  account_name: string;
  account_id?: string;
  access_token?: string;
  refresh_token?: string;
}

export interface CreateContentRequest {
  social_account_id?: number;
  platform: string;
  link: string;
  original_text?: string;
  description?: string;
  tags?: string[];
}

export interface SyncResponse {
  account_id: number;
  platform: string;
  account_name: string;
  synced_count: number;
  skipped_count: number;
  errors?: string[];
  message: string;
}
