part of 'import_bloc.dart';

abstract class ImportEvent extends Equatable {
  const ImportEvent();

  @override
  List<Object?> get props => [];
}

class ImportLoadSources extends ImportEvent {}

class ImportSourceSelected extends ImportEvent {
  final ImportSource source;

  const ImportSourceSelected(this.source);

  @override
  List<Object?> get props => [source];
}

class ImportFileSelected extends ImportEvent {
  final String filePath;
  final String fileName;
  final List<int> fileBytes;

  const ImportFileSelected({
    required this.filePath,
    required this.fileName,
    required this.fileBytes,
  });

  @override
  List<Object?> get props => [filePath, fileName, fileBytes];
}

class ImportTransactionUpdated extends ImportEvent {
  final int index;
  final ParsedTransaction transaction;

  const ImportTransactionUpdated({
    required this.index,
    required this.transaction,
  });

  @override
  List<Object?> get props => [index, transaction];
}

class ImportAllAccountsSet extends ImportEvent {
  final String accountId;

  const ImportAllAccountsSet(this.accountId);

  @override
  List<Object?> get props => [accountId];
}

class ImportConfirmed extends ImportEvent {}

class ImportBackPressed extends ImportEvent {}

class ImportReset extends ImportEvent {}
