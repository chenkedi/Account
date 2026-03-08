import {
  Box,
  Button,
  FormControl,
  FormLabel,
  IconButton,
  Input,
  Select,
  Textarea,
  VStack,
  Heading,
  HStack,
  useToast,
  Card,
  CardBody,
  Flex,
  Radio,
  RadioGroup,
  Stack,
} from '@chakra-ui/react';
import { useEffect } from 'react';
import { useNavigate, useSearchParams, useLocation } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import type { TransactionType, TransactionCreateInput } from '../../../data/models/transaction.model';
import {
  useAccounts,
  useAccountActions,
  useCategories,
  useCategoryActions,
  useTransactionActions,
} from '../../../store';
import { formatAmount } from '../../../core/utils/amount.utils';

function SettingsIcon() {
  return (
    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
    </svg>
  );
}

const transactionFormSchema = z.object({
  type: z.enum(['income', 'expense', 'transfer']),
  amount: z.string().refine((val) => !isNaN(parseFloat(val)) && parseFloat(val) > 0, {
    message: '请输入有效的金额',
  }),
  accountId: z.string().uuid('请选择账户'),
  categoryId: z.string().uuid().optional(),
  transactionDate: z.string(),
  note: z.string().optional(),
});

type TransactionFormData = z.infer<typeof transactionFormSchema>;

export function TransactionFormPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const toast = useToast();

  const typeFromUrl = searchParams.get('type') as TransactionType | null;

  const accounts = useAccounts();
  const categories = useCategories();
  const { fetchAccounts } = useAccountActions();
  const { fetchCategoriesByType } = useCategoryActions();
  const { createTransaction } = useTransactionActions();

  const {
    control,
    handleSubmit,
    watch,
    formState: { isSubmitting, errors },
  } = useForm<TransactionFormData>({
    resolver: zodResolver(transactionFormSchema),
    defaultValues: {
      type: typeFromUrl || 'expense',
      amount: '',
      accountId: '',
      categoryId: '',
      transactionDate: new Date().toISOString().split('T')[0],
      note: '',
    },
  });

  const watchedType = watch('type');

  useEffect(() => {
    fetchAccounts();
  }, [fetchAccounts]);

  useEffect(() => {
    if (watchedType === 'income' || watchedType === 'expense') {
      fetchCategoriesByType(watchedType);
    }
  }, [watchedType, fetchCategoriesByType]);

  const onSubmit = async (data: TransactionFormData) => {
    try {
      const input: TransactionCreateInput = {
        type: data.type,
        amount: parseFloat(data.amount),
        account_id: data.accountId,
        category_id: data.categoryId || null,
        transaction_date: new Date(data.transactionDate),
        note: data.note || null,
        currency: 'CNY',
      };

      await createTransaction(input);

      toast({
        title: '交易创建成功',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });

      navigate('/transactions');
    } catch (error) {
      toast({
        title: '创建失败',
        description: (error as Error).message,
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  return (
    <VStack spacing={6} align="stretch">
      <Box>
        <Heading
          size="2xl"
          bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
          bgClip="text"
          letterSpacing="-0.03em"
        >
          添加交易
        </Heading>
      </Box>

      <Card variant="elevated">
        <CardBody p={6}>
          <form onSubmit={handleSubmit(onSubmit)}>
            <VStack spacing={6}>
              {/* 交易类型 */}
              <FormControl>
                <FormLabel fontWeight="700">交易类型</FormLabel>
                <Controller
                  name="type"
                  control={control}
                  render={({ field }) => (
                    <RadioGroup {...field}>
                      <Stack direction="row" spacing={4}>
                        <Radio value="income" colorScheme="green">
                          收入
                        </Radio>
                        <Radio value="expense" colorScheme="red">
                          支出
                        </Radio>
                        <Radio value="transfer" colorScheme="blue">
                          转账
                        </Radio>
                      </Stack>
                    </RadioGroup>
                  )}
                />
              </FormControl>

              {/* 金额 */}
              <FormControl isInvalid={!!errors.amount}>
                <FormLabel fontWeight="700">金额</FormLabel>
                <Controller
                  name="amount"
                  control={control}
                  render={({ field }) => (
                    <Input
                      {...field}
                      type="number"
                      step="0.01"
                      placeholder="0.00"
                      size="lg"
                      fontSize="2xl"
                      fontWeight="800"
                    />
                  )}
                />
              </FormControl>

              {/* 账户 */}
              <FormControl isInvalid={!!errors.accountId}>
                <FormLabel fontWeight="700">账户</FormLabel>
                <Controller
                  name="accountId"
                  control={control}
                  render={({ field }) => (
                    <Select {...field} placeholder="选择账户" size="lg">
                      {accounts.map((account) => (
                        <option key={account.id} value={account.id}>
                          {account.name} - {formatAmount(account.balance)}
                        </option>
                      ))}
                    </Select>
                  )}
                />
              </FormControl>

              {/* 分类 (仅收入/支出显示) */}
              {(watchedType === 'income' || watchedType === 'expense') && (
                <FormControl>
                  <FormLabel fontWeight="700">
                    <Flex justify="space-between" align="center">
                      <span>分类</span>
                      <IconButton
                        aria-label="管理分类"
                        icon={<SettingsIcon />}
                        size="xs"
                        variant="ghost"
                        colorScheme="brand"
                        onClick={() => {
                          navigate('/categories', {
                            state: { from: location.pathname, type: watchedType },
                          });
                        }}
                      />
                    </Flex>
                  </FormLabel>
                  <Controller
                    name="categoryId"
                    control={control}
                    render={({ field }) => (
                      <Select {...field} placeholder="选择分类" size="lg">
                        {categories.map((category) => (
                          <option key={category.id} value={category.id}>
                            {category.name}
                          </option>
                        ))}
                      </Select>
                    )}
                  />
                </FormControl>
              )}

              {/* 日期 */}
              <FormControl>
                <FormLabel fontWeight="700">日期</FormLabel>
                <Controller
                  name="transactionDate"
                  control={control}
                  render={({ field }) => (
                    <Input {...field} type="date" size="lg" />
                  )}
                />
              </FormControl>

              {/* 备注 */}
              <FormControl>
                <FormLabel fontWeight="700">备注</FormLabel>
                <Controller
                  name="note"
                  control={control}
                  render={({ field }) => (
                    <Textarea
                      {...field}
                      value={field.value || ''}
                      placeholder="添加备注..."
                      size="lg"
                      rows={3}
                    />
                  )}
                />
              </FormControl>

              {/* 按钮 */}
              <HStack spacing={4} w="full" pt={4}>
                <Button
                  size="lg"
                  height="14"
                  flex={1}
                  borderRadius="2xl"
                  variant="outline"
                  borderWidth="2px"
                  onClick={() => navigate('/transactions')}
                >
                  取消
                </Button>
                <Button
                  type="submit"
                  size="lg"
                  height="14"
                  flex={1}
                  borderRadius="2xl"
                  bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
                  color="white"
                  isLoading={isSubmitting}
                  loadingText="保存中..."
                  _hover={{
                    bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
                    transform: 'translateY(-2px)',
                    boxShadow: 'lg',
                  }}
                  transition="all 0.2s ease"
                >
                  保存
                </Button>
              </HStack>
            </VStack>
          </form>
        </CardBody>
      </Card>
    </VStack>
  );
}
