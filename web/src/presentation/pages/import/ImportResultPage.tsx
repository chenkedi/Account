import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  SimpleGrid,
  Alert,
  AlertIcon,
} from '@chakra-ui/react';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  useImportState,
  useImportActions,
  useTransactionActions,
} from '../../../store';

function SuccessIcon() {
  return (
    <svg viewBox="0 0 24 24" width="80" height="80" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
      <polyline points="22 4 12 14.01 9 11.01" />
    </svg>
  );
}

function ErrorIcon() {
  return (
    <svg viewBox="0 0 24 24" width="80" height="80" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="8" x2="12" y2="12" />
      <line x1="12" y1="16" x2="12.01" y2="16" />
    </svg>
  );
}

export function ImportResultPage() {
  const navigate = useNavigate();

  const { result } = useImportState();
  const { resetImport } = useImportActions();
  const { fetchTransactions } = useTransactionActions();

  useEffect(() => {
    if (!result) {
      navigate('/import');
    }
  }, [result, navigate]);

  if (!result) {
    return null;
  }

  const handleViewTransactions = async () => {
    await fetchTransactions();
    navigate('/transactions');
  };

  const handleImportMore = () => {
    resetImport();
    navigate('/import');
  };

  return (
    <VStack spacing={6} align="stretch">
      {/* 头部 */}
      <Box>
        <Heading
          size="2xl"
          bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
          bgClip="text"
          letterSpacing="-0.03em"
        >
          导入完成
        </Heading>
      </Box>

      {/* 结果状态 */}
      <Card variant="elevated">
        <CardBody py={12}>
          <VStack spacing={6}>
            <Box color={result.success ? 'green.500' : 'red.500'}>
              {result.success ? <SuccessIcon /> : <ErrorIcon />}
            </Box>
            <VStack spacing={2}>
              <Text fontSize="2xl" fontWeight="800" color="gray.800">
                {result.success ? '导入成功！' : '导入失败'}
              </Text>
              <Text color="gray.500" textAlign="center" maxW="sm">
                {result.success
                  ? '您的交易记录已成功导入'
                  : '导入过程中出现错误，请重试'}
              </Text>
            </VStack>
          </VStack>
        </CardBody>
      </Card>

      {/* 统计信息 */}
      <SimpleGrid columns={{ base: 2, md: 4 }} spacing={4}>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                成功导入
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="green.600">
                {result.importedCount}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                跳过
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="orange.600">
                {result.skippedCount}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                失败
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="red.600">
                {result.errorCount}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                总计
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="gray.700">
                {result.importedCount + result.skippedCount + result.errorCount}
              </Text>
            </VStack>
          </CardBody>
        </Card>
      </SimpleGrid>

      {/* 错误信息 */}
      {result.errors && result.errors.length > 0 && (
        <Card variant="elevated">
          <CardBody>
            <Heading size="md" mb={4} color="gray.700">
              错误详情
            </Heading>
            <VStack spacing={3} align="stretch">
              {result.errors.slice(0, 10).map((error, index) => (
                <Alert key={index} status="error" borderRadius="lg">
                  <AlertIcon />
                  <Text>
                    第 {error.index + 1} 行: {error.message}
                  </Text>
                </Alert>
              ))}
              {result.errors.length > 10 && (
                <Text color="gray.500" textAlign="center">
                  还有 {result.errors.length - 10} 条错误...
                </Text>
              )}
            </VStack>
          </CardBody>
        </Card>
      )}

      {/* 操作按钮 */}
      <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4}>
        <Button
          size="lg"
          height="14"
          borderRadius="2xl"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          onClick={handleViewTransactions}
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          transition="all 0.2s ease"
        >
          查看交易
        </Button>
        <Button
          size="lg"
          height="14"
          borderRadius="2xl"
          variant="outline"
          borderWidth="2px"
          onClick={handleImportMore}
          _hover={{
            borderColor: 'brand.300',
            bg: 'brand.50',
            transform: 'translateY(-2px)',
            boxShadow: 'md',
          }}
          transition="all 0.2s ease"
        >
          继续导入
        </Button>
      </SimpleGrid>
    </VStack>
  );
}
