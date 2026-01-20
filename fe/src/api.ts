import type { User, SocialAccount, Content, ContentWithUser, CreateSocialAccountRequest, CreateContentRequest, SyncResponse } from './types';

const API_BASE_URL = '/api';

// Helper to format retry time in a human-readable way
const formatRetryTime = (seconds: number): string => {
  if (seconds < 60) {
    return `${seconds} second${seconds !== 1 ? 's' : ''}`;
  }
  const minutes = Math.ceil(seconds / 60);
  return `${minutes} minute${minutes !== 1 ? 's' : ''}`;
};

// Custom error class for rate limit errors
export class RateLimitError extends Error {
  retryAfter: number;
  isRateLimitError = true as const;
  
  constructor(message: string, retryAfter: number) {
    super(message);
    this.name = 'RateLimitError';
    this.retryAfter = retryAfter;
  }
}

// Type guard for rate limit errors (more reliable than instanceof)
export function isRateLimitError(err: unknown): err is RateLimitError {
  return (
    err !== null &&
    typeof err === 'object' &&
    'isRateLimitError' in err &&
    (err as RateLimitError).isRateLimitError === true
  );
}

// Helper to make fetch requests with credentials included
const fetchWithCredentials = (url: string, options: RequestInit = {}): Promise<Response> => {
  return fetch(url, {
    ...options,
    credentials: 'include',
  });
};

export const api = {
  // User
  getCurrentUser: async (): Promise<User> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/user`);
    if (!res.ok) throw new Error('Failed to fetch user');
    return res.json();
  },

  logout: (): void => {
    // Redirect to OAuth2 proxy sign out endpoint
    window.location.href = '/oauth2/sign_out?rd=/';
  },

  // Social Accounts
  getSocialAccounts: async (): Promise<SocialAccount[]> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/social-accounts`);
    if (!res.ok) throw new Error('Failed to fetch social accounts');
    return res.json();
  },

  createSocialAccount: async (data: CreateSocialAccountRequest): Promise<SocialAccount> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/social-accounts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create social account');
    return res.json();
  },

  deleteSocialAccount: async (id: number): Promise<void> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/social-accounts/${id}`, {
      method: 'DELETE',
    });
    if (!res.ok) throw new Error('Failed to delete social account');
  },

  pullContent: async (accountId: number): Promise<SyncResponse> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/social-accounts/${accountId}/pull`, {
      method: 'POST',
    });
    if (res.status === 429) {
      const data = await res.json();
      const retryAfter = data.retry_after || 900;
      throw new RateLimitError(
        `X/Twitter rate limit exceeded. Please try again in ${formatRetryTime(retryAfter)}.`,
        retryAfter
      );
    }
    if (!res.ok) {
      // Try to get error details from response
      try {
        const data = await res.json();
        const errorMsg = data.error || 'Failed to pull content';
        // Check if error message indicates rate limiting
        if (errorMsg.toLowerCase().includes('rate limit') || errorMsg.includes('429') || errorMsg.includes('too many')) {
          throw new RateLimitError(
            `X/Twitter rate limit exceeded. Please try again later.`,
            900
          );
        }
        throw new Error(errorMsg);
      } catch (parseErr) {
        if (parseErr instanceof RateLimitError) throw parseErr;
        if (isRateLimitError(parseErr)) throw parseErr;
        throw new Error('Failed to pull content');
      }
    }
    return res.json();
  },

  // Content
  getContent: async (): Promise<Content[]> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/content`);
    if (!res.ok) throw new Error('Failed to fetch content');
    return res.json();
  },

  createContent: async (data: CreateContentRequest): Promise<Content> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/content`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create content');
    return res.json();
  },

  deleteContent: async (id: number): Promise<void> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/content/${id}`, {
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
    const res = await fetchWithCredentials(url);
    if (!res.ok) throw new Error('Failed to fetch all content');
    return res.json();
  },

  // Twitter OAuth
  getTwitterOAuthStatus: async (): Promise<{ configured: boolean }> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/auth/twitter/status`);
    if (!res.ok) throw new Error('Failed to get Twitter OAuth status');
    return res.json();
  },

  getTwitterOAuthURL: async (): Promise<{ url: string }> => {
    const res = await fetchWithCredentials(`${API_BASE_URL}/auth/twitter`);
    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.error || 'Failed to get Twitter OAuth URL');
    }
    return res.json();
  },

  connectTwitter: async (): Promise<void> => {
    const { url } = await api.getTwitterOAuthURL();
    window.location.href = url;
  },
};
