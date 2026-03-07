import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import type { Account, AccountType, AccountCreateInput, AccountUpdateInput } from '../../../data/models/account.model';

const accountSchema = z.object({
  name: z.string().min(1, '账户名称不能为空'),
  type: z.enum(['cash', 'bank', 'credit_card', 'investment', 'other'] as const),
  balance: z.string().refine((val) => !isNaN(parseFloat(val)) && parseFloat(val) >= 0, '余额必须是有效数字'),
  currency: z.string().min(1, '货币不能为空'),
});

type AccountFormValues = z.infer<typeof accountSchema>;

interface AccountFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  account?: Account | null;
  onSubmit: (data: AccountCreateInput | AccountUpdateInput) => Promise<void>;
  isLoading?: boolean;
}

const accountTypeOptions: { value: AccountType; label: string }[] = [
  { value: 'cash', label: '现金' },
  { value: 'bank', label: '银行卡' },
  { value: 'credit_card', label: '信用卡' },
  { value: 'investment', label: '投资账户' },
  { value: 'other', label: '其他' },
];

const currencyOptions = [
  { value: 'CNY', label: '人民币 (¥)' },
  { value: 'USD', label: '美元 ($)' },
  { value: 'EUR', label: '欧元 (€)' },
];

export function AccountFormModal({
  isOpen,
  onClose,
  account,
  onSubmit,
  isLoading = false,
}: AccountFormModalProps) {
  const toast = useToast();
  const isEdit = !!account;

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AccountFormValues>({
    resolver: zodResolver(accountSchema),
    defaultValues: {
      name: account?.name || '',
      type: account?.type || 'cash',
      balance: account?.balance.toString() || '0',
      currency: account?.currency || 'CNY',
    },
  });

  const handleFormSubmit = async (values: AccountFormValues) => {
    try {
      const data: AccountCreateInput | AccountUpdateInput = {
        name: values.name,
        type: values.type,
        balance: parseFloat(values.balance),
        currency: values.currency,
      };
      await onSubmit(data);
      toast({
        title: isEdit ? '账户已更新' : '账户已创建',
        status: 'success',
        duration: 3000,
      });
      reset();
      onClose();
    } catch (error) {
      toast({
        title: '操作失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="md">
      <ModalOverlay />
      <ModalContent borderRadius="2xl">
        <ModalHeader>
          {isEdit ? '编辑账户' : '添加账户'}
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <VStack spacing={4}>
            <FormControl isInvalid={!!errors.name}>
              <FormLabel fontWeight="600">账户名称</FormLabel>
              <Controller
                name="name"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="例如：工商银行储蓄卡" borderRadius="xl" />
                )}
              />
              {errors.name && (
                <Box color="red.500" fontSize="sm" mt={1}>
                  {errors.name.message}
                </Box>
              )}
            </FormControl>

            <FormControl isInvalid={!!errors.type}>
              <FormLabel fontWeight="600">账户类型</FormLabel>
              <Controller
                name="type"
                control={control}
                render={({ field }) => (
                  <Select {...field} borderRadius="xl">
                    {accountTypeOptions.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </Select>
                )}
              />
              {errors.type && (
                <Box color="red.500" fontSize="sm" mt={1}>
                  {errors.type.message}
                </Box>
              )}
            </FormControl>

            <FormControl isInvalid={!!errors.balance}>
              <FormLabel fontWeight="600">初始余额</FormLabel>
              <Controller
                name="balance"
                control={control}
                render={({ field }) => (
                  <Input {...field} type="number" step="0.01" placeholder="0.00" borderRadius="xl" />
                )}
              />
              {errors.balance && (
                <Box color="red.500" fontSize="sm" mt={1}>
                  {errors.balance.message}
                </Box>
              )}
            </FormControl>

            <FormControl isInvalid={!!errors.currency}>
              <FormLabel fontWeight="600">货币</FormLabel>
              <Controller
                name="currency"
                control={control}
                render={({ field }) => (
                  <Select {...field} borderRadius="xl">
                    {currencyOptions.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </Select>
                )}
              />
              {errors.currency && (
                <Box color="red.500" fontSize="sm" mt={1}>
                  {errors.currency.message}
                </Box>
              )}
            </FormControl>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button onClick={onClose} variant="outline" borderRadius="xl" borderWidth="2px">
            取消
          </Button>
          <Button
            onClick={handleSubmit(handleFormSubmit)}
            colorScheme="brand"
            borderRadius="xl"
            isLoading={isLoading}
            bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
            color="white"
          >
            {isEdit ? '保存修改' : '创建账户'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
