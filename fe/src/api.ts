import type { User, SocialAccount, Content, ContentWithUser, CreateSocialAccountRequest, CreateContentRequest } from './types';

const API_BASE_URL = 'http://localhost:8080/api';

export const api = {
  // User
  getCurrentUser: async (): Promise<User> => {
    const res = await fetch(`${API_BASE_URL}/user`);
    if (!res.ok) throw new Error('Failed to fetch user');
    return res.json();
  },

  // Social Accounts
  getSocialAccounts: async (): Promise<SocialAccount[]> => {
    const res = await fetch(`${API_BASE_URL}/social-accounts`);
    if (!res.ok) throw new Error('Failed to fetch social accounts');
    return res.json();
  },

  createSocialAccount: async (data: CreateSocialAccountRequest): Promise<SocialAccount> => {
    const res = await fetch(`${API_BASE_URL}/social-accounts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create social account');
    return res.json();
  },

  deleteSocialAccount: async (id: number): Promise<void> => {
    const res = await fetch(`${API_BASE_URL}/social-accounts/${id}`, {
      method: 'DELETE',
    });
    if (!res.ok) throw new Error('Failed to delete social account');
  },

  pullContent: async (accountId: number): Promise<{ message: string }> => {
    const res = await fetch(`${API_BASE_URL}/social-accounts/${accountId}/pull`, {
      method: 'POST',
    });
    if (!res.ok) throw new Error('Failed to pull content');
    return res.json();
  },

  // Content
  getContent: async (): Promise<Content[]> => {
    const res = await fetch(`${API_BASE_URL}/content`);
    if (!res.ok) throw new Error('Failed to fetch content');
    return res.json();
  },

  createContent: async (data: CreateContentRequest): Promise<Content> => {
    const res = await fetch(`${API_BASE_URL}/content`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create content');
    return res.json();
  },

  deleteContent: async (id: number): Promise<void> => {
    const res = await fetch(`${API_BASE_URL}/content/${id}`, {
      method: 'DELETE',
    });
    if (!res.ok) throw new Error('Failed to delete content');
  },

  // Admin
  getAllContent: async (filters?: { platform?: string; username?: string }): Promise<ContentWithUser[]> => {
    const params = new URLSearchParams();
    if (filters?.platform) params.append('platform', filters.platform);
    if (filters?.username) params.append('username', filters.username);
    
    const url = `${API_BASE_URL}/admin/content${params.toString() ? '?' + params.toString() : ''}`;
    const res = await fetch(url);
    if (!res.ok) throw new Error('Failed to fetch all content');
    return res.json();
  },
};
