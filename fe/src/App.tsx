import { useState, useEffect } from 'react';
import { api } from './api';
import type { User } from './types';
import { AdminDashboard } from './AdminDashboard';
import { CreatorDashboard } from './CreatorDashboard';
import "./index.css";

export function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const loadUser = async () => {
      try {
        const userData = await api.getCurrentUser();
        setUser(userData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load user');
      } finally {
        setLoading(false);
      }
    };
    loadUser();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="bg-red-900/50 border-2 border-red-500 rounded-lg p-6 max-w-md">
          <h2 className="text-xl font-bold mb-2">Error</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">No user found</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <header className="bg-[#1a1a1a] border-b-2 border-[#fbf0df] p-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">Social Tracker</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm">
              {user.username} ({user.role})
            </span>
          </div>
        </div>
      </header>
      
      <main className="py-6">
        {user.role === 'admin' ? <AdminDashboard /> : <CreatorDashboard />}
      </main>
    </div>
  );
}

export default App;
