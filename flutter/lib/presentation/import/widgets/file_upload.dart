import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../data/models/import.dart';
import '../bloc/import_bloc.dart';

class FileUpload extends StatelessWidget {
  const FileUpload({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        final isLoading = state.uploadStatus == ImportStepStatus.loading;

        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                '上传文件',
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 8),
              Text(
                '请选择要导入的${_getSourceName(state.selectedSource)}账单文件',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                ),
              ),
              const SizedBox(height: 32),
              Expanded(
                child: _FileUploadArea(
                  isLoading: isLoading,
                  onFileSelected: (filePath, fileName, fileBytes) {
                    context.read<ImportBloc>().add(
                      ImportFileSelected(
                        filePath: filePath,
                        fileName: fileName,
                        fileBytes: fileBytes,
                      ),
                    );
                  },
                ),
              ),
              if (state.uploadStatus == ImportStepStatus.error) ...[
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.errorContainer,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    state.errorMessage ?? '上传失败',
                    style: TextStyle(
                      color: Theme.of(context).colorScheme.onErrorContainer,
                    ),
                  ),
                ),
              ],
            ],
          ),
        );
      },
    );
  }

  String _getSourceName(ImportSource? source) {
    switch (source) {
      case ImportSource.alipay:
        return '支付宝';
      case ImportSource.wechat:
        return '微信';
      case ImportSource.bank:
        return '银行';
      case ImportSource.generic:
      default:
        return '';
    }
  }
}

class _FileUploadArea extends StatefulWidget {
  final bool isLoading;
  final Function(String filePath, String fileName, List<int> fileBytes) onFileSelected;

  const _FileUploadArea({
    required this.isLoading,
    required this.onFileSelected,
  });

  @override
  State<_FileUploadArea> createState() => _FileUploadAreaState();
}

class _FileUploadAreaState extends State<_FileUploadArea> {
  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Container(
      decoration: BoxDecoration(
        border: Border.all(
          color: colorScheme.outline,
          width: 2,
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: widget.isLoading ? null : _pickFile,
          borderRadius: BorderRadius.circular(16),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              if (widget.isLoading) ...[
                const CircularProgressIndicator(),
                const SizedBox(height: 16),
                Text(
                  '解析中...',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
              ] else ...[
                Icon(
                  Icons.cloud_upload_outlined,
                  size: 80,
                  color: colorScheme.primary,
                ),
                const SizedBox(height: 16),
                Text(
                  '点击选择文件',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
                const SizedBox(height: 8),
                Text(
                  '支持 CSV 格式',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(height: 24),
                FilledButton.icon(
                  onPressed: _pickFile,
                  icon: const Icon(Icons.file_open),
                  label: const Text('选择文件'),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _pickFile() async {
    // For now, show a placeholder - file picking would require
    // file_picker package integration
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('请集成 file_picker 包以选择文件'),
        ),
      );
    }

    // TODO: Integrate file_picker
    // Example code for file_picker:
    /*
    final result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['csv'],
    );

    if (result != null && result.files.single.bytes != null) {
      widget.onFileSelected(
        result.files.single.path ?? '',
        result.files.single.name,
        result.files.single.bytes!,
      );
    }
    */
  }
}
