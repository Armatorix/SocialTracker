import { useState, useEffect } from 'react';
import { api } from './api';
import type { ContentWithUser } from './types';

export function AdminDashboard() {
  const [content, setContent] = useState<ContentWithUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [platformFilter, setPlatformFilter] = useState('');
  const [usernameFilter, setUsernameFilter] = useState('');

  const loadContent = async () => {
    try {
      setLoading(true);
      setError('');
      const filters: any = {};
      if (platformFilter) filters.platform = platformFilter;
      if (usernameFilter) filters.username = usernameFilter;
      const data = await api.getAllContent(filters);
      setContent(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load content');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadContent();
  }, [platformFilter, usernameFilter]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className="w-full max-w-7xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Admin Dashboard - All Content</h1>
      
      {/* Filters */}
      <div className="mb-6 flex gap-4">
        <div className="flex-1">
          <label className="block text-sm font-medium mb-2">Filter by Platform</label>
          <select
            value={platformFilter}
            onChange={(e) => setPlatformFilter(e.target.value)}
            className="w-full bg-[#1a1a1a] border-2 border-[#fbf0df] rounded-lg p-2 text-white"
          >
            <option value="">All Platforms</option>
            <option value="twitter">Twitter/X</option>
            <option value="facebook">Facebook</option>
            <option value="instagram">Instagram</option>
            <option value="youtube">YouTube</option>
            <option value="tiktok">TikTok</option>
          </select>
        </div>
        
        <div className="flex-1">
          <label className="block text-sm font-medium mb-2">Filter by Username</label>
          <input
            type="text"
            value={usernameFilter}
            onChange={(e) => setUsernameFilter(e.target.value)}
            placeholder="Search username..."
            className="w-full bg-[#1a1a1a] border-2 border-[#fbf0df] rounded-lg p-2 text-white"
          />
        </div>
      </div>

      {error && (
        <div className="bg-red-900/50 border-2 border-red-500 rounded-lg p-4 mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-[#1a1a1a]">
                <th className="border-2 border-[#fbf0df] p-3 text-left">ID</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Username</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Email</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Platform</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Link</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Description</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Tags</th>
                <th className="border-2 border-[#fbf0df] p-3 text-left">Created At</th>
              </tr>
            </thead>
            <tbody>
              {content.length === 0 ? (
                <tr>
                  <td colSpan={8} className="border-2 border-[#fbf0df] p-4 text-center">
                    No content found
                  </td>
                </tr>
              ) : (
                content.map((item) => (
                  <tr key={item.id} className="hover:bg-[#1a1a1a]/50">
                    <td className="border-2 border-[#fbf0df] p-3">{item.id}</td>
                    <td className="border-2 border-[#fbf0df] p-3">{item.username}</td>
                    <td className="border-2 border-[#fbf0df] p-3">{item.email}</td>
                    <td className="border-2 border-[#fbf0df] p-3 capitalize">{item.platform}</td>
                    <td className="border-2 border-[#fbf0df] p-3">
                      <a 
                        href={item.link} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        className="text-blue-400 hover:underline truncate block max-w-xs"
                      >
                        {item.link}
                      </a>
                    </td>
                    <td className="border-2 border-[#fbf0df] p-3">
                      {item.description || item.original_text || '-'}
                    </td>
                    <td className="border-2 border-[#fbf0df] p-3">
                      {item.tags && item.tags.length > 0 ? item.tags.join(', ') : '-'}
                    </td>
                    <td className="border-2 border-[#fbf0df] p-3 text-sm">
                      {formatDate(item.created_at)}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
