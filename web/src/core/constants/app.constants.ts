export const APP_CONSTANTS = {
  appName: 'Account',
  storagePrefix: 'account_',
  defaultCurrency: 'CNY',
  dateFormat: 'yyyy-MM-dd',
  dateTimeFormat: 'yyyy-MM-dd HH:mm:ss',
} as const;

export const TRANSACTION_TYPES = {
  INCOME: 'income',
  EXPENSE: 'expense',
  TRANSFER: 'transfer',
} as const;

export const IMPORT_SOURCES = {
  ALIPAY: 'alipay',
  WECHAT: 'wechat',
  BANK: 'bank',
  GENERIC: 'generic',
} as const;
