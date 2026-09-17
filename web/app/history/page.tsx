'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import {
  ClipboardDocumentIcon,
  ChartBarIcon,
  SparklesIcon,
  ArrowDownTrayIcon,
  PencilSquareIcon,
  CommandLineIcon,
  EyeIcon,
  PlusIcon,
  PlayCircleIcon,
} from '@heroicons/react/24/outline';

interface SummaryRecord {
  id: string;
  video_id: string;
  summary: string;
  key_points: string;
  timestamps: string;
  tags: string;
  source: string;
  created_at: string;
}

interface Stats {
  total_videos: number;
  total_summaries: number;
  ai_summaries: number;
}

export default function History() {
  const [summaries, setSummaries] = useState<SummaryRecord[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedSummary, setSelectedSummary] = useState<SummaryRecord | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    Promise.all([
      fetch('/api/history').then(res => res.json()),
      fetch('/api/stats').then(res => res.json()),
    ]).then(([data, statsData]) => {
      setSummaries(data);
      setStats(statsData);
    }).finally(() => setLoading(false));
  }, []);

  const handleExport = (id: string, format: string) => {
    window.open(`/api/export/${id}/${format}`, '_blank');
  };

  const filteredSummaries = summaries.filter(s => 
    s.summary.toLowerCase().includes(searchQuery.toLowerCase()) ||
    s.tags?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-6xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <ClipboardDocumentIcon className="w-8 h-8 text-purple-400" />
            Summary History
          </h1>
          <Link href="/" className="px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all flex items-center gap-2">
            <PlusIcon className="w-5 h-5" />
            New Summary
          </Link>
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-purple-300 text-sm">Total Videos</p>
                  <p className="text-3xl font-bold text-white">{stats.total_videos}</p>
                </div>
                <PlayCircleIcon className="w-10 h-10 text-purple-400" />
              </div>
            </div>
            <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-purple-300 text-sm">Total Summaries</p>
                  <p className="text-3xl font-bold text-white">{stats.total_summaries}</p>
                </div>
                <ClipboardDocumentIcon className="w-10 h-10 text-green-400" />
              </div>
            </div>
            <div className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-purple-300 text-sm">AI Generated</p>
                  <p className="text-3xl font-bold text-white">{stats.ai_summaries}</p>
                </div>
                <SparklesIcon className="w-10 h-10 text-yellow-400" />
              </div>
            </div>
          </div>
        )}

        {/* Search */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="Search summaries..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full px-6 py-3 bg-slate-800/50 border border-purple-500/30 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:border-purple-500"
          />
        </div>

        {loading ? (
          <div className="text-center py-20">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-400"></div>
          </div>
        ) : filteredSummaries.length === 0 ? (
          <div className="text-center py-20 text-gray-400">
            <SparklesIcon className="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p>No summaries yet. Start by extracting a video transcript!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredSummaries.map((s) => (
              <div key={s.id} className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20 hover:border-purple-500/40 transition-colors">
                <div className="flex justify-between items-start mb-3">
                  <div className="flex items-center gap-3">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                      s.source === 'deepseek' ? 'bg-purple-500/20 text-purple-300' : 'bg-gray-500/20 text-gray-400'
                    }`}>
                      {s.source}
                    </span>
                    <span className="text-gray-500 text-sm">
                      {new Date(s.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => setSelectedSummary(selectedSummary?.id === s.id ? null : s)}
                      className="px-3 py-1 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm transition-all flex items-center gap-1"
                    >
                      <EyeIcon className="w-4 h-4" />
                      View
                    </button>
                    <button
                      onClick={() => handleExport(s.id, 'markdown')}
                      className="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm transition-all flex items-center gap-1"
                    >
                      <PencilSquareIcon className="w-4 h-4" />
                      MD
                    </button>
                    <button
                      onClick={() => handleExport(s.id, 'json')}
                      className="px-3 py-1 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm transition-all flex items-center gap-1"
                    >
                      <CommandLineIcon className="w-4 h-4" />
                      JSON
                    </button>
                    <button
                      onClick={() => handleExport(s.id, 'text')}
                      className="px-3 py-1 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm transition-all flex items-center gap-1"
                    >
                      <ArrowDownTrayIcon className="w-4 h-4" />
                      TXT
                    </button>
                  </div>
                </div>
                <p className="text-gray-300 mb-3 line-clamp-2">{s.summary}</p>
                <div className="flex gap-2 flex-wrap">
                  {(() => {
                    try {
                      const tags = JSON.parse(s.tags || '[]');
                      return tags.map((tag: string, idx: number) => (
                        <span key={idx} className="px-2 py-1 bg-slate-700 text-gray-300 rounded text-xs">
                          {tag}
                        </span>
                      ));
                    } catch {
                      return null;
                    }
                  })()}
                </div>

                {/* Expanded View */}
                {selectedSummary?.id === s.id && (
                  <div className="mt-4 p-4 bg-slate-900/50 rounded-lg border border-purple-500/10">
                    <h4 className="text-white font-semibold mb-2">Full Summary</h4>
                    <p className="text-gray-300 whitespace-pre-wrap">{s.summary}</p>
                    {s.key_points && (() => {
                      try {
                        const points = JSON.parse(s.key_points);
                        return (
                          <div className="mt-3">
                            <h5 className="text-purple-300 font-medium">Key Points</h5>
                            <ul className="text-gray-400 text-sm space-y-1">
                              {points.map((p: string, i: number) => <li key={i}>• {p}</li>)}
                            </ul>
                          </div>
                        );
                      } catch { return null; }
                    })()}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
