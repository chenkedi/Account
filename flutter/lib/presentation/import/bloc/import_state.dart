part of 'import_bloc.dart';

enum ImportScreen { source, upload, preview, result }

enum ImportStepStatus { initial, loading, success, error }

class ImportState extends Equatable {
  final ImportScreen currentScreen;
  final ImportStepStatus sourceStatus;
  final ImportStepStatus uploadStatus;
  final ImportStepStatus previewStatus;
  final ImportStepStatus importStatus;

  final List<ImportSourceInfo>? availableSources;
  final ImportSource? selectedSource;
  final ImportPreview? preview;
  final ImportResult? result;
  final String? errorMessage;

  const ImportState({
    this.currentScreen = ImportScreen.source,
    this.sourceStatus = ImportStepStatus.initial,
    this.uploadStatus = ImportStepStatus.initial,
    this.previewStatus = ImportStepStatus.initial,
    this.importStatus = ImportStepStatus.initial,
    this.availableSources,
    this.selectedSource,
    this.preview,
    this.result,
    this.errorMessage,
  });

  ImportState copyWith({
    ImportScreen? currentScreen,
    ImportStepStatus? sourceStatus,
    ImportStepStatus? uploadStatus,
    ImportStepStatus? previewStatus,
    ImportStepStatus? importStatus,
    List<ImportSourceInfo>? availableSources,
    ImportSource? selectedSource,
    ImportPreview? preview,
    ImportResult? result,
    String? errorMessage,
  }) {
    return ImportState(
      currentScreen: currentScreen ?? this.currentScreen,
      sourceStatus: sourceStatus ?? this.sourceStatus,
      uploadStatus: uploadStatus ?? this.uploadStatus,
      previewStatus: previewStatus ?? this.previewStatus,
      importStatus: importStatus ?? this.importStatus,
      availableSources: availableSources ?? this.availableSources,
      selectedSource: selectedSource ?? this.selectedSource,
      preview: preview ?? this.preview,
      result: result ?? this.result,
      errorMessage: errorMessage,
    );
  }

  @override
  List<Object?> get props => [
    currentScreen,
    sourceStatus,
    uploadStatus,
    previewStatus,
    importStatus,
    availableSources,
    selectedSource,
    preview,
    result,
    errorMessage,
  ];
}
