import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  HStack,
  SimpleGrid,
  useToast,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Badge,
  Checkbox,
  Select,
} from '@chakra-ui/react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  useImportState,
  useImportActions,
  useAccounts,
  useCategories,
} from '../../../store';
import { formatAmount } from '../../../core/utils/amount.utils';

export function ImportPreviewPage() {
  const navigate = useNavigate();
  const toast = useToast();

  const { preview, isLoading } = useImportState();
  const { confirmImport, resetImport } = useImportActions();
  const accounts = useAccounts();
  const categories = useCategories();

  const [selectedTransactions, setSelectedTransactions] = useState<
    Record<number, { accountId: string | null; categoryId: string | null; note: string | null }>
  >({});

  if (!preview) {
    return (
      <VStack spacing={6} align="stretch">
        <Box>
          <Heading
            size="2xl"
            bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
            bgClip="text"
            letterSpacing="-0.03em"
          >
            预览导入
          </Heading>
        </Box>
        <Card variant="elevated">
          <CardBody py={12}>
            <VStack spacing={6}>
              <Text color="gray.500">
                没有可预览的数据
              </Text>
              <Button
                onClick={() => navigate('/import')}
                size="lg"
                bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
                color="white"
              >
                返回导入
              </Button>
            </VStack>
          </CardBody>
        </Card>
      </VStack>
    );
  }

  const handleConfirm = async () => {
    try {
      const transactions = preview.transactions.map((t) => ({
        rawIndex: t.rawIndex,
        accountId: selectedTransactions[t.rawIndex]?.accountId || null,
        categoryId: selectedTransactions[t.rawIndex]?.categoryId || null,
        note: selectedTransactions[t.rawIndex]?.note || null,
      }));

      const result = await confirmImport(transactions);

      toast({
        title: '导入成功',
        description: `成功导入 ${result.importedCount} 条记录`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });

      navigate('/import/result');
    } catch (error) {
      toast({
        title: '导入失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const handleCancel = () => {
    resetImport();
    navigate('/import');
  };

  const toggleTransaction = (index: number, checked: boolean) => {
    if (!checked) {
      const newSelected = { ...selectedTransactions };
      delete newSelected[index];
      setSelectedTransactions(newSelected);
    } else {
      setSelectedTransactions({
        ...selectedTransactions,
        [index]: {
          accountId: null,
          categoryId: null,
          note: null,
        },
      });
    }
  };

  const isAllSelected = preview.transactions.every(
    (t) => selectedTransactions[t.rawIndex] !== undefined
  );

  const toggleAll = () => {
    if (isAllSelected) {
      setSelectedTransactions({});
    } else {
      const newSelected: Record<
        number,
        { accountId: string | null; categoryId: string | null; note: string | null }
      > = {};
      preview.transactions.forEach((t) => {
        newSelected[t.rawIndex] = {
          accountId: null,
          categoryId: null,
          note: null,
        };
      });
      setSelectedTransactions(newSelected);
    }
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
          预览导入
        </Heading>
        <Text color="gray.500" fontWeight="500" mt={1}>
          文件名: {preview.fileName}
        </Text>
      </Box>

      {/* 统计信息 */}
      <SimpleGrid columns={{ base: 2, md: 4 }} spacing={4}>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                总记录
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="gray.700">
                {preview.totalRecords}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                有效记录
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="green.600">
                {preview.validRecords}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                重复记录
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="orange.600">
                {preview.duplicateRecords}
              </Text>
            </VStack>
          </CardBody>
        </Card>
        <Card variant="elevated">
          <CardBody p={4}>
            <VStack align="start">
              <Text fontSize="sm" fontWeight="600" color="gray.500">
                已选择
              </Text>
              <Text fontSize="2xl" fontWeight="800" color="brand.600">
                {Object.keys(selectedTransactions).length}
              </Text>
            </VStack>
          </CardBody>
        </Card>
      </SimpleGrid>

      {/* 交易预览表格 */}
      <Card variant="elevated" overflow="hidden">
        <CardBody p={0}>
          <Box overflowX="auto">
            <Table variant="simple">
              <Thead bg="gray.50">
                <Tr>
                  <Th>
                    <Checkbox
                      isChecked={isAllSelected}
                      onChange={toggleAll}
                    />
                  </Th>
                  <Th>日期</Th>
                  <Th>类型</Th>
                  <Th>金额</Th>
                  <Th>描述</Th>
                  <Th>对方</Th>
                  <Th>账户</Th>
                  <Th>分类</Th>
                  <Th>状态</Th>
                </Tr>
              </Thead>
              <Tbody>
                {preview.transactions.map((transaction) => (
                  <Tr key={transaction.rawIndex}>
                    <Td>
                      <Checkbox
                        isChecked={selectedTransactions[transaction.rawIndex] !== undefined}
                        onChange={(e) =>
                          toggleTransaction(transaction.rawIndex, e.target.checked)
                        }
                      />
                    </Td>
                    <Td>{transaction.date}</Td>
                    <Td>
                      <Badge
                        colorScheme={transaction.type === 'income' ? 'green' : 'red'}
                      >
                        {transaction.type === 'income' ? '收入' : '支出'}
                      </Badge>
                    </Td>
                    <Td
                      fontWeight="700"
                      color={transaction.type === 'income' ? 'green.600' : 'red.600'}
                    >
                      {formatAmount(transaction.amount)}
                    </Td>
                    <Td maxW="200px" isTruncated>
                      {transaction.description || '-'}
                    </Td>
                    <Td maxW="150px" isTruncated>
                      {transaction.counterparty || '-'}
                    </Td>
                    <Td>
                      <Select
                        placeholder="选择账户"
                        size="sm"
                        value={selectedTransactions[transaction.rawIndex]?.accountId || ''}
                        onChange={(e) =>
                          setSelectedTransactions({
                            ...selectedTransactions,
                            [transaction.rawIndex]: {
                              ...selectedTransactions[transaction.rawIndex],
                              accountId: e.target.value || null,
                            },
                          })
                        }
                      >
                        {accounts.map((account) => (
                          <option key={account.id} value={account.id}>
                            {account.name}
                          </option>
                        ))}
                      </Select>
                    </Td>
                    <Td>
                      <Select
                        placeholder="选择分类"
                        size="sm"
                        value={selectedTransactions[transaction.rawIndex]?.categoryId || ''}
                        onChange={(e) =>
                          setSelectedTransactions({
                            ...selectedTransactions,
                            [transaction.rawIndex]: {
                              ...selectedTransactions[transaction.rawIndex],
                              categoryId: e.target.value || null,
                            },
                          })
                        }
                      >
                        {categories.map((category) => (
                          <option key={category.id} value={category.id}>
                            {category.name}
                          </option>
                        ))}
                      </Select>
                    </Td>
                    <Td>
                      {transaction.isDuplicate ? (
                        <Badge colorScheme="orange">可能重复</Badge>
                      ) : (
                        <Badge colorScheme="green">正常</Badge>
                      )}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Box>
        </CardBody>
      </Card>

      {/* 操作按钮 */}
      <HStack spacing={4} justify="flex-end">
        <Button
          size="lg"
          height="14"
          variant="outline"
          borderWidth="2px"
          onClick={handleCancel}
        >
          取消
        </Button>
        <Button
          size="lg"
          height="14"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          onClick={handleConfirm}
          isLoading={isLoading}
          loadingText="导入中..."
          isDisabled={Object.keys(selectedTransactions).length === 0}
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          transition="all 0.2s ease"
        >
          确认导入 ({Object.keys(selectedTransactions).length} 条)
        </Button>
      </HStack>
    </VStack>
  );
}
