import { format, parseISO, startOfDay, endOfDay, startOfMonth, endOfMonth, startOfWeek, endOfWeek, startOfYear, endOfYear, subDays, subMonths, subWeeks, subYears } from 'date-fns';
import { zhCN } from 'date-fns/locale';

const LOCALE = zhCN;

export const dateUtils = {
  format: (date: Date | string, fmt: string = 'yyyy-MM-dd'): string => {
    const d = typeof date === 'string' ? parseISO(date) : date;
    return format(d, fmt, { locale: LOCALE });
  },

  formatDateTime: (date: Date | string): string => {
    return dateUtils.format(date, 'yyyy-MM-dd HH:mm:ss');
  },

  formatDisplay: (date: Date | string): string => {
    return dateUtils.format(date, 'yyyy年M月d日');
  },

  now: (): Date => new Date(),

  nowISO: (): string => new Date().toISOString(),

  startOfDay,
  endOfDay,
  startOfMonth,
  endOfMonth,
  startOfWeek,
  endOfWeek,
  startOfYear,
  endOfYear,
  subDays,
  subMonths,
  subWeeks,
  subYears,

  getDateRange: (type: 'week' | 'month' | 'quarter' | 'year' | '6months'): { start: Date; end: Date } => {
    const now = new Date();
    switch (type) {
      case 'week':
        return { start: startOfWeek(now, { weekStartsOn: 1 }), end: endOfWeek(now, { weekStartsOn: 1 }) };
      case 'month':
        return { start: startOfMonth(now), end: endOfMonth(now) };
      case 'quarter':
        const qStart = startOfMonth(subMonths(now, 2));
        return { start: qStart, end: endOfMonth(now) };
      case '6months':
        const hStart = startOfMonth(subMonths(now, 5));
        return { start: hStart, end: endOfMonth(now) };
      case 'year':
        return { start: startOfYear(now), end: endOfYear(now) };
      default:
        return { start: startOfMonth(now), end: endOfMonth(now) };
    }
  },
};
