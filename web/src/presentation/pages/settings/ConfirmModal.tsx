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
  Alert,
  AlertIcon,
} from '@chakra-ui/react';

interface ConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  isDangerous?: boolean;
  onConfirm: () => Promise<void>;
  isLoading?: boolean;
  warning?: string;
}

export function ConfirmModal({
  isOpen,
  onClose,
  title,
  message,
  confirmText = '确认',
  cancelText = '取消',
  isDangerous = false,
  onConfirm,
  isLoading = false,
  warning,
}: ConfirmModalProps) {
  const handleConfirm = async () => {
    try {
      await onConfirm();
      onClose();
    } catch (error) {
      // Error is handled by caller
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="md">
      <ModalOverlay />
      <ModalContent borderRadius="2xl">
        <ModalHeader>{title}</ModalHeader>
        <ModalBody>
          <VStack spacing={4} align="stretch">
            {warning && (
              <Alert status="warning" borderRadius="xl">
                <AlertIcon />
                <Text>{warning}</Text>
              </Alert>
            )}
            <Text>{message}</Text>
          </VStack>
        </ModalBody>

        <ModalFooter gap={3}>
          <Button onClick={onClose} variant="outline" borderRadius="xl" borderWidth="2px">
            {cancelText}
          </Button>
          <Button
            onClick={handleConfirm}
            colorScheme={isDangerous ? 'red' : 'brand'}
            borderRadius="xl"
            isLoading={isLoading}
            bgGradient={!isDangerous ? 'linear(135deg, brand.500 0%, brand.600 100%)' : undefined}
            color={!isDangerous ? 'white' : undefined}
          >
            {confirmText}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
