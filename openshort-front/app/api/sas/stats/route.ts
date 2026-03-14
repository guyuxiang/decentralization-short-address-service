import { NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:80';

export async function GET() {
  try {
    const res = await fetch(`${BACKEND_URL}/sas/stats`);
    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to fetch stats' }, { status: res.status });
    }
    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Stats API error:', error);
    return NextResponse.json({ error: 'Backend service unavailable' }, { status: 500 });
  }
}
