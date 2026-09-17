'use client';
import { useState } from 'react';
import Link from 'next/link';
import {
  PlayCircleIcon,
  DocumentTextIcon,
  ClipboardDocumentIcon,
  ArrowDownTrayIcon,
  SparklesIcon,
  CommandLineIcon,
  PencilSquareIcon,
  PlusIcon,
  TrashIcon,
  CheckIcon,
  LanguageIcon,
  FilmIcon,
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
  thumbnail?: string;
}

interface BatchResult {
  url: string;
  status: 'success' | 'error' | 'pending';
  data?: Result;
  error?: string;
}

interface Language {
  code: string;
  name: string;
}

export default function Home() {
  const [url, setUrl] = useState('');
  const [batchUrls, setBatchUrls] = useState('');
  const [batchMode, setBatchMode] = useState(false);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<Result | null>(null);
  const [summary, setSummary] = useState<Summary | null>(null);
  const [summarizing, setSummarizing] = useState(false);
  const [error, setError] = useState('');
  const [source, setSource] = useState('');
  const [batchResults, setBatchResults] = useState<BatchResult[]>([]);
  const [copied, setCopied] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState('en');
  const [languages, setLanguages] = useState<Language[]>([]);

  useState(() => {
    fetch('/api/languages')
      .then(res => res.json())
      .then(data => setLanguages(data))
      .catch(() => {});
  });

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

  const handleBatchExtract = async () => {
    const urls = batchUrls.split('\n').filter(u => u.trim());
    if (urls.length === 0) return;

    setLoading(true);
    setError('');
    setBatchResults(urls.map(u => ({ url: u, status: 'pending' as const })));

    try {
      const results: BatchResult[] = [];
      for (const u of urls) {
        try {
          const res = await fetch('/api/transcript', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: u.trim() }),
          });
          
          if (res.ok) {
            const data = await res.json();
            results.push({ url: u, status: 'success', data });
          } else {
            results.push({ url: u, status: 'error', error: 'Failed to extract' });
          }
        } catch {
          results.push({ url: u, status: 'error', error: 'Request failed' });
        }
      }
      setBatchResults(results);
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
          full_text: result.full_text || result.transcript.map(l => l.text).join(' '),
          language: selectedLanguage,
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to generate summary');
      }
      
      const data = await res.json();
      setSummary(data.summary);
      setSource(data.source || data.model || 'AI');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setSummarizing(false);
    }
  };

  const handleExport = async (format: string) => {
    if (!result) return;
    window.open(`/api/export/${result.id}/${format}`, '_blank');
  };

  const copyToClipboard = async (text: string) => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const getYouTubeThumbnail = (url: string) => {
    const match = url.match(/(?:v=|\/)([a-zA-Z0-9_-]{11})/);
    if (match) {
      return `https://img.youtube.com/vi/${match[1]}/maxresdefault.jpg`;
    }
    return null;
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

        {/* Mode Toggle */}
        <div className="flex justify-center gap-4 mb-8">
          <button
            onClick={() => setBatchMode(false)}
            className={`px-6 py-3 rounded-lg font-medium transition-all ${
              !batchMode ? 'bg-purple-600 text-white' : 'bg-slate-800/50 text-gray-400 hover:text-white'
            }`}
          >
            <DocumentTextIcon className="w-5 h-5 inline mr-2" />
            Single Video
          </button>
          <button
            onClick={() => setBatchMode(true)}
            className={`px-6 py-3 rounded-lg font-medium transition-all ${
              batchMode ? 'bg-purple-600 text-white' : 'bg-slate-800/50 text-gray-400 hover:text-white'
            }`}
          >
            <PlusIcon className="w-5 h-5 inline mr-2" />
            Batch Processing
          </button>
        </div>

        {/* Input Section */}
        <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-8 mb-8 border border-purple-500/20">
          {!batchMode ? (
            <div className="space-y-4">
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
              
              {/* Language Selector */}
              <div className="flex items-center gap-4">
                <LanguageIcon className="w-5 h-5 text-purple-400" />
                <select
                  value={selectedLanguage}
                  onChange={(e) => setSelectedLanguage(e.target.value)}
                  className="px-4 py-2 bg-slate-900/50 border border-purple-500/30 rounded-lg text-white focus:outline-none focus:border-purple-500"
                >
                  {languages.map(lang => (
                    <option key={lang.code} value={lang.code}>{lang.name}</option>
                  ))}
                </select>
                <span className="text-gray-400 text-sm">Summary language</span>
              </div>
            </div>
          ) : (
            <div>
              <textarea
                placeholder="Paste YouTube URLs (one per line)..."
                value={batchUrls}
                onChange={(e) => setBatchUrls(e.target.value)}
                rows={6}
                className="w-full px-6 py-4 bg-slate-900/50 border border-purple-500/30 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:border-purple-500 resize-none"
              />
              <button
                onClick={handleBatchExtract}
                disabled={loading || !batchUrls.trim()}
                className="mt-4 w-full px-8 py-4 bg-purple-600 hover:bg-purple-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-xl transition-all font-semibold text-lg flex items-center justify-center gap-2"
              >
                <PlusIcon className="w-6 h-6" />
                {loading ? 'Processing...' : `Process ${batchUrls.split('\n').filter(u => u.trim()).length} Videos`}
              </button>
            </div>
          )}
          {error && <p className="mt-4 text-red-400">{error}</p>}
        </div>

        {/* Batch Results */}
        {batchMode && batchResults.length > 0 && (
          <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 mb-8 border border-purple-500/20">
            <h3 className="text-xl font-semibold text-white mb-4">Batch Results</h3>
            <div className="space-y-3">
              {batchResults.map((r, idx) => (
                <div key={idx} className={`p-4 rounded-lg flex items-center gap-4 ${
                  r.status === 'success' ? 'bg-green-900/20 border border-green-500/30' :
                  r.status === 'pending' ? 'bg-yellow-900/20 border border-yellow-500/30' :
                  'bg-red-900/20 border border-red-500/30'
                }`}>
                  {r.status === 'success' ? (
                    <CheckIcon className="w-6 h-6 text-green-400 flex-shrink-0" />
                  ) : r.status === 'pending' ? (
                    <div className="w-6 h-6 border-2 border-yellow-400 border-t-transparent rounded-full animate-spin flex-shrink-0"></div>
                  ) : (
                    <TrashIcon className="w-6 h-6 text-red-400 flex-shrink-0" />
                  )}
                  <div className="flex-1 min-w-0">
                    <p className="text-white font-medium truncate">{r.url}</p>
                    {r.data && <p className="text-gray-400 text-sm">{r.data.title}</p>}
                    {r.error && <p className="text-red-400 text-sm">{r.error}</p>}
                  </div>
                  {r.data && (
                    <button
                      onClick={() => {
                        if (r.data) {
                          setResult(r.data);
                          setBatchMode(false);
                        }
                      }}
                      className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm flex-shrink-0"
                    >
                      View
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Single Video Results */}
        {result && !batchMode && (
          <div className="space-y-6">
            {/* Video Info with Thumbnail */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <div className="flex gap-6">
                {/* Thumbnail */}
                <div className="flex-shrink-0">
                  {getYouTubeThumbnail(result.transcript?.[0]?.text || '') ? (
                    <img 
                      src={`https://img.youtube.com/vi/${result.id.slice(0,11)}/maxresdefault.jpg`}
                      alt={result.title}
                      className="w-48 h-27 rounded-lg object-cover"
                      onError={(e) => {
                        (e.target as HTMLImageElement).src = 'https://via.placeholder.com/480x270/1e293b/64748b?text=Video';
                      }}
                    />
                  ) : (
                    <div className="w-48 h-27 bg-slate-700 rounded-lg flex items-center justify-center">
                      <FilmIcon className="w-12 h-12 text-gray-500" />
                    </div>
                  )}
                </div>
                
                {/* Video Info */}
                <div className="flex-1">
                  <h2 className="text-2xl font-bold text-white mb-2">{result.title}</h2>
                  <div className="flex gap-6 text-gray-400 mb-4">
                    <span>{result.author}</span>
                    <span>•</span>
                    <span>{result.length}</span>
                  </div>
                  <div className="flex gap-2">
                    <Link href="/history" className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm transition-all">
                      History
                    </Link>
                  </div>
                </div>
              </div>
            </div>

            {/* Transcript */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-xl font-semibold text-white flex items-center gap-2">
                  <ClipboardDocumentIcon className="w-6 h-6 text-purple-400" />
                  Transcript
                </h3>
                <div className="flex gap-2">
                  <button
                    onClick={() => handleExport('srt')}
                    className="px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                  >
                    <FilmIcon className="w-4 h-4" />
                    SRT
                  </button>
                  <button
                    onClick={() => handleExport('vtt')}
                    className="px-4 py-2 bg-cyan-600 hover:bg-cyan-700 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                  >
                    <FilmIcon className="w-4 h-4" />
                    VTT
                  </button>
                  <button
                    onClick={() => copyToClipboard(result.transcript.map(l => l.text).join('\n'))}
                    className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                  >
                    {copied ? <CheckIcon className="w-4 h-4" /> : <ClipboardDocumentIcon className="w-4 h-4" />}
                    {copied ? 'Copied!' : 'Copy'}
                  </button>
                  <button
                    onClick={handleSummarize}
                    disabled={summarizing}
                    className="px-6 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-600 text-white rounded-lg transition-all flex items-center gap-2"
                  >
                    <SparklesIcon className="w-5 h-5" />
                    {summarizing ? 'Summarizing...' : 'Generate Summary'}
                  </button>
                </div>
              </div>
              <div className="max-h-96 overflow-y-auto space-y-2">
                {result.transcript.map((line, idx) => (
                  <div key={idx} className="flex gap-4 p-3 bg-slate-900/30 rounded-lg hover:bg-slate-900/50 transition-colors">
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
        {!result && batchResults.length === 0 && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mt-12">
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <DocumentTextIcon className="w-12 h-12 text-purple-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Auto Transcription</h3>
              <p className="text-gray-400">Extract transcripts from any YouTube video automatically</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <SparklesIcon className="w-12 h-12 text-yellow-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">AI Summarization</h3>
              <p className="text-gray-400">Generate summaries with key points and timestamps</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <FilmIcon className="w-12 h-12 text-orange-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Subtitle Export</h3>
              <p className="text-gray-400">Download transcripts as SRT or VTT subtitle files</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <LanguageIcon className="w-12 h-12 text-green-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Multi-Language</h3>
              <p className="text-gray-400">Support for 16+ languages in summaries</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
