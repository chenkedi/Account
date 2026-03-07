import { ChakraProvider } from '@chakra-ui/react';
import { appTheme } from '../../core/theme';

interface ThemeProviderProps {
  children: React.ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  return <ChakraProvider theme={appTheme}>{children}</ChakraProvider>;
}
