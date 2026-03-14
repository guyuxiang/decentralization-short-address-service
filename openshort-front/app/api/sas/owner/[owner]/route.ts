import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:80';

export async function GET(request: NextRequest, { params }: { params: Promise<{ owner: string }> }) {
  try {
    const { owner } = await params;
    const res = await fetch(`${BACKEND_URL}/sas/adress/owner/${owner}`);
    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to fetch owner links' }, { status: res.status });
    }
    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Owner links API error:', error);
    return NextResponse.json({ error: 'Backend service unavailable' }, { status: 500 });
  }
}
