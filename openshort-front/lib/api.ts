import { REST_ENDPOINT, connectKeplrWallet, buyShortLink, setLongUrl, setSellFlag, setPrice } from '@/lib/cosmos';

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
  },

  async createShortLink(lUrl: string, customSUrl: string, owner: string): Promise<{ success: boolean; txHash?: string; message?: string; sUrl?: string }> {
    try {
      const connection = await connectKeplrWallet();
      const length = customSUrl ? 0 : Math.floor(Math.random() * 3) + 4;
      const result = await buyShortLink(connection.client, owner, '1', {
        sUrl: customSUrl || undefined,
        length: length,
      });
      if (customSUrl) {
        await setLongUrl(connection.client, owner, customSUrl, lUrl);
      }
      return { success: true, txHash: result.txHash, sUrl: customSUrl || `s${length}` };
    } catch (e: any) {
      console.error('Failed to create short link:', e);
      return { success: false, message: e.message || 'Failed to create link on blockchain' };
    }
  },

  async updateLongUrl(sUrl: string, lUrl: string, owner: string): Promise<{ success: boolean; txHash?: string; message?: string }> {
    try {
      const connection = await connectKeplrWallet();
      const result = await setLongUrl(connection.client, owner, sUrl, lUrl);
      return { success: true, txHash: result.txHash };
    } catch (e: any) {
      console.error('Failed to update long URL:', e);
      return { success: false, message: e.message || 'Failed to update URL on blockchain' };
    }
  },

  async listForSale(sUrl: string, price: string, owner: string): Promise<{ success: boolean; txHash?: string; message?: string }> {
    try {
      const connection = await connectKeplrWallet();
      await setSellFlag(connection.client, owner, sUrl, true);
      const result = await setPrice(connection.client, owner, sUrl, price);
      return { success: true, txHash: result.txHash };
    } catch (e: any) {
      console.error('Failed to list for sale:', e);
      return { success: false, message: e.message || 'Failed to list for sale' };
    }
  },

  async buyFromMarketplace(sUrl: string, price: string, buyer: string): Promise<{ success: boolean; txHash?: string; message?: string }> {
    try {
      const connection = await connectKeplrWallet();
      const result = await buyShortLink(connection.client, buyer, price, { sUrl });
      return { success: true, txHash: result.txHash };
    } catch (e: any) {
      console.error('Failed to buy link:', e);
      return { success: false, message: e.message || 'Failed to buy link' };
    }
  }
};
