'use client';

import { useEffect, useMemo, useState } from 'react';
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import {
  Check,
  Database,
  ExternalLink,
  Link2,
  RefreshCw,
  ShoppingBag,
  Sparkles,
  Wallet,
  X,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { api, ShortLink, Stats } from '@/lib/api';
import {
  buyShortLink,
  connectKeplrWallet,
  disconnectWallet,
  setLongUrl,
  setPrice,
  setSellFlag,
  WalletConnection,
} from '@/lib/cosmos';

const DEFAULT_STATS: Stats = {
  totalLinks: 0,
  totalClicks: 0,
  topLinks: [],
};

type Notice = {
  type: 'info' | 'error';
  text: string;
};

function StatTile({
  label,
  value,
  helper,
  accent,
}: {
  label: string;
  value: string | number;
  helper?: string;
  accent?: string;
}) {
  return (
    <div
      className={cn(
        'rounded-2xl border p-5 bg-white/10 shadow-sm backdrop-blur transition-all',
        accent || 'border-transparent'
      )}
    >
      <p className="text-xs uppercase tracking-[0.3em] text-white/60">{label}</p>
      <p className="mt-1 text-3xl font-semibold leading-tight text-white">{value}</p>
      {helper && <p className="mt-1 text-sm text-white/70">{helper}</p>}
    </div>
  );
}

function NoticeBanner({ notice }: { notice: Notice | null }) {
  if (!notice) return null;
  return (
    <div
      className={cn(
        'rounded-xl p-4 text-sm font-semibold flex items-center gap-2',
        notice.type === 'error' ? 'bg-red-500/20 text-red-200' : 'bg-sky-500/20 text-sky-100'
      )}
    >
      {notice.type === 'error' ? <X size={16} /> : <Check size={16} />}
      <span>{notice.text}</span>
    </div>
  );
}

function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    return error.message;
  }
  return typeof error === 'string' ? error : 'Unexpected error';
}

export default function OpenShortApp() {
  const [shortLinks, setShortLinks] = useState<ShortLink[]>([]);
  const [stats, setStats] = useState<Stats>(DEFAULT_STATS);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);
  const [walletConnection, setWalletConnection] = useState<WalletConnection | null>(null);
  const [notice, setNotice] = useState<Notice | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const [buyForm, setBuyForm] = useState({
    destination: 'https://',
    custom: '',
    length: 3,
    bid: '1',
  });

  const [selectedLinkId, setSelectedLinkId] = useState('');
  const [linkLongUrl, setLinkLongUrl] = useState('');
  const [salePrice, setSalePrice] = useState('');

  useEffect(() => {
    let active = true;
    setLoading(true);
    (async () => {
      try {
        const [links, statsData] = await Promise.all([api.getAllShortLinks(), api.getStats()]);
        if (!active) return;
        setShortLinks(links);
        setStats(statsData);
        setError(null);
      } catch (err) {
        if (!active) return;
        setError(getErrorMessage(err));
      } finally {
        if (!active) return;
        setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [refreshKey]);

  const marketLinks = useMemo(
    () => shortLinks.filter((link) => link.sell),
    [shortLinks]
  );

  const ownedLinks = useMemo(() => {
    if (!walletConnection) return [];
    return shortLinks.filter((link) => link.owner === walletConnection.address);
  }, [shortLinks, walletConnection]);

  useEffect(() => {
    if (!ownedLinks.length) {
      setSelectedLinkId('');
      return;
    }
    setSelectedLinkId((prev) => {
      if (ownedLinks.find((link) => link.sUrl === prev)) {
        return prev;
      }
      return ownedLinks[0].sUrl;
    });
  }, [ownedLinks]);

  const selectedLink = ownedLinks.find((link) => link.sUrl === selectedLinkId) ?? null;

  useEffect(() => {
    if (selectedLink) {
      setLinkLongUrl(selectedLink.lUrl);
      setSalePrice(selectedLink.price);
    } else {
      setLinkLongUrl('');
      setSalePrice('');
    }
  }, [selectedLink]);

  const chartData = useMemo(() => {
    const normalized =
      stats.topLinks.length > 0
        ? stats.topLinks
        : shortLinks
            .slice()
            .sort((a, b) => (b.clicks ?? 0) - (a.clicks ?? 0))
            .slice(0, 5)
            .map((link) => ({ sUrl: link.sUrl, clicks: link.clicks ?? 0 }));
    return normalized.map((link) => ({
      name: `/${link.sUrl}`,
      clicks: link.clicks ?? 0,
    }));
  }, [shortLinks, stats.topLinks]);

  const handleConnect = async () => {
    setNotice({ type: 'info', text: 'Connecting to Keplr...' });
    try {
      const connection = await connectKeplrWallet();
      setWalletConnection(connection);
      setNotice({ type: 'info', text: 'Wallet connected. You can now sign transactions.' });
    } catch (err) {
      setNotice({ type: 'error', text: getErrorMessage(err) });
    }
  };

  const handleDisconnect = () => {
    disconnectWallet();
    setWalletConnection(null);
    setNotice({ type: 'info', text: 'Wallet disconnected.' });
  };

  const handleBuy = async () => {
    if (!walletConnection) {
      setNotice({ type: 'error', text: 'Connect your wallet before purchasing a link.' });
      return;
    }
    if (!buyForm.bid) {
      setNotice({ type: 'error', text: 'Set a bid amount in OS.' });
      return;
    }
    setIsProcessing(true);
    setNotice({ type: 'info', text: 'Signing purchase... please approve in Keplr.' });
    try {
      await buyShortLink(walletConnection.client, walletConnection.address, buyForm.bid, {
        sUrl: buyForm.custom.trim() || undefined,
        length: buyForm.custom.trim() ? 0 : buyForm.length,
        memo: 'OpenShort quick buy',
      });
      setNotice({ type: 'info', text: 'Purchase broadcast. Awaiting confirmation...' });
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      setNotice({ type: 'error', text: getErrorMessage(err) });
    } finally {
      setIsProcessing(false);
    }
  };

  const handleSetLongUrl = async () => {
    if (!walletConnection || !selectedLink) {
      setNotice({ type: 'error', text: 'Select your link and connect your wallet first.' });
      return;
    }
    if (!linkLongUrl) {
      setNotice({ type: 'error', text: 'Provide a destination URL first.' });
      return;
    }
    setIsProcessing(true);
    setNotice({ type: 'info', text: 'Updating destination...' });
    try {
      await setLongUrl(walletConnection.client, walletConnection.address, selectedLink.sUrl, linkLongUrl);
      setNotice({ type: 'info', text: 'Destination updated. Refreshing data...' });
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      setNotice({ type: 'error', text: getErrorMessage(err) });
    } finally {
      setIsProcessing(false);
    }
  };

  const handleSetPrice = async () => {
    if (!walletConnection || !selectedLink) {
      setNotice({ type: 'error', text: 'Select your link and connect your wallet first.' });
      return;
    }
    if (!salePrice) {
      setNotice({ type: 'error', text: 'Set a price before updating the listing.' });
      return;
    }
    setIsProcessing(true);
    setNotice({ type: 'info', text: 'Publishing sale price...' });
    try {
      await setPrice(walletConnection.client, walletConnection.address, selectedLink.sUrl, salePrice);
      setNotice({ type: 'info', text: 'Price updated. Refreshing data...' });
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      setNotice({ type: 'error', text: getErrorMessage(err) });
    } finally {
      setIsProcessing(false);
    }
  };

  const handleToggleSell = async () => {
    if (!walletConnection || !selectedLink) {
      setNotice({ type: 'error', text: 'Select your link and connect your wallet first.' });
      return;
    }
    setIsProcessing(true);
    setNotice({ type: 'info', text: 'Toggling listing state...' });
    try {
      await setSellFlag(
        walletConnection.client,
        walletConnection.address,
        selectedLink.sUrl,
        !selectedLink.sell
      );
      setNotice({ type: 'info', text: 'Listing updated. Refreshing data...' });
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      setNotice({ type: 'error', text: getErrorMessage(err) });
    } finally {
      setIsProcessing(false);
    }
  };

  const refreshPrompt = () => setRefreshKey((prev) => prev + 1);

  return (
    <div className="min-h-screen bg-[#05040b] text-white">
      <div className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-8">
        <header className="rounded-3xl bg-gradient-to-br from-[#111b36] to-[#0d1a2a] p-8 shadow-[0_20px_60px_rgba(5,4,11,0.5)]">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm uppercase tracking-[0.4em] text-white/60">OpenShort</p>
              <h1 className="text-4xl font-black text-white">Cosmos-native decentralized short links</h1>
              <p className="mt-3 text-base text-white/70">
                Purchase, manage, and trade short addresses that live on the OpenShort app chain.
              </p>
            </div>
            <div className="flex items-center gap-3">
              <button
                onClick={refreshPrompt}
                className="flex items-center gap-2 rounded-full border border-white/10 px-4 py-2 text-xs font-semibold uppercase tracking-[0.4em] text-white/70 transition hover:border-white/40"
              >
                <RefreshCw size={16} /> Refresh
              </button>
              {walletConnection ? (
                <button
                  onClick={handleDisconnect}
                  className="flex items-center gap-2 rounded-full bg-white px-4 py-2 text-xs font-semibold uppercase tracking-[0.4em] text-[#0d1a2a]"
                >
                  <Wallet size={16} /> Disconnect
                </button>
              ) : (
                <button
                  onClick={handleConnect}
                  className="flex items-center gap-2 rounded-full bg-[#1f8cff] px-4 py-2 text-xs font-semibold uppercase tracking-[0.4em] text-white shadow-lg transition hover:bg-[#1a78d3]"
                >
                  <Wallet size={16} /> Connect Keplr
                </button>
              )}
            </div>
          </div>
          <div className="mt-6 grid gap-3 sm:grid-cols-3">
            <StatTile
              label="Total links"
              value={loading ? '—' : stats.totalLinks.toLocaleString()}
              helper="All minted short URLs"
            />
            <StatTile
              label="Clicks"
              value={loading ? '—' : stats.totalClicks.toLocaleString()}
              helper="Recorded visits"
            />
            <StatTile
              label="On sale"
              value={loading ? '—' : marketLinks.length}
              helper="Available for purchase"
            />
          </div>
          {notice && (
            <div className="mt-6">
              <NoticeBanner notice={notice} />
            </div>
          )}
        </header>

        <section className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className="space-y-6 rounded-3xl bg-white/[0.03] p-6 shadow-[0_10px_40px_rgba(5,4,11,0.6)]">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.4em] text-white/50">Quick mint</p>
                <h2 className="text-2xl font-bold">Reserve a short link</h2>
              </div>
              <span className="inline-flex items-center gap-1 rounded-full bg-slate-900/70 px-3 py-1 text-xs uppercase tracking-[0.4em] text-white/80">
                <Link2 size={14} /> Cosmos SDK
              </span>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <label className="space-y-1 text-sm text-white/60">
                Destination URL
                <input
                  className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none transition focus:border-[#1f8cff] focus:bg-white/[0.05]"
                  placeholder="https://example.com"
                  value={buyForm.destination}
                  onChange={(event) =>
                    setBuyForm((prev) => ({ ...prev, destination: event.target.value }))
                  }
                />
              </label>
              <label className="space-y-1 text-sm text-white/60">
                Custom alias
                <input
                  className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none transition focus:border-[#1f8cff] focus:bg-white/[0.05]"
                  placeholder="enter up to 6 characters"
                  value={buyForm.custom}
                  maxLength={6}
                  onChange={(event) =>
                    setBuyForm((prev) => ({ ...prev, custom: event.target.value.trim() }))
                  }
                />
              </label>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
              <label className="space-y-1 text-xs uppercase tracking-[0.4em] text-white/50">
                Length
                <input
                  type="range"
                  min={1}
                  max={6}
                  value={buyForm.length}
                  onChange={(event) =>
                    setBuyForm((prev) => ({ ...prev, length: Number(event.target.value) }))
                  }
                  disabled={!!buyForm.custom}
                  className="h-2 w-full cursor-pointer accent-[#1f8cff]"
                />
                <span className="text-sm text-white">
                  {buyForm.custom ? 'Locked for custom alias' : `${buyForm.length} characters`}
                </span>
              </label>
              <label className="space-y-1 text-sm text-white/60">
                Bid amount (os)
                <input
                  className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none transition focus:border-[#1f8cff] focus:bg-white/[0.05]"
                  placeholder="1"
                  value={buyForm.bid}
                  onChange={(event) => setBuyForm((prev) => ({ ...prev, bid: event.target.value }))}
                />
              </label>
              <div className="rounded-2xl border border-dashed border-white/20 p-3 text-xs uppercase tracking-[0.4em] text-white/60">
                {walletConnection ? (
                  <p className="text-white">Connected as {walletConnection.address.slice(0, 6)}…</p>
                ) : (
                  <p>Wallet needed to mint</p>
                )}
              </div>
            </div>
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <p className="text-sm text-white/60">
                After minting, use the Manage panel to map the short link to your destination URL.
              </p>
              <button
                onClick={handleBuy}
                disabled={isProcessing}
                className="inline-flex items-center justify-center gap-2 rounded-2xl bg-gradient-to-r from-[#ffb347] to-[#ffcc33] px-6 py-3 text-sm font-semibold uppercase tracking-[0.3em] text-[#0d1a2a] shadow-lg shadow-[#ffb347]/40 transition disabled:cursor-not-allowed disabled:opacity-60"
              >
                <Sparkles size={18} /> Purchase short link
              </button>
            </div>
          </div>

          <div className="space-y-6">
            <div className="rounded-3xl border border-white/10 bg-black/40 p-6 shadow-[0_10px_40px_rgba(5,4,11,0.5)]">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs uppercase tracking-[0.4em] text-white/50">Wallet</p>
                  <h3 className="text-xl font-semibold">Keplr Integration</h3>
                </div>
                <Wallet size={24} />
              </div>
              <p className="mt-4 text-sm text-white/70">
                Interact directly with the openshort app chain via Cosmos SDK messages and Amino signatures.
              </p>
              {walletConnection ? (
                <div className="mt-4 rounded-2xl bg-white/5 p-4 text-sm">
                  <p className="text-white/80">Connected address</p>
                  <p className="font-bold text-white">{walletConnection.address}</p>
                </div>
              ) : (
                <p className="mt-4 text-sm text-white/60">Keplr is required to submit transactions.</p>
              )}
            </div>

            <div className="rounded-3xl border border-white/10 bg-white/[0.02] p-6 shadow-[0_10px_40px_rgba(5,4,11,0.5)]">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs uppercase tracking-[0.4em] text-white/50">Marketplace</p>
                  <h3 className="text-xl font-semibold">Links for sale</h3>
                </div>
                <ShoppingBag size={24} />
              </div>
              <p className="mt-3 text-sm text-white/60">
                These listings can be purchased using a Cosmos transaction.
              </p>
              <div className="mt-4 space-y-4">
                {marketLinks.slice(0, 4).map((link) => (
                  <div
                    key={link.sUrl}
                    className="flex items-center justify-between rounded-2xl border border-white/10 bg-white/5 px-4 py-3"
                  >
                    <div>
                      <p className="text-lg font-semibold text-white">/{link.sUrl}</p>
                      <p className="text-sm text-white/60">{link.lUrl}</p>
                      <p className="text-xs uppercase tracking-[0.3em] text-white/40">
                        {link.clicks?.toLocaleString() ?? 0} clicks
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-white/80">Price</p>
                      <p className="text-xl font-semibold text-white">{link.price}</p>
                    </div>
                  </div>
                ))}
                {!marketLinks.length && (
                  <p className="rounded-2xl border border-dashed border-white/20 p-4 text-center text-sm text-white/60">
                    No listings available right now.
                  </p>
                )}
              </div>
            </div>
          </div>
        </section>

        <section className="grid gap-6 lg:grid-cols-[1.5fr_1fr]">
          <div className="rounded-3xl bg-white/[0.04] p-6 shadow-[0_10px_40px_rgba(5,4,11,0.55)]">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.4em] text-white/50">Top demand</p>
                <h3 className="text-2xl font-semibold">Live clicks chart</h3>
              </div>
              <Database size={24} />
            </div>
            <div className="mt-4 h-48">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} margin={{ top: 10, right: 16, left: 0, bottom: 0 }}>
                  <XAxis dataKey="name" tickLine={false} axisLine={false} tick={{ fill: 'rgba(255,255,255,0.7)', fontSize: 12 }} />
                  <YAxis
                    tick={{ fill: 'rgba(255,255,255,0.7)', fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <Tooltip
                    contentStyle={{ backgroundColor: 'rgba(13, 26, 42, 0.9)', border: 'none', borderRadius: 12 }}
                    itemStyle={{ color: '#fff' }}
                  />
                  <Bar dataKey="clicks" fill="#1f8cff" radius={[8, 8, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="space-y-6">
            <div className="rounded-3xl border border-white/10 bg-black/60 p-6 shadow-[0_10px_40px_rgba(5,4,11,0.6)]">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs uppercase tracking-[0.4em] text-white/50">Owned links</p>
                  <h3 className="text-xl font-semibold">Your portfolio</h3>
                </div>
                <Check size={24} />
              </div>
              {ownedLinks.length ? (
                <div className="mt-4 w-full space-y-3">
                  {ownedLinks.map((link) => (
                    <div
                      key={link.sUrl}
                      className={cn(
                        'rounded-2xl border px-4 py-3 transition',
                        selectedLinkId === link.sUrl
                          ? 'border-[#1f8cff] bg-[#1f8cff]/10'
                          : 'border-white/10 bg-white/5 hover:border-white/30'
                      )}
                      onClick={() => setSelectedLinkId(link.sUrl)}
                    >
                      <p className="text-lg font-semibold text-white">/{link.sUrl}</p>
                      <p className="text-xs text-white/60">{link.lUrl || 'Not mapped yet'}</p>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="mt-4 rounded-2xl border border-dashed border-white/20 p-4 text-center text-sm text-white/60">
                  Own links will appear here after you mint them on-chain.
                </p>
              )}
            </div>

            {selectedLink && (
              <div className="rounded-3xl border border-white/10 bg-white/[0.03] p-6 shadow-[0_10px_40px_rgba(5,4,11,0.55)]">
                <p className="text-xs uppercase tracking-[0.4em] text-white/50">Manage link</p>
                <h3 className="text-2xl font-semibold">/{selectedLink.sUrl}</h3>
                <label className="mt-4 block text-sm text-white/60">
                  Destination URL
                  <input
                    className="mt-2 w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none transition focus:border-[#1f8cff]"
                    value={linkLongUrl}
                    onChange={(event) => setLinkLongUrl(event.target.value)}
                    placeholder="https://example.com"
                  />
                </label>
                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                  <button
                    onClick={handleSetLongUrl}
                    disabled={isProcessing}
                    className="rounded-2xl bg-gradient-to-r from-[#34d399] to-[#10b981] px-4 py-3 text-sm font-semibold uppercase tracking-[0.4em] text-[#03140b] transition disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    Update destination
                  </button>
                  <label className="space-y-1 text-sm text-white/60">
                    Sale price
                    <input
                      className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none transition focus:border-[#1f8cff]"
                      value={salePrice}
                      onChange={(event) => setSalePrice(event.target.value)}
                      placeholder="e.g. 50"
                    />
                  </label>
                </div>
                <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                  <button
                    onClick={handleSetPrice}
                    disabled={isProcessing}
                    className="flex-1 rounded-2xl border border-white/20 px-4 py-3 text-sm font-semibold uppercase tracking-[0.4em] text-white transition hover:border-white/40 disabled:opacity-60"
                  >
                    Set listing price
                  </button>
                  <button
                    onClick={handleToggleSell}
                    disabled={isProcessing}
                    className={cn(
                      'flex-1 rounded-2xl px-4 py-3 text-sm font-semibold uppercase tracking-[0.4em] transition disabled:opacity-60',
                      selectedLink.sell
                        ? 'bg-red-500 text-white'
                        : 'bg-[#1f8cff] text-white'
                    )}
                  >
                    {selectedLink.sell ? 'Unlist link' : 'Put on sale'}
                  </button>
                </div>
              </div>
            )}
          </div>
        </section>

        {loading && (
          <p className="text-center text-sm text-white/60">
            Fetching the latest data from the REST API and Cosmos app chain...
          </p>
        )}
        {error && (
          <p className="text-center text-sm text-red-400">
            Unable to load data: {error}
          </p>
        )}
      </div>
    </div>
  );
}
