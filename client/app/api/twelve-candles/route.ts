import { NextRequest, NextResponse } from "next/server";

const SYMBOL_MAP: Record<string, string> = {
  EURUSD: "EUR/USD",
  GBPUSD: "GBP/USD",
  XAUUSD: "XAU/USD",
};

export async function GET(req: NextRequest) {
  const apiKey = process.env.TWELVE_DATA_API_KEY;

  if (!apiKey) {
    return NextResponse.json(
      { error: "Missing TWELVE_DATA_API_KEY" },
      { status: 500 },
    );
  }

  const { searchParams } = new URL(req.url);
  const rawSymbol = (searchParams.get("symbol") || "").toUpperCase();
  const symbol = SYMBOL_MAP[rawSymbol];

  if (!symbol) {
    return NextResponse.json({ error: "Unsupported symbol" }, { status: 400 });
  }

  const url = new URL("https://api.twelvedata.com/time_series");
  url.searchParams.set("symbol", symbol);
  url.searchParams.set("interval", "5min");
  url.searchParams.set("outputsize", "180");
  url.searchParams.set("timezone", "UTC");
  url.searchParams.set("format", "JSON");

  const response = await fetch(url.toString(), {
    headers: {
      Authorization: `apikey ${apiKey}`,
    },
    cache: "no-store",
  });

  const data = await response.json();

  if (!response.ok || data.status === "error") {
    return NextResponse.json(
      {
        error: "Failed to fetch Twelve Data candles",
        details: data,
      },
      { status: 500 },
    );
  }

  const values = Array.isArray(data.values) ? data.values : [];

  const candles = values
    .slice()
    .reverse()
    .map((c: any) => ({
      time: Math.floor(new Date(`${c.datetime}Z`).getTime() / 1000),
      open: Number(c.open),
      high: Number(c.high),
      low: Number(c.low),
      close: Number(c.close),
    }));

  return NextResponse.json({
    symbol: rawSymbol,
    candles,
  });
}
