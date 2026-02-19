import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../data/models/import.dart';
import '../../../data/datasources/remote/api_service.dart';
import '../../../injection_container.dart' as di;

part 'import_event.dart';
part 'import_state.dart';

class ImportBloc extends Bloc<ImportEvent, ImportState> {
  final ApiService _apiService;

  ImportBloc({ApiService? apiService})
      : _apiService = apiService ?? di.sl<ApiService>(),
        super(const ImportState()) {
    on<ImportLoadSources>(_onLoadSources);
    on<ImportSourceSelected>(_onSourceSelected);
    on<ImportFileSelected>(_onFileSelected);
    on<ImportTransactionUpdated>(_onTransactionUpdated);
    on<ImportAllAccountsSet>(_onAllAccountsSet);
    on<ImportConfirmed>(_onImportConfirmed);
    on<ImportBackPressed>(_onBackPressed);
    on<ImportReset>(_onReset);
  }

  Future<void> _onLoadSources(
    ImportLoadSources event,
    Emitter<ImportState> emit,
  ) async {
    emit(state.copyWith(sourceStatus: ImportStepStatus.loading));

    try {
      final sources = await _apiService.getImportSources();
      emit(state.copyWith(
        sourceStatus: ImportStepStatus.success,
        availableSources: sources,
      ));
    } catch (e) {
      emit(state.copyWith(
        sourceStatus: ImportStepStatus.error,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onSourceSelected(
    ImportSourceSelected event,
    Emitter<ImportState> emit,
  ) async {
    emit(state.copyWith(
      selectedSource: event.source,
      currentScreen: ImportScreen.upload,
    ));
  }

  Future<void> _onFileSelected(
    ImportFileSelected event,
    Emitter<ImportState> emit,
  ) async {
    emit(state.copyWith(uploadStatus: ImportStepStatus.loading));

    try {
      final preview = await _apiService.uploadAndParseFile(
        state.selectedSource!,
        event.fileName,
        event.fileBytes,
      );

      // Auto-select first matching account for each transaction
      final updatedTransactions = preview.transactions.map((tx) {
        if (tx.canBeImported && tx.selectedAccountId == null) {
          // Check if we have suggestions
          if (preview.accountSuggestions != null && tx.accountName != null) {
            final suggestions = preview.accountSuggestions![tx.accountName!];
            if (suggestions != null && suggestions.isNotEmpty) {
              return tx.copyWith(selectedAccountId: suggestions.first.id);
            }
          }
          // If no suggestions, use first available account
          if (preview.categories.isNotEmpty) {
            // Try to find first account from somewhere
          }
        }
        return tx;
      }).toList();

      final updatedPreview = ImportPreview(
        jobId: preview.jobId,
        source: preview.source,
        totalRows: preview.totalRows,
        validRows: preview.validRows,
        duplicateRows: preview.duplicateRows,
        transactions: updatedTransactions,
        accountSuggestions: preview.accountSuggestions,
        categories: preview.categories,
      );

      emit(state.copyWith(
        uploadStatus: ImportStepStatus.success,
        preview: updatedPreview,
        currentScreen: ImportScreen.preview,
      ));
    } catch (e) {
      emit(state.copyWith(
        uploadStatus: ImportStepStatus.error,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onTransactionUpdated(
    ImportTransactionUpdated event,
    Emitter<ImportState> emit,
  ) async {
    if (state.preview == null) return;

    final updatedTransactions = List<ParsedTransaction>.from(state.preview!.transactions);
    updatedTransactions[event.index] = event.transaction;

    final updatedPreview = ImportPreview(
      jobId: state.preview!.jobId,
      source: state.preview!.source,
      totalRows: state.preview!.totalRows,
      validRows: state.preview!.validRows,
      duplicateRows: state.preview!.duplicateRows,
      transactions: updatedTransactions,
      accountSuggestions: state.preview!.accountSuggestions,
      categories: state.preview!.categories,
    );

    emit(state.copyWith(preview: updatedPreview));
  }

  Future<void> _onAllAccountsSet(
    ImportAllAccountsSet event,
    Emitter<ImportState> emit,
  ) async {
    if (state.preview == null) return;

    final updatedTransactions = state.preview!.transactions.map((tx) {
      if (tx.canBeImported) {
        return tx.copyWith(selectedAccountId: event.accountId);
      }
      return tx;
    }).toList();

    final updatedPreview = ImportPreview(
      jobId: state.preview!.jobId,
      source: state.preview!.source,
      totalRows: state.preview!.totalRows,
      validRows: state.preview!.validRows,
      duplicateRows: state.preview!.duplicateRows,
      transactions: updatedTransactions,
      accountSuggestions: state.preview!.accountSuggestions,
      categories: state.preview!.categories,
    );

    emit(state.copyWith(preview: updatedPreview));
  }

  Future<void> _onImportConfirmed(
    ImportConfirmed event,
    Emitter<ImportState> emit,
  ) async {
    emit(state.copyWith(importStatus: ImportStepStatus.loading));

    try {
      final result = await _apiService.executeImport(
        state.preview!.jobId,
        state.preview!.transactions,
      );

      emit(state.copyWith(
        importStatus: ImportStepStatus.success,
        result: result,
        currentScreen: ImportScreen.result,
      ));
    } catch (e) {
      emit(state.copyWith(
        importStatus: ImportStepStatus.error,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onBackPressed(
    ImportBackPressed event,
    Emitter<ImportState> emit,
  ) async {
    switch (state.currentScreen) {
      case ImportScreen.source:
        // Already at first screen, do nothing
        break;
      case ImportScreen.upload:
        emit(state.copyWith(currentScreen: ImportScreen.source));
        break;
      case ImportScreen.preview:
        emit(state.copyWith(currentScreen: ImportScreen.upload));
        break;
      case ImportScreen.result:
        emit(state.copyWith(currentScreen: ImportScreen.preview));
        break;
    }
  }

  Future<void> _onReset(
    ImportReset event,
    Emitter<ImportState> emit,
  ) async {
    emit(const ImportState());
  }
}
