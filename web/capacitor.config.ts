import { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.account.app',
  appName: 'Account',
  webDir: 'dist',
  server: {
    androidScheme: 'https'
  },
  plugins: {
    Storage: {
      name: 'AccountStorage'
    }
  }
};

export default config;
