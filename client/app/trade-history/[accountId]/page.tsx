"use client";

import { webApi } from "@/app/lib/api/client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import { useParams } from "next/navigation";
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ReferenceLine,
} from "recharts";
import LogoutButton from "@/app/components/LogoutButton";
import BackButton from "@/app/components/BackButton";
import { formatMoney } from "@/app/lib/format";

type Side = "BUY" | "SELL";

type Trade = {
  id: number | string;
  symbol: string;
  side: Side;
  volume: number;
  openTime: string;
  closeTime: string;
  openPrice: number;
  closePrice: number;
  profit: number;
  commission: number;
  swap: number;
};

type ApiResponse = {
  username: string;
  trades: any[]; // keep as-is (backend shape), we map safely
};

type Row = Trade & { total: number };

function formatDateTime(iso: string) {
  const d = new Date(iso);
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  })}`;
}

function SideBadge({ side }: { side: Side }) {
  const isBuy = side === "BUY";
  return (
    <span
      className={[
        "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-semibold",
        "border border-white/10 bg-white/5",
        isBuy ? "text-[#60A5FA]" : "text-[#A855F7]",
      ].join(" ")}
    >
      <span
        className={[
          "h-2 w-2 rounded-full",
          isBuy ? "bg-[#60A5FA]" : "bg-[#A855F7]",
        ].join(" ")}
      />
      {side}
    </span>
  );
}

function Th({ children }: { children: React.ReactNode }) {
  return (
    <th className="px-3 py-3 text-xs font-semibold uppercase tracking-wide text-[#A1A1AA]">
      {children}
    </th>
  );
}

function Td({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <td
      className={["px-3 py-3 whitespace-nowrap", className]
        .filter(Boolean)
        .join(" ")}
    >
      {children}
    </td>
  );
}

function Stat({
  label,
  value,
  accent,
}: {
  label: string;
  value: string;
  accent: "blue" | "purple";
}) {
  const ring =
    accent === "blue"
      ? "from-[#60A5FA]/35 to-transparent"
      : "from-[#A855F7]/35 to-transparent";

  return (
    <div className="relative rounded-xl border border-white/10 bg-white/5 px-3 py-2 overflow-hidden">
      <div
        className={[
          "pointer-events-none absolute inset-0 bg-gradient-to-r",
          ring,
        ].join(" ")}
      />
      <div className="relative">
        <div className="text-[11px] uppercase tracking-wide text-[#A1A1AA]">
          {label}
        </div>
        <div className="mt-0.5 text-base font-semibold">{value}</div>
      </div>
    </div>
  );
}

function mapTrade(t: any): Trade {
  return {
    id: t.positionId ?? t.id,
    symbol: t.symbol,
    side: t.side,
    volume: Number(t.volume),
    openTime: t.openTime ?? t.open_time,
    closeTime: t.closeTime ?? t.close_time,
    openPrice: Number(t.openPrice ?? t.open_price),
    closePrice: Number(t.closePrice ?? t.close_price),
    profit: Number(t.profit ?? 0),
    commission: Number(t.commission ?? 0),
    swap: Number(t.swap ?? 0),
  };
}

function PnlTooltip({
  active,
  payload,
}: {
  active?: boolean;
  payload?: Array<any>;
}) {
  if (!active || !payload?.length) return null;

  const p = payload[0]?.payload as { time: string; cumPnl: number } | undefined;
  if (!p) return null;

  return (
    <div
      style={{
        background: "#0B0C10",
        border: "1px solid rgba(255,255,255,0.12)",
        borderRadius: 14,
        padding: "10px 12px",
        color: "white",
      }}
    >
      <div style={{ fontSize: 12, opacity: 0.85 }}>
        Close: {formatDateTime(p.time)}
      </div>
      <div style={{ marginTop: 4, fontWeight: 600 }}>
        {formatMoney(Number(p.cumPnl), { sign: true })}
      </div>
    </div>
  );
}

export default function TradingAccountTradeHistoryPage() {
  const params = useParams<{ accountId: string }>();
  const accountId = params.accountId;

  const [username, setUsername] = useState("");
  const [trades, setTrades] = useState<Trade[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // ✅ Makes tooltip easy to hide on mobile when tapping outside
  const chartWrapRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    let alive = true;

    async function load() {
      setLoading(true);
      setError(null);

      try {
        const { data } = await webApi.get<ApiResponse>(
          `/trading-accounts/${accountId}/trade-history`,
        );

        const mapped = (data?.trades ?? []).map(mapTrade);

        if (!alive) return;
        setUsername(data?.username ?? "");
        setTrades(mapped);
      } catch (e: any) {
        if (!alive) return;
        setError(
          e?.response?.data?.error ||
            e?.message ||
            "Failed to load trade history.",
        );
      } finally {
        if (!alive) return;
        setLoading(false);
      }
    }

    load();
    return () => {
      alive = false;
    };
  }, [accountId]);

  const rows: Row[] = useMemo(() => {
    return trades.map((t) => ({
      ...t,
      total: t.profit + t.commission + t.swap,
    }));
  }, [trades]);

  // Table is descending (backend already sends it like that), but chart must be ascending
  const chartRows = useMemo(() => {
    return rows.slice().reverse();
  }, [rows]);

  const chartData = useMemo(() => {
    let cumPnl = 0;
    return chartRows.map((t, i) => {
      cumPnl += t.total;
      return { i, time: t.closeTime, cumPnl };
    });
  }, [chartRows]);

  const stats = useMemo(() => {
    const net = rows.reduce((sum, t) => sum + t.total, 0);
    const wins = rows.filter((t) => t.total > 0).length;
    const winRate = rows.length ? (wins / rows.length) * 100 : 0;

    const best = rows.reduce(
      (m, t) => Math.max(m, t.total),
      Number.NEGATIVE_INFINITY,
    );
    const worst = rows.reduce(
      (m, t) => Math.min(m, t.total),
      Number.POSITIVE_INFINITY,
    );

    return {
      tradesCount: rows.length,
      net,
      winRate,
      best: isFinite(best) ? best : 0,
      worst: isFinite(worst) ? worst : 0,
    };
  }, [rows]);

  const netIsPos = stats.net >= 0;

  return (
    <div className="min-h-screen bg-[#0B0C10] text-white">
      {/* top bar */}
      <div className="flex justify-between pt-4 px-4">
        <BackButton text="Back to competition" />
        <LogoutButton />
      </div>

      {/* background glow */}
      <div className="pointer-events-none fixed inset-0">
        <div className="absolute -top-24 left-1/2 h-[520px] w-[520px] -translate-x-1/2 rounded-full bg-[#A855F7]/20 blur-3xl" />
        <div className="absolute -top-10 left-[55%] h-[520px] w-[520px] -translate-x-1/2 rounded-full bg-[#60A5FA]/15 blur-3xl" />
      </div>

      <div className="relative mx-auto w-full max-w-6xl px-4 py-8 space-y-6">
        {/* Header */}
        <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur px-6 py-5">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div className="space-y-1">
              <div className="text-xs text-[#A1A1AA]">
                Trading account • Trade history
              </div>
              <h1 className="text-2xl sm:text-3xl font-semibold tracking-tight">
                {username ? `${username}'s Trade History` : "Trade History"}
              </h1>
              <p className="text-sm text-[#A1A1AA]">
                Net PnL = profit - commission - swap
              </p>
            </div>

            <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
              <Stat
                label="Net PnL"
                value={formatMoney(stats.net, { sign: true })}
                accent={netIsPos ? "blue" : "purple"}
              />
              <Stat
                label="Win rate"
                value={`${stats.winRate.toFixed(0)}%`}
                accent="blue"
              />
              <Stat
                label="Best trade"
                value={formatMoney(stats.best, { sign: true })}
                accent="blue"
              />
              <Stat
                label="Worst trade"
                value={formatMoney(stats.worst)}
                accent="purple"
              />
            </div>
          </div>
        </div>

        {/* Loading / Error */}
        {loading ? (
          <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6 text-center text-[#A1A1AA]">
            Loading…
          </div>
        ) : error ? (
          <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6 text-center text-red-400">
            {error}
          </div>
        ) : (
          <>
            {/* Chart */}
            <div className="relative rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-4 sm:p-6 overflow-hidden">
              {/* soft neon ring */}
              <div className="pointer-events-none absolute inset-0">
                <div className="absolute inset-[-2px] rounded-2xl bg-gradient-to-r from-[#A855F7] to-[#60A5FA] opacity-25 blur-[10px]" />
                <div className="absolute inset-[1px] rounded-2xl bg-[#151621]/90" />
              </div>

              <div className="relative">
                <div className="space-y-1">
                  <h2 className="text-lg font-semibold">PnL Curve</h2>
                  <p className="text-xs text-[#A1A1AA]">
                    Cumulative net PnL by close time
                  </p>
                </div>

                <div ref={chartWrapRef} className="mt-4 h-[300px] sm:h-[420px]">
                  <ResponsiveContainer
                    width="100%"
                    height="100%"
                    style={{
                      WebkitTapHighlightColor: "transparent",
                      WebkitUserSelect: "none",
                      userSelect: "none",
                      outline: "none",
                    }}
                  >
                    <AreaChart
                      tabIndex={-1}
                      data={chartData}
                      style={{ outline: "none" }}
                      margin={{ top: 10, right: 16, left: 0, bottom: 0 }}
                      accessibilityLayer={false}
                    >
                      <CartesianGrid
                        strokeDasharray="4 4"
                        stroke="rgba(255,255,255,0.08)"
                      />

                      <XAxis
                        dataKey="i"
                        type="number"
                        scale="linear"
                        domain={[0, "dataMax"]}
                        tickFormatter={(idx) => {
                          const d = chartData[idx]?.time;
                          return d ? new Date(d).toLocaleDateString() : "";
                        }}
                        padding={{ left: 10, right: 10 }}
                        stroke="rgba(255,255,255,0.55)"
                        tick={{ fontSize: 12 }}
                        interval="preserveStartEnd"
                        allowDuplicatedCategory={false}
                      />

                      <YAxis
                        tickFormatter={(v) => Number(v).toFixed(0)}
                        stroke="rgba(255,255,255,0.55)"
                        tick={{ fontSize: 12 }}
                        width={44}
                      />
                      <ReferenceLine y={0} stroke="rgba(255,255,255,0.22)" />

                      <Tooltip
                        content={<PnlTooltip />}
                        cursor={{
                          stroke: "rgba(255,255,255,0.12)",
                          strokeWidth: 1,
                        }}
                        isAnimationActive={false}
                      />

                      <defs>
                        <linearGradient
                          id="neonFill"
                          x1="0"
                          y1="0"
                          x2="1"
                          y2="0"
                        >
                          <stop
                            offset="0%"
                            stopColor="#A855F7"
                            stopOpacity={0.28}
                          />
                          <stop
                            offset="100%"
                            stopColor="#60A5FA"
                            stopOpacity={0.28}
                          />
                        </linearGradient>
                      </defs>

                      <Area
                        type="linear"
                        dataKey="cumPnl"
                        stroke="#60A5FA"
                        strokeWidth={2.2}
                        fill="url(#neonFill)"
                        fillOpacity={1}
                        dot={false}
                        activeDot={{ r: 5 }}
                        isAnimationActive={false}
                        focusable={false}
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>

            {/* Table */}
            <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-4 sm:p-6">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-lg font-semibold">Trades</h2>
                <span className="text-xs text-[#A1A1AA]">
                  {stats.tradesCount} trade{stats.tradesCount === 1 ? "" : "s"}
                </span>
              </div>

              <div className="mt-4 overflow-x-auto">
                <table className="min-w-[1180px] w-full text-sm">
                  <thead>
                    <tr className="text-left border-b border-white/10">
                      <Th>id</Th>
                      <Th>symbol</Th>
                      <Th>side</Th>
                      <Th>volume</Th>
                      <Th>open time</Th>
                      <Th>close time</Th>
                      <Th>open</Th>
                      <Th>close</Th>
                      <Th>net</Th>
                    </tr>
                  </thead>

                  <tbody className="divide-y divide-white/10">
                    {rows.length === 0 ? (
                      <tr>
                        <td
                          colSpan={9}
                          className="px-3 py-6 text-center text-[#A1A1AA]"
                        >
                          No trades yet.
                        </td>
                      </tr>
                    ) : (
                      rows.map((t) => {
                        const isWin = t.total >= 0;

                        return (
                          <tr
                            key={t.id}
                            className="hover:bg-white/5 transition"
                          >
                            <Td className="text-white/90">{t.id}</Td>
                            <Td className="font-semibold">{t.symbol}</Td>
                            <Td>
                              <SideBadge side={t.side} />
                            </Td>
                            <Td className="text-white/90">{t.volume}</Td>
                            <Td className="text-white/90">
                              {formatDateTime(t.openTime)}
                            </Td>
                            <Td className="text-white/90">
                              {formatDateTime(t.closeTime)}
                            </Td>
                            <Td className="text-white/90">{t.openPrice}</Td>
                            <Td className="text-white/90">{t.closePrice}</Td>
                            <Td
                              className={[
                                "font-semibold",
                                isWin ? "text-green-400" : "text-red-400",
                              ].join(" ")}
                            >
                              {formatMoney(t.total, { sign: true })}
                              <div className="text-[11px] font-normal text-[#A1A1AA]">
                                p:{formatMoney(t.profit)} c:
                                {formatMoney(t.commission)} s:
                                {formatMoney(t.swap)}
                              </div>
                            </Td>
                          </tr>
                        );
                      })
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
