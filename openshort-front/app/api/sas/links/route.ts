import { NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:80';

export async function GET() {
  try {
    const res = await fetch(`${BACKEND_URL}/sas/adress/sUrls/detail`);
    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to fetch links' }, { status: res.status });
    }
    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Links API error:', error);
    return NextResponse.json({ error: 'Backend service unavailable' }, { status: 500 });
  }
}
