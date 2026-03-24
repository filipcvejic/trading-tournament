export type Candle = {
  time: number;
  open: number;
  high: number;
  low: number;
  close: number;
};

const SYMBOL_MAP: Record<"EURUSD" | "GBPUSD" | "XAUUSD", string> = {
  EURUSD: "EUR/USD",
  GBPUSD: "GBP/USD",
  XAUUSD: "XAU/USD",
};

export async function getTwelveCandles(
  rawSymbol: "EURUSD" | "GBPUSD" | "XAUUSD",
): Promise<Candle[]> {
  const apiKey = process.env.TWELVE_DATA_API_KEY;

  console.log("[twelve] symbol:", rawSymbol);
  console.log("[twelve] api key exists:", Boolean(apiKey));

  if (!apiKey) {
    throw new Error("Missing TWELVE_DATA_API_KEY");
  }

  const symbol = SYMBOL_MAP[rawSymbol];

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

  console.log("[twelve] rawSymbol:", rawSymbol);
  console.log("[twelve] status:", response.status);
  console.log("[twelve] data:", JSON.stringify(data));

  if (!response.ok || data.status === "error") {
    throw new Error(`TwelveData ${rawSymbol} failed: ${JSON.stringify(data)}`);
  }

  const values = Array.isArray(data.values) ? data.values : [];

  return values
    .slice()
    .reverse()
    .map((c: any) => ({
      time: Math.floor(new Date(`${c.datetime}Z`).getTime() / 1000),
      open: Number(c.open),
      high: Number(c.high),
      low: Number(c.low),
      close: Number(c.close),
    }));
}
