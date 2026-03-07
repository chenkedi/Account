import {
  Box,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  SimpleGrid,
  Flex,
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import type { ImportSource } from '../../../data/models/import.model';

function AlipayIcon() {
  return (
    <svg viewBox="0 0 24 24" width="40" height="40" fill="currentColor">
      <path d="M19.5 3h-15A1.5 1.5 0 0 0 3 4.5v15A1.5 1.5 0 0 0 4.5 21h15a1.5 1.5 0 0 0 1.5-1.5v-15A1.5 1.5 0 0 0 19.5 3zm-2.5 14.5c-.3.1-.6.1-.9.2-.5.1-.9.2-1.4.3-.2.1-.4.1-.6.1-.3 0-.5-.1-.8-.2-.1 0-.1-.1-.2-.1-.6-.4-1.1-.8-1.6-1.3-.3-.3-.6-.7-.8-1.1 0-.1-.1-.1-.1-.2-.1-.2-.2-.3-.2-.5 0-.1.1-.2.1-.3.4-1.1.9-2.2 1.3-3.3.1-.2.2-.4.3-.6.1-.2.1-.4.1-.6 0-.3-.1-.6-.3-.8-.2-.2-.4-.3-.7-.4-.2-.1-.4-.1-.6 0h-2c-.1 0-.2 0-.3.1-.2.1-.3.2-.3.4 0 .1 0 .2.1.3.3.5.6 1 .8 1.5.2.4.4.8.6 1.2.1.2.2.3.2.5 0 .2 0 .4-.1.6-.4 1.1-.8 2.1-1.2 3.2-.1.3-.2.5-.3.7-.1.2-.1.4-.1.6 0 .3.1.6.3.8.2.2.4.3.7.4.5.2 1 .4 1.5.5.3.1.6.1.9.2h.3c.4.1.8.2 1.2.3.2.1.4.1.6.1.4 0 .7-.1 1-.3.2-.1.3-.3.3-.5 0-.1 0-.2-.1-.3-.1-.1-.2-.3-.3-.4z" />
    </svg>
  );
}

function WechatIcon() {
  return (
    <svg viewBox="0 0 24 24" width="40" height="40" fill="currentColor">
      <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.045c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z" />
    </svg>
  );
}

function BankIcon() {
  return (
    <svg viewBox="0 0 24 24" width="40" height="40" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 21h18" />
      <path d="M3 10h18" />
      <path d="M5 6l7-3 7 3" />
      <path d="M4 10v11" />
      <path d="M20 10v11" />
      <path d="M8 10v11" />
      <path d="M12 10v11" />
      <path d="M16 10v11" />
    </svg>
  );
}

function FileIcon() {
  return (
    <svg viewBox="0 0 24 24" width="40" height="40" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
      <polyline points="14 2 14 8 20 8" />
      <line x1="16" y1="13" x2="8" y2="13" />
      <line x1="16" y1="17" x2="8" y2="17" />
      <polyline points="10 9 9 9 8 9" />
    </svg>
  );
}

interface ImportSourceCardProps {
  icon: React.ReactNode;
  name: string;
  description: string;
  colorScheme: 'blue' | 'green' | 'purple' | 'orange';
  source: ImportSource;
}

function ImportSourceCard({
  icon,
  name,
  description,
  colorScheme,
  source,
}: ImportSourceCardProps) {
  const navigate = useNavigate();

  const colorMap = {
    blue: { bg: 'brand.50', color: 'brand.600', border: 'brand.100' },
    green: { bg: 'income.50', color: 'income.600', border: 'income.100' },
    purple: { bg: 'purple.50', color: 'purple.600', border: 'purple.100' },
    orange: { bg: 'transfer.50', color: 'transfer.600', border: 'transfer.100' },
  };

  const colors = colorMap[colorScheme as keyof typeof colorMap] || colorMap.blue;

  const handleClick = () => {
    navigate(`/import/upload?source=${source}`);
  };

  return (
    <Card
      variant="elevated"
      cursor="pointer"
      onClick={handleClick}
      overflow="hidden"
      _hover={{
        transform: 'translateY(-4px)',
        boxShadow: 'lg',
      }}
      transition="all 0.3s ease"
    >
      <CardBody p={6}>
        <VStack spacing={4}>
          <Box
            w={16}
            h={16}
            borderRadius="2xl"
            bg={colors.bg}
            color={colors.color}
            borderWidth="2px"
            borderColor={colors.border}
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            {icon}
          </Box>
          <VStack spacing={1} textAlign="center">
            <Text fontSize="xl" fontWeight="800" color="gray.800">
              {name}
            </Text>
            <Text fontSize="sm" color="gray.500" fontWeight="500">
              {description}
            </Text>
          </VStack>
        </VStack>
      </CardBody>
    </Card>
  );
}

export function ImportPage() {
  return (
    <VStack spacing={6} align="stretch">
      {/* 头部 */}
      <Box>
        <Heading
          size="2xl"
          bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
          bgClip="text"
          letterSpacing="-0.03em"
        >
          导入账单
        </Heading>
        <Text color="gray.500" fontWeight="500" mt={1}>
          从支付宝、微信或银行导入交易记录
        </Text>
      </Box>

      {/* 导入来源 */}
      <SimpleGrid columns={{ base: 1, sm: 2 }} spacing={5}>
        <ImportSourceCard
          icon={<AlipayIcon />}
          name="支付宝"
          description="支持支付宝 CSV 账单"
          colorScheme="blue"
          source="alipay"
        />
        <ImportSourceCard
          icon={<WechatIcon />}
          name="微信支付"
          description="支持微信支付 CSV 账单"
          colorScheme="green"
          source="wechat"
        />
        <ImportSourceCard
          icon={<BankIcon />}
          name="银行对账单"
          description="支持银行 CSV 格式"
          colorScheme="blue"
          source="bank"
        />
        <ImportSourceCard
          icon={<FileIcon />}
          name="通用 CSV"
          description="自定义格式导入"
          colorScheme="orange"
          source="generic"
        />
      </SimpleGrid>

      {/* 导入说明 */}
      <Card variant="elevated">
        <CardBody p={6}>
          <VStack spacing={4} align="stretch">
            <Flex align="center" gap={3}>
              <Box
                w={10}
                h={10}
                borderRadius="xl"
                bg="brand.50"
                color="brand.600"
                display="flex"
                alignItems="center"
                justifyContent="center"
              >
                <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="16" x2="12" y2="12" />
                  <line x1="12" y1="8" x2="12.01" y2="8" />
                </svg>
              </Box>
              <Text fontSize="lg" fontWeight="800" color="gray.800">
                导入说明
              </Text>
            </Flex>
            <VStack spacing={3} align="stretch" pl={1}>
              {[
                '选择导入来源',
                '上传对应格式的账单文件',
                '预览并确认导入数据',
                '完成导入',
              ].map((step, index) => (
                <Flex key={index} align="center" gap={4}>
                  <Box
                    w={8}
                    h={8}
                    borderRadius="full"
                    bg="gray.100"
                    color="gray.600"
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    fontSize="sm"
                    fontWeight="800"
                  >
                    {index + 1}
                  </Box>
                  <Text color="gray.600" fontWeight="500">
                    {step}
                  </Text>
                </Flex>
              ))}
            </VStack>
          </VStack>
        </CardBody>
      </Card>
    </VStack>
  );
}
