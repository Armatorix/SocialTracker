import { useState, useEffect } from 'react';
import { api } from './api';
import type { SocialAccount, Content } from './types';

const platformStyles: Record<string, string> = {
  twitter: 'bg-sky-500',
  facebook: 'bg-blue-600',
  instagram: 'bg-gradient-to-tr from-yellow-400 via-pink-500 to-purple-600',
  youtube: 'bg-red-600',
  tiktok: 'bg-black',
};

const platformIcons: Record<string, JSX.Element> = {
  twitter: (
    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
    </svg>
  ),
  facebook: (
    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
      <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
    </svg>
  ),
  instagram: (
    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
      <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
    </svg>
  ),
  youtube: (
    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
      <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
    </svg>
  ),
  tiktok: (
    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
      <path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z"/>
    </svg>
  ),
};

export function CreatorDashboard() {
  const [socialAccounts, setSocialAccounts] = useState<SocialAccount[]>([]);
  const [content, setContent] = useState<Content[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  const [showAccountForm, setShowAccountForm] = useState(false);
  const [showContentForm, setShowContentForm] = useState(false);
  const [accountForm, setAccountForm] = useState({ platform: 'twitter', account_name: '' });
  const [contentForm, setContentForm] = useState({
    platform: 'twitter',
    link: '',
    original_text: '',
    description: '',
    tags: '',
  });

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');
      const [accountsData, contentData] = await Promise.all([
        api.getSocialAccounts(),
        api.getContent(),
      ]);
      setSocialAccounts(accountsData);
      setContent(contentData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleAddAccount = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createSocialAccount(accountForm);
      setAccountForm({ platform: 'twitter', account_name: '' });
      setShowAccountForm(false);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add account');
    }
  };

  const handleDeleteAccount = async (id: number) => {
    if (!confirm('Are you sure you want to delete this account?')) return;
    try {
      await api.deleteSocialAccount(id);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete account');
    }
  };

  const handlePullContent = async (accountId: number) => {
    try {
      const result = await api.pullContent(accountId);
      alert(result.message);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to pull content');
    }
  };

  const handleAddContent = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const tags = contentForm.tags ? contentForm.tags.split(',').map(t => t.trim()) : [];
      await api.createContent({
        ...contentForm,
        tags,
        original_text: contentForm.original_text || undefined,
        description: contentForm.description || undefined,
      });
      setContentForm({ platform: 'twitter', link: '', original_text: '', description: '', tags: '' });
      setShowContentForm(false);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add content');
    }
  };

  const handleDeleteContent = async (id: number) => {
    if (!confirm('Are you sure you want to delete this content?')) return;
    try {
      await api.deleteContent(id);
      loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete content');
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleString();
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="w-10 h-10 border-4 border-purple-500/30 border-t-purple-500 rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      {/* Header */}
      <div className="mb-8 animate-fade-in">
        <h1 className="text-3xl sm:text-4xl font-bold mb-2 bg-gradient-to-r from-purple-400 to-cyan-400 bg-clip-text text-transparent">
          Creator Dashboard
        </h1>
        <p className="text-slate-300">Manage your social accounts and track your content</p>
      </div>
      
      {error && (
        <div className="bg-red-500/10 border border-red-500/30 rounded-2xl p-4 mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-red-500/20 flex items-center justify-center shrink-0">
              <svg className="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <p className="text-red-400">{error}</p>
          </div>
          <button onClick={() => setError('')} className="text-slate-200 hover:text-white transition-colors p-1">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      )}

      {/* Social Accounts Section */}
      <section className="mb-10">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-purple-700 flex items-center justify-center shadow-lg">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
            </div>
            <div>
              <h2 className="text-xl font-bold text-white">Connected Accounts</h2>
              <p className="text-sm text-slate-600">{socialAccounts.length} accounts linked</p>
            </div>
          </div>
          <button
            onClick={() => setShowAccountForm(!showAccountForm)}
            className={`px-5 py-2.5 rounded-xl font-semibold transition-all ${
              showAccountForm 
                ? 'bg-slate-700 text-white border border-slate-600 hover:bg-slate-600' 
                : 'bg-gradient-to-r from-purple-500 to-purple-700 text-white shadow-lg shadow-purple-500/25 hover:shadow-purple-500/40'
            }`}
          >
            {showAccountForm ? '✕ Cancel' : '+ Add Account'}
          </button>
        </div>

        {showAccountForm && (
          <form onSubmit={handleAddAccount} className="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-6 mb-6 border border-slate-700/50 animate-fade-in">
            <h3 className="text-lg font-semibold mb-4 text-white">Link New Account</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-100 mb-2">Platform</label>
                <select
                  value={accountForm.platform}
                  onChange={(e) => setAccountForm({ ...accountForm, platform: e.target.value })}
                  className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  required
                >
                  <option value="twitter">Twitter/X</option>
                  <option value="facebook">Facebook</option>
                  <option value="instagram">Instagram</option>
                  <option value="youtube">YouTube</option>
                  <option value="tiktok">TikTok</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-100 mb-2">Account Name</label>
                <input
                  type="text"
                  value={accountForm.account_name}
                  onChange={(e) => setAccountForm({ ...accountForm, account_name: e.target.value })}
                  className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  placeholder="@username"
                  required
                />
              </div>
            </div>
            <button type="submit" className="mt-4 px-5 py-2.5 rounded-xl font-semibold bg-gradient-to-r from-purple-500 to-purple-700 text-white shadow-lg shadow-purple-500/25 hover:shadow-purple-500/40 transition-all">
              Link Account
            </button>
          </form>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {socialAccounts.length === 0 ? (
            <div className="col-span-full bg-slate-800/50 backdrop-blur-sm rounded-2xl p-8 text-center border border-slate-700/50">
              <div className="w-16 h-16 rounded-full bg-slate-700/50 flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-slate-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                </svg>
              </div>
              <p className="text-slate-300 mb-2">No accounts connected yet</p>
              <p className="text-sm text-slate-100">Click "Add Account" to link your social media profiles</p>
            </div>
          ) : (
            socialAccounts.map((account) => (
              <div key={account.id} className="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-5 border border-slate-700/50 hover:border-purple-500/50 transition-all hover:shadow-lg hover:shadow-purple-500/10">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex items-center gap-3">
                    <div className={`w-12 h-12 rounded-xl flex items-center justify-center text-white shadow-lg ${platformStyles[account.platform] || 'bg-slate-600'}`}>
                      {platformIcons[account.platform]}
                    </div>
                    <div>
                      <h3 className="font-bold capitalize text-white">{account.platform}</h3>
                      <p className="text-sm text-slate-100">{account.account_name}</p>
                    </div>
                  </div>
                  <button
                    onClick={() => handleDeleteAccount(account.id)}
                    className="p-2 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 hover:bg-red-500/20 transition-all"
                    title="Delete"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
                <div className="flex items-center gap-2 text-xs text-slate-100 mb-4">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Last synced: {formatDate(account.last_pull_at)}
                </div>
                <button
                  onClick={() => handlePullContent(account.id)}
                  className="w-full py-2.5 rounded-xl font-medium bg-slate-700 text-white border border-slate-600 hover:bg-slate-600 hover:border-purple-500/50 transition-all flex items-center justify-center gap-2"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  Sync Content
                </button>
              </div>
            ))
          )}
        </div>
      </section>

      {/* Content Section */}
      <section>
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-cyan-500 to-cyan-700 flex items-center justify-center shadow-lg">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <h2 className="text-xl font-bold text-white">My Content</h2>
              <p className="text-sm text-slate-100">{content.length} items tracked</p>
            </div>
          </div>
          <button
            onClick={() => setShowContentForm(!showContentForm)}
            className={`px-5 py-2.5 rounded-xl font-semibold transition-all ${
              showContentForm 
                ? 'bg-slate-700 text-white border border-slate-600 hover:bg-slate-600' 
                : 'bg-gradient-to-r from-cyan-500 to-cyan-700 text-white shadow-lg shadow-cyan-500/25 hover:shadow-cyan-500/40'
            }`}
          >
            {showContentForm ? '✕ Cancel' : '+ Add Content'}
          </button>
        </div>

        {showContentForm && (
          <form onSubmit={handleAddContent} className="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-6 mb-6 border border-slate-700/50 animate-fade-in">
            <h3 className="text-lg font-semibold mb-4 text-white">Add New Content</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium text-slate-100 mb-2">Platform</label>
                <select
                  value={contentForm.platform}
                  onChange={(e) => setContentForm({ ...contentForm, platform: e.target.value })}
                  className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                  required
                >
                  <option value="twitter">Twitter/X</option>
                  <option value="facebook">Facebook</option>
                  <option value="instagram">Instagram</option>
                  <option value="youtube">YouTube</option>
                  <option value="tiktok">TikTok</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-100 mb-2">Link (URL) *</label>
                <input
                  type="url"
                  value={contentForm.link}
                  onChange={(e) => setContentForm({ ...contentForm, link: e.target.value })}
                  className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                  placeholder="https://..."
                  required
                />
              </div>
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-slate-100 mb-2">Original Text/Caption</label>
              <textarea
                value={contentForm.original_text}
                onChange={(e) => setContentForm({ ...contentForm, original_text: e.target.value })}
                className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent resize-none"
                rows={3}
                placeholder="The original post text or caption..."
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-slate-100 mb-2">Short Description</label>
              <textarea
                value={contentForm.description}
                onChange={(e) => setContentForm({ ...contentForm, description: e.target.value })}
                className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent resize-none"
                rows={2}
                placeholder="What is this content about..."
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-slate-100 mb-2">Tags (comma-separated)</label>
              <input
                type="text"
                value={contentForm.tags}
                onChange={(e) => setContentForm({ ...contentForm, tags: e.target.value })}
                className="w-full bg-slate-900/50 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent"
                placeholder="tag1, tag2, tag3"
              />
            </div>
            <button type="submit" className="px-5 py-2.5 rounded-xl font-semibold bg-gradient-to-r from-cyan-500 to-cyan-700 text-white shadow-lg shadow-cyan-500/25 hover:shadow-cyan-500/40 transition-all">
              Add Content
            </button>
          </form>
        )}

        <div className="space-y-4">
          {content.length === 0 ? (
            <div className="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-8 text-center border border-slate-700/50">
              <div className="w-16 h-16 rounded-full bg-slate-700/50 flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-slate-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <p className="text-slate-300 mb-2">No content yet</p>
              <p className="text-sm text-slate-100">Add content manually or sync from your connected accounts</p>
            </div>
          ) : (
            content.map((item) => (
              <div key={item.id} className="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-5 border border-slate-700/50 hover:border-cyan-500/50 transition-all hover:shadow-lg hover:shadow-cyan-500/10">
                <div className="flex justify-between items-start gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-wrap items-center gap-2 mb-3">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold text-white ${platformStyles[item.platform] || 'bg-slate-600'}`}>
                        {item.platform}
                      </span>
                      <span className="text-xs text-slate-100 flex items-center gap-1">
                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {formatDate(item.created_at)}
                      </span>
                    </div>
                    <a
                      href={item.link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-cyan-400 hover:text-cyan-300 hover:underline block mb-2 truncate"
                    >
                      {item.link}
                    </a>
                    {item.description && (
                      <p className="text-sm text-slate-300 mb-2">{item.description}</p>
                    )}
                    {item.original_text && (
                      <p className="text-sm text-slate-100 mb-2 italic line-clamp-2">"{item.original_text}"</p>
                    )}
                    {item.tags && item.tags.length > 0 && (
                      <div className="flex flex-wrap gap-2 mt-3">
                        {item.tags.map((tag, i) => (
                          <span key={i} className="px-2 py-1 rounded-full text-xs bg-slate-700 text-slate-300 border border-slate-600">
                            #{tag}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                  <button
                    onClick={() => handleDeleteContent(item.id)}
                    className="shrink-0 px-3 py-2 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 hover:bg-red-500/20 transition-all text-sm font-medium flex items-center gap-1"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </section>
    </div>
  );
}
