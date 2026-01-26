"use client";

import { webApi } from "@/app/lib/api/client";
import { useRouter } from "next/navigation";
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";

type Row = {
  tradingAccountLogin: number;
  rank: number;
  username: string;
  accountSize: number;
  profit: number;
  equity: number;
  gainPercent: number;
};

function medal(rank: number) {
  if (rank === 1) return "🥇";
  if (rank === 2) return "🥈";
  if (rank === 3) return "🥉";
  return null;
}

// FLIP animation helper
function animateReorder(
  prevRects: Map<string, DOMRect>,
  rowEls: Map<string, HTMLTableRowElement>,
) {
  rowEls.forEach((el, key) => {
    const prev = prevRects.get(key);
    if (!prev) return;

    const next = el.getBoundingClientRect();
    const dy = prev.top - next.top;

    if (Math.abs(dy) < 1) return;

    el.animate(
      [{ transform: `translateY(${dy}px)` }, { transform: "translateY(0px)" }],
      { duration: 260, easing: "cubic-bezier(0.2, 0.8, 0.2, 1)" },
    );
  });
}

export default function LeaderboardTable({
  competitionId,
  endsAt,
}: {
  competitionId: string; // ✅ UUID
  endsAt: string;
}) {
  const router = useRouter();

  const [rows, setRows] = useState<Row[]>([]);
  const [loading, setLoading] = useState(true);

  // ✅ avoid Date.now in render: compute ended in state
  const [ended, setEnded] = useState(false);

  // keep row refs for FLIP
  const rowElsRef = useRef<Map<string, HTMLTableRowElement>>(new Map());
  const prevRectsRef = useRef<Map<string, DOMRect> | null>(null);

  // stable row key: username is unique enough for leaderboard
  const rowKey = (r: Row) => r.username;

  async function fetchLeaderboard() {
    const { data } = await webApi.get(
      `/competitions/${competitionId}/leaderboard`,
    );
    setRows(Array.isArray(data) ? data : []);
    setLoading(false);
  }

  // ✅ compute ended, update once per second (cheap) until ended true
  useEffect(() => {
    const endsAtMs = new Date(endsAt).getTime();

    const tick = () => setEnded(Date.now() > endsAtMs);
    tick();

    const t = setInterval(tick, 1000);
    return () => clearInterval(t);
  }, [endsAt]);

  // capture previous positions before DOM updates after rows change
  useLayoutEffect(() => {
    const rects = new Map<string, DOMRect>();
    rowElsRef.current.forEach((el, key) =>
      rects.set(key, el.getBoundingClientRect()),
    );
    prevRectsRef.current = rects;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rows.length]);

  // run FLIP animation after rows change
  useLayoutEffect(() => {
    if (!prevRectsRef.current) return;
    animateReorder(prevRectsRef.current, rowElsRef.current);
    prevRectsRef.current = null;
  }, [rows]);

  useEffect(() => {
    let alive = true;

    const run = async () => {
      try {
        await fetchLeaderboard();
      } catch {
        // later UX
      }
    };

    run();

    if (!ended) {
      const t = setInterval(() => {
        if (!alive) return;
        run();
      }, 60 * 1000);

      return () => {
        alive = false;
        clearInterval(t);
      };
    }

    return () => {
      alive = false;
    };
  }, [competitionId, ended]);

  const subtitle = useMemo(
    () => (ended ? "Final results" : "Updates every minute"),
    [ended],
  );

  return (
    <div className="rounded-2xl border border-white/10 bg-[#151621]/80 backdrop-blur p-6">
      <div className="flex items-center justify-between">
        <h2 className="text-lg sm:text-xl font-semibold">Leaderboard</h2>
        <span className="text-xs text-[#A1A1AA]">{subtitle}</span>
      </div>

      <div className="mt-4 overflow-x-auto">
        <table className="min-w-[780px] w-full text-md">
          <thead>
            <tr className="text-left text-[#A1A1AA] border-b border-white/10">
              <th className="px-3 py-2">#</th>
              <th className="px-3 py-2">User</th>
              <th className="px-3 py-2">Account</th>
              <th className="px-3 py-2">Profit</th>
              <th className="px-3 py-2">Equity</th>
              <th className="px-3 py-2">Gain %</th>
            </tr>
          </thead>

          <tbody className="divide-y divide-white/10">
            {loading ? (
              <tr>
                <td
                  colSpan={6}
                  className="px-3 py-6 text-center text-[#A1A1AA]"
                >
                  Loading…
                </td>
              </tr>
            ) : rows.length === 0 ? (
              <tr>
                <td
                  colSpan={6}
                  className="px-3 py-6 text-center text-[#A1A1AA]"
                >
                  No data yet.
                </td>
              </tr>
            ) : (
              rows.map((r) => {
                const pos = r.gainPercent >= 0;

                return (
                  <tr
                    key={rowKey(r)}
                    className={[
                      "hover:bg-white/5 transition cursor-pointer",
                      r.rank <= 3 ? "bg-white/[0.03]" : "",
                    ].join(" ")}
                    onClick={() =>
                      router.push(`/trade-history/${r.tradingAccountLogin}`)
                    }
                  >
                    <td className="px-3 py-2 font-semibold">
                      {r.rank <= 3 ? (
                        <span className="text-xl">{medal(r.rank)}</span>
                      ) : (
                        <span className="text-[#E5E7EB]">{r.rank}</span>
                      )}
                    </td>

                    <td className="px-3 py-2 font-medium text-[#BFDBFE]">
                      {r.username}
                    </td>
                    <td className="px-3 py-2">
                      ${r.accountSize.toLocaleString()}
                    </td>
                    <td
                      className={`px-3 py-2 font-semibold ${pos ? "text-green-400" : "text-red-400"}`}
                    >
                      {pos ? "+" : ""}${r.profit.toFixed(2)}
                    </td>
                    <td className="px-3 py-2">${r.equity.toFixed(2)}</td>
                    <td
                      className={`px-3 py-2 font-semibold ${pos ? "text-green-400" : "text-red-400"}`}
                    >
                      {pos ? "+" : ""}
                      {r.gainPercent.toFixed(2)}%
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
