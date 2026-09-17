import type { Metadata, Viewport } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import SWRegister from '@/components/sw-register';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'AI Video Summarizer - YouTube to Summary',
  description: 'Extract transcripts and generate AI summaries from YouTube videos',
  manifest: '/manifest.json',
};

export const viewport: Viewport = {
  themeColor: '#533afd',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SWRegister />
        {children}
      </body>
    </html>
  );
}
