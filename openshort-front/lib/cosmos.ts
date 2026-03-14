import { AminoMsg, OfflineAminoSigner } from '@cosmjs/amino';
import { SigningStargateClient, StdFee, SigningStargateClientOptions } from '@cosmjs/stargate';
import { ChainInfo, Keplr } from '@keplr-wallet/types';

interface BroadcastTxResponse {
  readonly code: number;
  readonly hash: string;
  readonly height: number;
  readonly rawLog: string;
}

export const CHAIN_ID = process.env.NEXT_PUBLIC_CHAIN_ID || 'openshort';
export const RPC_ENDPOINT = process.env.NEXT_PUBLIC_RPC_ENDPOINT || '/api/cosmos';
export const RPC_WS_ENDPOINT = process.env.NEXT_PUBLIC_RPC_WS_ENDPOINT || 'ws://43.167.195.109:26657';
export const REST_ENDPOINT = process.env.NEXT_PUBLIC_REST_ENDPOINT || 'http://localhost:80';

const STAKE_CURRENCY = {
  coinDenom: 'OS',
  coinMinimalDenom: 'uos',
  coinDecimals: 6,
};

export const CHAIN_INFO = {
  chainId: CHAIN_ID,
  chainName: 'OpenShort App Chain',
  rpc: RPC_ENDPOINT,
  rest: REST_ENDPOINT,
  bip44: {
    coinType: 118,
  },
  bech32Config: {
    bech32PrefixAccAddr: 'cosmos',
    bech32PrefixAccPub: 'cosmospub',
    bech32PrefixValAddr: 'cosmosvaloper',
    bech32PrefixValPub: 'cosmosvaloperpub',
    bech32PrefixConsAddr: 'cosmosvalcons',
    bech32PrefixConsPub: 'cosmosvalconspub',
  },
  currencies: [STAKE_CURRENCY],
  feeCurrencies: [STAKE_CURRENCY],
  stakeCurrency: STAKE_CURRENCY,
} as unknown as ChainInfo;

const BUY_FEE: StdFee = {
  amount: [{ denom: STAKE_CURRENCY.coinMinimalDenom, amount: '200' }],
  gas: '200000',
};

const SET_FEE: StdFee = {
  amount: [{ denom: STAKE_CURRENCY.coinMinimalDenom, amount: '80' }],
  gas: '120000',
};

const SELL_FEE: StdFee = {
  amount: [{ denom: STAKE_CURRENCY.coinMinimalDenom, amount: '60' }],
  gas: '100000',
};

export interface WalletConnection {
  signer: OfflineAminoSigner;
  address: string;
}

let cachedConnection: WalletConnection | null = null;

function normalizeAmount(amount: string) {
  const cleaned = amount.trim();
  if (!cleaned) {
    throw new Error('Amount is required');
  }
  if (!/^[0-9]+(\.[0-9]+)?$/.test(cleaned)) {
    throw new Error('Amount must be a number');
  }
  return cleaned;
}

async function sendSasTx(
  signer: OfflineAminoSigner,
  signerAddress: string,
  rpcEndpoint: string,
  msg: AminoMsg,
  fee: StdFee,
  memo = ''
): Promise<{ txHash: string }> {
  const client = await SigningStargateClient.offline(signer) as any;
  
  const signedTx = await client.sign(signerAddress, [msg as any], fee, memo);
  
  const txBytes = client.encode(signedTx);
  
  const broadcastRes = await fetch(rpcEndpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: 1,
      method: 'broadcast_tx_commit',
      params: { tx: Buffer.from(txBytes).toString('base64') },
    }),
  });
  
  const broadcastData = await broadcastRes.json();
  
  if (broadcastData.result && broadcastData.result.check_tx && broadcastData.result.check_tx.code !== 0) {
    throw new Error(broadcastData.result.check_tx.log || `Tx failed (code ${broadcastData.result.check_tx.code})`);
  }
  
  return {
    txHash: broadcastData.result.hash,
  };
}

export async function connectKeplrWallet(): Promise<WalletConnection> {
  if (typeof window === 'undefined') {
    throw new Error('Keplr wallet must be used in a browser environment');
  }
  const keplr = window.keplr;
  if (!keplr) {
    throw new Error('Keplr extension is not installed');
  }
  await keplr.experimentalSuggestChain(CHAIN_INFO);
  const signer: OfflineAminoSigner = await keplr.getOfflineSignerOnlyAmino(CHAIN_ID);
  const accounts = await signer.getAccounts();
  if (!accounts.length) {
    throw new Error('No accounts available from Keplr');
  }
  
  if (cachedConnection && cachedConnection.address === accounts[0].address) {
    return cachedConnection;
  }
  
  cachedConnection = {
    signer,
    address: accounts[0].address,
  };
  return cachedConnection;
}

export function disconnectWallet() {
  cachedConnection = null;
}

const MSG_TYPE_PREFIX = 'openshort/sas';

export async function buyShortLink(
  signer: OfflineAminoSigner,
  signerAddress: string,
  rpcEndpoint: string,
  bidAmount: string,
  options: { sUrl?: string; length?: number; memo?: string } = {}
) {
  const msg: AminoMsg = {
    type: `${MSG_TYPE_PREFIX}/BuySUrl`,
    value: {
      SUrl: options.sUrl ?? '',
      Bid: [
        {
          denom: STAKE_CURRENCY.coinMinimalDenom,
          amount: normalizeAmount(bidAmount),
        },
      ],
      Buyer: signerAddress,
      Length: options.length ?? 0,
    },
  };
  return sendSasTx(signer, signerAddress, rpcEndpoint, msg, BUY_FEE, options.memo);
}

export async function setLongUrl(
  signer: OfflineAminoSigner,
  signerAddress: string,
  rpcEndpoint: string,
  sUrl: string,
  lUrl: string,
  memo?: string
) {
  const msg: AminoMsg = {
    type: `${MSG_TYPE_PREFIX}/SetLUrl`,
    value: {
      SUrl: sUrl,
      LUrl: lUrl,
      Owner: signerAddress,
    },
  };
  return sendSasTx(signer, signerAddress, rpcEndpoint, msg, SET_FEE, memo);
}

export async function setSellFlag(
  signer: OfflineAminoSigner,
  signerAddress: string,
  rpcEndpoint: string,
  sUrl: string,
  isSell: boolean,
  memo?: string
) {
  const msg: AminoMsg = {
    type: `${MSG_TYPE_PREFIX}/SetSell`,
    value: {
      SUrl: sUrl,
      IsSell: isSell,
      Owner: signerAddress,
    },
  };
  return sendSasTx(signer, signerAddress, rpcEndpoint, msg, SELL_FEE, memo);
}

export async function setPrice(
  signer: OfflineAminoSigner,
  signerAddress: string,
  rpcEndpoint: string,
  sUrl: string,
  price: string,
  memo?: string
) {
  const msg: AminoMsg = {
    type: `${MSG_TYPE_PREFIX}/SetPrice`,
    value: {
      SUrl: sUrl,
      Bid: [
        {
          denom: STAKE_CURRENCY.coinMinimalDenom,
          amount: normalizeAmount(price),
        },
      ],
      Owner: signerAddress,
    },
  };
  return sendSasTx(signer, signerAddress, rpcEndpoint, msg, SET_FEE, memo);
}

declare global {
  interface Window {
    keplr?: Keplr;
  }
}
