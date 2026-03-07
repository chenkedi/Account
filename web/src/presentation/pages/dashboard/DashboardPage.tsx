import {
  Box,
  Button,
  Card,
  CardBody,
  Flex,
  Heading,
  SimpleGrid,
  Text,
  VStack,
} from '@chakra-ui/react';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { startOfMonth, endOfMonth } from 'date-fns';
import { useStats, useStatsState, useStatsActions, useAccounts, useAccountsState, useAccountActions } from '../../../store';
import { formatAmount } from '../../../core/utils/amount.utils';

// 问候语组件
function Greeting() {
  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return '早上好';
    if (hour < 18) return '下午好';
    return '晚上好';
  };

  return (
    <Box>
      <Text
        color="gray.500"
        fontWeight="500"
        fontSize="md"
        mb={1}
      >
        {getGreeting()}！
      </Text>
      <Heading
        size="2xl"
        bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
        bgClip="text"
        letterSpacing="-0.03em"
      >
        财务概览
      </Heading>
    </Box>
  );
}

// 统计卡片组件
interface StatCardProps {
  label: string;
  value: string;
  type: 'income' | 'expense' | 'balance' | 'savings';
  icon: React.ReactNode;
}

function StatCard({ label, value, type, icon }: StatCardProps) {
  const gradientMap = {
    income: 'linear(135deg, income.500 0%, income.600 100%)',
    expense: 'linear(135deg, expense.500 0%, expense.600 100%)',
    balance: 'linear(135deg, brand.500 0%, brand.600 100%)',
    savings: 'linear(135deg, transfer.500 0%, transfer.600 100%)',
  };

  const bgMap = {
    income: 'income.50',
    expense: 'expense.50',
    balance: 'brand.50',
    savings: 'transfer.50',
  };

  const colorMap = {
    income: 'income.600',
    expense: 'expense.600',
    balance: 'brand.600',
    savings: 'transfer.600',
  };

  return (
    <Card
      variant="elevated"
      overflow="hidden"
      _hover={{
        transform: 'translateY(-4px)',
        boxShadow: 'lg',
      }}
      transition="all 0.3s ease"
    >
      <CardBody p={6}>
        <Flex justify="space-between" align="flex-start">
          <VStack align="flex-start" spacing={3}>
            <Text
              color="gray.500"
              fontWeight="600"
              fontSize="sm"
            >
              {label}
            </Text>
            <Text
              fontSize="3xl"
              fontWeight="800"
              color={colorMap[type]}
              letterSpacing="-0.02em"
            >
              {value}
            </Text>
          </VStack>
          <Box
            w={14}
            h={14}
            borderRadius="2xl"
            bg={bgMap[type]}
            display="flex"
            alignItems="center"
            justifyContent="center"
            color={colorMap[type]}
            fontSize="2xl"
          >
            {icon}
          </Box>
        </Flex>
        {/* 底部装饰条 */}
        <Box
          mt={4}
          w="full"
          h="1"
          bgGradient={gradientMap[type]}
          borderRadius="full"
          opacity={0.3}
        />
      </CardBody>
    </Card>
  );
}

// 快捷操作按钮
interface QuickActionProps {
  label: string;
  icon: React.ReactNode;
  colorScheme?: 'income' | 'expense' | 'brand' | 'transfer';
  onClick?: () => void;
}

function QuickAction({ label, icon, colorScheme = 'brand', onClick }: QuickActionProps) {
  const colorMap = {
    income: 'income',
    expense: 'expense',
    brand: 'brand',
    transfer: 'transfer',
  };

  return (
    <Button
      onClick={onClick}
      variant="outline"
      borderWidth="2px"
      height="24"
      borderRadius="2xl"
      flexDirection="column"
      gap={2}
      colorScheme={colorMap[colorScheme]}
      _hover={{
        transform: 'translateY(-2px)',
        boxShadow: 'md',
      }}
      transition="all 0.2s ease"
    >
      <Box fontSize="2xl">{icon}</Box>
      <Text fontSize="sm" fontWeight="600">{label}</Text>
    </Button>
  );
}

// 图标组件
function IncomeIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="12" y1="1" x2="12" y2="23" />
      <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
    </svg>
  );
}

function ExpenseIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="5" y1="12" x2="19" y2="12" />
      <path d="M12 5v14" />
    </svg>
  );
}

function BalanceIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
    </svg>
  );
}

function SavingsIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 10h18V6a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7" />
      <polyline points="16 14 21 9 16 4" />
      <line x1="21" y1="9" x2="9" y2="9" />
    </svg>
  );
}

function PlusIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

function MinusIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

function ArrowRightLeftIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="17 1 21 5 17 9" />
      <path d="M3 11V9a4 4 0 0 1 4-4h14" />
      <polyline points="7 23 3 19 7 15" />
      <path d="M21 13v2a4 4 0 0 1-4 4H3" />
    </svg>
  );
}

function ChartIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 21H4.6c-.56 0-.84 0-1.054-.109a1 1 0 0 1-.437-.437C3 20.24 3 19.96 3 19.4V3" />
      <path d="m7 14 4-4 4 4 6-6" />
    </svg>
  );
}

function TransactionIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <path d="M7 10h10" />
      <path d="M7 14h6" />
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

export function DashboardPage() {
  const navigate = useNavigate();

  const stats = useStats();
  const { isLoading: statsLoading } = useStatsState();
  const accounts = useAccounts();
  const { isLoading: accountsLoading } = useAccountsState();
  const { fetchStats } = useStatsActions();
  const { fetchAccounts } = useAccountActions();

  useEffect(() => {
    const now = new Date();
    const start = startOfMonth(now);
    const end = endOfMonth(now);

    Promise.all([
      fetchStats(start, end).catch((err) => {
        console.error('Failed to fetch stats:', err);
      }),
      fetchAccounts().catch((err) => {
        console.error('Failed to fetch accounts:', err);
      }),
    ]);
  }, [fetchStats, fetchAccounts]);

  const totalIncome = stats?.totalIncome || 0;
  const totalExpense = stats?.totalExpense || 0;
  const netIncome = stats?.netIncome || 0;

  const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0);

  return (
    <VStack spacing={8} align="stretch">
      {/* 问候语 */}
      <Greeting />

      {/* 统计卡片 */}
      <SimpleGrid columns={{ base: 1, sm: 2 }} spacing={5}>
        <StatCard
          label="本月收入"
          value={statsLoading ? '...' : formatAmount(totalIncome)}
          type="income"
          icon={<IncomeIcon />}
        />
        <StatCard
          label="本月支出"
          value={statsLoading ? '...' : formatAmount(totalExpense)}
          type="expense"
          icon={<ExpenseIcon />}
        />
        <StatCard
          label="当前余额"
          value={accountsLoading ? '...' : formatAmount(totalBalance)}
          type="balance"
          icon={<BalanceIcon />}
        />
        <StatCard
          label="本月结余"
          value={statsLoading ? '...' : formatAmount(netIncome)}
          type="savings"
          icon={<SavingsIcon />}
        />
      </SimpleGrid>

      {/* 快捷操作 */}
      <Card variant="elevated">
        <CardBody p={6}>
          <Flex justify="space-between" align="center" mb={5}>
            <Heading size="md" color="gray.800" letterSpacing="-0.02em">
              快捷操作
            </Heading>
            <Text color="gray.400" fontSize="sm" fontWeight="500">
              快速记账
            </Text>
          </Flex>
          <SimpleGrid columns={2} spacing={4}>
            <QuickAction
              label="添加收入"
              icon={<PlusIcon />}
              colorScheme="income"
              onClick={() => navigate('/transactions/new?type=income')}
            />
            <QuickAction
              label="添加支出"
              icon={<MinusIcon />}
              colorScheme="expense"
              onClick={() => navigate('/transactions/new?type=expense')}
            />
            <QuickAction
              label="转账"
              icon={<ArrowRightLeftIcon />}
              colorScheme="transfer"
              onClick={() => navigate('/transactions/new?type=transfer')}
            />
            <QuickAction
              label="查看统计"
              icon={<ChartIcon />}
              colorScheme="brand"
              onClick={() => navigate('/stats')}
            />
          </SimpleGrid>
        </CardBody>
      </Card>

      {/* 主要操作按钮 */}
      <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4}>
        <Button
          as={RouterLink}
          to="/transactions"
          size="lg"
          height="16"
          borderRadius="2xl"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          fontSize="lg"
          fontWeight="700"
          leftIcon={<TransactionIcon />}
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          _active={{
            transform: 'translateY(0)',
          }}
          transition="all 0.2s ease"
        >
          查看交易
        </Button>
        <Button
          as={RouterLink}
          to="/import"
          size="lg"
          height="16"
          borderRadius="2xl"
          variant="outline"
          borderWidth="2px"
          borderColor="gray.200"
          fontSize="lg"
          fontWeight="700"
          leftIcon={<ImportIcon />}
          _hover={{
            borderColor: 'brand.300',
            bg: 'brand.50',
            transform: 'translateY(-2px)',
            boxShadow: 'md',
          }}
          _active={{
            transform: 'translateY(0)',
          }}
          transition="all 0.2s ease"
        >
          导入账单
        </Button>
      </SimpleGrid>
    </VStack>
  );
}
