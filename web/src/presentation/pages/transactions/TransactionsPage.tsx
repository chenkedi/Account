import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  Flex,
  SimpleGrid,
  Spinner,
  useToast,
  HStack,
  Badge,
} from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';
import { useEffect } from 'react';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import type { Transaction } from '../../../data/models/transaction.model';
import { useTransactions, useTransactionsState, useTransactionActions } from '../../../store';
import { formatAmount } from '../../../core/utils/amount.utils';

function PlusIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

function ImportIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="17 8 12 3 7 8" />
      <line x1="12" y1="3" x2="12" y2="15" />
    </svg>
  );
}

function DownloadIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="7 10 12 15 17 10" />
      <line x1="12" y1="15" x2="12" y2="3" />
    </svg>
  );
}

function EmptyStateIcon() {
  return (
    <svg viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <path d="M2 10h20" />
      <path d="M6 16h4" />
      <path d="M14 16h4" />
    </svg>
  );
}

interface TransactionItemProps {
  transaction: Transaction;
}

function TransactionItem({ transaction }: TransactionItemProps) {

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'income':
        return 'income';
      case 'expense':
        return 'expense';
      case 'transfer':
        return 'brand';
      default:
        return 'gray';
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'income':
        return '收入';
      case 'expense':
        return '支出';
      case 'transfer':
        return '转账';
      default:
        return type;
    }
  };

  const getAmountPrefix = (type: string) => {
    switch (type) {
      case 'income':
        return '+';
      case 'expense':
        return '-';
      default:
        return '';
    }
  };

  return (
    <Card
      variant="outline"
      cursor="pointer"
      _hover={{
        transform: 'translateY(-2px)',
        boxShadow: 'md',
      }}
      transition="all 0.2s ease"
    >
      <CardBody p={4}>
        <Flex justify="space-between" align="center">
          <Flex align="center" gap={4}>
            <Box>
              <Text fontWeight="700" color="gray.800">
                {transaction.note || getTypeLabel(transaction.type)}
              </Text>
              <HStack mt={1} spacing={2}>
                <Badge colorScheme={getTypeColor(transaction.type)} variant="subtle">
                  {getTypeLabel(transaction.type)}
                </Badge>
                <Text fontSize="sm" color="gray.500">
                  {format(new Date(transaction.transaction_date), 'yyyy-MM-dd', {
                    locale: zhCN,
                  })}
                </Text>
              </HStack>
            </Box>
          </Flex>
          <Box>
            <Text
              fontSize="xl"
              fontWeight="800"
              color={
                transaction.type === 'income'
                  ? 'income.600'
                  : transaction.type === 'expense'
                  ? 'expense.600'
                  : 'brand.600'
              }
            >
              {getAmountPrefix(transaction.type)}
              {formatAmount(transaction.amount)}
            </Text>
          </Box>
        </Flex>
      </CardBody>
    </Card>
  );
}

export function TransactionsPage() {
  const toast = useToast();
  const transactions = useTransactions();
  const { isLoading } = useTransactionsState();
  const { fetchTransactions } = useTransactionActions();

  useEffect(() => {
    fetchTransactions().catch((err) => {
      toast({
        title: '加载失败',
        description: (err as Error).message,
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    });
  }, [fetchTransactions, toast]);

  return (
    <VStack spacing={6} align="stretch">
      {/* 头部 */}
      <Flex justify="space-between" align="center">
        <Box>
          <Heading
            size="2xl"
            bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
            bgClip="text"
            letterSpacing="-0.03em"
          >
            交易记录
          </Heading>
          <Text color="gray.500" fontWeight="500" mt={1}>
            管理您的收入和支出
          </Text>
        </Box>
        <Button
          as={RouterLink}
          to="/transactions/new"
          size="lg"
          borderRadius="2xl"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          leftIcon={<PlusIcon />}
          fontWeight="700"
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          transition="all 0.2s ease"
        >
          添加
        </Button>
      </Flex>

      {/* 加载状态 */}
      {isLoading && (
        <Card variant="elevated">
          <CardBody py={16}>
            <VStack spacing={4}>
              <Spinner size="xl" color="brand.500" />
              <Text color="gray.500">加载中...</Text>
            </VStack>
          </CardBody>
        </Card>
      )}

      {/* 空状态 */}
      {!isLoading && transactions.length === 0 && (
        <Card variant="elevated">
          <CardBody py={16}>
            <VStack spacing={6}>
              <Box color="gray.300">
                <EmptyStateIcon />
              </Box>
              <VStack spacing={2}>
                <Text fontSize="xl" fontWeight="700" color="gray.700">
                  暂无交易记录
                </Text>
                <Text color="gray.500" textAlign="center" maxW="sm">
                  开始记录您的第一笔交易，或导入历史账单
                </Text>
              </VStack>
            </VStack>
          </CardBody>
        </Card>
      )}

      {/* 交易列表 */}
      {!isLoading && transactions.length > 0 && (
        <VStack spacing={3} align="stretch">
          {transactions.map((transaction) => (
            <TransactionItem key={transaction.id} transaction={transaction} />
          ))}
        </VStack>
      )}

      {/* 操作按钮 */}
      <SimpleGrid columns={{ base: 1, sm: 2 }} spacing={4}>
        <Button
          as={RouterLink}
          to="/import"
          size="lg"
          height="14"
          borderRadius="2xl"
          variant="outline"
          borderWidth="2px"
          borderColor="gray.200"
          leftIcon={<ImportIcon />}
          fontWeight="700"
          _hover={{
            borderColor: 'brand.300',
            bg: 'brand.50',
            transform: 'translateY(-2px)',
            boxShadow: 'md',
          }}
          transition="all 0.2s ease"
        >
          导入账单
        </Button>
        <Button
          size="lg"
          height="14"
          borderRadius="2xl"
          variant="outline"
          borderWidth="2px"
          borderColor="gray.200"
          leftIcon={<DownloadIcon />}
          fontWeight="700"
          _hover={{
            borderColor: 'gray.300',
            bg: 'gray.50',
            transform: 'translateY(-2px)',
            boxShadow: 'md',
          }}
          transition="all 0.2s ease"
        >
          导出
        </Button>
      </SimpleGrid>
    </VStack>
  );
}
