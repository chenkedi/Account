export const API_CONSTANTS = {
  baseUrl: 'http://localhost:8080/api/v1',
  apiPrefix: '/api/v1',
  wsBaseUrl: 'ws://localhost:8080',

  // Auth endpoints
  auth: {
    register: '/auth/register',
    login: '/auth/login',
  },

  // Account endpoints
  accounts: '/accounts',

  // Category endpoints
  categories: '/categories',

  // Transaction endpoints
  transactions: '/transactions',
  transactionsRange: '/transactions/range',
  transactionsStats: '/transactions/stats',
  transactionsStatsDetailed: '/transactions/stats/detailed',

  // Sync endpoints
  syncPull: '/sync/pull',
  syncPush: '/sync/push',

  // Import endpoints
  importSources: '/import/sources',
  importTemplate: '/import/template',
  importUpload: '/import/upload',
  importExecute: '/import/execute',

  // WebSocket
  wsSync: '/ws/sync',
} as const;
