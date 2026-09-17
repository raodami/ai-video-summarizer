'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  PlayCircleIcon,
  DocumentTextIcon,
  ClipboardDocumentIcon,
  ArrowDownTrayIcon,
  SparklesIcon,
  CommandLineIcon,
  PhotoIcon,
  PencilSquareIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';

interface TranscriptLine {
  start: number;
  text: string;
}

interface Summary {
  summary: string;
  key_points: string[];
  timestamps: { time: string; text: string }[];
  tags: string[];
}

interface Result {
  id: string;
  title: string;
  author: string;
  length: string;
  transcript: TranscriptLine[];
  full_text?: string;
}

export default function Home() {
  const router = useRouter();
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<Result | null>(null);
  const [summary, setSummary] = useState<Summary | null>(null);
  const [summarizing, setSummarizing] = useState(false);
  const [error, setError] = useState('');
  const [source, setSource] = useState('');

  const handleExtract = async () => {
    if (!url) return;
    setLoading(true);
    setError('');
    setResult(null);
    setSummary(null);
    
    try {
      const res = await fetch('/api/transcript', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to extract transcript');
      }
      
      const data = await res.json();
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const handleSummarize = async () => {
    if (!result) return;
    setSummarizing(true);
    setError('');
    
    try {
      const res = await fetch('/api/summarize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          transcript_id: result.id,
          full_text: result.full_text || result.transcript.map(l => l.text).join(' ')
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to generate summary');
      }
      
      const data = await res.json();
      setSummary(data.summary);
      setSource(data.source);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setSummarizing(false);
    }
  };

  const handleExport = async (format: string) => {
    if (!result || !summary) return;
    const windowRef = window.open(`/api/export/${result.id}/${format}`);
    if (windowRef) {
      windowRef.focus();
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-4xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-white mb-4 flex items-center justify-center gap-3">
            <PlayCircleIcon className="w-12 h-12 text-purple-400" />
            AI Video Summarizer
          </h1>
          <p className="text-xl text-purple-200">Transform YouTube videos into concise summaries</p>
        </div>

        {/* Input Section */}
        <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-8 mb-8 border border-purple-500/20">
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Paste YouTube URL here..."
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleExtract()}
              className="flex-1 px-6 py-4 bg-slate-900/50 border border-purple-500/30 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:border-purple-500 text-lg"
            />
            <button
              onClick={handleExtract}
              disabled={loading || !url}
              className="px-8 py-4 bg-purple-600 hover:bg-purple-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-xl transition-all font-semibold text-lg flex items-center gap-2"
            >
              <DocumentTextIcon className="w-6 h-6" />
              {loading ? 'Extracting...' : 'Extract'}
            </button>
          </div>
          {error && <p className="mt-4 text-red-400">{error}</p>}
        </div>

        {/* Results Section */}
        {result && (
          <div className="space-y-6">
            {/* Video Info */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <h2 className="text-2xl font-bold text-white mb-2">{result.title}</h2>
              <div className="flex gap-6 text-gray-400">
                <span>{result.author}</span>
                <span>•</span>
                <span>{result.length}</span>
              </div>
            </div>

            {/* Transcript */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-xl font-semibold text-white flex items-center gap-2">
                  <ClipboardDocumentIcon className="w-6 h-6 text-purple-400" />
                  Transcript
                </h3>
                <button
                  onClick={handleSummarize}
                  disabled={summarizing}
                  className="px-6 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-600 text-white rounded-lg transition-all flex items-center gap-2"
                >
                  <SparklesIcon className="w-5 h-5" />
                  {summarizing ? 'Summarizing...' : 'Generate Summary'}
                </button>
              </div>
              <div className="max-h-96 overflow-y-auto space-y-2">
                {result.transcript.map((line, idx) => (
                  <div key={idx} className="flex gap-4 p-3 bg-slate-900/30 rounded-lg">
                    <span className="text-purple-400 font-mono text-sm min-w-[60px]">
                      {formatTime(line.start)}
                    </span>
                    <p className="text-gray-300">{line.text}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Summary */}
            {summary && (
              <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-xl font-semibold text-white flex items-center gap-2">
                    <SparklesIcon className="w-6 h-6 text-yellow-400" />
                    AI Summary ({source})
                  </h3>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleExport('markdown')}
                      className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                    >
                      <PencilSquareIcon className="w-4 h-4" />
                      Markdown
                    </button>
                    <button
                      onClick={() => handleExport('json')}
                      className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                    >
                      <CommandLineIcon className="w-4 h-4" />
                      JSON
                    </button>
                    <button
                      onClick={() => handleExport('text')}
                      className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                    >
                      <DocumentTextIcon className="w-4 h-4" />
                      Text
                    </button>
                  </div>
                </div>
                
                {/* Summary Text */}
                <div className="mb-6 p-4 bg-slate-900/50 rounded-xl">
                  <p className="text-gray-200 leading-relaxed">{summary.summary}</p>
                </div>

                {/* Key Points */}
                <div className="mb-6">
                  <h4 className="text-lg font-semibold text-purple-300 mb-3">Key Points</h4>
                  <ul className="space-y-2">
                    {summary.key_points.map((point, idx) => (
                      <li key={idx} className="flex items-start gap-3 text-gray-300">
                        <span className="w-2 h-2 bg-purple-400 rounded-full mt-2 flex-shrink-0"></span>
                        {point}
                      </li>
                    ))}
                  </ul>
                </div>

                {/* Timestamps */}
                <div className="mb-6">
                  <h4 className="text-lg font-semibold text-purple-300 mb-3">Key Moments</h4>
                  <div className="space-y-2">
                    {summary.timestamps.map((ts, idx) => (
                      <div key={idx} className="flex items-center gap-4 p-3 bg-slate-900/30 rounded-lg">
                        <span className="text-purple-400 font-mono font-bold">{ts.time}</span>
                        <span className="text-gray-300">{ts.text}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Tags */}
                <div>
                  <h4 className="text-lg font-semibold text-purple-300 mb-3">Tags</h4>
                  <div className="flex flex-wrap gap-2">
                    {summary.tags.map((tag, idx) => (
                      <span key={idx} className="px-3 py-1 bg-purple-500/20 text-purple-300 rounded-full text-sm">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Features */}
        {!result && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <DocumentTextIcon className="w-12 h-12 text-purple-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Auto Transcription</h3>
              <p className="text-gray-400">Extract transcripts from any YouTube video automatically</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <SparklesIcon className="w-12 h-12 text-yellow-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">AI Summarization</h3>
              <p className="text-gray-400">Generate concise summaries with key points and timestamps</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <ArrowDownTrayIcon className="w-12 h-12 text-green-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Export & Share</h3>
              <p className="text-gray-400">Download summaries in Markdown, JSON, or plain text</p>
            </div>
          </div>
        )}

        {/* History Link */}
        {result && (
          <div className="mt-8 text-center">
            <Link href="/history" className="text-purple-400 hover:text-purple-300 underline">
              View History
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}
