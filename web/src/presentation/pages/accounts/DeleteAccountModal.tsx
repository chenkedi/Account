import {
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Text,
  VStack,
  useToast,
  Alert,
  AlertIcon,
} from '@chakra-ui/react';
import type { Account } from '../../../data/models/account.model';

interface DeleteAccountModalProps {
  isOpen: boolean;
  onClose: () => void;
  account: Account | null;
  onDelete: (id: string) => Promise<void>;
  isLoading?: boolean;
}

export function DeleteAccountModal({
  isOpen,
  onClose,
  account,
  onDelete,
  isLoading = false,
}: DeleteAccountModalProps) {
  const toast = useToast();

  const handleDelete = async () => {
    if (!account) return;

    try {
      await onDelete(account.id);
      toast({
        title: '账户已删除',
        status: 'success',
        duration: 3000,
      });
      onClose();
    } catch (error) {
      toast({
        title: '删除失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  if (!account) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="md">
      <ModalOverlay />
      <ModalContent borderRadius="2xl">
        <ModalHeader>删除账户</ModalHeader>
        <ModalBody>
          <VStack spacing={4} align="stretch">
            <Alert status="warning" borderRadius="xl">
              <AlertIcon />
              <Text>此操作不可撤销，删除后账户相关数据将无法恢复。</Text>
            </Alert>
            <Text>
              确定要删除账户 <Text as="span" fontWeight="bold">{account.name}</Text> 吗？
            </Text>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button onClick={onClose} variant="outline" borderRadius="xl" borderWidth="2px">
            取消
          </Button>
          <Button
            onClick={handleDelete}
            colorScheme="red"
            borderRadius="xl"
            isLoading={isLoading}
          >
            删除
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
