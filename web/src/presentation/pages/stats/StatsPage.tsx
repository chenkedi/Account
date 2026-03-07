import {
  Box,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  HStack,
  Spinner,
  SimpleGrid,
} from '@chakra-ui/react';
import { useEffect, useCallback } from 'react';
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import {
  startOfMonth,
  endOfMonth,
  startOfQuarter,
  endOfQuarter,
  startOfYear,
  endOfYear,
  format,
} from 'date-fns';
import { zhCN } from 'date-fns/locale';
import {
  useStats,
  useStatsState,
  useStatsActions,
} from '../../../store';
import { formatAmount } from '../../../core/utils/amount.utils';

type TimeRange = 'month' | 'quarter' | 'year';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

function ChartIcon() {
  return (
    <svg viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 21H4.6c-.56 0-.84 0-1.054-.109a1 1 0 0 1-.437-.437C3 20.24 3 19.96 3 19.4V3" />
      <path d="m7 14 4-4 4 4 6-6" />
    </svg>
  );
}

function StatCard({ label, value, type }: { label: string; value: string; type: 'income' | 'expense' | 'net' }) {
  const colorMap = {
    income: { bg: 'income.50', color: 'income.600' },
    expense: { bg: 'expense.50', color: 'expense.600' },
    net: { bg: 'brand.50', color: 'brand.600' },
  };

  const colors = colorMap[type];

  return (
    <Card variant="elevated">
      <CardBody p={4}>
        <VStack align="start" spacing={2}>
          <Text fontSize="sm" fontWeight="600" color="gray.500">
            {label}
          </Text>
          <Text fontSize="2xl" fontWeight="800" color={colors.color}>
            {value}
          </Text>
        </VStack>
      </CardBody>
    </Card>
  );
}

function getDateRange(range: TimeRange): { start: Date; end: Date } {
  const now = new Date();
  switch (range) {
    case 'month':
      return { start: startOfMonth(now), end: endOfMonth(now) };
    case 'quarter':
      return { start: startOfQuarter(now), end: endOfQuarter(now) };
    case 'year':
      return { start: startOfYear(now), end: endOfYear(now) };
  }
}

export function StatsPage() {
  const stats = useStats();
  const { isLoading, timeRange } = useStatsState();
  const { fetchStats, setTimeRange } = useStatsActions();

  const loadStats = useCallback(
    (range: TimeRange) => {
      const { start, end } = getDateRange(range);
      fetchStats(start, end).catch((err) => {
        console.error('Failed to fetch stats:', err);
      });
    },
    [fetchStats]
  );

  useEffect(() => {
    loadStats(timeRange as TimeRange);
  }, [timeRange, loadStats]);

  const handleTabChange = (index: number) => {
    const ranges: TimeRange[] = ['month', 'quarter', 'year'];
    const newRange = ranges[index];
    setTimeRange(newRange);
  };

  const incomeByCategory =
    stats?.byCategory?.filter((c) => c.type === 'income') || [];
  const expenseByCategory =
    stats?.byCategory?.filter((c) => c.type === 'expense') || [];
  const monthlyTrend = stats?.monthlyTrend || [];

  const totalIncome = stats?.totalIncome || 0;
  const totalExpense = stats?.totalExpense || 0;
  const netIncome = stats?.netIncome || 0;

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
          统计报表
        </Heading>
        <Text color="gray.500" fontWeight="500" mt={1}>
          查看您的收支分析
        </Text>
      </Box>

      {/* 概览卡片 */}
      <SimpleGrid columns={{ base: 1, sm: 3 }} spacing={4}>
        <StatCard
          label="总收入"
          value={isLoading ? '...' : formatAmount(totalIncome)}
          type="income"
        />
        <StatCard
          label="总支出"
          value={isLoading ? '...' : formatAmount(totalExpense)}
          type="expense"
        />
        <StatCard
          label="净收入"
          value={isLoading ? '...' : formatAmount(netIncome)}
          type="net"
        />
      </SimpleGrid>

      {/* Tabs */}
      <Card variant="elevated" overflow="hidden">
        <Tabs
          variant="soft-rounded"
          colorScheme="brand"
          px={6}
          pt={6}
          onChange={handleTabChange}
        >
          <TabList bg="gray.50" p={1.5} borderRadius="2xl">
            <HStack spacing={2}>
              <Tab
                fontWeight="700"
                borderRadius="xl"
                _selected={{
                  bg: 'white',
                  color: 'brand.600',
                  boxShadow: 'sm',
                }}
              >
                本月
              </Tab>
              <Tab
                fontWeight="700"
                borderRadius="xl"
                _selected={{
                  bg: 'white',
                  color: 'brand.600',
                  boxShadow: 'sm',
                }}
              >
                本季度
              </Tab>
              <Tab
                fontWeight="700"
                borderRadius="xl"
                _selected={{
                  bg: 'white',
                  color: 'brand.600',
                  boxShadow: 'sm',
                }}
              >
                本年
              </Tab>
            </HStack>
          </TabList>

          <TabPanels>
            {(['month', 'quarter', 'year'] as const).map((range) => (
              <TabPanel key={range} px={0}>
                {isLoading ? (
                  <Card variant="outline" borderWidth="0">
                    <CardBody py={12}>
                      <VStack spacing={6}>
                        <Spinner size="xl" color="brand.500" />
                        <Text color="gray.500">加载中...</Text>
                      </VStack>
                    </CardBody>
                  </Card>
                ) : !stats ||
                  (incomeByCategory.length === 0 &&
                    expenseByCategory.length === 0) ? (
                  <Card variant="outline" borderWidth="0">
                    <CardBody py={12}>
                      <VStack spacing={6}>
                        <Box color="gray.300">
                          <ChartIcon />
                        </Box>
                        <VStack spacing={2}>
                          <Text fontSize="xl" fontWeight="700" color="gray.700">
                            暂无数据
                          </Text>
                          <Text color="gray.500" textAlign="center" maxW="sm">
                            添加交易后即可查看统计数据
                          </Text>
                        </VStack>
                      </VStack>
                    </CardBody>
                  </Card>
                ) : (
                  <VStack spacing={6} py={4}>
                    {/* 支出分类饼图 */}
                    {expenseByCategory.length > 0 && (
                      <Card variant="outline" w="full">
                        <CardBody>
                          <Heading size="md" mb={4} color="gray.700">
                            支出分类
                          </Heading>
                          <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                              <Pie
                                data={expenseByCategory}
                                cx="50%"
                                cy="50%"
                                labelLine={false}
                                label={({ categoryName, amount }) =>
                                  `${categoryName} ${formatAmount(amount)}`
                                }
                                outerRadius={100}
                                fill="#8884d8"
                                dataKey="amount"
                              >
                                {expenseByCategory.map((_, index) => (
                                  <Cell
                                    key={`cell-${index}`}
                                    fill={COLORS[index % COLORS.length]}
                                  />
                                ))}
                              </Pie>
                              <Tooltip
                                formatter={(value: number) => formatAmount(value)}
                              />
                            </PieChart>
                          </ResponsiveContainer>
                        </CardBody>
                      </Card>
                    )}

                    {/* 收入分类饼图 */}
                    {incomeByCategory.length > 0 && (
                      <Card variant="outline" w="full">
                        <CardBody>
                          <Heading size="md" mb={4} color="gray.700">
                            收入分类
                          </Heading>
                          <ResponsiveContainer width="100%" height={300}>
                            <PieChart>
                              <Pie
                                data={incomeByCategory}
                                cx="50%"
                                cy="50%"
                                labelLine={false}
                                label={({ categoryName, amount }) =>
                                  `${categoryName} ${formatAmount(amount)}`
                                }
                                outerRadius={100}
                                fill="#8884d8"
                                dataKey="amount"
                              >
                                {incomeByCategory.map((_, index) => (
                                  <Cell
                                    key={`cell-${index}`}
                                    fill={COLORS[index % COLORS.length]}
                                  />
                                ))}
                              </Pie>
                              <Tooltip
                                formatter={(value: number) => formatAmount(value)}
                              />
                            </PieChart>
                          </ResponsiveContainer>
                        </CardBody>
                      </Card>
                    )}

                    {/* 月度趋势图 */}
                    {monthlyTrend.length > 0 && (
                      <Card variant="outline" w="full">
                        <CardBody>
                          <Heading size="md" mb={4} color="gray.700">
                            月度趋势
                          </Heading>
                          <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={monthlyTrend}>
                              <CartesianGrid strokeDasharray="3 3" />
                              <XAxis
                                dataKey="month"
                                tickFormatter={(month) =>
                                  format(new Date(month + '-01'), 'M月', {
                                    locale: zhCN,
                                  })
                                }
                              />
                              <YAxis />
                              <Tooltip
                                formatter={(value: number) => formatAmount(value)}
                              />
                              <Legend />
                              <Bar
                                dataKey="income"
                                fill="#10b981"
                                name="收入"
                              />
                              <Bar
                                dataKey="expense"
                                fill="#ef4444"
                                name="支出"
                              />
                            </BarChart>
                          </ResponsiveContainer>
                        </CardBody>
                      </Card>
                    )}
                  </VStack>
                )}
              </TabPanel>
            ))}
          </TabPanels>
        </Tabs>
      </Card>
    </VStack>
  );
}
