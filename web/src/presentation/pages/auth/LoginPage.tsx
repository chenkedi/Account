import { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Flex,
  FormControl,
  FormLabel,
  Heading,
  Input,
  InputGroup,
  InputRightElement,
  Link,
  Text,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useAuthActions } from '../../../store';

// Logo 组件
function AppLogo() {
  return (
    <Box mb={8}>
      <Box
        w="80px"
        h="80px"
        bgGradient="linear(135deg, brand.500 0%, brand.600 50%, brand.700 100%)"
        borderRadius="2xl"
        display="flex"
        alignItems="center"
        justifyContent="center"
        boxShadow="2xl"
        position="relative"
        overflow="hidden"
      >
        {/* 装饰性圆环 */}
        <Box
          position="absolute"
          w="120px"
          h="120px"
          border="2px solid"
          borderColor="whiteAlpha.300"
          borderRadius="full"
          top="-20px"
          right="-20px"
        />
        <Box
          position="absolute"
          w="80px"
          h="80px"
          border="2px solid"
          borderColor="whiteAlpha.200"
          borderRadius="full"
          bottom="-10px"
          left="-10px"
        />
        {/* 图标 */}
        <Box color="white" fontSize="40px" fontWeight="800" zIndex={1}>
          ¥
        </Box>
      </Box>
    </Box>
  );
}

// 装饰性背景元素
function BackgroundDecoration() {
  return (
    <Box
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      zIndex={-1}
      bgGradient="linear(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%)"
      overflow="hidden"
    >
      {/* 装饰性圆形 */}
      <Box
        position="absolute"
        w="500px"
        h="500px"
        bg="whiteAlpha.100"
        borderRadius="full"
        top="-200px"
        right="-100px"
        filter="blur(40px)"
      />
      <Box
        position="absolute"
        w="400px"
        h="400px"
        bg="whiteAlpha.100"
        borderRadius="full"
        bottom="-150px"
        left="-100px"
        filter="blur(40px)"
      />
      {/* 网格装饰 */}
      <Box
        position="absolute"
        inset={0}
        opacity={0.03}
        backgroundImage="linear-gradient(rgba(255,255,255,0.5) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.5) 1px, transparent 1px)"
        backgroundSize="50px 50px"
      />
    </Box>
  );
}

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const toast = useToast();
  const { login } = useAuthActions();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      toast({
        title: '请填写完整信息',
        status: 'warning',
        duration: 3000,
        position: 'top',
      });
      return;
    }

    setIsLoading(true);
    try {
      await login(email, password);
      toast({
        title: '登录成功',
        description: '欢迎回来！',
        status: 'success',
        duration: 3000,
        position: 'top',
        containerStyle: {
          borderRadius: '16px',
        },
      });
      navigate('/');
    } catch (error) {
      toast({
        title: '登录失败',
        description: (error as Error).message || '请检查您的邮箱和密码',
        status: 'error',
        duration: 5000,
        position: 'top',
        containerStyle: {
          borderRadius: '16px',
        },
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <BackgroundDecoration />
      <Box minH="100vh" display="flex" alignItems="center" justifyContent="center" p={4}>
        <Container maxW="md" w="full">
          <Box
            as="form"
            onSubmit={handleSubmit}
            bg="white"
            borderRadius="3xl"
            p={{ base: 8, md: 12 }}
            boxShadow="2xl"
            position="relative"
            overflow="hidden"
          >
            {/* 顶部装饰条 */}
            <Box
              position="absolute"
              top={0}
              left={0}
              right={0}
              h="4px"
              bgGradient="linear(90deg, brand.500 0%, brand.600 50%, brand.700 100%)"
            />

            <VStack spacing={8} align="stretch">
              {/* Logo 和标题 */}
              <VStack spacing={4} textAlign="center">
                <AppLogo />
                <Box>
                  <Heading
                    size="2xl"
                    bgGradient="linear(135deg, brand.600 0%, brand.700 100%)"
                    bgClip="text"
                    letterSpacing="-0.03em"
                  >
                    Account
                  </Heading>
                  <Text mt={2} fontSize="lg" color="gray.500" fontWeight="500">
                    个人财务管理
                  </Text>
                </Box>
              </VStack>

              {/* 表单 */}
              <VStack spacing={5}>
                <FormControl isRequired>
                  <FormLabel
                    color="gray.700"
                    fontWeight="600"
                    fontSize="sm"
                    mb={2}
                    ml={1}
                  >
                    邮箱地址
                  </FormLabel>
                  <Input
                    type="email"
                    placeholder="your@email.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    autoComplete="email"
                    size="lg"
                    _placeholder={{ color: 'gray.400' }}
                  />
                </FormControl>

                <FormControl isRequired>
                  <FormLabel
                    color="gray.700"
                    fontWeight="600"
                    fontSize="sm"
                    mb={2}
                    ml={1}
                  >
                    密码
                  </FormLabel>
                  <InputGroup size="lg">
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      autoComplete="current-password"
                      _placeholder={{ color: 'gray.400' }}
                    />
                    <InputRightElement h="full" pr={3}>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => setShowPassword(!showPassword)}
                        color="gray.500"
                        fontWeight="600"
                        h="full"
                        px={3}
                        _hover={{
                          color: 'brand.600',
                          bg: 'transparent',
                        }}
                      >
                        {showPassword ? '隐藏' : '显示'}
                      </Button>
                    </InputRightElement>
                  </InputGroup>
                </FormControl>

                {/* 登录按钮 */}
                <Button
                  type="submit"
                  size="lg"
                  width="full"
                  isLoading={isLoading}
                  loadingText="登录中"
                  mt={4}
                  bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
                  color="white"
                  _hover={{
                    bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
                    transform: 'translateY(-2px)',
                    boxShadow: 'lg',
                  }}
                  _active={{
                    transform: 'translateY(0)',
                  }}
                  transition="all 0.2s ease"
                >
                  登录
                </Button>
              </VStack>

              {/* 注册链接 */}
              <Flex justify="center" mt={2}>
                <Text color="gray.600" fontSize="md">
                  还没有账户？{' '}
                  <Link
                    as={RouterLink}
                    to="/register"
                    color="brand.600"
                    fontWeight="700"
                    _hover={{
                      color: 'brand.700',
                      textDecoration: 'none',
                    }}
                    transition="color 0.2s ease"
                  >
                    立即注册
                  </Link>
                </Text>
              </Flex>
            </VStack>
          </Box>
        </Container>
      </Box>
    </>
  );
}
