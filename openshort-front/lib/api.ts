import { REST_ENDPOINT } from '@/lib/cosmos';

const BASE_URL = REST_ENDPOINT;

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
  topLinks: { sUrl: string; clicks: number }[];
}

export const api = {
  async getAllShortLinks(): Promise<ShortLink[]> {
    const res = await fetch(`${BASE_URL}/sas/adress/sUrls`);
    if (!res.ok) return [];
    const data = await res.json();
    if (Array.isArray(data.result)) {
      return data.result;
    }
    return Array.isArray(data.links) ? data.links : [];
  },

  async getStats(): Promise<Stats> {
    const res = await fetch(`${BASE_URL}/sas/stats`);
    if (!res.ok) return { totalLinks: 0, totalClicks: 0, topLinks: [] };
    return await res.json();
  },

  async getLongUrl(sUrl: string): Promise<string> {
    const res = await fetch(`${BASE_URL}/sas/adress/${sUrl}/lUrl`);
    if (!res.ok) return '';
    const data = await res.json();
    return data.result || '';
  },

  async getLinksByOwner(owner: string): Promise<ShortLink[]> {
    if (!owner) return [];
    const links = await this.getAllShortLinks();
    return links.filter((link) => link.owner === owner);
  },

  // Note: These would require wallet integration in a real app
  // For now we define the structure
  async buyShortLink(sUrl: string, amount: string, owner: string) {
    // POST /sas/adress
  },

  async setLongUrl(sUrl: string, lUrl: string, owner: string) {
    // PUT /sas/adress/lUrl
  }
};
