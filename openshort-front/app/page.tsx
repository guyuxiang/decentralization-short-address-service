'use client';

import React, { useState, useEffect } from 'react';
import { 
  Home, 
  Link as LinkIcon, 
  ShoppingCart, 
  LayoutDashboard, 
  Plus, 
  Search, 
  Bell, 
  User,
  ChevronLeft,
  ExternalLink,
  BarChart3,
  Clock,
  Tag,
  ArrowRight,
  Copy,
  Check,
  Wallet,
  X
} from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@/lib/utils';
import { api } from '@/lib/api';
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  BarChart,
  Bar
} from 'recharts';

type View = 'home' | 'manage' | 'marketplace' | 'dashboard';

const SidebarItem = ({ 
  icon: Icon, 
  label, 
  active, 
  onClick 
}: { 
  icon: any, 
  label: string, 
  active: boolean, 
  onClick: () => void 
}) => (
  <button
    onClick={onClick}
    className={cn(
      "flex items-center w-full gap-3 px-4 py-3 text-sm font-medium transition-colors rounded-lg group",
      active 
        ? "bg-[#e8f0fe] text-[#0055d4]" 
        : "text-[#64748b] hover:bg-gray-100 hover:text-[#1e293b]"
    )}
  >
    <Icon size={20} className={cn(active ? "text-[#0055d4]" : "text-[#94a3b8] group-hover:text-[#64748b]")} />
    {label}
  </button>
);

const Card = ({ children, className }: { children: React.ReactNode, className?: string }) => (
  <div className={cn("bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden", className)}>
    {children}
  </div>
);

const ChainInfoPanel = () => {
  const [chainStats, setChainStats] = useState({
    latestBlock: 0,
    validators: 4,
    avgBlockTime: '5.2s',
    totalAddresses: 0,
    totalTxns: 0,
    gasPrice: '0.002 os',
  });

  useEffect(() => {
    const fetchChainData = async () => {
      try {
        const [blockRes, statsRes] = await Promise.all([
          fetch('http://localhost:26657/status'),
          fetch('/api/sas/stats')
        ]);

        const blockData = await blockRes.json();
        const statsData = await statsRes.json();

        setChainStats(prev => ({
          ...prev,
          latestBlock: parseInt(blockData.result?.sync_info?.latest_block_height) || 0,
          totalAddresses: (statsData.totalLinks || 0) + 100,
          totalTxns: parseInt(statsData.totalClicks) || 0,
        }));
      } catch (e) {
        console.error('Failed to fetch chain data:', e);
      }
    };

    fetchChainData();
    const interval = setInterval(fetchChainData, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Card className="p-6 h-full">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-10 h-10 bg-[#0055d4] rounded-lg flex items-center justify-center">
          <svg viewBox="0 0 24 24" fill="none" className="w-6 h-6 text-white">
            <path d="M12 2L2 7l10 5 10-5-10-5z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M2 17l10 5 10-5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M2 12l10 5 10-5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
        </div>
        <div>
          <h3 className="font-bold text-[#1e293b]">OpenShort Chain</h3>
          <p className="text-xs text-[#64748b]">Cosmos SDK Blockchain</p>
        </div>
      </div>

      <div className="space-y-4">
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Latest Block</span>
          <span className="font-bold text-[#1e293b]">#{chainStats.latestBlock.toLocaleString()}</span>
        </div>
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Avg Block Time</span>
          <span className="font-bold text-[#1e293b]">{chainStats.avgBlockTime}</span>
        </div>
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Total Addresses</span>
          <span className="font-bold text-[#1e293b]">{chainStats.totalAddresses.toLocaleString()}</span>
        </div>
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Total Txns</span>
          <span className="font-bold text-[#1e293b]">{chainStats.totalTxns.toLocaleString()}</span>
        </div>
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Gas Tracker</span>
          <span className="font-bold text-[#1e293b]">{chainStats.gasPrice}</span>
        </div>
        <div className="flex justify-between items-center py-2 border-b border-gray-100">
          <span className="text-sm text-[#64748b]">Validators</span>
          <span className="font-bold text-[#1e293b]">{chainStats.validators}</span>
        </div>
        <div className="flex justify-between items-center py-2">
          <span className="text-sm text-[#64748b]">Chain ID</span>
          <span className="font-mono text-xs text-[#0055d4]">openshort</span>
        </div>
      </div>

      <div className="mt-6 pt-4 border-t border-gray-100">
        <div className="flex items-center gap-2 text-xs text-[#64748b]">
          <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
          <span>Connected to OpenShort Network</span>
        </div>
      </div>
    </Card>
  );
};

const HomeView = () => {
  const [url, setUrl] = useState('');
  const [customSUrl, setCustomSUrl] = useState('');
  const [shortLinks, setShortLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const links = await api.getAllShortLinks();
        setShortLinks(links);
      } catch (e) {
        console.error('Failed to fetch links:', e);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const forSaleLinks = shortLinks.filter(link => link.sell);
  const totalClicks = shortLinks.reduce((sum, link) => sum + (link.clicks || 0), 0);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div className="lg:col-span-2 flex flex-col gap-8">
        <div className="grid grid-cols-3 gap-4">
          <Card className="p-4 text-center">
            <p className="text-xs text-[#64748b] uppercase tracking-wider">Total Links</p>
            <p className="text-2xl font-bold text-[#1e293b]">{shortLinks.length}</p>
          </Card>
          <Card className="p-4 text-center">
            <p className="text-xs text-[#64748b] uppercase tracking-wider">On Sale</p>
            <p className="text-2xl font-bold text-[#1e293b]">{forSaleLinks.length}</p>
          </Card>
          <Card className="p-4 text-center">
            <p className="text-xs text-[#64748b] uppercase tracking-wider">Total Clicks</p>
            <p className="text-2xl font-bold text-[#1e293b]">{totalClicks.toLocaleString()}</p>
          </Card>
        </div>
        
        <h1 className="text-2xl font-bold text-[#1e293b]">Quick create: Short link</h1>
        
        <Card className="p-6 flex-1">
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-semibold text-[#1e293b] mb-2">Enter your destination URL</label>
              <div className="flex gap-2">
                <input 
                  type="text" 
                  placeholder="https://example.com/my-long-url"
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#0055d4]/20 focus:border-[#0055d4]"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                />
                <button className="px-6 py-2 bg-[#0055d4] text-white font-semibold rounded-lg hover:bg-[#0044aa] transition-colors">
                  Create link
                </button>
              </div>
            </div>

            <div className="pt-4 border-t border-gray-100">
              <label className="block text-sm font-semibold text-[#1e293b] mb-2">Custom back-half (optional)</label>
              <div className="flex items-center gap-2">
                <span className="text-gray-400 font-medium">openshort.io /</span>
                <input 
                  type="text" 
                  placeholder="my-custom-link"
                  className="w-48 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#0055d4]/20 focus:border-[#0055d4]"
                  value={customSUrl}
                  onChange={(e) => setCustomSUrl(e.target.value)}
                />
              </div>
            </div>

            <div className="flex items-center justify-between pt-2">
              <div className="flex items-center gap-2 text-sm text-[#64748b]">
                <input type="checkbox" className="rounded border-gray-300 text-[#0055d4] focus:ring-[#0055d4]" />
                <span>Also create a QR Code for this link</span>
              </div>
              <div className="text-right">
                <span className="text-[10px] font-bold text-[#94a3b8] uppercase tracking-wider block">Price</span>
                <span className="text-xl font-bold text-[#1e293b]">1 os</span>
              </div>
            </div>
          </div>
        </Card>

        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="p-6 bg-gradient-to-br from-[#0055d4] to-[#003388] text-white">
            <div className="flex justify-between items-start mb-4">
              <div className="p-2 bg-white/10 rounded-lg">
                <Tag size={24} />
              </div>
              <span className="text-xs font-bold uppercase tracking-wider opacity-70">Premium</span>
            </div>
            <h3 className="text-xl font-bold mb-2">Get custom links</h3>
            <p className="text-sm opacity-80 mb-4">Stand out with branded links and a complimentary domain.</p>
            <button className="text-sm font-bold flex items-center gap-1 hover:gap-2 transition-all">
              Upgrade now <ArrowRight size={16} />
            </button>
          </Card>

          <Card className="p-6 bg-[#fff8f0] border-[#ffe4c4]">
            <div className="flex justify-between items-start mb-4">
              <div className="p-2 bg-[#f27d26]/10 rounded-lg text-[#f27d26]">
                <ShoppingCart size={24} />
              </div>
            </div>
            <h3 className="text-xl font-bold text-[#1e293b] mb-2">Marketplace</h3>
            <p className="text-sm text-[#64748b] mb-4">Buy and sell unique short links in our decentralized market.</p>
            <button className="text-sm font-bold text-[#f27d26] flex items-center gap-1 hover:gap-2 transition-all">
              Explore market <ArrowRight size={16} />
            </button>
          </Card>
        </div>
      </div>

      <div>
        <ChainInfoPanel />
      </div>
    </div>
  );
};

const ManageView = ({ walletAddress }: { walletAddress: string | null }) => {
  const [copied, setCopied] = useState<string | null>(null);
  const [links, setLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLinks = async () => {
      try {
        let fetchedLinks;
        if (walletAddress) {
          fetchedLinks = await api.getLinksByOwner(walletAddress);
        } else {
          fetchedLinks = await api.getAllShortLinks();
        }
        setLinks(fetchedLinks);
      } catch (e) {
        console.error('Failed to fetch links:', e);
      } finally {
        setLoading(false);
      }
    };
    fetchLinks();
  }, [walletAddress]);

  const handleCopy = (url: string) => {
    navigator.clipboard.writeText(`openshort.io/${url}`);
    setCopied(url);
    setTimeout(() => setCopied(null), 2000);
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-[#1e293b]">Bitly Links</h1>
        <button className="px-4 py-2 bg-[#0055d4] text-white font-semibold rounded-lg hover:bg-[#0044aa] transition-colors flex items-center gap-2">
          <Plus size={18} /> Create link
        </button>
      </div>

      <div className="flex gap-4 mb-6">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
          <input 
            type="text" 
            placeholder="Search links"
            className="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#0055d4]/20"
          />
        </div>
        <button className="px-4 py-2 bg-white border border-gray-200 rounded-lg text-sm font-medium text-[#64748b] hover:bg-gray-50 flex items-center gap-2">
          <Clock size={16} /> Filter by date
        </button>
      </div>

      <div className="space-y-4">
        {loading && (
          <Card className="p-8 text-center text-[#64748b]">Loading...</Card>
        )}
        {!loading && links.length === 0 && (
          <Card className="p-8 text-center text-[#64748b]">No links found. Create your first short link!</Card>
        )}
        {links.map((link) => (
          <Card key={link.sUrl} className="p-5 hover:border-[#0055d4] transition-colors group">
            <div className="flex gap-4">
              <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center text-[#0055d4]">
                <LinkIcon size={20} />
              </div>
              <div className="flex-1">
                <div className="flex justify-between items-start mb-1">
                  <h3 className="font-bold text-lg text-[#1e293b] group-hover:text-[#0055d4] transition-colors">
                    {link.sUrl}
                  </h3>
                  <div className="flex gap-2">
                    <button 
                      onClick={() => handleCopy(link.sUrl)}
                      className="p-2 text-gray-400 hover:text-[#0055d4] hover:bg-[#e8f0fe] rounded-lg transition-all"
                    >
                      {copied === link.sUrl ? <Check size={18} className="text-green-500" /> : <Copy size={18} />}
                    </button>
                    <button className="p-2 text-gray-400 hover:text-[#0055d4] hover:bg-[#e8f0fe] rounded-lg transition-all">
                      <BarChart3 size={18} />
                    </button>
                  </div>
                </div>
                <p className="text-sm text-[#0055d4] font-medium mb-2 flex items-center gap-1">
                  openshort.io/{link.sUrl} <ExternalLink size={12} />
                </p>
                <p className="text-sm text-[#64748b] truncate mb-4">
                  {link.lUrl}
                </p>
                <div className="flex items-center gap-6 text-xs text-[#94a3b8] font-medium uppercase tracking-wider">
                  <span className="flex items-center gap-1"><BarChart3 size={14} /> {link.clicks} clicks</span>
                  <span className="flex items-center gap-1"><Clock size={14} /> Expires {link.expiresAt}</span>
                  {link.sell && <span className="px-2 py-0.5 bg-green-100 text-green-700 rounded-full">For Sale: {link.price}</span>}
                </div>
              </div>
            </div>
              </Card>
            ))}
          </div>
    </div>
  );
};

const MarketplaceView = ({ walletAddress }: { walletAddress: string | null }) => {
  const [isSellModalOpen, setIsSellModalOpen] = useState(false);
  const [selectedLink, setSelectedLink] = useState<string | null>(null);
  const [price, setPrice] = useState('');
  const [links, setLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLinks = async () => {
      try {
        const allLinks = await api.getAllShortLinks();
        setLinks(allLinks);
      } catch (e) {
        console.error('Failed to fetch links:', e);
      } finally {
        setLoading(false);
      }
    };
    fetchLinks();
  }, []);

  const forSaleLinks = links.filter(l => l.sell);
  const myLinks = walletAddress ? links.filter(l => l.owner === walletAddress && !l.sell) : [];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-[#1e293b]">Marketplace</h1>
          <p className="text-sm text-[#64748b]">Discover and trade unique decentralized short links.</p>
        </div>
        <button 
          onClick={() => setIsSellModalOpen(true)}
          className="px-6 py-2 bg-[#f27d26] text-white font-bold rounded-lg hover:bg-[#e06c15] transition-colors flex items-center gap-2"
        >
          <Tag size={18} />
          Sell My Link
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? (
          <Card className="p-6 col-span-full text-center text-[#64748b]">Loading...</Card>
        ) : forSaleLinks.length === 0 ? (
          <Card className="p-6 col-span-full text-center text-[#64748b]">No links for sale at the moment.</Card>
        ) : (
          forSaleLinks.map((link) => (
          <Card key={link.sUrl} className="p-6 flex flex-col h-full">
            <div className="flex justify-between items-start mb-4">
              <div className="px-3 py-1 bg-[#e8f0fe] text-[#0055d4] rounded-full text-xs font-bold uppercase tracking-wider">
                Available
              </div>
              <div className="text-right">
                <p className="text-xs text-[#94a3b8] font-bold uppercase">Price</p>
                <p className="text-lg font-bold text-[#1e293b]">{link.price}</p>
              </div>
            </div>
            
            <h3 className="text-2xl font-bold text-[#1e293b] mb-1">/{link.sUrl}</h3>
            <p className="text-sm text-[#64748b] mb-6 flex-1">
              A premium short link perfect for your brand. Currently pointing to {new URL(link.lUrl).hostname}.
            </p>

            <div className="space-y-3">
              <div className="flex justify-between text-xs font-medium text-[#94a3b8]">
                <span>Total Clicks</span>
                <span className="text-[#1e293b]">{link.clicks}</span>
              </div>
              <div className="flex justify-between text-xs font-medium text-[#94a3b8]">
                <span>Expires</span>
                <span className="text-[#1e293b]">{link.expiresAt}</span>
              </div>
              <button className="w-full py-3 bg-[#0055d4] text-white font-bold rounded-lg hover:bg-[#0044aa] transition-colors mt-4">
                Buy Now
              </button>
            </div>
          </Card>
        )))}
      </div>

      <AnimatePresence>
        {isSellModalOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsSellModalOpen(false)}
              className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            />
            <motion.div 
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="relative bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden"
            >
              <div className="p-6 border-b border-gray-100 flex justify-between items-center">
                <h2 className="text-xl font-bold text-[#1e293b]">List Link for Sale</h2>
                <button onClick={() => setIsSellModalOpen(false)} className="text-gray-400 hover:text-[#1e293b]">
                  <X size={24} />
                </button>
              </div>
              
              <div className="p-6 space-y-6">
                <div>
                  <label className="block text-sm font-bold text-[#1e293b] mb-3">Select a link to sell</label>
                  <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto pr-2">
                    {myLinks.map(link => (
                      <button
                        key={link.sUrl}
                        onClick={() => setSelectedLink(link.sUrl)}
                        className={cn(
                          "flex items-center justify-between p-3 rounded-xl border transition-all text-left",
                          selectedLink === link.sUrl 
                            ? "border-[#0055d4] bg-[#e8f0fe] ring-1 ring-[#0055d4]" 
                            : "border-gray-200 hover:border-gray-300"
                        )}
                      >
                        <div>
                          <p className="font-bold text-[#1e293b]">/{link.sUrl}</p>
                          <p className="text-xs text-[#64748b] truncate w-48">{link.lUrl}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-xs font-bold text-[#94a3b8] uppercase">{link.clicks} clicks</p>
                        </div>
                      </button>
                    ))}
                  </div>
                </div>

                {selectedLink && (
                  <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
                    <label className="block text-sm font-bold text-[#1e293b] mb-2">Set your price (os)</label>
                    <div className="relative">
                      <input 
                        type="number" 
                        placeholder="e.g. 50"
                        className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#0055d4]/20 focus:border-[#0055d4] font-bold text-lg"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                      />
                      <span className="absolute right-4 top-1/2 -translate-y-1/2 font-bold text-[#64748b]">os</span>
                    </div>
                  </motion.div>
                )}

                <button 
                  disabled={!selectedLink || !price}
                  className="w-full py-4 bg-[#f27d26] text-white font-bold rounded-xl hover:bg-[#e06c15] transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-[#f27d26]/20"
                  onClick={() => {
                    alert(`Listing /${selectedLink} for ${price} os`);
                    setIsSellModalOpen(false);
                  }}
                >
                  Confirm Listing
                </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
};

const DashboardView = ({ walletAddress }: { walletAddress: string | null }) => {
  const [stats, setStats] = useState<{ totalLinks: number; totalClicks: number; topShortUrls: { key: string; value: number }[] }>({ totalLinks: 0, totalClicks: 0, topShortUrls: [] });
  const [links, setLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [statsData, linksData] = await Promise.all([
          api.getStats(),
          walletAddress ? api.getLinksByOwner(walletAddress) : api.getAllShortLinks()
        ]);
        setStats(statsData);
        setLinks(linksData);
      } catch (e) {
        console.error('Failed to fetch dashboard data:', e);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [walletAddress]);

  const forSaleCount = links.filter(l => l.sell).length;
  const totalClicks = links.reduce((sum, l) => sum + (l.clicks || 0), 0);
  const topLinks = (stats.topShortUrls || []).slice(0, 5);

  const chartData = topLinks.length > 0 
    ? topLinks.map((item: any) => ({ name: item.key, clicks: item.value }))
    : links.slice(0, 5).map(l => ({ name: l.sUrl, clicks: l.clicks || 0 }));

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold text-[#1e293b]">Dashboard</h1>

      {loading ? (
        <p className="text-[#64748b]">Loading...</p>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card className="p-6">
              <p className="text-sm font-bold text-[#94a3b8] uppercase tracking-wider mb-1">Total Links</p>
              <p className="text-3xl font-bold text-[#1e293b]">{links.length.toLocaleString()}</p>
            </Card>
            <Card className="p-6">
              <p className="text-sm font-bold text-[#94a3b8] uppercase tracking-wider mb-1">Total Clicks</p>
              <p className="text-3xl font-bold text-[#1e293b]">{totalClicks.toLocaleString()}</p>
            </Card>
            <Card className="p-6">
              <p className="text-sm font-bold text-[#94a3b8] uppercase tracking-wider mb-1">Links for Sale</p>
              <p className="text-3xl font-bold text-[#1e293b]">{forSaleCount}</p>
            </Card>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card className="p-6">
              <h3 className="font-bold text-[#1e293b] mb-6">Click Performance</h3>
          <div className="h-[300px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{fill: '#94a3b8', fontSize: 12}} dy={10} />
                <YAxis axisLine={false} tickLine={false} tick={{fill: '#94a3b8', fontSize: 12}} />
                <Tooltip 
                  contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 10px 15px -3px rgb(0 0 0 / 0.1)' }}
                />
                <Line type="monotone" dataKey="clicks" stroke="#0055d4" strokeWidth={3} dot={{ r: 4, fill: '#0055d4', strokeWidth: 2, stroke: '#fff' }} activeDot={{ r: 6 }} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>

        <Card className="p-6">
          <h3 className="font-bold text-[#1e293b] mb-6">Top Performing Links</h3>
          <div className="space-y-4">
            {chartData.length === 0 ? (
              <p className="text-[#64748b] text-sm">No data available</p>
            ) : (
              chartData.map((item: any, i: number) => (
                <div key={item.name} className="flex items-center gap-4">
                  <div className="w-8 h-8 rounded-lg bg-gray-100 flex items-center justify-center text-xs font-bold text-[#64748b]">
                    {i + 1}
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-bold text-[#1e293b]">/{item.name}</p>
                    <div className="w-full bg-gray-100 h-1.5 rounded-full mt-1">
                      <div 
                        className="bg-[#0055d4] h-full rounded-full" 
                        style={{ width: `${(item.clicks / (chartData[0]?.clicks || 1)) * 100}%` }}
                      />
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-bold text-[#1e293b]">{item.clicks.toLocaleString()}</p>
                    <p className="text-[10px] text-[#94a3b8] font-bold uppercase">Clicks</p>
                  </div>
                </div>
              ))
              )}
          </div>
        </Card>
      </div>
        </>
      )}
    </div>
  );
};

export default function OpenShortApp() {
  const [activeView, setActiveView] = useState<View>('home');
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [walletAddress, setWalletAddress] = useState<string | null>(null);

  const connectWallet = () => {
    setWalletAddress('cosmos1rxh8pl8k3gea67t4uhw2387v9hqpgz3u2awk2g');
  };

  const renderView = () => {
    switch (activeView) {
      case 'home': return <HomeView />;
      case 'manage': return <ManageView walletAddress={walletAddress} />;
      case 'marketplace': return <MarketplaceView walletAddress={walletAddress} />;
      case 'dashboard': return <DashboardView walletAddress={walletAddress} />;
      default: return <HomeView />;
    }
  };

  return (
    <div className="flex h-screen overflow-hidden">
      <motion.aside 
        initial={false}
        animate={{ width: isSidebarOpen ? 260 : 80 }}
        className="bg-white border-r border-gray-200 flex flex-col z-20"
      >
        <div className="p-6 flex items-center gap-3">
          <div className="w-10 h-10 bg-[#0055d4] rounded-xl flex items-center justify-center text-white">
            <LinkIcon size={24} />
          </div>
          {isSidebarOpen && <span className="font-bold text-xl tracking-tight text-[#1e293b]">OpenShort</span>}
        </div>

        <nav className="flex-1 px-4 py-6 space-y-1">
          <SidebarItem 
            icon={Home} 
            label={isSidebarOpen ? "Home" : ""} 
            active={activeView === 'home'} 
            onClick={() => setActiveView('home')} 
          />
          <SidebarItem 
            icon={LinkIcon} 
            label={isSidebarOpen ? "Manage" : ""} 
            active={activeView === 'manage'} 
            onClick={() => setActiveView('manage')} 
          />
          <SidebarItem 
            icon={ShoppingCart} 
            label={isSidebarOpen ? "Marketplace" : ""} 
            active={activeView === 'marketplace'} 
            onClick={() => setActiveView('marketplace')} 
          />
          <SidebarItem 
            icon={LayoutDashboard} 
            label={isSidebarOpen ? "Dashboard" : ""} 
            active={activeView === 'dashboard'} 
            onClick={() => setActiveView('dashboard')} 
          />
        </nav>

        <div className="p-4 border-t border-gray-100">
          <button 
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
            className="w-full p-2 text-gray-400 hover:text-[#1e293b] hover:bg-gray-100 rounded-lg transition-all flex items-center justify-center"
          >
            <ChevronLeft className={cn("transition-transform", !isSidebarOpen && "rotate-180")} size={20} />
          </button>
        </div>
      </motion.aside>

      <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <header className="h-20 bg-white border-bottom border-gray-200 flex items-center justify-between px-8 z-10">
          <div className="flex-1 max-w-2xl relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
            <input 
              type="text" 
              placeholder="Search your links..."
              className="w-full pl-10 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#0055d4]/10 focus:bg-white transition-all"
            />
          </div>

          <div className="flex items-center gap-4">
            <button className="px-4 py-2 bg-[#008080] text-white text-sm font-bold rounded-lg hover:bg-[#006666] transition-colors">
              Upgrade
            </button>
            <button className="p-2 text-gray-400 hover:text-[#1e293b] transition-colors">
              <Bell size={20} />
            </button>
            <div className="h-8 w-px bg-gray-200 mx-2" />
            
            {walletAddress ? (
              <div className="flex items-center gap-3 px-4 py-2 bg-[#e8f0fe] border border-[#0055d4]/20 rounded-xl">
                <div className="w-8 h-8 bg-[#0055d4] rounded-full flex items-center justify-center text-white">
                  <Wallet size={16} />
                </div>
                <div className="hidden md:block">
                  <p className="text-sm font-bold text-[#1e293b]">
                    {walletAddress.slice(0, 6)}...{walletAddress.slice(-5)}
                  </p>
                  <p className="text-[10px] font-bold text-green-600 uppercase tracking-wider">Connected</p>
                </div>
              </div>
            ) : (
              <button 
                onClick={connectWallet}
                className="flex items-center gap-2 px-6 py-2 bg-[#1e293b] text-white font-bold rounded-xl hover:bg-black transition-all shadow-lg shadow-black/10"
              >
                <Wallet size={18} />
                Connect Wallet
              </button>
            )}
          </div>
        </header>

        <div className="flex-1 overflow-y-auto p-8">
          <AnimatePresence mode="wait">
            <motion.div
              key={activeView}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
            >
              {renderView()}
            </motion.div>
          </AnimatePresence>
        </div>
      </main>
    </div>
  );
}
