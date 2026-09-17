'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import {
  ClipboardDocumentIcon,
  ChartBarIcon,
  SparklesIcon,
  PlayCircleIcon,
  UserGroupIcon,
  ClockIcon,
  ArrowTrendingUpIcon,
} from '@heroicons/react/24/outline';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';

interface Stats {
  total_videos: number;
  total_summaries: number;
  ai_summaries: number;
}

interface DailyStat {
  date: string;
  count: number;
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [dailyStats, setDailyStats] = useState<DailyStat[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      fetch('/api/stats').then(res => res.json()),
      fetch('/api/history?limit=30').then(res => res.json()),
    ]).then(([statsData, history]) => {
      setStats(statsData);
      
      // Group by date
      const dateMap = new Map<string, number>();
      history.forEach((s: any) => {
        const date = new Date(s.created_at).toLocaleDateString();
        dateMap.set(date, (dateMap.get(date) || 0) + 1);
      });
      
      setDailyStats(Array.from(dateMap.entries()).map(([date, count]) => ({ date, count })));
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  const pieData = [
    { name: 'AI Generated', value: stats?.ai_summaries || 0 },
    { name: 'Mock Summaries', value: (stats?.total_summaries || 0) - (stats?.ai_summaries || 0) },
  ];

  const COLORS = ['#8b5cf6', '#6b7280'];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-6xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <ChartBarIcon className="w-8 h-8 text-purple-400" />
            Analytics Dashboard
          </h1>
          <Link href="/" className="px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all flex items-center gap-2">
            <PlayCircleIcon className="w-5 h-5" />
            New Summary
          </Link>
        </div>

        {loading ? (
          <div className="text-center py-20">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-400"></div>
          </div>
        ) : (
          <>
            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-purple-300 text-sm">Total Videos</p>
                    <p className="text-3xl font-bold text-white">{stats?.total_videos || 0}</p>
                  </div>
                  <PlayCircleIcon className="w-10 h-10 text-purple-400" />
                </div>
              </div>
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-purple-300 text-sm">Total Summaries</p>
                    <p className="text-3xl font-bold text-white">{stats?.total_summaries || 0}</p>
                  </div>
                  <ClipboardDocumentIcon className="w-10 h-10 text-green-400" />
                </div>
              </div>
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-purple-300 text-sm">AI Generated</p>
                    <p className="text-3xl font-bold text-white">{stats?.ai_summaries || 0}</p>
                  </div>
                  <SparklesIcon className="w-10 h-10 text-yellow-400" />
                </div>
              </div>
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-purple-300 text-sm">AI Rate</p>
                    <p className="text-3xl font-bold text-white">
                      {stats?.total_summaries ? Math.round((stats.ai_summaries / stats.total_summaries) * 100) : 0}%
                    </p>
                  </div>
                  <ArrowTrendingUpIcon className="w-10 h-10 text-blue-400" />
                </div>
              </div>
            </div>

            {/* Charts */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Daily Activity */}
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <ClockIcon className="w-5 h-5 text-purple-400" />
                  Daily Activity
                </h3>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={dailyStats}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                    <XAxis dataKey="date" tick={{ fill: '#9ca3af' }} />
                    <YAxis tick={{ fill: '#9ca3af' }} />
                    <Tooltip 
                      contentStyle={{ backgroundColor: '#1e293b', border: '1px solid #4c1d95' }}
                      labelStyle={{ color: '#e2e8f0' }}
                    />
                    <Bar dataKey="count" fill="#8b5cf6" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </div>

              {/* AI vs Mock */}
              <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                  <SparklesIcon className="w-5 h-5 text-yellow-400" />
                  Summary Sources
                </h3>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {pieData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ backgroundColor: '#1e293b', border: '1px solid #4c1d95' }}
                      labelStyle={{ color: '#e2e8f0' }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
