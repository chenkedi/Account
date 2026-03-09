import { useEffect, useState } from 'react';
import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  Flex,
  HStack,
  Badge,
  Spinner,
  SimpleGrid,
  IconButton,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
} from '@chakra-ui/react';
import {
  useAccounts,
  useAccountsState,
  useAccountActions,
} from '../../../store';
import type { Account, AccountType, AccountCreateInput, AccountUpdateInput } from '../../../data/models/account.model';
import { AccountFormModal } from './AccountFormModal';
import { DeleteAccountModal } from './DeleteAccountModal';

function PlusIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

function MoreIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="1" />
      <circle cx="19" cy="12" r="1" />
      <circle cx="5" cy="12" r="1" />
    </svg>
  );
}

function EditIcon() {
  return (
    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
    </svg>
  );
}

function DeleteIcon() {
  return (
    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="3 6 5 6 21 6" />
      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
  );
}

function WalletIcon() {
  return (
    <svg viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="7" width="20" height="14" rx="2" />
      <path d="M16 12h4" />
      <path d="M2 7V5a2 2 0 0 1 2-2h14" />
    </svg>
  );
}

function getAccountTypeLabel(type: AccountType): string {
  const labels: Record<AccountType, string> = {
    cash: '现金',
    bank: '银行卡',
    credit_card: '信用卡',
    investment: '投资账户',
    other: '其他',
  };
  return labels[type] || type;
}

function getAccountTypeColor(type: AccountType): string {
  const colors: Record<AccountType, string> = {
    cash: 'green',
    bank: 'blue',
    credit_card: 'purple',
    investment: 'orange',
    other: 'gray',
  };
  return colors[type] || 'gray';
}

export function AccountsPage() {
  const accounts = useAccounts();
  const { isLoading, error } = useAccountsState();
  const { fetchAccounts, createAccount, updateAccount, deleteAccount } = useAccountActions();

  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);

  useEffect(() => {
    fetchAccounts();
  }, [fetchAccounts]);

  const handleAddAccount = () => {
    setSelectedAccount(null);
    setIsFormModalOpen(true);
  };

  const handleEditAccount = (account: Account) => {
    setSelectedAccount(account);
    setIsFormModalOpen(true);
  };

  const handleDeleteAccount = (account: Account) => {
    setSelectedAccount(account);
    setIsDeleteModalOpen(true);
  };

  const handleFormSubmit = async (data: AccountCreateInput | AccountUpdateInput) => {
    if (selectedAccount) {
      await updateAccount(selectedAccount.id, data);
    } else {
      await createAccount(data as AccountCreateInput);
    }
  };

  return (
    <VStack spacing={6} align="stretch">
      {/* Header */}
      <Flex justify="space-between" align="center">
        <Box>
          <Heading
            size="2xl"
            bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
            bgClip="text"
            letterSpacing="-0.03em"
          >
            账户管理
          </Heading>
          <Text color="gray.500" fontWeight="500" mt={1}>
            管理您的支付账户
          </Text>
        </Box>
        <Button
          size="lg"
          borderRadius="2xl"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          leftIcon={<PlusIcon />}
          fontWeight="700"
          onClick={handleAddAccount}
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          transition="all 0.2s ease"
        >
          添加账户
        </Button>
      </Flex>

      {/* Error */}
      {error && (
        <Card bg="red.50" borderColor="red.200" borderWidth="1px">
          <CardBody>
            <Text color="red.600">{error}</Text>
          </CardBody>
        </Card>
      )}

      {/* Loading */}
      {isLoading && accounts.length === 0 && (
        <Flex justify="center" py={20}>
          <Spinner size="xl" color="brand.500" />
        </Flex>
      )}

      {/* Empty State */}
      {!isLoading && accounts.length === 0 && (
        <Card variant="elevated">
          <CardBody py={16}>
            <VStack spacing={6}>
              <Box color="gray.300">
                <WalletIcon />
              </Box>
              <VStack spacing={2}>
                <Text fontSize="xl" fontWeight="700" color="gray.700">
                  暂无账户
                </Text>
                <Text color="gray.500" textAlign="center" maxW="sm">
                  添加您的银行账户、电子钱包等开始记账
                </Text>
              </VStack>
            </VStack>
          </CardBody>
        </Card>
      )}

      {/* Account List */}
      {accounts.length > 0 && (
        <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
          {accounts.map((account) => (
            <Card
              key={account.id}
              variant="elevated"
              _hover={{
                transform: 'translateY(-4px)',
                boxShadow: 'lg',
              }}
              transition="all 0.2s ease"
            >
              <CardBody>
                <VStack align="stretch" spacing={3}>
                  <Flex justify="space-between" align="start">
                    <HStack spacing={2}>
                      <Badge
                        colorScheme={getAccountTypeColor(account.type)}
                        borderRadius="full"
                        px={3}
                        py={1}
                      >
                        {getAccountTypeLabel(account.type)}
                      </Badge>
                      <Text fontSize="sm" color="gray.500">
                        {account.currency}
                      </Text>
                    </HStack>
                    <Menu>
                      <MenuButton
                        as={IconButton}
                        icon={<MoreIcon />}
                        variant="ghost"
                        size="sm"
                        aria-label="账户操作"
                      />
                      <MenuList>
                        <MenuItem
                          icon={<EditIcon />}
                          onClick={() => handleEditAccount(account)}
                        >
                          编辑
                        </MenuItem>
                        <MenuItem
                          icon={<DeleteIcon />}
                          color="red.500"
                          onClick={() => handleDeleteAccount(account)}
                        >
                          删除
                        </MenuItem>
                      </MenuList>
                    </Menu>
                  </Flex>
                  <VStack align="start" spacing={0}>
                    <Text fontSize="lg" fontWeight="bold" noOfLines={1}>
                      {account.name}
                    </Text>
                    {account.tail_number && (
                      <Text fontSize="sm" color="gray.500" fontWeight="500">
                        尾号 {account.tail_number}
                      </Text>
                    )}
                  </VStack>
                  <Text
                    fontSize="2xl"
                    fontWeight="bold"
                    color={account.balance >= 0 ? 'green.500' : 'red.500'}
                  >
                    {account.balance.toLocaleString('zh-CN', {
                      style: 'currency',
                      currency: account.currency,
                    })}
                  </Text>
                </VStack>
              </CardBody>
            </Card>
          ))}
        </SimpleGrid>
      )}

      <AccountFormModal
        isOpen={isFormModalOpen}
        onClose={() => setIsFormModalOpen(false)}
        account={selectedAccount}
        onSubmit={handleFormSubmit}
        isLoading={isLoading}
      />

      <DeleteAccountModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        account={selectedAccount}
        onDelete={deleteAccount}
        isLoading={isLoading}
      />

    </VStack>
  );
}
