import { SigningAminoClient, StdFee, AminoMsg, OfflineAminoSigner } from '@cosmjs/amino';
import type { ChainInfo, Keplr } from '@keplr-wallet/types';
import type { BroadcastTxResponse } from '@cosmjs/stargate';

export const CHAIN_ID = 'openshort';
export const RPC_ENDPOINT = 'http://43.167.195.109:26657';
export const REST_ENDPOINT = 'http://43.167.195.109:80';

const STAKE_CURRENCY = {
  coinDenom: 'OS',
  coinMinimalDenom: 'os',
  coinDecimals: 6,
};

export const CHAIN_INFO: ChainInfo = {
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
  gasPriceStep: {
    low: 0.001,
    average: 0.002,
    high: 0.003,
  },
  features: ['stargate'],
};

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
  client: SigningAminoClient;
  address: string;
}

let cachedConnection: WalletConnection | null = null;

function normalizeAmount(amount: string) {
  const cleaned = amount.trim();
  if (!cleaned) {
    throw new Error('Amount is required');
  }
  if (!/^[0-9]+(\\.[0-9]+)?$/.test(cleaned)) {
    throw new Error('Amount must be a number');
  }
  return cleaned;
}

async function sendSasTx(
  client: SigningAminoClient,
  signer: string,
  msg: AminoMsg,
  fee: StdFee,
  memo = ''
): Promise<{ txHash: string }> {
  const response: BroadcastTxResponse = await client.signAndBroadcast(signer, [msg], fee, memo);
  if (response.code && response.code !== 0) {
    throw new Error(response.rawLog ?? `Tx failed (code ${response.code})`);
  }
  return {
    txHash: response.transactionHash,
  };
}

export async function connectKeplrWallet(): Promise<WalletConnection> {
  if (typeof window === 'undefined') {
    throw new Error('Keplr wallet must be used in a browser environment');
  }
  if (cachedConnection) {
    return cachedConnection;
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
  const client = await SigningAminoClient.connectWithSigner(RPC_ENDPOINT, signer);
  cachedConnection = {
    client,
    address: accounts[0].address,
  };
  return cachedConnection;
}

export function disconnectWallet() {
  if (cachedConnection) {
    cachedConnection.client.disconnect();
    cachedConnection = null;
  }
}

export async function buyShortLink(
  client: SigningAminoClient,
  buyer: string,
  bidAmount: string,
  options: { sUrl?: string; length?: number; memo?: string } = {}
) {
  const msg: AminoMsg = {
    type: 'sas/BuySUrl',
    value: {
      SUrl: options.sUrl ?? '',
      Bid: [
        {
          denom: STAKE_CURRENCY.coinMinimalDenom,
          amount: normalizeAmount(bidAmount),
        },
      ],
      Buyer: buyer,
      Length: options.length ?? 0,
    },
  };
  return sendSasTx(client, buyer, msg, BUY_FEE, options.memo);
}

export async function setLongUrl(
  client: SigningAminoClient,
  owner: string,
  sUrl: string,
  lUrl: string,
  memo?: string
) {
  const msg: AminoMsg = {
    type: 'sas/SetLUrl',
    value: {
      SUrl: sUrl,
      LUrl: lUrl,
      Owner: owner,
    },
  };
  return sendSasTx(client, owner, msg, SET_FEE, memo);
}

export async function setSellFlag(
  client: SigningAminoClient,
  owner: string,
  sUrl: string,
  isSell: boolean,
  memo?: string
) {
  const msg: AminoMsg = {
    type: 'sas/SetSell',
    value: {
      SUrl: sUrl,
      IsSell: isSell,
      Owner: owner,
    },
  };
  return sendSasTx(client, owner, msg, SELL_FEE, memo);
}

export async function setPrice(
  client: SigningAminoClient,
  owner: string,
  sUrl: string,
  price: string,
  memo?: string
) {
  const msg: AminoMsg = {
    type: 'sas/SetPrice',
    value: {
      SUrl: sUrl,
      Bid: [
        {
          denom: STAKE_CURRENCY.coinMinimalDenom,
          amount: normalizeAmount(price),
        },
      ],
      Owner: owner,
    },
  };
  return sendSasTx(client, owner, msg, SET_FEE, memo);
}

declare global {
  interface Window {
    keplr?: Keplr;
  }
}
