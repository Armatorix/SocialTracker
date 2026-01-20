import { useState, useEffect } from 'react';
import { api } from './api';
import type { SocialAccount, Content } from './types';

export function CreatorDashboard() {
  const [socialAccounts, setSocialAccounts] = useState<SocialAccount[]>([]);
  const [content, setContent] = useState<Content[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Form states
  const [showAccountForm, setShowAccountForm] = useState(false);
  const [showContentForm, setShowContentForm] = useState(false);
  const [accountForm, setAccountForm] = useState({
    platform: 'twitter',
    account_name: '',
  });
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

  useEffect(() => {
    loadData();
  }, []);

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

  if (loading) return <div className="text-center py-8">Loading...</div>;

  return (
    <div className="w-full max-w-7xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Creator Dashboard</h1>
      
      {error && (
        <div className="bg-red-900/50 border-2 border-red-500 rounded-lg p-4 mb-4">
          {error}
          <button onClick={() => setError('')} className="ml-4 underline">Dismiss</button>
        </div>
      )}

      {/* Social Accounts Section */}
      <section className="mb-8">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">Social Media Accounts</h2>
          <button
            onClick={() => setShowAccountForm(!showAccountForm)}
            className="bg-[#fbf0df] text-[#1a1a1a] px-4 py-2 rounded-lg font-bold hover:bg-[#f3d5a3]"
          >
            {showAccountForm ? 'Cancel' : 'Add Account'}
          </button>
        </div>

        {showAccountForm && (
          <form onSubmit={handleAddAccount} className="bg-[#1a1a1a] p-4 rounded-lg mb-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-2">Platform</label>
                <select
                  value={accountForm.platform}
                  onChange={(e) => setAccountForm({ ...accountForm, platform: e.target.value })}
                  className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
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
                <label className="block text-sm font-medium mb-2">Account Name</label>
                <input
                  type="text"
                  value={accountForm.account_name}
                  onChange={(e) => setAccountForm({ ...accountForm, account_name: e.target.value })}
                  className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
                  placeholder="@username or account name"
                  required
                />
              </div>
            </div>
            <button
              type="submit"
              className="mt-4 bg-[#fbf0df] text-[#1a1a1a] px-4 py-2 rounded-lg font-bold hover:bg-[#f3d5a3]"
            >
              Add Account
            </button>
          </form>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {socialAccounts.length === 0 ? (
            <p className="text-gray-400 col-span-full">No social accounts connected yet.</p>
          ) : (
            socialAccounts.map((account) => (
              <div key={account.id} className="bg-[#1a1a1a] border-2 border-[#fbf0df] rounded-lg p-4">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h3 className="font-bold capitalize">{account.platform}</h3>
                    <p className="text-sm text-gray-400">{account.account_name}</p>
                  </div>
                  <button
                    onClick={() => handleDeleteAccount(account.id)}
                    className="text-red-400 hover:text-red-300 text-sm"
                  >
                    Delete
                  </button>
                </div>
                <p className="text-xs text-gray-500 mb-3">
                  Last pulled: {formatDate(account.last_pull_at)}
                </p>
                <button
                  onClick={() => handlePullContent(account.id)}
                  className="w-full bg-[#fbf0df] text-[#1a1a1a] px-3 py-2 rounded-lg font-bold hover:bg-[#f3d5a3] text-sm"
                >
                  Pull Content
                </button>
              </div>
            ))
          )}
        </div>
      </section>

      {/* Content Section */}
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">My Content</h2>
          <button
            onClick={() => setShowContentForm(!showContentForm)}
            className="bg-[#fbf0df] text-[#1a1a1a] px-4 py-2 rounded-lg font-bold hover:bg-[#f3d5a3]"
          >
            {showContentForm ? 'Cancel' : 'Add Content Manually'}
          </button>
        </div>

        {showContentForm && (
          <form onSubmit={handleAddContent} className="bg-[#1a1a1a] p-4 rounded-lg mb-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium mb-2">Platform</label>
                <select
                  value={contentForm.platform}
                  onChange={(e) => setContentForm({ ...contentForm, platform: e.target.value })}
                  className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
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
                <label className="block text-sm font-medium mb-2">Link (URL) *</label>
                <input
                  type="url"
                  value={contentForm.link}
                  onChange={(e) => setContentForm({ ...contentForm, link: e.target.value })}
                  className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
                  placeholder="https://..."
                  required
                />
              </div>
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">Original Text/Caption</label>
              <textarea
                value={contentForm.original_text}
                onChange={(e) => setContentForm({ ...contentForm, original_text: e.target.value })}
                className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
                rows={3}
                placeholder="The original post text or caption..."
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">Short Description</label>
              <textarea
                value={contentForm.description}
                onChange={(e) => setContentForm({ ...contentForm, description: e.target.value })}
                className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
                rows={2}
                placeholder="What is this content about..."
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-2">Tags (comma-separated)</label>
              <input
                type="text"
                value={contentForm.tags}
                onChange={(e) => setContentForm({ ...contentForm, tags: e.target.value })}
                className="w-full bg-[#242424] border-2 border-[#fbf0df] rounded-lg p-2"
                placeholder="tag1, tag2, tag3"
              />
            </div>
            <button
              type="submit"
              className="bg-[#fbf0df] text-[#1a1a1a] px-4 py-2 rounded-lg font-bold hover:bg-[#f3d5a3]"
            >
              Add Content
            </button>
          </form>
        )}

        <div className="space-y-4">
          {content.length === 0 ? (
            <p className="text-gray-400">No content yet. Add content manually or pull from connected accounts.</p>
          ) : (
            content.map((item) => (
              <div key={item.id} className="bg-[#1a1a1a] border-2 border-[#fbf0df] rounded-lg p-4">
                <div className="flex justify-between items-start mb-2">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <span className="bg-[#fbf0df] text-[#1a1a1a] px-2 py-1 rounded text-xs font-bold capitalize">
                        {item.platform}
                      </span>
                      <span className="text-xs text-gray-500">
                        {formatDate(item.created_at)}
                      </span>
                    </div>
                    <a
                      href={item.link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-400 hover:underline block mb-2"
                    >
                      {item.link}
                    </a>
                    {item.description && (
                      <p className="text-sm mb-2">{item.description}</p>
                    )}
                    {item.original_text && (
                      <p className="text-sm text-gray-400 mb-2 italic">{item.original_text}</p>
                    )}
                    {item.tags && item.tags.length > 0 && (
                      <div className="flex flex-wrap gap-1">
                        {item.tags.map((tag, i) => (
                          <span key={i} className="bg-[#242424] px-2 py-1 rounded text-xs">
                            #{tag}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                  <button
                    onClick={() => handleDeleteContent(item.id)}
                    className="text-red-400 hover:text-red-300 text-sm ml-4"
                  >
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
