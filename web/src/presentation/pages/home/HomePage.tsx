import { Box, Flex, IconButton, Text, VStack, useColorModeValue } from '@chakra-ui/react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const navItems = [
  { path: '/', label: '首页', icon: <HomeIcon /> },
  { path: '/transactions', label: '交易', icon: <TransactionIcon /> },
  { path: '/accounts', label: '账户', icon: <AccountIcon /> },
  { path: '/stats', label: '统计', icon: <StatsIcon /> },
  { path: '/settings', label: '设置', icon: <SettingsIcon /> },
];

function HomeIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  );
}

function TransactionIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <path d="M7 10h10" />
      <path d="M7 14h6" />
    </svg>
  );
}

function StatsIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="18" y1="20" x2="18" y2="10" />
      <line x1="12" y1="20" x2="12" y2="4" />
      <line x1="6" y1="20" x2="6" y2="14" />
    </svg>
  );
}

function AccountIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="7" width="18" height="14" rx="2" ry="2" />
      <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
    </svg>
  );
}

function SettingsIcon() {
  return (
    <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
    </svg>
  );
}

function ImportIcon() {
  return (
    <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="17 8 12 3 7 8" />
      <line x1="12" y1="3" x2="12" y2="15" />
    </svg>
  );
}

export function HomePage() {
  const navigate = useNavigate();
  const location = useLocation();

  const bgColor = useColorModeValue('gray.50', 'gray.900');
  const headerBg = useColorModeValue('rgba(255, 255, 255, 0.85)', 'rgba(15, 23, 42, 0.85)');
  const navBg = useColorModeValue('rgba(255, 255, 255, 0.95)', 'rgba(15, 23, 42, 0.95)');
  const borderColor = useColorModeValue('rgba(99, 102, 241, 0.1)', 'rgba(99, 102, 241, 0.15)');

  return (
    <Flex direction="column" height="100%" bg={bgColor}>
      {/* 头部 */}
      <Box
        as="header"
        position="sticky"
        top={0}
        zIndex={10}
        px={4}
        pt={4}
      >
        <Box
          bg={headerBg}
          backdropFilter="blur(20px)"
          borderWidth="1px"
          borderColor={borderColor}
          borderRadius="2xl"
          px={4}
          py={3}
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          boxShadow="soft"
        >
          {/* Logo */}
          <Flex alignItems="center" gap={3}>
            <Box
              w="10"
              h="10"
              bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
              borderRadius="xl"
              display="flex"
              alignItems="center"
              justifyContent="center"
              color="white"
              fontSize="xl"
              fontWeight="800"
              boxShadow="md"
            >
              ¥
            </Box>
            <Text
              fontSize="xl"
              fontWeight="800"
              bgGradient="linear(135deg, brand.600 0%, brand.700 100%)"
              bgClip="text"
              letterSpacing="-0.02em"
            >
              Account
            </Text>
          </Flex>
          {/* 导入按钮 */}
          <IconButton
            icon={<ImportIcon />}
            aria-label="导入账单"
            variant="soft"
            colorScheme="brand"
            size="md"
            borderRadius="xl"
            onClick={() => navigate('/import')}
            _hover={{
              transform: 'translateY(-2px)',
              boxShadow: 'md',
            }}
            transition="all 0.2s ease"
          />
        </Box>
      </Box>

      {/* 内容区域 */}
      <Box flex={1} overflowY="auto" px={4} pt={4} pb={28}>
        <Outlet />
      </Box>

      {/* 底部导航 */}
      <Box
        as="nav"
        position="fixed"
        bottom={0}
        left={0}
        right={0}
        px={4}
        pb="env(safe-area-inset-bottom)"
        zIndex={20}
      >
        <Box
          bg={navBg}
          backdropFilter="blur(20px)"
          borderWidth="1px"
          borderColor={borderColor}
          borderRadius="2xl"
          py={2}
          px={4}
          mb={4}
          mx="auto"
          maxW="md"
          boxShadow="card"
        >
          <Flex justify="space-around" alignItems="center">
            {navItems.map((item) => {
              const isActive = location.pathname === item.path;
              return (
                <VStack
                  key={item.path}
                  spacing={1}
                  cursor="pointer"
                  onClick={() => navigate(item.path)}
                  py={2}
                  px={4}
                  borderRadius="xl"
                  transition="all 0.2s ease"
                  _hover={{
                    bg: isActive ? 'transparent' : 'gray.100',
                  }}
                  position="relative"
                >
                  <Box
                    fontSize="24px"
                    color={isActive ? 'brand.600' : 'gray.400'}
                    transition="all 0.2s ease"
                    transform={isActive ? 'translateY(-2px) scale(1.1)' : 'none'}
                    filter={isActive ? 'drop-shadow(0 2px 8px rgba(99, 102, 241, 0.4))' : 'none'}
                  >
                    {item.icon}
                  </Box>
                  <Text
                    fontSize="11px"
                    fontWeight={isActive ? '700' : '600'}
                    color={isActive ? 'brand.600' : 'gray.500'}
                    transition="all 0.2s ease"
                  >
                    {item.label}
                  </Text>
                </VStack>
              );
            })}
          </Flex>
        </Box>
      </Box>
    </Flex>
  );
}
