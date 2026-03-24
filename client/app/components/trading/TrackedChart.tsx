"use client";

import { useEffect, useMemo, useRef } from "react";
import {
  CandlestickSeries,
  ColorType,
  CrosshairMode,
  LineStyle,
  createChart,
  type IPriceLine,
  type LogicalRange,
  type Time,
  type UTCTimestamp,
} from "lightweight-charts";

type TradeSide = "BUY" | "SELL";

export type Candle = {
  time: number;
  open: number;
  high: number;
  low: number;
  close: number;
};

export type Trade = {
  positionId: number;
  symbol: string;
  side: TradeSide;
  openPrice: number;
  volume: number;
  stopLoss: number | null;
  openedAt: string;
  closedAt: string | null;
};

type Props = {
  title: "EURUSD" | "GBPUSD" | "XAUUSD";
  candles: Candle[];
  trades: Trade[];
};

const COLORS = {
  bg: "#070814",
  panel: "#0b1020",
  text: "#d7d2ff",
  grid: "rgba(139, 92, 246, 0.10)",
  border: "rgba(139, 92, 246, 0.24)",
  glow: "rgba(139, 92, 246, 0.18)",
  candleUp: "#cfc8ff",
  candleDown: "#7c5cff",
  wickUp: "#ebe7ff",
  wickDown: "#9d87ff",
  buy: "#22c55e",
  sell: "#ef4444",
  stopLoss: "rgba(167, 139, 250, 0.85)",
};

function normalizeSymbol(raw: string) {
  if (!raw) return;

  const s = raw?.toUpperCase().replace(/\//g, "").replace(/\s+/g, "");
  if (s.startsWith("EURUSD")) return "EURUSD";
  if (s.startsWith("GBPUSD")) return "GBPUSD";
  if (s.startsWith("XAUUSD") || s.startsWith("GOLD")) return "XAUUSD";
  return s;
}

function toUnix(iso: string) {
  return Math.floor(new Date(iso).getTime() / 1000);
}

function getPriceConfig(symbol: "EURUSD" | "GBPUSD" | "XAUUSD") {
  if (symbol === "XAUUSD") {
    return {
      precision: 2,
      minMove: 0.01,
      formatter: (value: number) => value.toFixed(2),
    };
  }

  return {
    precision: 5,
    minMove: 0.00001,
    formatter: (value: number) => value.toFixed(5),
  };
}

function getZoomAwareRadius(params: {
  volume: number;
  visibleRange: LogicalRange | null;
  hostWidth: number;
}) {
  const { volume, visibleRange, hostWidth } = params;

  const minRadius = 4;
  const maxRadius = 14;

  const visibleBars =
    visibleRange &&
    Number.isFinite(visibleRange.from) &&
    Number.isFinite(visibleRange.to)
      ? Math.max(visibleRange.to - visibleRange.from, 1)
      : 144;

  const pxPerBar = hostWidth / visibleBars;

  const volumeClamped = Math.max(0.01, Math.min(volume, 3));
  const volumeRatio = Math.log10(volumeClamped + 1) / Math.log10(4);

  const baseFromZoom = pxPerBar * 0.32;
  const radius = baseFromZoom + volumeRatio * 4;

  return 2 * Math.max(minRadius, Math.min(radius, maxRadius));
}

function findNearestCandleTime(
  candles: Candle[],
  tradeUnix: number,
): number | null {
  if (candles.length === 0) return null;

  let nearest = candles[0].time;
  let nearestDiff = Math.abs(candles[0].time - tradeUnix);

  for (let i = 1; i < candles.length; i += 1) {
    const diff = Math.abs(candles[i].time - tradeUnix);
    if (diff < nearestDiff) {
      nearest = candles[i].time;
      nearestDiff = diff;
    }
  }

  return nearest;
}

export default function TrackedChart({ title, candles, trades }: Props) {
  const hostRef = useRef<HTMLDivElement | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
  const seriesRef = useRef<any>(null);
  const priceLinesRef = useRef<IPriceLine[]>([]);

  const filteredTrades = useMemo(
    () => trades.filter((trade) => normalizeSymbol(trade?.symbol) === title),
    [trades, title],
  );

  const priceConfig = useMemo(() => getPriceConfig(title), [title]);

  useEffect(() => {
    if (!hostRef.current) return;

    const chart = createChart(hostRef.current, {
      width: hostRef.current.clientWidth,
      height: hostRef.current.clientHeight,
      layout: {
        background: { type: ColorType.Solid, color: COLORS.bg },
        textColor: COLORS.text,
        attributionLogo: false,
      },
      grid: {
        vertLines: { color: COLORS.grid },
        horzLines: { color: COLORS.grid },
      },
      crosshair: {
        mode: CrosshairMode.Normal,
      },
      rightPriceScale: {
        borderColor: COLORS.border,
        ticksVisible: true,
        minimumWidth: 74,
        scaleMargins: {
          top: 0.08,
          bottom: 0.08,
        },
      },
      timeScale: {
        borderColor: COLORS.border,
        timeVisible: true,
        secondsVisible: false,
        rightOffset: 4,
        barSpacing: 8,
      },
      handleScroll: {
        mouseWheel: true,
        pressedMouseMove: true,
        horzTouchDrag: true,
        vertTouchDrag: true,
      },
      handleScale: {
        mouseWheel: true,
        pinch: true,
        axisPressedMouseMove: true,
      },
    });

    const series = chart.addSeries(CandlestickSeries, {
      upColor: COLORS.candleUp,
      downColor: COLORS.candleDown,
      borderUpColor: COLORS.candleUp,
      borderDownColor: COLORS.candleDown,
      wickUpColor: COLORS.wickUp,
      wickDownColor: COLORS.wickDown,
      priceLineVisible: false,
      lastValueVisible: false,
      priceFormat: {
        type: "price",
        precision: priceConfig.precision,
        minMove: priceConfig.minMove,
      },
    });

    chartRef.current = chart;
    seriesRef.current = series;

    const drawOverlay = () => {
      const canvas = canvasRef.current;
      const host = hostRef.current;
      const currentChart = chartRef.current;
      const currentSeries = seriesRef.current;

      if (!canvas || !host || !currentChart || !currentSeries) return;

      const dpr = window.devicePixelRatio || 1;
      const width = host.clientWidth;
      const height = host.clientHeight;

      canvas.width = Math.floor(width * dpr);
      canvas.height = Math.floor(height * dpr);
      canvas.style.width = `${width}px`;
      canvas.style.height = `${height}px`;

      const ctx = canvas.getContext("2d");
      if (!ctx) return;

      ctx.setTransform(1, 0, 0, 1, 0, 0);
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.scale(dpr, dpr);

      const visibleRange = currentChart.timeScale().getVisibleLogicalRange();

      for (const trade of filteredTrades) {
        const tradeUnix = toUnix(trade.openedAt);
        const alignedTradeTime = findNearestCandleTime(candles, tradeUnix);

        if (alignedTradeTime === null) continue;

        const x = currentChart
          .timeScale()
          .timeToCoordinate(alignedTradeTime as Time);
        const y = currentSeries.priceToCoordinate(trade.openPrice);

        if (x == null || y == null) continue;

        const color = trade.side === "BUY" ? COLORS.buy : COLORS.sell;
        const baseRadius = getZoomAwareRadius({
          volume: trade.volume,
          visibleRange,
          hostWidth: width,
        });

        const isClosed = Boolean(trade.closedAt);
        const radius = isClosed ? baseRadius * 0.88 : baseRadius;

        // open trade dobija blagi outer glow
        if (!isClosed) {
          ctx.beginPath();
          ctx.arc(x, y, radius + 2, 0, Math.PI * 2);
          ctx.strokeStyle = color;
          ctx.globalAlpha = 0.18;
          ctx.lineWidth = 1.5;
          ctx.stroke();
          ctx.globalAlpha = 1;
        }

        // main circle
        ctx.beginPath();
        ctx.arc(x, y, radius, 0, Math.PI * 2);
        ctx.strokeStyle = color;
        ctx.lineWidth = isClosed ? 2.2 : 1.5;

        // open = pun
        // closed = tamni centar + obojena ivica
        ctx.fillStyle = isClosed ? COLORS.bg : color;
        ctx.fill();
        ctx.stroke();

        // anchor marker for open SL on same trade x-position
        if (!trade.closedAt && trade.stopLoss !== null) {
          const slY = currentSeries.priceToCoordinate(trade.stopLoss);
          if (slY != null) {
            const slRadius = Math.max(2.5, Math.min(radius * 0.45, 5));

            ctx.beginPath();
            ctx.arc(x, slY, slRadius, 0, Math.PI * 2);
            ctx.fillStyle = COLORS.bg;
            ctx.strokeStyle = COLORS.stopLoss;
            ctx.lineWidth = 1.5;
            ctx.fill();
            ctx.stroke();
          }
        }
      }
    };

    const resizeObserver = new ResizeObserver(() => {
      if (!hostRef.current || !chartRef.current) return;

      chartRef.current.applyOptions({
        width: hostRef.current.clientWidth,
        height: hostRef.current.clientHeight,
      });

      drawOverlay();
    });

    resizeObserver.observe(hostRef.current);

    chart.timeScale().subscribeVisibleLogicalRangeChange(drawOverlay);

    requestAnimationFrame(drawOverlay);

    return () => {
      resizeObserver.disconnect();
      chart.timeScale().unsubscribeVisibleLogicalRangeChange(drawOverlay);
      chart.remove();
    };
  }, [filteredTrades, priceConfig, title]);

  useEffect(() => {
    const chart = chartRef.current;
    const series = seriesRef.current;

    if (!chart || !series || candles.length === 0) return;

    series.setData(
      candles.map((c) => ({
        time: c.time as UTCTimestamp,
        open: c.open,
        high: c.high,
        low: c.low,
        close: c.close,
      })),
    );

    // last 12h on 5m = 144 bars
    const lastIndex = candles.length - 1;
    const from = Math.max(lastIndex - 143, 0);
    const to = lastIndex + 4;
    chart.timeScale().setVisibleLogicalRange({ from, to });

    for (const line of priceLinesRef.current) {
      series.removePriceLine(line);
    }
    priceLinesRef.current = [];

    filteredTrades
      .filter((trade) => !trade.closedAt && trade.stopLoss !== null)
      .forEach((trade) => {
        const line = series.createPriceLine({
          price: trade.stopLoss as number,
          color: COLORS.stopLoss,
          lineWidth: 1,
          lineStyle: LineStyle.Dashed,
          axisLabelVisible: false,
          title: "",
        });

        priceLinesRef.current.push(line);
      });

    requestAnimationFrame(() => {
      window.dispatchEvent(new Event("resize"));
    });
  }, [candles, filteredTrades]);

  return (
    <section
      className="rounded-3xl border p-4"
      style={{
        background: `linear-gradient(180deg, ${COLORS.panel} 0%, ${COLORS.bg} 100%)`,
        borderColor: COLORS.border,
        boxShadow: `0 0 24px ${COLORS.glow}`,
      }}
    >
      <h2 className="mb-3 text-lg font-semibold text-white">{title}</h2>

      <div className="relative h-[420px] w-full overflow-hidden rounded-2xl">
        <div ref={hostRef} className="h-full w-full" />
        <canvas
          ref={canvasRef}
          className="pointer-events-none absolute inset-0 z-10"
        />
      </div>
    </section>
  );
}
