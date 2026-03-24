import { redirect } from "next/navigation";
import TrackedChart, {
  type Candle,
  type Trade,
} from "../../components/trading/TrackedChart";
import { getServerApi } from "@/app/lib/api/server";

async function fetchCandles(
  symbol: "EURUSD" | "GBPUSD" | "XAUUSD",
): Promise<Candle[]> {
  const baseUrl = process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000";

  const res = await fetch(`${baseUrl}/api/twelve-candles?symbol=${symbol}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    const err = await res.json().catch(() => null);
    throw new Error(
      err?.error
        ? `${err.error}: ${JSON.stringify(err.details ?? {})}`
        : `Failed to fetch candles for ${symbol}`,
    );
  }

  const data = await res.json();
  return data.candles;
}

async function fetchTrackedTrades(): Promise<Trade[]> {
  const api = await getServerApi();
  try {
    const { data } = await api.get("/admin/tracked-trades");

    if (Array.isArray(data)) {
      return data;
    }
  } catch (err: any) {
    if (err.response?.status === 403) {
      redirect("/competition");
    }

    throw err;
  }
}

export default async function TrackedTradesPage() {
  const [eurusdCandles, gbpusdCandles, xauusdCandles, trackedTrades] =
    await Promise.all([
      fetchCandles("EURUSD"),
      fetchCandles("GBPUSD"),
      fetchCandles("XAUUSD"),
      fetchTrackedTrades(),
    ]);

  return (
    <main
      className="min-h-screen p-6 md:p-8"
      style={{
        background:
          "radial-gradient(circle at top, rgba(139,92,246,0.18) 0%, rgba(7,8,20,1) 32%, rgba(3,5,14,1) 100%)",
      }}
    >
      <div className="mx-auto max-w-[1600px] space-y-6">
        <header>
          <h1 className="text-3xl font-bold text-white md:text-4xl">
            Tracked Trades
          </h1>
        </header>

        <div className="flex flex-col gap-6">
          <TrackedChart
            title="EURUSD"
            candles={eurusdCandles}
            trades={trackedTrades}
          />
          <TrackedChart
            title="GBPUSD"
            candles={gbpusdCandles}
            trades={trackedTrades}
          />
          <TrackedChart
            title="XAUUSD"
            candles={xauusdCandles}
            trades={trackedTrades}
          />
        </div>
      </div>
    </main>
  );
}
