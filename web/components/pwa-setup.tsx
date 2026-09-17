'use client';

import { useEffect } from 'react';
import Script from 'next/script';

export default function PWASetup() {
  useEffect(() => {
    // iOS home screen support
    if (window.navigator.standalone) {
      document.documentElement.style.setProperty('--viewport-height', `${window.innerHeight}px`);
    }
    
    // Add to home screen prompt (simplified)
    let deferredPrompt;
    window.addEventListener('beforeinstallprompt', (e) => {
      e.preventDefault();
      deferredPrompt = e;
    });
  }, []);

  return (
    <>
      <meta name="theme-color" content="#533afd" />
      <meta name="apple-mobile-web-app-capable" content="yes" />
      <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
      <Script src="/manifest.json" type="application/manifest+json" />
    </>
  );
}
