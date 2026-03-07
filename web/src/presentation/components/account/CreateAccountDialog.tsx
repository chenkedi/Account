import { useState } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  VStack,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
} from '@chakra-ui/react';
import type { AccountType, AccountCreateInput } from '../../../data/models/account.model';

interface CreateAccountDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: AccountCreateInput) => Promise<void>;
  isLoading?: boolean;
}

const ACCOUNT_TYPES: { value: AccountType; label: string }[] = [
  { value: 'cash', label: '现金' },
  { value: 'bank', label: '银行卡' },
  { value: 'credit_card', label: '信用卡' },
  { value: 'investment', label: '投资账户' },
  { value: 'other', label: '其他' },
];

export function CreateAccountDialog({
  isOpen,
  onClose,
  onSubmit,
  isLoading = false,
}: CreateAccountDialogProps) {
  const [formData, setFormData] = useState<AccountCreateInput>({
    name: '',
    type: 'cash',
    currency: 'CNY',
    balance: 0,
  });

  const [errors, setErrors] = useState<Partial<Record<keyof AccountCreateInput, string>>>({});

  const handleChange = (field: keyof AccountCreateInput, value: string | number) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof AccountCreateInput, string>> = {};

    if (!formData.name.trim()) {
      newErrors.name = '请输入账户名称';
    }

    if (formData.balance === undefined || formData.balance === null) {
      newErrors.balance = '请输入初始余额';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validate()) return;

    try {
      await onSubmit(formData);
      // 重置表单
      setFormData({
        name: '',
        type: 'cash',
        currency: 'CNY',
        balance: 0,
      });
      onClose();
    } catch (error) {
      // 错误已由 store 处理
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="lg">
      <ModalOverlay backdropFilter="blur(4px)" />
      <ModalContent borderRadius="2xl">
        <ModalHeader fontSize="xl" fontWeight="bold">
          新增账户
        </ModalHeader>
        <ModalCloseButton />

        <ModalBody>
          <VStack spacing={4}>
            <FormControl isRequired isInvalid={!!errors.name}>
              <FormLabel>账户名称</FormLabel>
              <Input
                placeholder="例如：我的银行卡"
                value={formData.name}
                onChange={(e) => handleChange('name', e.target.value)}
              />
            </FormControl>

            <FormControl isRequired>
              <FormLabel>账户类型</FormLabel>
              <Select
                value={formData.type}
                onChange={(e) => handleChange('type', e.target.value as AccountType)}
              >
                {ACCOUNT_TYPES.map((type) => (
                  <option key={type.value} value={type.value}>
                    {type.label}
                  </option>
                ))}
              </Select>
            </FormControl>

            <FormControl isRequired isInvalid={!!errors.balance}>
              <FormLabel>初始余额</FormLabel>
              <NumberInput
                value={formData.balance}
                onChange={(valueString) => handleChange('balance', parseFloat(valueString) || 0)}
                min={-999999999}
                max={999999999}
                step={0.01}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
            </FormControl>

            <FormControl isRequired>
              <FormLabel>币种</FormLabel>
              <Select
                value={formData.currency}
                onChange={(e) => handleChange('currency', e.target.value)}
              >
                <option value="CNY">人民币 (CNY)</option>
                <option value="USD">美元 (USD)</option>
                <option value="EUR">欧元 (EUR)</option>
                <option value="JPY">日元 (JPY)</option>
                <option value="HKD">港币 (HKD)</option>
              </Select>
            </FormControl>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button variant="ghost" onClick={onClose}>
            取消
          </Button>
          <Button
            colorScheme="blue"
            onClick={handleSubmit}
            isLoading={isLoading}
            loadingText="保存中..."
          >
            保存
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
