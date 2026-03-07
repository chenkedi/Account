import { useState } from 'react';
import {
  Box,
  Heading,
  Text,
  VStack,
  Button,
  Card,
  CardBody,
  Flex,
  useToast,
} from '@chakra-ui/react';
import { useAuth, useAuthActions } from '../../../store';
import { ConfirmModal } from './ConfirmModal';

function ChevronRightIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
    </svg>
  );
}

function UserIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

function LockIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0 1 10 0v4" />
    </svg>
  );
}

function DownloadIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="7 10 12 15 17 10" />
      <line x1="12" y1="15" x2="12" y2="3" />
    </svg>
  );
}

function TrashIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="3 6 5 6 21 6" />
      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
  );
}

function InfoIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="16" x2="12" y2="12" />
      <line x1="12" y1="8" x2="12.01" y2="8" />
    </svg>
  );
}

function LogoutIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  );
}

interface SettingItemProps {
  icon: React.ReactNode;
  label: string;
  subLabel?: string;
  onClick?: () => void;
  colorScheme?: 'brand' | 'gray' | 'red';
  showChevron?: boolean;
}

function SettingItem({ icon, label, subLabel, onClick, colorScheme = 'gray', showChevron = true }: SettingItemProps) {
  const colorMap = {
    brand: { bg: 'brand.50', color: 'brand.600' },
    gray: { bg: 'gray.50', color: 'gray.600' },
    red: { bg: 'expense.50', color: 'expense.600' },
  };

  const colors = colorMap[colorScheme];

  return (
    <Flex
      align="center"
      gap={4}
      py={4}
      px={2}
      cursor="pointer"
      onClick={onClick}
      borderRadius="xl"
      _hover={{ bg: 'gray.50' }}
      transition="all 0.2s ease"
    >
      <Box
        w={12}
        h={12}
        borderRadius="xl"
        bg={colors.bg}
        color={colors.color}
        display="flex"
        alignItems="center"
        justifyContent="center"
      >
        {icon}
      </Box>
      <Flex flex={1} direction="column">
        <Text fontWeight="700" color="gray.800">
          {label}
        </Text>
        {subLabel && (
          <Text fontSize="sm" color="gray.500" fontWeight="500">
            {subLabel}
          </Text>
        )}
      </Flex>
      {showChevron && (
        <Box color="gray.400">
          <ChevronRightIcon />
        </Box>
      )}
    </Flex>
  );
}

export function SettingsPage() {
  const { user } = useAuth();
  const { logout } = useAuthActions();
  const toast = useToast();

  const [logoutModalOpen, setLogoutModalOpen] = useState(false);
  const [clearDataModalOpen, setClearDataModalOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await logout();
      toast({
        title: '已退出登录',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: '退出失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  const handleClearData = async () => {
    try {
      // Clear local storage
      localStorage.clear();

      toast({
        title: '本地数据已清除',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: '清除失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  return (
    <>
      <VStack spacing={6} align="stretch">
        {/* 头部 */}
        <Box>
          <Heading
            size="2xl"
            bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
            bgClip="text"
            letterSpacing="-0.03em"
          >
            设置
          </Heading>
          <Text color="gray.500" fontWeight="500" mt={1}>
            管理您的账户和偏好
          </Text>
        </Box>

        {/* 账户信息 */}
        <Card variant="elevated">
          <CardBody p={2}>
            <SettingItem
              icon={<UserIcon />}
              label="账户信息"
              subLabel={user?.email || '未登录'}
              colorScheme="brand"
            />
            <SettingItem
              icon={<LockIcon />}
              label="修改密码"
              colorScheme="gray"
            />
          </CardBody>
        </Card>

        {/* 数据管理 */}
        <Card variant="elevated">
          <CardBody p={2}>
            <Text fontWeight="700" color="gray.800" px={2} pt={2} pb={1}>
              数据管理
            </Text>
            <SettingItem
              icon={<DownloadIcon />}
              label="导出所有数据"
              colorScheme="gray"
            />
            <SettingItem
              icon={<TrashIcon />}
              label="清除本地数据"
              colorScheme="red"
              showChevron={false}
              onClick={() => setClearDataModalOpen(true)}
            />
          </CardBody>
        </Card>

        {/* 关于 */}
        <Card variant="elevated">
          <CardBody p={2}>
            <Text fontWeight="700" color="gray.800" px={2} pt={2} pb={1}>
              关于
            </Text>
            <SettingItem
              icon={<InfoIcon />}
              label="Account"
              subLabel="v0.0.1 · 个人财务管理应用"
              colorScheme="brand"
              showChevron={false}
            />
          </CardBody>
        </Card>

        {/* 退出登录 */}
        <Button
          onClick={() => setLogoutModalOpen(true)}
          size="lg"
          height="14"
          borderRadius="2xl"
          variant="outline"
          borderWidth="2px"
          borderColor="expense.200"
          color="expense.600"
          leftIcon={<LogoutIcon />}
          fontWeight="700"
          _hover={{
            borderColor: 'expense.300',
            bg: 'expense.50',
            transform: 'translateY(-2px)',
            boxShadow: 'md',
          }}
          transition="all 0.2s ease"
        >
          退出登录
        </Button>
      </VStack>

      <ConfirmModal
        isOpen={logoutModalOpen}
        onClose={() => setLogoutModalOpen(false)}
        title="退出登录"
        message="确定要退出登录吗？"
        confirmText="退出"
        isDangerous={true}
        onConfirm={handleLogout}
      />

      <ConfirmModal
        isOpen={clearDataModalOpen}
        onClose={() => setClearDataModalOpen(false)}
        title="清除本地数据"
        message="确定要清除所有本地数据吗？此操作不可撤销。"
        confirmText="清除"
        isDangerous={true}
        onConfirm={handleClearData}
        warning="清除后，您的本地数据将无法恢复。"
      />
    </>
  );
}
