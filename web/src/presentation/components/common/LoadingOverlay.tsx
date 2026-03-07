import { Flex, Spinner, Text, VStack } from '@chakra-ui/react';

interface LoadingOverlayProps {
  message?: string;
  isLoading?: boolean;
}

export function LoadingOverlay({ message = '加载中...', isLoading = true }: LoadingOverlayProps) {
  if (!isLoading) return null;

  return (
    <Flex
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      bg="rgba(0,0,0,0.5)"
      alignItems="center"
      justifyContent="center"
      zIndex={9999}
    >
      <VStack spacing={4} bg="white" p={8} borderRadius="lg">
        <Spinner size="xl" color="brand.500" />
        <Text>{message}</Text>
      </VStack>
    </Flex>
  );
}
