import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  VStack,
  SimpleGrid,
  Box,
  Text,
  Radio,
  RadioGroup,
  Stack,
  FormErrorMessage,
} from '@chakra-ui/react';
import { useEffect, useMemo } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import type {
  Category,
  CategoryType,
  CategoryCreateInput,
  CategoryUpdateInput,
} from '../../../data/models/category.model';

const categoryFormSchema = z.object({
  name: z.string().min(1, '分类名称不能为空').max(50, '分类名称不能超过50个字符'),
  type: z.enum(['income', 'expense']),
  parent_id: z.string().optional(),
  icon: z.string().min(1, '请选择一个图标'),
});

type CategoryFormData = z.infer<typeof categoryFormSchema>;

interface CategoryFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  category: Category | null;
  defaultType?: CategoryType;
  allCategories: Category[];
  onSubmit: (data: CategoryCreateInput | CategoryUpdateInput) => Promise<void>;
  isLoading: boolean;
}

const PRESET_ICONS = [
  { value: 'shopping', label: '购物', emoji: '🛍️' },
  { value: 'food', label: '餐饮', emoji: '🍔' },
  { value: 'transport', label: '交通', emoji: '🚗' },
  { value: 'entertainment', label: '娱乐', emoji: '🎮' },
  { value: 'health', label: '医疗', emoji: '🏥' },
  { value: 'education', label: '教育', emoji: '📚' },
  { value: 'housing', label: '住房', emoji: '🏠' },
  { value: 'utilities', label: '水电', emoji: '💡' },
  { value: 'salary', label: '工资', emoji: '💰' },
  { value: 'bonus', label: '奖金', emoji: '🧧' },
  { value: 'investment', label: '投资', emoji: '📈' },
  { value: 'gift', label: '礼金', emoji: '🎁' },
  { value: 'other', label: '其他', emoji: '📦' },
];

export function CategoryFormModal({
  isOpen,
  onClose,
  category,
  defaultType = 'expense',
  allCategories,
  onSubmit,
}: CategoryFormModalProps) {
  const isEditMode = !!category;

  const {
    control,
    handleSubmit,
    watch,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CategoryFormData>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues: {
      name: '',
      type: defaultType,
      parent_id: undefined,
      icon: 'other',
    },
  });

  const watchedType = watch('type');

  // Reset form when modal opens or category changes
  useEffect(() => {
    if (isOpen) {
      if (category) {
        reset({
          name: category.name,
          type: category.type,
          parent_id: category.parent_id || undefined,
          icon: category.icon || 'other',
        });
      } else {
        reset({
          name: '',
          type: defaultType,
          parent_id: undefined,
          icon: 'other',
        });
      }
    }
  }, [isOpen, category, defaultType, reset]);

  // Get available parent categories (same type, not self or descendants)
  const availableParentCategories = useMemo(() => {
    if (!watchedType) return [];

    return allCategories.filter((c) => {
      // Same type only
      if (c.type !== watchedType) return false;
      // In edit mode, exclude self and descendants
      if (category) {
        if (c.id === category.id) return false;
        // Check if c is a descendant of category
        let current = allCategories.find((cat) => cat.id === c.parent_id);
        while (current) {
          if (current.id === category.id) return false;
          current = allCategories.find((cat) => cat.id === current!.parent_id);
        }
      }
      return true;
    });
  }, [allCategories, watchedType, category]);

  const onFormSubmit = async (data: CategoryFormData) => {
    const input: CategoryCreateInput | CategoryUpdateInput = {
      name: data.name,
      type: data.type,
      parent_id: data.parent_id,
      icon: data.icon,
    };
    await onSubmit(input);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="xl">
      <ModalOverlay />
      <ModalContent borderRadius="2xl">
        <ModalHeader fontSize="xl" fontWeight="700">
          {isEditMode ? '编辑分类' : '添加分类'}
        </ModalHeader>
        <ModalBody>
          <VStack spacing={6} align="stretch">
            {/* Name */}
            <FormControl isInvalid={!!errors.name}>
              <FormLabel fontWeight="600">分类名称</FormLabel>
              <Controller
                name="name"
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="请输入分类名称" size="lg" />
                )}
              />
              <FormErrorMessage>{errors.name?.message}</FormErrorMessage>
            </FormControl>

            {/* Type */}
            <FormControl isInvalid={!!errors.type}>
              <FormLabel fontWeight="600">分类类型</FormLabel>
              <Controller
                name="type"
                control={control}
                render={({ field }) => (
                  <RadioGroup {...field}>
                    <Stack direction="row" spacing={6}>
                      <Radio value="expense" colorScheme="red">支出</Radio>
                      <Radio value="income" colorScheme="green">收入</Radio>
                    </Stack>
                  </RadioGroup>
                )}
              />
            </FormControl>

            {/* Parent Category */}
            <FormControl>
              <FormLabel fontWeight="600">父分类（可选）</FormLabel>
              <Controller
                name="parent_id"
                control={control}
                render={({ field }) => (
                  <Select {...field} placeholder="选择父分类（可选）" size="lg">
                    <option value="">无</option>
                    {availableParentCategories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </Select>
                )}
              />
              <Text fontSize="sm" color="gray.500" mt={1}>
                选择父分类可将此分类作为子分类
              </Text>
            </FormControl>

            {/* Icon Selection */}
            <FormControl isInvalid={!!errors.icon}>
              <FormLabel fontWeight="600">图标</FormLabel>
              <Controller
                name="icon"
                control={control}
                render={({ field }) => (
                  <SimpleGrid columns={7} spacing={2}>
                    {PRESET_ICONS.map((icon) => (
                      <Box
                        key={icon.value}
                        onClick={() => field.onChange(icon.value)}
                        cursor="pointer"
                        p={2}
                        borderRadius="xl"
                        borderWidth="2px"
                        borderColor={field.value === icon.value ? 'brand.500' : 'transparent'}
                        bg={field.value === icon.value ? 'brand.50' : 'gray.50'}
                        _hover={{
                          bg: field.value === icon.value ? 'brand.50' : 'gray.100',
                        }}
                        transition="all 0.2s"
                        title={icon.label}
                      >
                        <Text fontSize="2xl" textAlign="center">
                          {icon.emoji}
                        </Text>
                      </Box>
                    ))}
                  </SimpleGrid>
                )}
              />
              <FormErrorMessage>{errors.icon?.message}</FormErrorMessage>
            </FormControl>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button
            variant="outline"
            borderRadius="xl"
            onClick={onClose}
            isDisabled={isSubmitting}
          >
            取消
          </Button>
          <Button
            borderRadius="xl"
            bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
            color="white"
            onClick={handleSubmit(onFormSubmit)}
            isLoading={isSubmitting}
            loadingText="保存中..."
            _hover={{
              bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            }}
          >
            保存
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
