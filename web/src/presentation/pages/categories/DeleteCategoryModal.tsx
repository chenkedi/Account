import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Text,
  VStack,
  Alert,
  AlertIcon,
  Box,
  Flex,
} from '@chakra-ui/react';
import { useMemo } from 'react';
import type { Category } from '../../../data/models/category.model';

interface DeleteCategoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  category: Category | null;
  allCategories: Category[];
  onDelete: () => Promise<void>;
  isLoading: boolean;
}

function TagIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
      <line x1="7" y1="7" x2="7.01" y2="7" />
    </svg>
  );
}

export function DeleteCategoryModal({
  isOpen,
  onClose,
  category,
  allCategories,
  onDelete,
  isLoading,
}: DeleteCategoryModalProps) {
  // Check if category has children
  const childCategories = useMemo(() => {
    if (!category) return [];
    return allCategories.filter((c) => c.parent_id === category.id);
  }, [category, allCategories]);

  const hasChildren = childCategories.length > 0;

  if (!category) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="lg">
      <ModalOverlay />
      <ModalContent borderRadius="2xl">
        <ModalHeader fontSize="xl" fontWeight="700" color="red.600">
          删除分类
        </ModalHeader>

        <ModalBody>
          <VStack spacing={4} align="stretch">
            <Text>
              确定要删除分类 <Text as="strong">"{category.name}"</Text> 吗？此操作不可撤销。
            </Text>

            {hasChildren && (
              <Alert status="warning" borderRadius="xl">
                <AlertIcon />
                <Box flex={1}>
                  <Text fontWeight="600">此分类包含子分类</Text>
                  <Text fontSize="sm">
                    删除此分类将同时删除以下 {childCategories.length} 个子分类：
                  </Text>
                  <VStack align="stretch" mt={2} spacing={1}>
                    {childCategories.map((child) => (
                      <Flex key={child.id} align="center" gap={2}>
                        <Box color="gray.400">
                          <TagIcon />
                        </Box>
                        <Text fontSize="sm" fontWeight="500">{child.name}</Text>
                      </Flex>
                    ))}
                  </VStack>
                </Box>
              </Alert>
            )}

            <Alert status="error" borderRadius="xl">
              <AlertIcon />
              <Text fontSize="sm">
                如果该分类已被用于交易记录，删除后这些交易将不再关联此分类。
              </Text>
            </Alert>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button
            variant="outline"
            borderRadius="xl"
            onClick={onClose}
            isDisabled={isLoading}
          >
            取消
          </Button>
          <Button
            borderRadius="xl"
            bg="red.500"
            color="white"
            onClick={onDelete}
            isLoading={isLoading}
            loadingText="删除中..."
            _hover={{
              bg: 'red.600',
            }}
          >
            确认删除
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
