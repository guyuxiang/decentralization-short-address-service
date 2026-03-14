import { REST_ENDPOINT } from '@/lib/cosmos';

const BASE_URL = '/api/sas';

export interface ShortLink {
  sUrl: string;
  lUrl: string;
  owner: string;
  price: string;
  sell: boolean;
  expiresAt: string;
  clicks?: number;
}

export interface Stats {
  totalLinks: number;
  totalClicks: number;
  topShortUrls: { key: string; value: number }[];
}

export const api = {
  async getAllShortLinks(): Promise<ShortLink[]> {
    const res = await fetch(`${BASE_URL}/links`);
    if (!res.ok) return [];
    const data = await res.json();
    if (data.result && Array.isArray(data.result)) {
      return data.result.map((link: any) => ({
        sUrl: link.sUrl,
        lUrl: link.lUrl || '',
        owner: link.owner,
        price: link.price,
        sell: link.isSell,
        expiresAt: link.expirationTime,
        clicks: link.clicks,
      }));
    }
    return [];
  },

  async getStats(): Promise<Stats> {
    const res = await fetch(`${BASE_URL}/stats`);
    if (!res.ok) return { totalLinks: 0, totalClicks: 0, topShortUrls: [] };
    const data = await res.json();
    return {
      totalLinks: data.totalLinks || 0,
      totalClicks: data.totalClicks || 0,
      topShortUrls: data.topShortUrls || [],
    };
  },

  async getLongUrl(sUrl: string): Promise<string> {
    const res = await fetch(`${BASE_URL}/links/${sUrl}/lUrl`);
    if (!res.ok) return '';
    const data = await res.json();
    return data.lUrl || '';
  },

  async getLinksByOwner(owner: string): Promise<ShortLink[]> {
    if (!owner) return [];
    const res = await fetch(`${BASE_URL}/owner/${owner}`);
    if (!res.ok) return [];
    const data = await res.json();
    if (data.result && Array.isArray(data.result)) {
      return data.result.map((link: any) => ({
        sUrl: link.sUrl,
        lUrl: link.lUrl || '',
        owner: link.owner,
        price: link.price,
        sell: link.isSell,
        expiresAt: link.expirationTime,
        clicks: link.clicks,
      }));
    }
    return [];
  },

  async buyShortLink(sUrl: string, amount: string, owner: string) {
  },

  async setLongUrl(sUrl: string, lUrl: string, owner: string) {
  },

  async claimFaucet(address: string): Promise<{ success: boolean; message: string }> {
    const res = await fetch(`${BASE_URL}/faucet`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ address }),
    });
    if (!res.ok) {
      return { success: false, message: 'Failed to claim tokens' };
    }
    return res.json();
  }
};
