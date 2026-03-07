export const amountUtils = {
  format: (amount: number, currency: string = 'CNY'): string => {
    return new Intl.NumberFormat('zh-CN', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2,
    }).format(amount);
  },

  formatSimple: (amount: number): string => {
    return amount.toFixed(2);
  },

  parse: (value: string): number => {
    const num = parseFloat(value.replace(/[^\d.-]/g, ''));
    return isNaN(num) ? 0 : num;
  },

  isPositive: (amount: number): boolean => amount > 0,
  isNegative: (amount: number): boolean => amount < 0,

  sum: (amounts: number[]): number => {
    return amounts.reduce((acc, curr) => acc + curr, 0);
  },
};

// Named export for convenience
export const formatAmount = (amount: number, currency: string = 'CNY'): string => {
  return amountUtils.format(amount, currency);
};
