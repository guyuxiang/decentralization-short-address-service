import { NextRequest, NextResponse } from 'next/server';

const RPC_ENDPOINT = process.env.NEXT_PUBLIC_RPC_ENDPOINT || 'http://43.167.195.109:26657';

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const params = searchParams.get('params') || '';
    
    const body: any = {
      jsonrpc: '2.0',
      id: 1,
    };
    
    if (params) {
      try {
        body.method = searchParams.get('method') || 'status';
        body.params = JSON.parse(params);
      } catch {
        body.method = searchParams.get('method') || 'status';
        body.params = [];
      }
    } else {
      body.method = 'status';
      body.params = [];
    }

    const res = await fetch(RPC_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return NextResponse.json({ error: 'RPC request failed' }, { status: res.status });
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Cosmos RPC error:', error);
    return NextResponse.json({ error: 'Blockchain service unavailable' }, { status: 500 });
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const res = await fetch(RPC_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        id: 1,
        ...body,
      }),
    });

    if (!res.ok) {
      return NextResponse.json({ error: 'RPC request failed' }, { status: res.status });
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Cosmos RPC error:', error);
    return NextResponse.json({ error: 'Blockchain service unavailable' }, { status: 500 });
  }
}
