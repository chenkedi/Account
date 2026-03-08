import { useEffect, useState, useMemo } from 'react';
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
  Spinner,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  IconButton,
  useToast,
  Collapse,
} from '@chakra-ui/react';
import {
  useCategories,
  useCategoriesState,
  useCategoryActions,
} from '../../../store';
import type { Category, CategoryType, CategoryCreateInput, CategoryUpdateInput } from '../../../data/models/category.model';
import { CategoryFormModal } from './CategoryFormModal';
import { DeleteCategoryModal } from './DeleteCategoryModal';

function PlusIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

function ChevronDownIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <svg
      viewBox="0 0 24 24"
      width="16"
      height="16"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      style={{ transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s' }}
    >
      <polyline points="6 9 12 15 18 9" />
    </svg>
  );
}

function FolderIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
    </svg>
  );
}

function TagIcon() {
  return (
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
      <line x1="7" y1="7" x2="7.01" y2="7" />
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

interface CategoryTreeItemProps {
  category: Category;
  allCategories: Category[];
  level: number;
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  expandedCategories: Set<string>;
  toggleExpanded: (id: string) => void;
}

function CategoryTreeItem({
  category,
  allCategories,
  level,
  onEdit,
  onDelete,
  expandedCategories,
  toggleExpanded,
}: CategoryTreeItemProps) {
  const children = useMemo(
    () => allCategories.filter((c) => c.parent_id === category.id),
    [allCategories, category.id]
  );

  const hasChildren = children.length > 0;
  const isExpanded = expandedCategories.has(category.id);

  return (
    <Box>
      <Flex
        align="center"
        justify="space-between"
        py={3}
        px={4}
        ml={level * 4}
        borderRadius="xl"
        _hover={{ bg: 'gray.50' }}
        transition="all 0.2s"
      >
        <Flex align="center" gap={3} flex={1}>
          {hasChildren ? (
            <IconButton
              aria-label={isExpanded ? '收起' : '展开'}
              icon={<ChevronDownIcon isOpen={isExpanded} />}
              size="xs"
              variant="ghost"
              onClick={() => toggleExpanded(category.id)}
            />
          ) : (
            <Box w="24px" />
          )}
          <Box color="brand.500">
            {hasChildren ? <FolderIcon /> : <TagIcon />}
          </Box>
          <Box>
            <Text fontWeight="600" color="gray.800">
              {category.name}
            </Text>
            {hasChildren && (
              <Text fontSize="xs" color="gray.500">
                {children.length} 个子分类
              </Text>
            )}
          </Box>
        </Flex>

        <HStack spacing={1}>
          <IconButton
            aria-label="编辑"
            icon={<EditIcon />}
            size="sm"
            variant="ghost"
            colorScheme="brand"
            onClick={() => onEdit(category)}
          />
          <IconButton
            aria-label="删除"
            icon={<DeleteIcon />}
            size="sm"
            variant="ghost"
            colorScheme="red"
            onClick={() => onDelete(category)}
          />
        </HStack>
      </Flex>

      <Collapse in={isExpanded}>
        {children.map((child) => (
          <CategoryTreeItem
            key={child.id}
            category={child}
            allCategories={allCategories}
            level={level + 1}
            onEdit={onEdit}
            onDelete={onDelete}
            expandedCategories={expandedCategories}
            toggleExpanded={toggleExpanded}
          />
        ))}
      </Collapse>
    </Box>
  );
}

export function CategoriesPage() {
  const categories = useCategories();
  const { isLoading, error } = useCategoriesState();
  const { fetchCategories, createCategory, updateCategory, deleteCategory } =
    useCategoryActions();

  const [activeTab, setActiveTab] = useState<CategoryType>('expense');
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set());

  const toast = useToast();

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  const filteredCategories = useMemo(() => {
    return categories.filter((c) => c.type === activeTab);
  }, [categories, activeTab]);

  const rootCategories = useMemo(() => {
    return filteredCategories.filter((c) => c.parent_id === null);
  }, [filteredCategories]);

  const toggleExpanded = (id: string) => {
    setExpandedCategories((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const handleAddCategory = () => {
    setSelectedCategory(null);
    setIsFormModalOpen(true);
  };

  const handleEditCategory = (category: Category) => {
    setSelectedCategory(category);
    setIsFormModalOpen(true);
  };

  const handleDeleteCategory = (category: Category) => {
    setSelectedCategory(category);
    setIsDeleteModalOpen(true);
  };

  const handleFormSubmit = async (data: CategoryCreateInput | CategoryUpdateInput) => {
    try {
      // 处理 parent_id: 空字符串转为 null
      const submitData = {
        ...data,
        parent_id: data.parent_id === '' ? null : data.parent_id,
      };

      if (selectedCategory) {
        await updateCategory(selectedCategory.id, submitData);
        toast({
          title: '分类已更新',
          status: 'success',
          duration: 3000,
        });
      } else {
        await createCategory(submitData as CategoryCreateInput);
        toast({
          title: '分类已创建',
          status: 'success',
          duration: 3000,
        });
      }
      setIsFormModalOpen(false);
    } catch (error) {
      toast({
        title: selectedCategory ? '更新失败' : '创建失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
      throw error;
    }
  };

  const handleDeleteConfirm = async () => {
    if (!selectedCategory) return;

    try {
      await deleteCategory(selectedCategory.id);
      toast({
        title: '分类已删除',
        status: 'success',
        duration: 3000,
      });
      setIsDeleteModalOpen(false);
    } catch (error) {
      toast({
        title: '删除失败',
        description: (error as Error).message,
        status: 'error',
        duration: 5000,
      });
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
            分类管理
          </Heading>
          <Text color="gray.500" fontWeight="500" mt={1}>
            管理收入和支出分类
          </Text>
        </Box>
        <Button
          size="lg"
          borderRadius="2xl"
          bgGradient="linear(135deg, brand.500 0%, brand.600 100%)"
          color="white"
          leftIcon={<PlusIcon />}
          fontWeight="700"
          onClick={handleAddCategory}
          _hover={{
            bgGradient: 'linear(135deg, brand.600 0%, brand.700 100%)',
            transform: 'translateY(-2px)',
            boxShadow: 'lg',
          }}
          transition="all 0.2s ease"
        >
          添加分类
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
      {isLoading && categories.length === 0 && (
        <Flex justify="center" py={20}>
          <Spinner size="xl" color="brand.500" />
        </Flex>
      )}

      {/* Tabs */}
      <Tabs
        variant="soft-rounded"
        colorScheme="brand"
        index={activeTab === 'expense' ? 0 : 1}
        onChange={(index) => setActiveTab(index === 0 ? 'expense' : 'income')}
      >
        <TabList>
          <Tab fontWeight="600">支出分类</Tab>
          <Tab fontWeight="600">收入分类</Tab>
        </TabList>

        <TabPanels>
          <TabPanel px={0}>
            <CategoryList
              categories={filteredCategories}
              rootCategories={rootCategories}
              isLoading={isLoading}
              onEdit={handleEditCategory}
              onDelete={handleDeleteCategory}
              expandedCategories={expandedCategories}
              toggleExpanded={toggleExpanded}
            />
          </TabPanel>
          <TabPanel px={0}>
            <CategoryList
              categories={filteredCategories}
              rootCategories={rootCategories}
              isLoading={isLoading}
              onEdit={handleEditCategory}
              onDelete={handleDeleteCategory}
              expandedCategories={expandedCategories}
              toggleExpanded={toggleExpanded}
            />
          </TabPanel>
        </TabPanels>
      </Tabs>

      <CategoryFormModal
        isOpen={isFormModalOpen}
        onClose={() => setIsFormModalOpen(false)}
        category={selectedCategory}
        defaultType={activeTab}
        allCategories={categories}
        onSubmit={handleFormSubmit}
        isLoading={isLoading}
      />

      <DeleteCategoryModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        category={selectedCategory}
        allCategories={categories}
        onDelete={handleDeleteConfirm}
        isLoading={isLoading}
      />
    </VStack>
  );
}

interface CategoryListProps {
  categories: Category[];
  rootCategories: Category[];
  isLoading: boolean;
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  expandedCategories: Set<string>;
  toggleExpanded: (id: string) => void;
}

function CategoryList({
  categories,
  rootCategories,
  isLoading,
  onEdit,
  onDelete,
  expandedCategories,
  toggleExpanded,
}: CategoryListProps) {
  if (isLoading && categories.length === 0) {
    return (
      <Flex justify="center" py={12}>
        <Spinner size="lg" color="brand.500" />
      </Flex>
    );
  }

  if (categories.length === 0) {
    return (
      <Card variant="elevated">
        <CardBody py={12}>
          <VStack spacing={4}>
            <Box color="gray.300">
              <FolderIcon />
            </Box>
            <Text fontSize="lg" fontWeight="600" color="gray.600">
              暂无分类
            </Text>
            <Text fontSize="sm" color="gray.500">
              点击上方按钮添加新分类
            </Text>
          </VStack>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card variant="elevated">
      <CardBody p={0}>
        {rootCategories.map((category) => (
          <CategoryTreeItem
            key={category.id}
            category={category}
            allCategories={categories}
            level={0}
            onEdit={onEdit}
            onDelete={onDelete}
            expandedCategories={expandedCategories}
            toggleExpanded={toggleExpanded}
          />
        ))}
      </CardBody>
    </Card>
  );
}

