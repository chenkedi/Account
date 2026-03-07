import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Card,
  CardBody,
  useToast,
  Progress,
} from '@chakra-ui/react';
import { useState, useCallback } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDropzone } from 'react-dropzone';
import type { ImportSource } from '../../../data/models/import.model';
import { useImportActions, useImportState } from '../../../store';

function UploadIcon() {
  return (
    <svg viewBox="0 0 24 24" width="80" height="80" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="17 8 12 3 7 8" />
      <line x1="12" y1="3" x2="12" y2="15" />
    </svg>
  );
}

const sourceNames: Record<ImportSource, string> = {
  alipay: '支付宝',
  wechat: '微信支付',
  bank: '银行对账单',
  generic: '通用 CSV',
};

export function ImportUploadPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const toast = useToast();

  const source = searchParams.get('source') as ImportSource | null;
  const { isLoading } = useImportState();
  const { uploadFile, selectSource } = useImportActions();

  const [uploadProgress, setUploadProgress] = useState(0);

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      if (!source) {
        toast({
          title: '错误',
          description: '请先选择导入来源',
          status: 'error',
          duration: 3000,
          isClosable: true,
        });
        navigate('/import');
        return;
      }

      const file = acceptedFiles[0];
      if (!file) return;

      try {
        selectSource(source);
        setUploadProgress(30);

        const preview = await uploadFile(source, file);
        setUploadProgress(100);

        toast({
          title: '解析成功',
          description: `成功解析 ${preview.totalRecords} 条记录`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });

        navigate('/import/preview');
      } catch (error) {
        setUploadProgress(0);
        toast({
          title: '解析失败',
          description: (error as Error).message,
          status: 'error',
          duration: 5000,
          isClosable: true,
        });
      }
    },
    [source, selectSource, uploadFile, navigate, toast]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv'],
    },
    maxFiles: 1,
    disabled: isLoading,
  });

  if (!source) {
    return (
      <VStack spacing={6} align="stretch">
        <Box>
          <Heading
            size="2xl"
            bgGradient="linear(135deg, gray.800 0%, gray.900 100%)"
            bgClip="text"
            letterSpacing="-0.03em"
          >
            上传文件
          </Heading>
        </Box>
        <Card variant="elevated">
          <CardBody py={12}>
            <VStack spacing={6}>
              <Text color="gray.500">
                请先选择导入来源
              </Text>
              <Button
                onClick={() => navigate('/import')}
                size="lg"
                bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
                color="white"
              >
                返回选择来源
              </Button>
            </VStack>
          </CardBody>
        </Card>
      </VStack>
    );
  }

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
          上传文件
        </Heading>
        <Text color="gray.500" fontWeight="500" mt={1}>
          {sourceNames[source]} - 请上传 CSV 文件
        </Text>
      </Box>

      {/* 上传区域 */}
      <Card variant="elevated">
        <CardBody py={12}>
          <VStack spacing={6}>
            <Box
              {...getRootProps()}
              w="full"
              p={12}
              border="2px dashed"
              borderColor={isDragActive ? 'brand.400' : 'gray.200'}
              borderRadius="2xl"
              bg={isDragActive ? 'brand.50' : 'gray.50'}
              cursor={isLoading ? 'not-allowed' : 'pointer'}
              transition="all 0.2s ease"
              _hover={{
                borderColor: isLoading ? 'gray.200' : 'brand.400',
                bg: isLoading ? 'gray.50' : 'brand.50',
              }}
            >
              <input {...getInputProps()} />
              <VStack spacing={4}>
                <Box color={isDragActive ? 'brand.500' : 'gray.300'}>
                  <UploadIcon />
                </Box>
                <VStack spacing={2}>
                  <Text fontSize="xl" fontWeight="700" color="gray.700">
                    {isDragActive ? '释放文件开始上传' : '点击或拖拽文件到这里'}
                  </Text>
                  <Text color="gray.500" textAlign="center" maxW="sm">
                    支持 CSV 格式文件
                  </Text>
                </VStack>
              </VStack>
            </Box>

            {isLoading && (
              <VStack w="full" spacing={4}>
                <Text color="gray.600" fontWeight="600">
                  解析中...
                </Text>
                <Progress
                  value={uploadProgress}
                  size="lg"
                  w="full"
                  colorScheme="brand"
                  borderRadius="full"
                />
              </VStack>
            )}

            <Button
              onClick={() => navigate('/import')}
              size="lg"
              variant="outline"
              borderWidth="2px"
            >
              返回选择来源
            </Button>
          </VStack>
        </CardBody>
      </Card>
    </VStack>
  );
}
