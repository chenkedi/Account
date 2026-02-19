import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../injection_container.dart' as di;
import '../bloc/import_bloc.dart';
import '../widgets/source_selection.dart';
import '../widgets/file_upload.dart';
import '../widgets/import_preview.dart';
import '../widgets/import_result.dart';

class ImportPage extends StatelessWidget {
  const ImportPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => di.sl<ImportBloc>()..add(const ImportLoadSources()),
      child: const ImportView(),
    );
  }
}

class ImportView extends StatelessWidget {
  const ImportView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('导入账单'),
        leading: const _LeadingBackButton(),
      ),
      body: BlocBuilder<ImportBloc, ImportState>(
        builder: (context, state) {
          switch (state.currentScreen) {
            case ImportScreen.source:
              return const SourceSelection();
            case ImportScreen.upload:
              return const FileUpload();
            case ImportScreen.preview:
              return const ImportPreview();
            case ImportScreen.result:
              return const ImportResult();
          }
        },
      ),
    );
  }
}

class _LeadingBackButton extends StatelessWidget {
  const _LeadingBackButton();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ImportBloc, ImportState>(
      builder: (context, state) {
        if (state.currentScreen == ImportScreen.source) {
          return const SizedBox.shrink();
        }
        return IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            context.read<ImportBloc>().add(const ImportBackPressed());
          },
        );
      },
    );
  }
}
