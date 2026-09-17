'use client';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import {
  ClipboardDocumentIcon,
  ChartBarIcon,
  SparklesIcon,
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

export default function History() {
  const [summaries, setSummaries] = useState<SummaryRecord[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/api/history')
      .then(res => res.json())
      .then(data => setSummaries(data))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-4xl mx-auto px-4 py-12">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <ClipboardDocumentIcon className="w-8 h-8 text-purple-400" />
            Summary History
          </h1>
          <Link href="/" className="px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all">
            New Summary
          </Link>
        </div>

        {loading ? (
          <div className="text-center py-20">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-400"></div>
          </div>
        ) : summaries.length === 0 ? (
          <div className="text-center py-20 text-gray-400">
            <SparklesIcon className="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p>No summaries yet. Start by extracting a video transcript!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {summaries.map((s) => (
              <div key={s.id} className="bg-slate-800/50 backdrop-blur rounded-xl p-6 border border-purple-500/20">
                <div className="flex justify-between items-start mb-3">
                  <div>
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                      s.source === 'deepseek' ? 'bg-purple-500/20 text-purple-300' : 'bg-gray-500/20 text-gray-400'
                    }`}>
                      {s.source}
                    </span>
                    <span className="ml-3 text-gray-500 text-sm">
                      {new Date(s.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>
                <p className="text-gray-300 mb-3">{s.summary}</p>
                <div className="flex gap-2">
                  {s.tags && JSON.parse(s.tags || '[]').map((tag: string, idx: number) => (
                    <span key={idx} className="px-2 py-1 bg-slate-700 text-gray-300 rounded text-xs">
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
