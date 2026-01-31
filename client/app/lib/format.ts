type FormatMoneyOptions = {
  sign?: boolean;
  currency?: boolean;
};

export function formatMoney(value: number, options?: FormatMoneyOptions) {
  const { sign = false, currency = true } = options || {};

  const abs = Math.abs(value).toLocaleString("en-US", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

  const prefix = value < 0 ? "-" : value > 0 && sign ? "+" : "";

  const symbol = currency ? "$" : "";

  return `${prefix}${symbol}${abs}`;
}
