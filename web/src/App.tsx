import { useEffect } from 'react';
import { ThemeProvider } from './presentation/providers/ThemeProvider';
import { AppRouter } from './presentation/router/AppRouter';
import { useAuthActions } from './store';

function App() {
  const { checkAuth } = useAuthActions();

  useEffect(() => {
    // 检查认证状态
    checkAuth();
  }, [checkAuth]);

  return (
    <ThemeProvider>
      <AppRouter />
    </ThemeProvider>
  );
}

export default App;
