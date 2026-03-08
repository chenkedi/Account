import { useEffect, useCallback } from 'react';
import { ThemeProvider } from './presentation/providers/ThemeProvider';
import { AppRouter } from './presentation/router/AppRouter';
import { useAuthActions } from './store';
import { authEvents, AUTH_EVENTS } from './core/events/auth.events';

function App() {
  const { checkAuth, logout } = useAuthActions();

  const handleSessionExpired = useCallback(() => {
    console.log('[App] Session expired, logging out');
    logout();
  }, [logout]);

  useEffect(() => {
    // 检查认证状态
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    // 监听会话过期事件
    authEvents.on(AUTH_EVENTS.SESSION_EXPIRED, handleSessionExpired);

    // 清理时取消监听
    return () => {
      authEvents.off(AUTH_EVENTS.SESSION_EXPIRED, handleSessionExpired);
    };
  }, [handleSessionExpired]);

  return (
    <ThemeProvider>
      <AppRouter />
    </ThemeProvider>
  );
}

export default App;
